Tipos de Vulnerabilidades a Introducir
Inyección SQL (SQL Injection): La más clásica y peligrosa para APIs que interactúan con bases de datos.
Exposición de Información Sensible (Sensitive Data Exposure): Mostrar datos que no deberían ser públicos.
Manejo Inseguro de Errores: Revelar detalles internos del sistema.
Autenticación/Autorización Débil o Ausente: Control de acceso deficiente.
Validación de Entrada Insuficiente: Confiar ciegamente en los datos del usuario.
Fuerza Bruta / Enumeración de Usuarios: Facilitar adivinanzas de credenciales o IDs.
Dependencias con Vulnerabilidades Conocidas: Aunque go mod download ya las baja, es bueno mencionarlo.
Modificaciones para Introducir Vulnerabilidades
1. Inyección SQL (SQL Injection)
Esta es la joya de la corona para escaners. La forma más sencilla de introducirla es construyendo consultas SQL directamente con cadenas de entrada del usuario sin parametrización.

Archivo a Modificar: /handlers/instrument_handler.go

Vulnerabilidad: Utilizar fmt.Sprintf o concatenación de cadenas para incluir entradas de usuario directamente en la consulta SQL.

Go

package handlers

import (
    "context"
    "encoding/json"
    "fmt" // Necesario para fmt.Sprintf
    "net/http"
    "strconv"
    "time"

    "instruments-api/db"
    "instruments-api/models"

    "github.com/go-chi/chi/v5"
)

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

// --- VULNERABILIDAD: SQL INJECTION en GetAllInstruments (con filtro inseguro) ---
// Simula un filtro dinámico que podría ser explotado
func GetAllInstruments(w http.ResponseWriter, r *http.Request) {
    // Parámetro de consulta opcional para "filtrar"
    filterName := r.URL.Query().Get("name") // Obtiene el parámetro 'name' de la URL

    var query string
    if filterName != "" {
        // VULNERABLE: Concatena directamente el valor del filtro en la cláusula WHERE.
        // Un atacante podría usar ?name=' OR 1=1 -- para inyectar.
        query = fmt.Sprintf("SELECT id, name, description, price, created_at, updated_at FROM instruments WHERE name = '%s'", filterName)
    } else {
        query = "SELECT id, name, description, price, created_at, updated_at FROM instruments"
    }

    rows, err := db.DBConn.QueryContext(context.Background(), query)
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

// --- VULNERABILIDAD: SQL INJECTION en DeleteInstrument ---
// Similar al caso de GetInstrumentByID, se concatena el ID.
func DeleteInstrument(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")

    // VULNERABLE: Concatenación directa de ID en la consulta SQL.
    // Un atacante podría usar "1 OR 1=1" para eliminar todos los registros.
    query := fmt.Sprintf("DELETE FROM instruments WHERE id = %s", id)

    result, err := db.DBConn.ExecContext(context.Background(), query)
    if err != nil {
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

// Los demás handlers (Create, Update) se mantienen seguros usando parámetros.
// Puedes vulnerarlos de forma similar si lo deseas, pero con estos ejemplos ya es suficiente.

// (Mantén los handlers CreateInstrument y UpdateInstrument como están, ya que usan parámetros seguros)
func CreateInstrument(w http.ResponseWriter, r *http.Request) {
    var ins models.Instrument
    if err := json.NewDecoder(r.Body).Decode(&ins); err != nil {
        http.Error(w, "JSON inválido", http.StatusBadRequest)
        return
    }

    now := time.Now()
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
    result, err := db.DBConn.ExecContext(context.Background(), `
        UPDATE instruments 
        SET name = $1, description = $2, price = $3, updated_at = $4 
        WHERE id = $5`,
        ins.Name, ins.Description, ins.Price, now, id)

    if err != nil {
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
2. Manejo Inseguro de Errores / Exposición de Información Sensible
Revelar mensajes de error detallados puede dar pistas a un atacante sobre la estructura interna de tu aplicación, base de datos o sistema operativo.

Archivo a Modificar: /db/db.go y /handlers/instrument_handler.go

Vulnerabilidad: Imprimir errores técnicos directamente o usar log.Fatal que sale de la aplicación.

Go

// db/db.go (Vulnerabilidad: log.Fatalf revela detalles)
package db

import (
    "context"
    "database/sql"
    "log"
    "os"
    "time"

    _ "github.com/jackc/pgconn"
    _ "github.com/jackc/pgx/v4"
    _ "github.com/jackc/pgx/v4/stdlib"
)

var DBConn *sql.DB

func InitDB() {
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        log.Fatal("DATABASE_URL no definida en el entorno. (¡Vulnerabilidad: mensaje revelador!)") // Mensaje más explícito, pero aún malo
    }

    var err error
    DBConn, err = sql.Open("pgx", dsn)
    if err != nil {
        // VULNERABLE: log.Fatalf expone el error técnico al log (si es accesible) y termina la aplicación.
        // En un entorno real, esto podría ser un mensaje genérico y registrar el error completo.
        log.Fatalf("No se pudo abrir la conexión a la base de datos: %v (¡Vulnerabilidad: detalles internos!)", err) 
    }

    DBConn.SetMaxOpenConns(25)
    DBConn.SetMaxIdleConns(25)
    DBConn.SetConnMaxLifetime(5 * time.Minute)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    err = DBConn.PingContext(ctx)
    if err != nil {
        // VULNERABLE: log.Fatalf expone el error técnico de la base de datos.
        log.Fatalf("Ping a la base de datos falló: %v (¡Vulnerabilidad: detalles técnicos de DB!)", err) 
    }

    log.Println("Conexión a la base de datos establecida exitosamente.")
}

Go

// handlers/instrument_handler.go (Vulnerabilidad: mensajes de error)
// En los handlers existentes, los mensajes de error ya son un poco genéricos,
// pero podríamos hacerlos más "útiles" para un atacante si revelamos qué salió mal.

// Ejemplo de GetInstrumentByID con error más específico (y malo)
// Reemplaza la línea: http.Error(w, "Instrumento no encontrado", http.StatusNotFound)
// Por esta:
// if err != nil {
//     // VULNERABLE: Revela el error exacto de la base de datos
//     http.Error(w, fmt.Sprintf("Error interno de DB: %v", err), http.StatusInternalServerError)
//     return
// }
// (Aunque ya tenemos el SQLi, esto añade una capa de vulnerabilidad si el SQLi falla pero hay un error de DB)
3. Autenticación/Autorización Débil o Ausente
Tu API actual no tiene ninguna autenticación. Esto ya es una vulnerabilidad en sí misma (cualquiera puede acceder a cualquier endpoint). Para hacerla explícitamente "débil", no se requiere cambio de código, solo notarlo.

Vulnerabilidad: No hay autenticación ni autorización para ninguna ruta. Cualquier usuario puede realizar cualquier operación (CRUD).

4. Validación de Entrada Insuficiente
Confiar en que el cliente enviará datos válidos es un error.

Archivo a Modificar: /handlers/instrument_handler.go

Vulnerabilidad: Permitir IDs no numéricos en rutas esperadas como números, o aceptar campos vacíos/malformados donde se esperan valores.

Go

// handlers/instrument_handler.go

// --- VULNERABILIDAD: Validación de ID Insuficiente en GetInstrumentByID ---
// Ya lo tiene vulnerable a SQLi, pero si no fuera SQLi, el ID podría ser inválido.
// En este caso, chi.URLParam() ya devuelve string.
// Si tuviéramos un middleware para validar que el ID es un número, podríamos quitarlo.
// (No necesita cambio de código, es la ausencia de validación extra)

// --- VULNERABILIDAD: Falta de validación de campos obligatorios en CreateInstrument ---
// Si un campo 'name' es obligatorio, el código actual lo aceptaría vacío.
// Podríamos dejarlo así y argumentar la falta de validación de negocio.
func CreateInstrument(w http.ResponseWriter, r *http.Request) {
    var ins models.Instrument
    if err := json.NewDecoder(r.Body).Decode(&ins); err != nil {
        http.Error(w, "JSON inválido", http.StatusBadRequest)
        return
    }

    // VULNERABLE: No se valida si ins.Name o ins.Price son valores sensatos o están presentes.
    // Un atacante podría enviar un instrumento con nombre vacío o precio negativo/cero,
    // o un precio excesivamente grande que cause un desbordamiento o problemas en la DB.
    if ins.Name == "" {
        // En una app segura: http.Error(w, "Nombre es obligatorio", http.StatusBadRequest)
        // Pero para ser vulnerable, lo permitimos.
    }
    if ins.Price <= 0 {
        // En una app segura: http.Error(w, "El precio debe ser positivo", http.StatusBadRequest)
        // Pero para ser vulnerable, lo permitimos.
    }

    now := time.Now()
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
5. Dependencias con Vulnerabilidades Conocidas
Esto es más difícil de "inyectar" a propósito sin manipular los go.mod y go.sum con versiones muy antiguas. Sin embargo, al usar go mod download, estás confiando en las versiones actuales. Para tu tesis, simplemente puedes argumentar que el escáner de vulnerabilidades de dependencias (como go list -m all combinado con govulncheck o Snyk) es crucial porque las bibliotecas de terceros pueden contener vulnerabilidades.

6. Credenciales y Configuraciones Sensibles
Tu docker-compose.yml ya expone POSTGRES_USER y POSTGRES_PASSWORD directamente. Esto es una mala práctica si este archivo cae en manos equivocadas o si se usa en un entorno de producción sin Secret Management.

Archivo a Modificar: /docker-compose.yml (ya está mal, solo se mantiene)

Vulnerabilidad: Credenciales de base de datos hardcodeadas y expuestas en el docker-compose.yml.

YAML

# docker-compose.yml
services:
  db:
    environment:
      POSTGRES_USER: user     # VULNERABLE: Credencial hardcodeada y expuesta
      POSTGRES_PASSWORD: password # VULNERABLE: Credencial hardcodeada y expuesta
      POSTGRES_DB: mydatabase
Resumen de Cambios y Cómo Explotarlos/Detectarlos
Inyección SQL (handlers/instrument_handler.go):

GetInstrumentByID:
Explotación: GET /instruments/1 OR 1=1 -- (debería devolver todos los instrumentos si el backend lo procesa) o GET /instruments/1; DROP TABLE instruments; -- (¡muy peligroso, eliminaría la tabla!).
Detección: Escáneres de seguridad dinámica (DAST) como OWASP ZAP, Burp Suite, SQLMap. Análisis estático de código (SAST) buscando concatenaciones de strings en consultas SQL.
GetAllInstruments (con filtro name):
Explotación: GET /instruments?name=' OR 1=1 --
Detección: DAST, SAST.
DeleteInstrument:
Explotación: DELETE /instruments/1 OR 1=1 (eliminaría todos los instrumentos).
Detección: DAST, SAST.
Manejo Inseguro de Errores (db/db.go):

Explotación: Forzar un error (ej. desconectando la DB o pasando un DSN inválido en el docker-compose.yml si tuvieras control) y observar los logs del contenedor para ver la información técnica que se revela.
Detección: Revisión manual de logs, análisis estático (SAST) buscando log.Fatal o fmt.Errorf con http.Error directo a w.
Falta de Autenticación/Autorización (General):

Explotación: Simplemente acceder a cualquier endpoint (GET, POST, PUT, DELETE) sin ninguna credencial.
Detección: Revisión manual de las rutas, SAST/DAST reportarán la ausencia de mecanismos de seguridad en los endpoints.
Validación de Entrada Insuficiente (handlers/instrument_handler.go - CreateInstrument):

Explotación: Enviar JSON con valores no válidos (ej. {"name": "", "price": 0}) si la lógica de negocio los requiere.
Detección: Pruebas de fuzzing, DAST (encontrarán que la API acepta datos "malos"), SAST.
Credenciales Expuestas (docker-compose.yml):

Explotación: Acceso directo al archivo docker-compose.yml.
Detección: Herramientas de análisis de secretos (ej. gitleaks, trufflehog) en el repositorio, revisión manual del código y configuración.
Pasos para Probar las Vulnerabilidades
Actualiza los archivos db/db.go y handlers/instrument_handler.go con los cambios propuestos.

Asegúrate de que tu docker-compose.yml esté usando el bind mount (como lo dejamos).

Reconstruye y levanta los contenedores:

Bash

docker compose down -v --rmi all
docker compose up --build -d
Usa herramientas como curl o Postman/Insomnia para intentar explotar las vulnerabilidades.

Inyección SQL (GET by ID):

Bash

# Intentará obtener un instrumento, pero la inyección puede listar todos o causar un error.
curl -X GET http://localhost:8080/instruments/1%20OR%201=1%20--
(Asegúrate de codificar la URL, %20 para espacio, -- para comentarios SQL)
También puedes probar con 1; SELECT pg_sleep(5); -- para un ataque de inyección SQL de tiempo.

Inyección SQL (GET All con filtro):

Bash

curl -X GET "http://localhost:8080/instruments?name=' OR 1=1 --"
Inyección SQL (DELETE):

Bash

# ¡CUIDADO! Esto intentará eliminar registros.
curl -X DELETE http://localhost:8080/instruments/1%20OR%201=1%20--
Validación de entrada débil (POST):

Bash

curl -X POST -H "Content-Type: application/json" -d '{"name": "", "description": "Instrumento sin nombre", "price": 0}' http://localhost:8080/instruments