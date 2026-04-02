package api

import (
	"net/http"

	"github.com/DevitoDbug/portfolio/internals/web/pages"
)

func (a *Api) IndexHandler(w http.ResponseWriter, r *http.Request) {
	_ = pages.IndexPage().Render(r.Context(), w)
}
