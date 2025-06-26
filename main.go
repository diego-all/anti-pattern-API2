package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"instruments-api/db"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No se pudo cargar el archivo .env, usando variables del entorno")
	}

	db.InitDB()

	// Inicializa el router de la aplicación llamando a la función Routes del nuevo archivo routes.go.
	// Esto centraliza la definición de las rutas y middleware en un solo lugar.
	router := AppRoutes() // Corregido: Ahora llama a Routes() para que coincida con la definición en routes.go

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
		Handler:   router,
		TLSConfig: tlsConfig, // Asignar la configuración TLS vulnerable
	}

	log.Printf("Servidor iniciado en http://localhost:%s\n", port)
	//http.ListenAndServe(":"+port, r)

	log.Fatal(srv.ListenAndServeTLS("cert.pem", "key.pem"))
}
