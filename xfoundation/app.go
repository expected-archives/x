package xfoundation

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

type (
	Provider interface {
		Register(app *App) error
	}

	Scopeable[T any] interface {
		Scope() (T, error)
	}

	App struct {
		Env        string
		Logger     *zap.Logger
		Providers  []Provider
		startHooks []func(ctx context.Context) error
		stopHooks  []func(ctx context.Context) error
		values     map[reflect.Type]reflect.Value
	}
)

const (
	AppEnvProduction  = "production"
	AppEnvDevelopment = "development"
)

func (app *App) Run() {
	app.values = make(map[reflect.Type]reflect.Value)
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
	app.values[reflect.TypeOf(v)] = reflect.ValueOf(v)
}

func (app *App) Invoke(f any) ([]reflect.Value, error) {
	fType := reflect.TypeOf(f)
	if fType.Kind() != reflect.Func {
		return nil, fmt.Errorf("app.Invoke: invalid func type %v", fType)
	}
	var dependencies []reflect.Value
	for i := 0; i < fType.NumIn(); i++ {
		depType := fType.In(i)
		value, ok := app.values[depType]
		if !ok {
			return nil, fmt.Errorf("app.Invoke: cannot find dependency %v", depType)
		}
		dependencies = append(dependencies, value)
	}

	returnValues := reflect.ValueOf(f).Call(dependencies)
	if returnValuesLen := len(returnValues); returnValuesLen > 0 {
		if err, ok := returnValues[returnValuesLen-1].Interface().(error); ok && err != nil {
			return nil, err
		}
		returnValues = returnValues[:returnValuesLen-1]
	}
	return returnValues, nil
}

func panicOnError[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}
