package main

import (
	"encoding/gob"
	"flag"
	"log"
	"net/http"

	"github.com/joho/godotenv"
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

	app := NewApplication(config, deps)
	configureOAuth(config.Origin)

	app.InfoLog.Println("ORIGIN: ", config.Origin)
	app.InfoLog.Printf("css: %s js: %s\n", app.Assets["css"], app.Assets["js"])

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: app.ErrorLog,
		Handler:  app.routes("./ui/static/", *serveStatic),
	}

	app.InfoLog.Println("Starting server on", *addr)
	err = srv.ListenAndServe()
	app.ErrorLog.Fatal(err)
}

func init() {
	// Register the struct so encoding/gob knows how to handle it
	gob.Register(flashMessage{})
	gob.Register(PendingRegistration{})
}
