package main

import (
	"crypto/tls"
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

		// "The XSS case will probably require using the GetInstrumentByID endpoint or overlapping some logic."
		// from XSS4
		// r.Get("/products/get-xss/{id}", GetProductXSS) // Nuevo endpoint vulnerable

		// NO FUNCIONO ANTERIORMENTE, VALIDAR DE NUEVO!!!
		//r.Put("/vulnerable-sqligetinst-put", handlers.GetInstrumentByIDSQLiPut) // Utiliza verbo, y funciona con curl

	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Configuración TLS vulnerable
	// Esto es específicamente lo que el escáner SAST debería identificar como "uso de un algoritmo criptográfico roto o riesgoso"
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			// Estos ciphersuites son conocidos por ser débiles o deprecados
			// e.g., CBC-SHA son vulnerables a ataques como BEAST
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,   // CipherSuite débil
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA, // CipherSuite débil
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,   // CipherSuite débil
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA, // CipherSuite débil
		},
		MinVersion: tls.VersionTLS12, // TLS 1.2 es el mínimo, pero los ciphersuites son el punto débil aquí
		// Otras configuraciones que pueden ser vulnerables:
		// MaxVersion: tls.VersionTLS12, // Limitar a TLS 1.2 puede ser una mala práctica si TLS 1.3 está disponible
		// InsecureSkipVerify: true, // NO USAR EN PRODUCCIÓN: deshabilita la verificación de certificados del cliente, muy riesgoso
	}

	// Crear un servidor HTTP con la configuración TLS personalizada
	srv := &http.Server{
		Addr:      ":" + port,
		Handler:   r,
		TLSConfig: tlsConfig, // Asignar la configuración TLS vulnerable
	}

	log.Printf("Servidor iniciado en http://localhost:%s\n", port)
	//http.ListenAndServe(":"+port, r)

	log.Fatal(srv.ListenAndServeTLS("cert.pem", "key.pem"))
}
