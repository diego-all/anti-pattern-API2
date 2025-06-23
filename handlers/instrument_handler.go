package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv" // Necesario para convertir string a int si se sigue usando en handlers

	"instruments-api/models" // Importar el paquete models

	"github.com/go-chi/chi/v5"
)

func GetAllInstruments(w http.ResponseWriter, r *http.Request) {
	instruments, err := models.GetAllInstruments(context.Background())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al obtener los instrumentos: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(instruments)
}

func GetInstrumentByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	ins, err := models.GetInstrumentByID(context.Background(), id)
	if err != nil {
		// La verificación de strconv.NumError ya no es directamente aplicable aquí
		// porque el error viene del modelo y puede ser más genérico.
		// El mensaje de error del modelo ahora es más descriptivo.
		http.Error(w, "Instrumento no encontrado o error de base de datos", http.StatusNotFound) // Mensaje genérico
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ins)
}

// CreateInstrument maneja la creación de un nuevo instrumento.
func CreateInstrument(w http.ResponseWriter, r *http.Request) {
	var ins models.Instrument
	if err := json.NewDecoder(r.Body).Decode(&ins); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	err := models.CreateInstrument(context.Background(), &ins)
	if err != nil {
		// 🚨 MALA PRÁCTICA: Se expone el error completo al cliente
		// Esto es un ejemplo claro de insecure error handling
		http.Error(w, fmt.Sprintf("Error al insertar el instrumento: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ins)
}

// UpdateInstrument maneja la actualización de un instrumento.
func UpdateInstrument(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var ins models.Instrument
	if err := json.NewDecoder(r.Body).Decode(&ins); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	rowsAffected, err := models.UpdateInstrument(context.Background(), id, &ins)
	if err != nil {
		// El error al no encontrar filas se maneja con RowsAffected en el handler
		http.Error(w, fmt.Sprintf("Error al actualizar el instrumento: %v", err), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No se pudo actualizar el instrumento o no se encontró", http.StatusInternalServerError)
		return
	}

	ins.ID, _ = strconv.Atoi(id)
	// ins.UpdatedAt se establece en el modelo, no es necesario reasignarlo aquí.
	// La línea comentada era: ins.UpdatedAt = now

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ins)
}

// DeleteInstrument maneja la eliminación de un instrumento.
func DeleteInstrument(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	rowsAffected, err := models.DeleteInstrument(context.Background(), id)
	if err != nil {
		// El error al no encontrar filas se maneja con RowsAffected en el handler
		http.Error(w, fmt.Sprintf("Error al eliminar: %v", err), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No se pudo eliminar el instrumento o no se encontró", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteInstrumentSQLi maneja la eliminación vulnerable de un instrumento por SQLi.
// Maybe it's for curl or r.URL.Query().Get("id")
func DeleteInstrumentSQLi(w http.ResponseWriter, r *http.Request) {
	// AHORA obtiene el ID como PARÁMETRO DE CONSULTA (ej. /endpoint?id=valor)
	id := r.URL.Query().Get("id")
	// id := chi.URLParam(r, "id") // Esta línea ya no es relevante aquí ya que el ID se obtiene de r.URL.Query()

	// Si no se proporciona ID, quizás quieras manejarlo
	if id == "" {
		http.Error(w, "El ID del instrumento es requerido", http.StatusBadRequest)
		return
	}

	rowsAffected, err := models.DeleteInstrumentSQLi(context.Background(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al eliminar el instrumento: %v", err), http.StatusInternalServerError) // Más detalle para debugging
		return
	}
	// if err != nil { // El error al no encontrar filas se maneja con RowsAffected en el modelo
	//  http.Error(w, "Error al eliminar", http.StatusInternalServerError)
	//  return
	// }

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

// GetInstrumentByIDSQLiURLParam obtiene un instrumento por ID vulnerable a SQLi vía URL param.
// QueryRowContext only return 1 row. Is not exploitable.
func GetInstrumentByIDSQLiURLParam(w http.ResponseWriter, r *http.Request) {

	//id := chi.URLParam(r, "id") will
	id := r.URL.Query().Get("id") // mario

	// var ins models.Instrument // La variable 'ins' ahora se declara dentro del modelo

	if id == "" {
		http.Error(w, "El ID del instrumento es requerido", http.StatusBadRequest)
		return
	}

	ins, err := models.GetInstrumentByIDSQLiURLParam(context.Background(), id)
	if err != nil {
		http.Error(w, "Instrumento no encontrado o error de base de datos", http.StatusNotFound) // Mensaje genérico
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ins)

}

// GetInstrumentByIDSQLi obtiene instrumentos por ID vulnerable a SQLi (puede devolver múltiples).
func GetInstrumentByIDSQLi(w http.ResponseWriter, r *http.Request) {
	// Obtiene el ID como PARÁMETRO DE CONSULTA (ej. /endpoint?id=valor)
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "El ID del instrumento es requerido", http.StatusBadRequest)
		return
	}

	instruments, err := models.GetInstrumentByIDSQLi(context.Background(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al consultar los instrumentos: %v", err), http.StatusInternalServerError)
		return
	}

	if len(instruments) == 0 {
		// La bandera 'found' se ha eliminado del modelo, se verifica aquí la longitud del slice.
		http.Error(w, "Instrumento(s) no encontrado(s) o error de base de datos", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(instruments) // Envía una lista de instrumentos
}
