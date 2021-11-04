package httpext

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/lpar/gzipped/v2"
	"github.com/sirupsen/logrus"
)

func ServeDirWithCompression(r chi.Router, path string, root gzipped.Dir) {
	filePath := path
	if strings.HasSuffix(filePath, "*") {
		filePath = filePath[:len(filePath)-1]
	}
	fs := http.StripPrefix(filePath, gzipped.FileServer(root))

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, ".js") {
			w.Header().Add("Content-Type", "application/javascript")
		}

		fs.ServeHTTP(w, r)
	})
}

func ServeDir(r chi.Router, path string, root http.Dir) {
	filePath := path
	if strings.HasSuffix(filePath, "*") {
		filePath = filePath[:len(filePath)-1]
	}
	fs := http.StripPrefix(filePath, http.FileServer(root))

	logrus.Infof("%+v", root)

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, ".js") {
			w.Header().Add("Content-Type", "application/javascript")
		}

		fs.ServeHTTP(w, r)
	})
}

func ServeFile(r chi.Router, path string, filePath string) {
	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filePath)
	})
}
