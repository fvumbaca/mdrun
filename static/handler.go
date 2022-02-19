package static

import (
	"embed"
	"net/http"
)

//go:embed js/* css/*
var staticFS embed.FS

func Static(prefixPath string) http.Handler {
	return http.StripPrefix(prefixPath, http.FileServer(http.FS(staticFS)))
}
