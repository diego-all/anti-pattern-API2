Tengo la siguiente api escrita en golang: Requiero cambiar el driver de:
github.com/jackc/pgx/v5/pgxpool" por esta:

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib" // Importar el adaptador stdlib


Aca esta la estructura de la API.

├── db
│   └── db.go
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── handlers
│   └── instrument_handler.go
├── main.go
├── models
│   └── instrument.go
├── README.md
├── request.md
└── sql
    └── init.sql

- /db/db.go
package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func InitDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL no definida en el entorno")
	}

	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	Pool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}

	err = Pool.Ping(ctx)
	if err != nil {
		log.Fatalf("Ping a la base de datos falló: %v", err)
	}
}

- /handlers/instrument_handler.go
package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"instruments-api/db"
	"instruments-api/models"

	"github.com/go-chi/chi/v5"
)

func GetAllInstruments(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Pool.Query(context.Background(), "SELECT id, name, description, price, created_at, updated_at FROM instruments")
	if err != nil {
		http.Error(w, "Error al obtener los instrumentos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var instruments []models.Instrument
	for rows.Next() {
		var ins models.Instrument
		if err := rows.Scan(&ins.ID, &ins.Name, &ins.Description, &ins.Price, &ins.CreatedAt, &ins.UpdatedAt); err != nil {
			http.Error(w, "Error al leer los datos", http.StatusInternalServerError)
			return
		}
		instruments = append(instruments, ins)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(instruments)
}

func GetInstrumentByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var ins models.Instrument
	err := db.Pool.QueryRow(context.Background(), `
		SELECT id, name, description, price, created_at, updated_at 
		FROM instruments WHERE id = $1`, id).
		Scan(&ins.ID, &ins.Name, &ins.Description, &ins.Price, &ins.CreatedAt, &ins.UpdatedAt)

	if err != nil {
		http.Error(w, "Instrumento no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ins)
}

func CreateInstrument(w http.ResponseWriter, r *http.Request) {
	var ins models.Instrument
	if err := json.NewDecoder(r.Body).Decode(&ins); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	now := time.Now()
	err := db.Pool.QueryRow(context.Background(), `
		INSERT INTO instruments (name, description, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`, ins.Name, ins.Description, ins.Price, now, now).
		Scan(&ins.ID)

	if err != nil {
		http.Error(w, "Error al insertar el instrumento", http.StatusInternalServerError)
		return
	}

	ins.CreatedAt = now
	ins.UpdatedAt = now

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ins)
}

func UpdateInstrument(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var ins models.Instrument
	if err := json.NewDecoder(r.Body).Decode(&ins); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	now := time.Now()
	result, err := db.Pool.Exec(context.Background(), `
		UPDATE instruments 
		SET name = $1, description = $2, price = $3, updated_at = $4 
		WHERE id = $5`,
		ins.Name, ins.Description, ins.Price, now, id)

	if err != nil || result.RowsAffected() == 0 {
		http.Error(w, "Error al actualizar el instrumento", http.StatusInternalServerError)
		return
	}

	ins.ID, _ = strconv.Atoi(id)
	ins.UpdatedAt = now

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ins)
}

func DeleteInstrument(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result, err := db.Pool.Exec(context.Background(), "DELETE FROM instruments WHERE id = $1", id)
	if err != nil || result.RowsAffected() == 0 {
		http.Error(w, "No se pudo eliminar", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/models/instrument.go
package models

import "time"

type Instrument struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

/sql/init.sql
CREATE TABLE IF NOT EXISTS instruments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(64),
    description VARCHAR(200),
    price INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL
);

INSERT INTO instruments (name, description, price, created_at, updated_at)
VALUES 
('Guitarra eléctrica', 'Guitarra Fender Stratocaster de seis cuerdas', 1200, NOW(), NOW()),
('Batería acústica', 'Set completo de batería Pearl con platillos', 2300, NOW(), NOW()),
('Teclado digital', 'Yamaha con 88 teclas contrapesadas', 850, NOW(), NOW()),
('Violín', 'Violín acústico hecho a mano con arco y estuche', 600, NOW(), NOW()),
('Saxofón alto', 'Saxofón profesional con boquilla y correa', 1500, NOW(), NOW());

docker-compose.yml

version: '3.8'

services:

  db:
    image: postgres:15-alpine # Usa una imagen oficial de PostgreSQL
    ports:
      - "5432:5432" # Opcional: Mapea el puerto de la DB al host (útil para herramientas externas)
    environment:
      POSTGRES_USER: user # Nombre de usuario de la base de datos
      POSTGRES_PASSWORD: password # Contraseña del usuario
      POSTGRES_DB: mydatabase # Nombre de la base de datos
    volumes:
      # Monta el archivo init.sql para que PostgreSQL lo ejecute al iniciar
      - ./sql/init.sql:/docker-entrypoint-initdb.d/init.sql
      # Persiste los datos de la base de datos en un volumen nombrado
      - pgdata:/var/lib/postgresql/data
  app:
    build:
      context: . # El contexto de construcción es el directorio actual
      dockerfile: Dockerfile # Usa el Dockerfile que acabamos de crear
    ports:
      - "8080:8080" # Mapea el puerto 8080 del host al puerto 8080 del contenedor
    environment:
      # Asegúrate de que estos valores coincidan con los de tu .env o sean los que esperas
      PORT: 8080
      DATABASE_URL: postgres://user:password@db:5432/mydatabase?sslmode=disable
    depends_on:
      - db # Asegura que la base de datos se inicie antes que la aplicación
    # Para desarrollo, puedes montar el código fuente para recargas en caliente si usas un observador
    # volumes:
    #   - .:/app # Descomenta para montar el volumen para desarrollo


volumes:
  pgdata: # Define el volumen nombrado para la persistencia de datos de PostgreSQL


- Dockerfile
# Usa una imagen base para la compilación
FROM golang:1.23.2-alpine AS builder

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia los archivos go.mod y go.sum para descargar las dependencias
COPY go.mod ./

# Descarga las dependencias del módulo
RUN go mod download

# Copia todo el código fuente de la aplicación
COPY . .

# Construye la aplicación Go
# CGO_ENABLED=0 deshabilita la vinculación de CGO para una compilación estática
# -o /go-app especifica el nombre del ejecutable y su ubicación
# main.go es el archivo de entrada principal
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go-app main.go

# Usa una imagen base ligera para la aplicación final
FROM alpine:latest

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /root/

# Copia el ejecutable compilado desde la etapa de construcción
COPY --from=builder /go-app .

# Copia el archivo init.sql para la inicialización de la base de datos (opcional, si quieres que el contenedor lo use)
# Aunque para un Docker Compose, es mejor que PostgreSQL lo maneje.
# Si el init.sql lo va a ejecutar la app, aquí es donde lo copiarías:
# COPY --from=builder /app/sql/init.sql /root/sql/init.sql

# Expone el puerto en el que la aplicación Go escuchará
EXPOSE 8080

# Comando para ejecutar la aplicación Go cuando el contenedor se inicie
CMD ["./go-app"]

- Aca esta el programa principal
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
		r.Get("/", handlers.GetAllInstruments)
		r.Get("/{id}", handlers.GetInstrumentByID)
		r.Post("/", handlers.CreateInstrument)
		r.Put("/{id}", handlers.UpdateInstrument)
		r.Delete("/{id}", handlers.DeleteInstrument)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor iniciado en http://localhost:%s\n", port)
	http.ListenAndServe(":"+port, r)
}


Podrias realizar la modificacion y entregar los archivos por completo y darme la respuesta en español. 







