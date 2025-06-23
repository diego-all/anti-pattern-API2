package models

import (
	"context"
	"database/sql" // Importamos el paquete estándar database/sql
	"fmt"
	"time"

	"instruments-api/db" // Importar el paquete db para acceder a DBConn
)

type Instrument struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GetAllInstruments obtiene todos los instrumentos de la base de datos.
func GetAllInstruments(ctx context.Context) ([]Instrument, error) {
	// Ahora usamos db.DBConn.Query() en lugar de db.Pool.Query()
	rows, err := db.DBConn.QueryContext(ctx, "SELECT id, name, description, price, created_at, updated_at FROM instruments")
	if err != nil {
		return nil, fmt.Errorf("error al obtener los instrumentos: %w", err)
	}
	defer rows.Close()

	var instruments []Instrument
	for rows.Next() {
		var ins Instrument
		if err := rows.Scan(&ins.ID, &ins.Name, &ins.Description, &ins.Price, &ins.CreatedAt, &ins.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error al leer los datos: %w", err)
		}
		instruments = append(instruments, ins)
	}

	// Verifica si hubo errores durante la iteración de las filas
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error en la iteración de resultados: %w", err)
	}

	return instruments, nil
}

// GetInstrumentByID obtiene un instrumento por su ID.
// --- VULNERABILIDAD: SQL INJECTION en GetInstrumentByID ---
// No se usa QueryRowContext con parámetros, se concatena la entrada directamente.
func GetInstrumentByID(ctx context.Context, id string) (*Instrument, error) {
	var ins Instrument
	// VULNERABLE: Concatenación directa de ID en la consulta SQL.
	// Un atacante podría pasar "1 OR 1=1 --" como ID para obtener todos los registros,
	// o "1; DROP TABLE instruments; --" para eliminar la tabla.
	query := fmt.Sprintf(`
        SELECT id, name, description, price, created_at, updated_at
        FROM instruments WHERE id = %s`, id) // ¡MUY PELIGROSO!

	// Ahora usamos db.DBConn.QueryRow() con la query vulnerable
	err := db.DBConn.QueryRowContext(ctx, query).
		Scan(&ins.ID, &ins.Name, &ins.Description, &ins.Price, &ins.CreatedAt, &ins.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("instrumento no encontrado")
		}
		return nil, fmt.Errorf("error de base de datos: %w", err)
	}
	return &ins, nil
}

// CreateInstrument crea un nuevo instrumento en la base de datos.
func CreateInstrument(ctx context.Context, ins *Instrument) error {
	now := time.Now()
	// Asigna CreatedAt y UpdatedAt aquí antes de la inserción
	ins.CreatedAt = now
	ins.UpdatedAt = now

	// Ahora usamos db.DBConn.QueryRow() para RETURNING
	err := db.DBConn.QueryRowContext(ctx, `
        INSERT INTO instruments (name, description, price, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id`, ins.Name, ins.Description, ins.Price, ins.CreatedAt, ins.UpdatedAt).
		Scan(&ins.ID)

	if err != nil {
		// El manejador se encargará de exponer el error completo al cliente (MALA PRÁCTICA)
		return fmt.Errorf("error al insertar el instrumento: %w", err)
	}
	return nil
}

// UpdateInstrument actualiza un instrumento existente por su ID.
func UpdateInstrument(ctx context.Context, id string, ins *Instrument) (int64, error) {
	now := time.Now()
	ins.UpdatedAt = now // Asigna UpdatedAt aquí antes de la actualización

	// Ahora usamos db.DBConn.Exec() en lugar de db.Pool.Exec()
	result, err := db.DBConn.ExecContext(ctx, `
        UPDATE instruments
        SET name = $1, description = $2, price = $3, updated_at = $4
        WHERE id = $5`,
		ins.Name, ins.Description, ins.Price, ins.UpdatedAt, id)

	if err != nil {
		return 0, fmt.Errorf("error al actualizar el instrumento: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error al obtener filas afectadas: %w", err)
	}
	return rowsAffected, nil
}

// DeleteInstrument elimina un instrumento por su ID.
func DeleteInstrument(ctx context.Context, id string) (int64, error) {
	// Ahora usamos db.DBConn.Exec() en lugar de db.Pool.Exec()
	result, err := db.DBConn.ExecContext(ctx, "DELETE FROM instruments WHERE id = $1", id)
	if err != nil {
		return 0, fmt.Errorf("error al eliminar: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error al obtener filas afectadas: %w", err)
	}
	return rowsAffected, nil
}

// DeleteInstrumentSQLi elimina un instrumento por ID (VULNERABLE A SQL INJECTION).
// Maybe it's for curl or r.URL.Query().Get("id")
func DeleteInstrumentSQLi(ctx context.Context, id string) (int64, error) {
	// AHORA obtiene el ID como PARÁMETRO DE CONSULTA (ej. /endpoint?id=valor)
	// id := chi.URLParam(r, "id")

	query := fmt.Sprintf("DELETE FROM instruments WHERE id = '%s'", id) // ¡VULNERABLE!

	fmt.Println("Consulta SQL ejecutada (vulnerable):", query) // Para ver la query inyectada en los logs

	result, err := db.DBConn.ExecContext(ctx, query)
	if err != nil {
		// Más detalle para debugging
		return 0, fmt.Errorf("error al eliminar el instrumento: %v", err)
	}
	// if err != nil { // El error al no encontrar filas se maneja con RowsAffected
	//  return 0, fmt.Errorf("Error al eliminar", http.StatusInternalServerError)
	// }

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Podría indicar un problema al obtener las filas afectadas después de una operación
		return 0, fmt.Errorf("error al verificar la eliminación: %w", err)
	}
	return rowsAffected, nil
}

// GetInstrumentByIDSQLiURLParam obtiene un instrumento por ID (VULNERABLE A SQL INJECTION en URL param).
// QueryRowContext only return 1 row. Is not exploitable.
func GetInstrumentByIDSQLiURLParam(ctx context.Context, id string) (*Instrument, error) {
	//id := chi.URLParam(r, "id") will
	// id := r.URL.Query().Get("id") // mario

	var ins Instrument

	// query := fmt.Sprintf("DELETE FROM instruments WHERE id = '%s'", id) // ¡VULNERABLE!
	query := fmt.Sprintf("SELECT id, name, description, price, created_at, updated_at FROM instruments WHERE id = '%s'", id) // ¡VULNERABLE!

	// db vs database

	// Will usa Query(query)

	// Ahora usamos db.DBConn.QueryRow() con las query vulnerable
	err := db.DBConn.QueryRowContext(ctx, query).
		Scan(&ins.ID, &ins.Name, &ins.Description, &ins.Price, &ins.CreatedAt, &ins.UpdatedAt)

	if err != nil {
		// Mensaje genérico
		return nil, fmt.Errorf("instrumento no encontrado o error de base de datos: %w", err)
	}
	return &ins, nil
}

// GetInstrumentByIDSQLi obtiene instrumentos por ID (VULNERABLE A SQL INJECTION, puede devolver múltiples).
func GetInstrumentByIDSQLi(ctx context.Context, id string) ([]Instrument, error) {
	// Obtiene el ID como PARÁMETRO DE CONSULTA (ej. /endpoint?id=valor)
	// id := r.URL.Query().Get("id")

	// Consulta SQL VULNERABLE: Concatenación directa de ID en la cláusula WHERE.
	// Un atacante podría usar '3' OR ''='' para que la condición WHERE sea siempre verdadera,
	// devolviendo todas las filas.
	query := fmt.Sprintf(`
        SELECT id, name, description, price, created_at, updated_at
        FROM instruments WHERE id = '%s'`, id) // ¡VULNERABLE!

	fmt.Println("Consulta SQL ejecutada (vulnerable):", query) // Para ver la query inyectada en los logs

	// CAMBIO CLAVE: Usar db.DBConn.QueryContext para esperar múltiples filas
	rows, err := db.DBConn.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al consultar la base de datos: %v", err)
	}
	defer rows.Close() // Es crucial cerrar las filas

	var instruments []Instrument
	// found := false // Bandera para saber si se encontró al menos un instrumento

	for rows.Next() {
		var ins Instrument
		// Asegúrate de que todos los campos del SELECT están siendo escaneados aquí.
		// Si Price, CreatedAt o UpdatedAt son nulos en la DB para alguna fila inyectada,
		// o si el payload es malicioso y altera el esquema, esto podría fallar.
		if err := rows.Scan(&ins.ID, &ins.Name, &ins.Description, &ins.Price, &ins.CreatedAt, &ins.UpdatedAt); err != nil {
			// Maneja el error de escaneo, podría ser por tipos de datos
			return nil, fmt.Errorf("error al leer los datos del instrumento: %v", err)
		}
		instruments = append(instruments, ins)
		// found = true // Ya no es necesaria la bandera aquí, se verifica len(instruments) en el handler
	}

	// Verifica si hubo errores durante la iteración de las filas
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error en la iteración de resultados: %v", err)
	}

	return instruments, nil
}
