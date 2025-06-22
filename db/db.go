package db

import (
	"context"
	"database/sql" // Importamos el paquete estándar database/sql
	"log"
	"time"

	// Importamos los drivers de pgx que permiten a database/sql interactuar con PostgreSQL
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// DBConn es la variable global que contendrá nuestra conexión a la base de datos
// Ahora es de tipo *sql.DB
var DBConn *sql.DB

// InitDB inicializa la conexión a la base de datos
func InitDB() {

	dsn := "host=db port=5432 user=user password=password dbname=mydatabase sslmode=disable timezone=UTC connect_timeout=5"

	var err error
	// Abre la conexión usando el driver "pgx" registrado por pgx/v4/stdlib
	DBConn, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("No se pudo abrir la conexión a la base de datos: %v", err)
	}

	// Configuración del pool de conexiones para database/sql
	// Esto es importante para un buen rendimiento en APIs web
	DBConn.SetMaxOpenConns(25)                 // Máximo número de conexiones abiertas
	DBConn.SetMaxIdleConns(25)                 // Máximo número de conexiones inactivas en el pool
	DBConn.SetConnMaxLifetime(5 * time.Minute) // Tiempo máximo que una conexión puede ser reutilizada

	// Intentamos hacer un Ping para verificar que la conexión es válida
	// Establecemos un contexto con timeout para el ping
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = DBConn.PingContext(ctx) // Usar PingContext para respetar el timeout
	if err != nil {
		log.Fatalf("Ping a la base de datos falló: %v", err)
	}

	log.Println("Conexión a la base de datos establecida exitosamente.")
}
