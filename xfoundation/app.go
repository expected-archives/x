package xfoundation

import (
	"context"
	"github.com/caumette-co/x/xfoundation/contracts"
	"github.com/caumette-co/x/xfoundation/template/gohtml"
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

		Env    Environment
		Logger *zap.Logger

		Providers []any

		startHooks []func(ctx context.Context) error
		stopHooks  []func(ctx context.Context) error
	}

	AppHandler interface {
		OnStartup(app *App) error
	}

	// AppHandlerOut is used to inject a AppHandler
	AppHandlerOut struct {
		dig.Out

		Handler AppHandler `group:"x.app_handler"`
	}

	Environment string
)

const (
	EnvironmentProduction  Environment = "production"
	EnvironmentDevelopment Environment = "development"
)

func (app *App) Run() {
	app.Container = dig.New()

	if app.Env == "" {
		app.Env = EnvironmentDevelopment
	}

	if app.Logger == nil {
		if app.Env == EnvironmentProduction {
			app.Logger = lo.Must(zap.NewProduction())
		} else {
			app.Logger = lo.Must(zap.NewDevelopment())
		}
	}

	zap.ReplaceGlobals(app.Logger)

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

	app.provideDefaultTemplateEngines()
	app.invokeAppHandlers()

	ctx, canceler := context.WithCancel(context.Background())

	app.callStartHooks(ctx)

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-s

	app.callStopHooks(ctx)

	canceler()
	app.Logger.Info("app stopped")
}

func (app *App) callStopHooks(ctx context.Context) {
	for _, hook := range app.stopHooks {
		if err := hook(ctx); err != nil {
			app.Logger.Fatal("failed to stop hook", zap.Error(err))
		}
	}
}

func (app *App) callStartHooks(ctx context.Context) {
	for _, hook := range app.startHooks {
		if err := hook(ctx); err != nil {
			app.Logger.Fatal("failed to start hook", zap.Error(err))
		}
	}
}

func (app *App) invokeAppHandlers() {
	type appHandlers struct {
		dig.In
		Handlers []AppHandler `group:"x.app_handler"`
	}

	_, err := app.Invoke(func(handlers appHandlers) error {
		app.Logger.Info("invoking app handlers", zap.Int("count", len(handlers.Handlers)))
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

func (app *App) provideDefaultTemplateEngines() {
	type templateEngines struct {
		dig.In
		Engines []contracts.TemplateEngine `group:"x.template_engine"`
	}

	_, err := app.Invoke(func(engines templateEngines) error {
		if len(engines.Engines) == 0 {
			app.Logger.Info("no template engine provided, using default gohtml")
			return app.Provide(gohtml.New(gohtml.Config{
				Folder:         "views",
				LayoutsFolder:  "layouts",
				PartialsFolder: "partials",
			}))
		}

		return nil
	})

	if err != nil {
		app.Logger.Fatal("failed to provide default template engine gohtml", zap.Error(err))
	}
}

func ProvideSingleValueFunc[T any](v T) func() T {
	return func() T { return v }
}

func NewAppHandlerOut(h AppHandler) AppHandlerOut {
	return AppHandlerOut{Handler: h}
}
