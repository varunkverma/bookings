package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/varunkverma/bookings/pkg/config"
	"github.com/varunkverma/bookings/pkg/handlers"
	"github.com/varunkverma/bookings/pkg/render"
)

const PORT = ":3000"

var appConfig config.AppConfig

var session *scs.SessionManager

// main is the main entry point of the application
func main() {

	// Change this to true when in production
	appConfig.InProduction = false

	// session creation
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true // persist cookie even after the window is closed by the user
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = appConfig.InProduction // cookies be encypted using https

	appConfig.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		fmt.Printf("error occured while creating template cache, err: %v\n", err)
	}
	appConfig.TemplateCache = tc
	appConfig.UseCache = true
	repo := handlers.NewRepo(&appConfig)
	handlers.NewHandlers(repo)

	render.NewTemplates(&appConfig)

	fmt.Println("Listening at Port", PORT)

	srv := &http.Server{
		Addr:    PORT,
		Handler: routes(&appConfig),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
