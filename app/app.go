package app

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/pedro-git-projects/nilptr-md/httpext"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	DefaultPort      = "8080"
	StaticPrefix     = "/static"
	TemplatesPattern = "templates/*.tmpl"
	StaticAssetsDir  = "assets/static"
	PagesAssetsDir   = "assets/pages"
	BaseTemplateName = "base.tmpl"
)

type TemplateData struct {
	Content template.HTML
	Title   string
	Now     time.Time
}

type App struct {
	server *http.Server
	logger *log.Logger
	router http.Handler
	pages  fs.FS
	md     goldmark.Markdown
	tmpl   *template.Template
}

func New(assetsFS embed.FS, templatesFS embed.FS, logger *log.Logger) *App {
	if logger == nil {
		logger = log.New(os.Stdout, "[app] ", log.LstdFlags)
	}

	tmpl, err := template.ParseFS(templatesFS, TemplatesPattern)
	if err != nil {
		logger.Fatalf("parsing templates: %v", err)
	}

	md := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)

	// Sub-tree for static assets
	staticSub, err := fs.Sub(assetsFS, StaticAssetsDir)
	if err != nil {
		logger.Fatalf("static assets directory %q not found: %v", StaticAssetsDir, err)
	}

	// Sub-tree for markdown pages
	pagesSub, err := fs.Sub(assetsFS, PagesAssetsDir)
	if err != nil {
		logger.Fatalf("pages assets directory %q not found: %v", PagesAssetsDir, err)
	}

	r := chi.NewRouter()
	r.Use(loggingMiddleware(logger))

	r.Handle(StaticPrefix+"/*", http.StripPrefix(StaticPrefix, http.FileServer(http.FS(staticSub))))

	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		handlePage(w, r, pagesSub, md, tmpl, logger)
	})

	addr := getPort()
	srv := &http.Server{
		Addr:           addr,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return &App{
		server: srv,
		logger: logger,
		router: r,
		pages:  pagesSub,
		md:     md,
		tmpl:   tmpl,
	}
}

func getPort() string {
	if p := os.Getenv("PORT"); p != "" {
		return ":" + p
	}
	return ":" + DefaultPort
}

func loggingMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Printf("%s %s completed in %v", r.Method, r.URL.Path, time.Since(start))
		})
	}
}

func handlePage(w http.ResponseWriter, r *http.Request, pages fs.FS, md goldmark.Markdown, tmpl *template.Template, logger *log.Logger) {
	p := r.URL.Path
	if p == "/" {
		p = "index"
	} else {
		p = strings.Trim(p, "/")
	}
	file := path.Clean(p + ".md")

	data, err := fs.ReadFile(pages, file)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var buf strings.Builder
	if err := md.Convert(data, &buf); err != nil {
		http.Error(w, "render error", http.StatusInternalServerError)
		logger.Printf("markdown convert error: %v", err)
		return
	}

	resp := TemplateData{
		Content: template.HTML(buf.String()),
		Title:   cases.Title(language.Und).String(path.Base(p)),
		Now:     time.Now(),
	}
	httpext.ContentType.Add(w.Header(), httpext.TextHTML)
	if err := tmpl.ExecuteTemplate(w, BaseTemplateName, resp); err != nil {
		http.Error(w, "template exec error", http.StatusInternalServerError)
		logger.Printf("template exec error: %v", err)
	}
}

func (a *App) Run() error {
	a.logger.Printf("Listening on %s", a.server.Addr)
	return a.server.ListenAndServe()
}
