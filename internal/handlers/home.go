package handlers

import (
	"context"
	"net/http"

	"ikurotime/code-engine/templates"
)

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	h.logger.Printf("%s %s", r.Method, r.URL.Path)
	templ := templates.Layout(templates.Home("World"))
	templ.Render(context.Background(), w)
}
