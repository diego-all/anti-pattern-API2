package main // Este archivo pertenece al mismo paquete 'main' que main.go

import (
	"instruments-api/handlers" // Importa tus handlers existentes
	"net/http"                 // Necesario para http.Handler

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors" // Importa el paquete CORS aquí
)

// AppRoutes configura y devuelve un router Chi con todas las rutas y middlewares.
// La inicialización del router y la configuración de CORS se realizan dentro de esta función.
func AppRoutes() http.Handler {
	r := chi.NewRouter() // Inicializa el router Chi aquí, como solicitaste.

	// Middleware para mitigar Spectre agregando la cabecera de protección
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Cabecera para mitigar Spectre (CORP)
			w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
			w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
			next.ServeHTTP(w, r)
		})
	})

	// Configuración de CORS
	// Esto es importante para permitir solicitudes desde diferentes orígenes,
	// lo cual es común en entornos de desarrollo o APIs públicas.
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"}, // Permite cualquier origen HTTP/HTTPS
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Duración de caché para pre-vuelos CORS en segundos
	}))

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
