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

// func DeleteUserSQLi(w http.ResponseWriter, r *http.Request) {

// 	var payload jsonResponse
// 	id := r.URL.Query().Get("id")

// 	fmt.Println("El id es:", id)

// 	app.infoLog.Println(r.URL, id)

// 	err := app.models.User.DeleteUserSQLi(id)
// 	if err != nil {
// 		app.errorLog.Println("Couldn't delete user")
// 		// send back a response
// 		payload.Error = true
// 		payload.Message = "Couldn't delete user"
// 		return
// 	}

// 	// send back a response
// 	payload.Error = false
// 	payload.Message = "User deleted correctly"

// 	err = app.writeJSON(w, http.StatusOK, payload)
// 	if err != nil {
// 		app.errorLog.Println(err)
// 	}
// }
