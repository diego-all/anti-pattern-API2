Como se realizo anteriormente, volver un endpoint vulnerable a SQL injection (DeleteInstrumentSQLi).
Que funciona y es explotable con el payload:

root@pho3nix:/home/diegoall/MAESTRIA_ING/anti-pattern-API/anti-pattern-API2# curl -X DELETE "http://localhost:8080/instruments/vulnerable-sqli?id=3%27%20OR%20%27%27=%27"
{"error":false}

root@pho3nix:/home/diegoall/MAESTRIA_ING/anti-pattern-API/anti-pattern-API2# docker logs -f anti-pattern-api2_app_1
2025/06/19 04:37:42 No se pudo cargar el archivo .env, usando variables del entorno
2025/06/19 04:37:42 Conexión a la base de datos establecida exitosamente.
2025/06/19 04:37:42 Servidor iniciado en http://localhost:8080

Consulta SQL ejecutada (vulnerable): DELETE FROM instruments WHERE id = '3' OR ''=''

Se evidencia que se eliminan los registros de la base de datos:
mydatabase=# select * from instruments;
 id |        name        |                   description                   | price |         created_at         |         updated_at         
----+--------------------+-------------------------------------------------+-------+----------------------------+----------------------------
  1 | Guitarra eléctrica | Guitarra Fender Stratocaster de seis cuerdas    |  1200 | 2025-06-19 04:34:03.760094 | 2025-06-19 04:34:03.760094
  2 | Batería acústica   | Set completo de batería Pearl con platillos     |  2300 | 2025-06-19 04:34:03.760094 | 2025-06-19 04:34:03.760094
  3 | Teclado digital    | Yamaha con 88 teclas contrapesadas              |   850 | 2025-06-19 04:34:03.760094 | 2025-06-19 04:34:03.760094
  4 | Violín             | Violín acústico hecho a mano con arco y estuche |   600 | 2025-06-19 04:34:03.760094 | 2025-06-19 04:34:03.760094
  5 | Saxofón alto       | Saxofón profesional con boquilla y correa       |  1500 | 2025-06-19 04:34:03.760094 | 2025-06-19 04:34:03.760094
(5 rows)

mydatabase=# select * from instruments;
 id | name | description | price | created_at | updated_at 
----+------+-------------+-------+------------+------------
(0 rows)


Ahora requiero algo similar pero para el endpoint de (GetInstrumentByIDSQLi) de tal forma que tambien se pueda explotar utilizando curl. ya lo implemente pero no borro los registros.

curl -X GET "http://localhost:8080/instruments/vulnerable-sqligetinst?id=3%27%20OR%20%27%27=%27"

root@pho3nix:/home/diegoall/MAESTRIA_ING/anti-pattern-API/anti-pattern-API2# curl -X GET "http://localhost:8080/instruments/vulnerable-sqligetinst?id=3%27%20OR%20%27%27=%27"
Instrumento no encontrado o error de base de datos

La funcion la implemente utilizando query params pero al parcer no funciona. Podrias ayudarme a corregirla de tal forma que pueda ser explotable.


- Aca esta /handlers/instrument_handlers.go

package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"instruments-api/db"
	"instruments-api/models"

	"github.com/go-chi/chi/v5"
)

func GetAllInstruments(w http.ResponseWriter, r *http.Request) {
	// Ahora usamos db.DBConn.Query() en lugar de db.Pool.Query()
	rows, err := db.DBConn.QueryContext(context.Background(), "SELECT id, name, description, price, created_at, updated_at FROM instruments")
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

// --- VULNERABILIDAD: SQL INJECTION en GetInstrumentByID ---
// No se usa QueryRowContext con parámetros, se concatena la entrada directamente.
func GetInstrumentByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var ins models.Instrument
	// VULNERABLE: Concatenación directa de ID en la consulta SQL.
	// Un atacante podría pasar "1 OR 1=1 --" como ID para obtener todos los registros,
	// o "1; DROP TABLE instruments; --" para eliminar la tabla.
	query := fmt.Sprintf(`
        SELECT id, name, description, price, created_at, updated_at
        FROM instruments WHERE id = %s`, id) // ¡MUY PELIGROSO!

	// Ahora usamos db.DBConn.QueryRow() con la query vulnerable
	err := db.DBConn.QueryRowContext(context.Background(), query).
		Scan(&ins.ID, &ins.Name, &ins.Description, &ins.Price, &ins.CreatedAt, &ins.UpdatedAt)

	if err != nil {
		http.Error(w, "Instrumento no encontrado o error de base de datos", http.StatusNotFound) // Mensaje genérico
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
	// Ahora usamos db.DBConn.QueryRow() en lugar de db.Pool.QueryRow() para RETURNING
	err := db.DBConn.QueryRowContext(context.Background(), `
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
	// Ahora usamos db.DBConn.Exec() en lugar de db.Pool.Exec()
	result, err := db.DBConn.ExecContext(context.Background(), `
        UPDATE instruments 
        SET name = $1, description = $2, price = $3, updated_at = $4 
        WHERE id = $5`,
		ins.Name, ins.Description, ins.Price, now, id)

	if err != nil { // El error al no encontrar filas se maneja con RowsAffected
		http.Error(w, "Error al actualizar el instrumento", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		http.Error(w, "No se pudo actualizar el instrumento o no se encontró", http.StatusInternalServerError)
		return
	}

	ins.ID, _ = strconv.Atoi(id)
	ins.UpdatedAt = now

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ins)
}

func DeleteInstrument(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Ahora usamos db.DBConn.Exec() en lugar de db.Pool.Exec()
	result, err := db.DBConn.ExecContext(context.Background(), "DELETE FROM instruments WHERE id = $1", id)
	if err != nil { // El error al no encontrar filas se maneja con RowsAffected
		http.Error(w, "Error al eliminar", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		http.Error(w, "No se pudo eliminar el instrumento o no se encontró", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Maybe it's for curl or r.URL.Query().Get("id")
func DeleteInstrumentSQLi(w http.ResponseWriter, r *http.Request) {
	// AHORA obtiene el ID como PARÁMETRO DE CONSULTA (ej. /endpoint?id=valor)
	id := r.URL.Query().Get("id")
	// id := chi.URLParam(r, "id")

	// Si no se proporciona ID, quizás quieras manejarlo
	if id == "" {
		http.Error(w, "El ID del instrumento es requerido", http.StatusBadRequest)
		return
	}

	query := fmt.Sprintf("DELETE FROM instruments WHERE id = '%s'", id) // ¡VULNERABLE!

	fmt.Println("Consulta SQL ejecutada (vulnerable):", query) // Para ver la query inyectada en los logs

	result, err := db.DBConn.ExecContext(context.Background(), query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al eliminar el instrumento: %v", err), http.StatusInternalServerError) // Más detalle para debugging
		return
	}
	// if err != nil { // El error al no encontrar filas se maneja con RowsAffected
	// 	http.Error(w, "Error al eliminar", http.StatusInternalServerError)
	// 	return
	// }

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Podría indicar un problema al obtener las filas afectadas después de una operación
		http.Error(w, "Error al verificar la eliminación", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		// Indica que no se encontró el instrumento o la inyección no eliminó nada
		http.Error(w, "No se pudo eliminar el instrumento o no se encontró", http.StatusNotFound)
		return
	}

	// w.WriteHeader(http.StatusNoContent)
	// Respuesta de éxito similar a tu ejemplo de DeleteUserSQLi
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"error": false}) // o un struct de payload

}

func GetInstrumentByIDSQLi(w http.ResponseWriter, r *http.Request) {

	//id := chi.URLParam(r, "id") will
	id := r.URL.Query().Get("id") // mario

	var ins models.Instrument

	if id == "" {
		http.Error(w, "El ID del instrumento es requerido", http.StatusBadRequest)
		return
	}

	// query := fmt.Sprintf("DELETE FROM instruments WHERE id = '%s'", id) // ¡VULNERABLE!
	query := fmt.Sprintf("SELECT id, name, description FROM instruments WHERE id = '%s'", id) // ¡VULNERABLE!

	// db vs database

	// Will usa Query(query)

	// Ahora usamos db.DBConn.QueryRow() con la query vulnerable
	err := db.DBConn.QueryRowContext(context.Background(), query).
		Scan(&ins.ID, &ins.Name, &ins.Description, &ins.Price, &ins.CreatedAt, &ins.UpdatedAt)

	if err != nil {
		http.Error(w, "Instrumento no encontrado o error de base de datos", http.StatusNotFound) // Mensaje genérico
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ins)

}

func GetAllInstrumentsSQLi(w http.ResponseWriter, r *http.Request) {

}


- Aca esta el programa principal main.go
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

		r.Get("/vulnerable-sqligetall", handlers.GetAllInstrumentsSQLi) // Utiliza verbo, y funciona con curl

	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor iniciado en http://localhost:%s\n", port)
	http.ListenAndServe(":"+port, r)
}

Respuesta en español

