package main // Este archivo pertenece al mismo paquete 'main' que main.go

import (
	"instruments-api/handlers" // Importa tus handlers existentes
	"net/http"                 // Necesario para http.Handler

	"github.com/go-chi/chi/v5"
	// Importa el paquete CORS aquí
)

// AppRoutes configura y devuelve un router Chi con todas las rutas.
func AppRoutes() http.Handler {
	r := chi.NewRouter() // Inicializa el router Chi aquí, como solicitaste.

	// Agrupa las rutas relacionadas con "/instruments"
	r.Route("/instruments", func(r chi.Router) {

		// Rutas CRUD estándar (las vulnerabilidades están en los handlers o modelos subyacentes)
		r.Get("/", handlers.GetAllInstruments)
		r.Get("/{id}", handlers.GetInstrumentByID)
		r.Post("/", handlers.CreateInstrument)
		r.Put("/{id}", handlers.UpdateInstrument)
		r.Delete("/{id}", handlers.DeleteInstrument)

		// Rutas Vulnerables (para propósitos académicos y de pruebas de seguridad)

		// Ruta DELETE vulnerable a SQLi (obtiene ID de query param)
		r.Delete("/vulnerable-sqli", handlers.DeleteInstrumentSQLi)

		// Ruta GET vulnerable a SQLi (obtiene ID de query param y puede devolver múltiples)
		r.Get("/vulnerable-sqligetinst", handlers.GetInstrumentByIDSQLi)

		// Ruta GET vulnerable a SQLi (obtiene ID de query param, pero originalmente diseñada para URL param)
		r.Get("/vulnerable-sqligetinsturlparam", handlers.GetInstrumentByIDSQLiURLParam)

		// Ruta PUT vulnerable a SQLi (obtiene ID de URL param y datos del JSON)
		//r.Put("/vulnerable-sqligetinst-put/{id}", handlers.GetInstrumentByIDSQLiPut)

		// Si en el futuro añades rutas para XSS, irían aquí.
		// r.Get("/products/get-xss/{id}", handlers.GetProductXSS)
	})

	return r
}
