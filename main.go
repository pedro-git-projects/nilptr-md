package main

import (
	"embed"
	_ "embed"
	"log"
	"os"

	"github.com/pedro-git-projects/nilptr-md/app"
)

//go:embed assets
var assetsFS embed.FS

//go:embed templates
var tmplFS embed.FS

func main() {
	app := app.New(
		assetsFS,
		tmplFS,
		log.New(os.Stdout, "[nilpt-md]", log.LstdFlags),
	)
	app.Run()
}
