package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// Handler is a default handler to serve up
// a home page.
func Handler(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("index.html"))
}

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("home.html"))
}
