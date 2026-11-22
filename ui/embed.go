package ui

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed dist/*
var distFS embed.FS

// Handler returns an http.Handler that serves the embedded Vue 3 SPA with fallback to index.html.
func Handler() http.Handler {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		return http.NotFoundHandler()
	}

	fileServer := http.FileServer(http.FS(sub))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		// Try opening the requested file
		f, err := sub.Open(path)
		if err == nil {
			_ = f.Close()
			fileServer.ServeHTTP(w, r)
			return
		}

		// SPA Fallback: serve index.html for client-side routing
		indexFile, err := sub.Open("index.html")
		if err == nil {
			defer indexFile.Close()
			stat, _ := indexFile.Stat()
			http.ServeContent(w, r, "index.html", stat.ModTime(), indexFile.(interface {
				ReadSeeker
			}))
			return
		}

		http.NotFound(w, r)
	})
}

type ReadSeeker interface {
	Read(p []byte) (n int, err error)
	Seek(offset int64, whence int) (int64, error)
}
