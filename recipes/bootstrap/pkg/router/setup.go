package router

import (
	"go.khulnasoft.com/velocity"
)

func InstallRouter(app *fiber.App) {
	setup(app, NewApiRouter(), NewHttpRouter())
}

func setup(app *fiber.App, router ...Router) {
	for _, r := range router {
		r.InstallRouter(app)
	}
}
