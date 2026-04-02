package api

import (
	"net/http"

	"github.com/DevitoDbug/portfolio/internals/web/pages/index"
)

func (a *Api) IndexHandler(w http.ResponseWriter, r *http.Request) {
	_ = index.IndexPage().Render(r.Context(), w)
}
