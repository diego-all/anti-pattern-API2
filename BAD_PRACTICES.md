# Bad practices


- SQL Injection in SELECT Sprintf


- SQL Injection in DELETE Sprintf (Payload codificado)

        curl -X DELETE "http://localhost:8080/instruments/vulnerable-sqli?id=3%27%20OR%20%27%27=%27"    FUNCIONO


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