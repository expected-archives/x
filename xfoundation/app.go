package xfoundation

import (
	"context"
	"go.uber.org/dig"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

type (
	Provider interface {
		Register(app *App) error
	}

	App struct {
		Env        string
		Logger     *zap.Logger
		Providers  []Provider
		container  *dig.Container
		startHooks []func(ctx context.Context) error
		stopHooks  []func(ctx context.Context) error
	}
)

const (
	AppEnvProduction  = "production"
	AppEnvDevelopment = "development"
)

func (app *App) Run() {
	app.container = dig.New()
	if app.Env == "" {
		app.Env = AppEnvDevelopment
	}

	if app.Logger == nil {
		if app.Env == AppEnvProduction {
			app.Logger = panicOnError(zap.NewProduction())
		} else {
			app.Logger = panicOnError(zap.NewDevelopment())
		}
	}
	app.Provide(app.Logger)

	for _, provider := range app.Providers {
		log := app.Logger.With(zap.Any("provider", provider))
		if err := provider.Register(app); err != nil {
			log.Fatal("failed to register provider", zap.Error(err))
		}
		log.Info("provider registered")
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
}

func (app *App) OnStart(hook func(ctx context.Context) error) {
	app.startHooks = append(app.startHooks, hook)
}

func (app *App) OnStop(hook func(ctx context.Context) error) {
	app.stopHooks = append(app.stopHooks, hook)
}

func (app *App) Provide(v any) {
	app.container.
		app.container.Provide(func() any {
		return v
	})
}

func panicOnError[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}
