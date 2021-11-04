package swagger

import (
	"net/http"
	"strings"

	httpSwagger "github.com/swaggo/http-swagger"
)

func WrapSwagger(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, ".js") {
		// will not be overwritten later in current version of package
		w.Header().Add("Content-Type", "application/javascript")
	}

	httpSwagger.WrapHandler.ServeHTTP(w, r)
}
