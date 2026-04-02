package api

import (
	"net/http"

	"github.com/DevitoDbug/portfolio/internals/web/pages"
)

func (a *Api) AboutHandler(w http.ResponseWriter, r *http.Request) {
	_ = pages.AboutPage().Render(r.Context(), w)
}
