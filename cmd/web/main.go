package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/cozycole/thesketchdb/internal/web"
)

func main() {
	addr := flag.String("addr", "localhost:8080", "HTTP network address")
	debug := flag.Bool("debug", false, "debug mode")
	dev := flag.Bool("dev", false, "use dev config and set debug true")
	serveStatic := flag.Bool("serve-static", false, "serve css, js and images")

	flag.Parse()

	err := godotenv.Load()
	if err != nil && *dev {
		log.Fatal("Error loading .env file")
	}

	config := loadConfig(*dev, *debug)
	deps, cleanup, err := newDependencies(config)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	app := web.NewApplication(config, deps)
	configureOAuth(config.Origin)

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: app.ErrorLog,
		Handler:  app.Routes("./ui/static/", *serveStatic),
	}

	app.InfoLog.Printf("Starting server on %s with origin %s", *addr, config.Origin)
	err = srv.ListenAndServe()
	app.ErrorLog.Fatal(err)
}
