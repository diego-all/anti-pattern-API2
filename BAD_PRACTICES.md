# Bad practices

> Burp Intruder

## SQL INjection in SELECT Sprintf (cometas+will)

**URL.Query vs URL.Param**

Permite retornar un usuario o varios. Segun sea la clausula.

        localhost:9090/instruments/vulnerable-sqligetinst?id=17' OR ''=' (Postman)

        curl -X GET "http://localhost:8080/instruments/vulnerable-sqligetinst?id=3"  (RECUPERA 1 FILA)

        curl -X GET "http://localhost:8080/instruments/vulnerable-sqligetinst?id=3%27%20OR%20%27%27=%27"  (RECUPERA TODAS LAS FILAS)


La cláusula WHERE id = %s normalmente devolvería un solo registro, si el campo id es clave primaria o único. Pero debido a que el id se está inyectando directamente con fmt.Sprintf, es posible manipular la consulta. Ejemplo:

SELECT id, name, email FROM users WHERE id = 1 OR 1=1   **(F) OR (V) = (V)**
⚠️ Esto devolvería todos los usuarios de la tabla.

Internamente ejecuta la consulta:

        SELECT id, name, description, price, created_at, updated_at
        FROM instruments
        WHERE id = '3' OR ''='';

Parece ser que no importa como se capture el parametro (Query param o path param), es la logica de la query y la falta de validacion en las entradas ya que con el Sprintf() se ejecuta directamente.

**DUDA:** dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

**Importante**: 
QueryContext (Cuando esperas cero, una o múltiples filas)
QueryRowContext (Como máximo una fila como resultado)


## SQL Injection in DELETE Sprintf (Payload codificado)

        curl -X DELETE "http://localhost:8080/instruments/vulnerable-sqli?id=3%27%20OR%20%27%27=%27"    FUNCIONO  (DELETE ALL)

       "3' OR ''='" esta debe ser or 1 = 1 la basica
       dynamic querys


## SQL Injection UNION SELECT null;--



## XSS in handler (productID, err := strconv.Atoi(chi.URLParam(r, "id"))) SANITIZER



## Exposición de Información Sensible: En docker compose



## Manejo Inseguro de Errores


## Autenticación/Autorización Débil o Ausente: Control de acceso deficiente.
Validación de Entrada Insuficiente: Confiar ciegamente en los datos del usuario.

## Fuerza Bruta / Enumeración de Usuarios: Facilitar adivinanzas de credenciales o IDs.

## Dependencias con Vulnerabilidades Conocidas: Aunque go mod download ya las baja, es bueno mencionarlo.


## cryptographic algorithm


## Harcoded credentials


Security Hotspots

- SQL Injection	High	Make sure using a dynamically formatted SQL query is safe here.
- Permission	Medium	Copying recursively might inadvertently add sensitive data to the container. Make sure it is safe here.
- Permission	Medium	The "golang" image runs with "root" as the default user. Make sure it is safe here
- Permission	Medium	The "postgres" image runs with "root" as the default user. Make sure it is safe here.