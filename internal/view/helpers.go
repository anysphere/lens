package view

import (
	"context"
	"errors"

	"github.com/anysphere/lens/internal"
	"github.com/anysphere/lens/internal/model"
	"github.com/anysphere/lens/internal/ui"
)

func extractApp(ctx context.Context) (*App, error) {
	app, ok := ctx.Value(internal.KeyApp).(*App)
	if !ok {
		return nil, errors.New("No application found in context")
	}

	return app, nil
}

func describeResource(app *App, m ui.Tabular, resource, path string) {
	v := NewLiveView(app, "Describe", model.NewDescribe(resource, path))
	if err := app.inject(v); err != nil {
		app.Flash().Err(err)
	}
}
