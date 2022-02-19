package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/fvumbaca/mdrun/browser"
	"github.com/fvumbaca/mdrun/rundoc"
	_ "github.com/fvumbaca/mdrun/rundoc"
	"github.com/fvumbaca/mdrun/static"
)

func main() {
	rootFS := os.DirFS("./")

	docHandler := rundoc.NewHandler(rootFS)

	browserHandler := browser.Handler{
		RootFS:      rootFS,
		FileHandler: docHandler,
	}

	fmt.Println("Starting up server on :3000")
	http.Handle("/-/static/", static.Static("/-/static"))
	http.Handle("/", &browserHandler)
	http.ListenAndServe(":3000", nil)
}
