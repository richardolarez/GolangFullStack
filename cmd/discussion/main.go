package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// ...
	ac := config.ParseAppConfig()

	db, err := sql.Open("postgres", ac.DBConnectionURI)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	session := session.NewSession(ac.SessionKey)
	templates, err := templates.GenerateTemplates("./web/views/*.gohtml")

	if err != nil {
		log.Fatal(err)
	}

	storage := storage.New(db)

	a := api.New(templates, storage, session)
	r := mux.NewRouter()
	a.RegisterRoutes(r)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/"))))
	r.Use(middleware.RequestLogger)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", ac.Port), r))

}
