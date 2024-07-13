package render

import (
	"context"
	"net/http"
	"time"

	"github.com/Glowman554/WebStarterGO/templates/components"
	"github.com/a-h/templ"
)

type Component func() templ.Component
type RequestComponent func(r *http.Request, w http.ResponseWriter) templ.Component

type Responder func(http.ResponseWriter, *http.Request)

func withContext(f func(context.Context) error) error {
	context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return f(context)
}

func ApplyLayout(component Component, title string) Component {
	return func() templ.Component {
		return components.Layout(title, component())
	}
}

func AsHandler(component Component) Responder {
	return AsRequestHandler(func(r *http.Request, w http.ResponseWriter) templ.Component {
		return component()
	})
}

func AsRequestHandler(component RequestComponent) Responder {
	return func(w http.ResponseWriter, r *http.Request) {
		err := withContext(func(ctx context.Context) error {
			return component(r, w).Render(ctx, w)
		})
		if err != nil {
			w.Write([]byte(err.Error()))
		}
	}
}
