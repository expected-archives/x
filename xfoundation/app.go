package xfoundation

import (
	"context"
	"github.com/expectedsh/dig"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

type (
	App struct {
		Container *dig.Container

		Env    AppEnv
		Logger *zap.Logger

		Providers []any

		startHooks []func(ctx context.Context) error
		stopHooks  []func(ctx context.Context) error
	}

	AppHandler interface {
		OnStartup(app *App) error
	}

	AppHandlerOut struct {
		dig.Out

		Handler AppHandler `group:"x.startup_handler"`
	}

	AppEnv string
)

const (
	AppEnvProduction  AppEnv = "production"
	AppEnvDevelopment AppEnv = "development"
)

func (app *App) Run() {
	app.Container = dig.New()

	if app.Env == "" {
		app.Env = AppEnvDevelopment
	}

	if app.Logger == nil {
		if app.Env == AppEnvProduction {
			app.Logger = lo.Must(zap.NewProduction())
		} else {
			app.Logger = lo.Must(zap.NewDevelopment())
		}
	}

	app.Logger.Info("app started", zap.String("env", string(app.Env)))

	lo.Must0(app.Provide(ProvideSingleValueFunc(app.Logger)))
	lo.Must0(app.Provide(ProvideSingleValueFunc(app.Env)))
	lo.Must0(app.Provide(ProvideSingleValueFunc(app)))

	for _, provider := range app.Providers {
		log := app.Logger.With(zap.Any("provider", reflect.TypeOf(provider)))
		if err := app.Provide(provider); err != nil {
			log.Fatal("failed to register provider", zap.Error(err))
		}
		log.Info("provider registered")
	}

	type appHandlers struct {
		dig.In
		Handlers []AppHandler `group:"x.startup_handler"`
	}

	_, err := app.Invoke(func(handlers appHandlers) error {
		for _, handler := range handlers.Handlers {
			if err := handler.OnStartup(app); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		app.Logger.Fatal("failed to start app", zap.Error(err))
	}

	ctx := context.Background()
	for _, hook := range app.startHooks {
		if err := hook(ctx); err != nil {
			app.Logger.Fatal("failed to start hook", zap.Error(err))
		}
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-s

	for _, hook := range app.stopHooks {
		if err := hook(ctx); err != nil {
			app.Logger.Fatal("failed to stop hook", zap.Error(err))
		}
	}

	app.Logger.Info("app stopped")
}

func (app *App) OnStart(hook func(ctx context.Context) error) {
	app.startHooks = append(app.startHooks, hook)
}

func (app *App) OnStop(hook func(ctx context.Context) error) {
	app.stopHooks = append(app.stopHooks, hook)
}

func (app *App) Provide(v any) error {
	return app.Container.Provide(v)
}

func (a *App) Invoke(f any) ([]reflect.Value, error) {
	invokeInfo := dig.InvokeInfo{}

	err := a.Container.Invoke(f, dig.FillInvokeInfo(&invokeInfo))
	if err != nil {
		return nil, err
	}

	return invokeInfo.Outputs, nil
}

func ProvideSingleValueFunc[T any](v T) func() T {
	return func() T { return v }
}

func NewAppHandlerOut(h AppHandler) AppHandlerOut {
	return AppHandlerOut{Handler: h}
}
