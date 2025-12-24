package internal

import "context"

type App struct{}

func New(cfg Config) (*App, error) {
	return nil, nil
}

func (a *App) Run(ctx context.Context) error {
	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	return nil
}
