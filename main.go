package main

/*
This is going to be sick!

Features:

- Support metadata at the top of the file
- Support required cli-tools to throw errors early if it can even run on the
system
- Support inline requests
- Support out-of-the-box containerization
- Interactive mode will print out runbook and prompt before each exec
- Script mode that will execute mardown as a script
- Infer language dependencies from code blocks

Example usage:

	$ mdrun my-runbook.md --set SOME_VAR=hello

	# -E pulls from the environment
	$ mdrun -E my-runbook.md
*/

import (
	// Keeping here for now to include tests while developing

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
