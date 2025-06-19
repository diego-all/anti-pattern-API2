package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"instruments-api/db"
	"instruments-api/handlers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No se pudo cargar el archivo .env, usando variables del entorno")
	}

	db.InitDB()

	r := chi.NewRouter()
	r.Route("/instruments", func(r chi.Router) {

		// URLparam
		r.Get("/", handlers.GetAllInstruments)
		r.Get("/{id}", handlers.GetInstrumentByID)
		r.Post("/", handlers.CreateInstrument)
		r.Put("/{id}", handlers.UpdateInstrument)
		r.Delete("/{id}", handlers.DeleteInstrument)

		// original r.URL.Query().Get("id")  {id}
		// r.Delete("/vulnerable/instruments", handlers.DeleteInstrumentSQLi)
		r.Delete("/vulnerable-sqli", handlers.DeleteInstrumentSQLi) // Utiliza verbo, y funciona con curl
		// URLparam
		// r.Delete("/vulnerable/instruments/{id}", handlers.DeleteInstrumentSQLi)

		r.Get("/vulnerable-sqligetinst", handlers.GetInstrumentByIDSQLi) // Utiliza verbo, y funciona con curl

		r.Get("/vulnerable-sqligetinsturlparam", handlers.GetInstrumentByIDSQLiURLParam) // Utiliza verbo, y funciona con curl

	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor iniciado en http://localhost:%s\n", port)
	http.ListenAndServe(":"+port, r)
}
