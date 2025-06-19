

- Paylodads de gemini

En resumen:

* Tu API es vulnerable, pero la forma en que lo es, no permite que los payloads clásicos de OR 1=1 (que son para columnas de texto o cuando el id se encierra entre comillas) funcionen directamente con una columna INTEGER sin comillas.

* PostgreSQL es demasiado "inteligente" (o estricto) para el tipo de inyección que intentamos hacer. Intenta convertir el texto inyectado (1 OR 1=1 --) a un entero antes de evaluarlo, y falla.

* Además, el uso de ; para encadenar comandos (DROP TABLE) generalmente no funciona con QueryRowContext o ExecContext por defecto en los drivers de Go, por razones de seguridad. Estos métodos están diseñados para ejecutar una sola declaración.


query := fmt.Sprintf(`
        SELECT id, name, description, price, created_at, updated_at 
        FROM instruments WHERE id = %s`, id)  **No funciona**

query := fmt.Sprintf(`
        SELECT id, name, description, price, created_at, updated_at
        FROM instruments WHERE id = '%s'`, id) // ¡AHORA SÍ CON COMILLAS!



La clave para que tu función DeleteInstrument sea vulnerable de la misma manera que DeleteUserSQLi es eliminar la parametrización de la consulta SQL y concatenar directamente la entrada del usuario.

El problema con tu función DeleteInstrument original es que usa $1 y pasa id como un argumento separado a ExecContext. Esto hace que la base de datos trate id como un valor seguro, impidiendo la inyección.

La función DeleteUserSQLi funciona porque (asumiendo que app.models.User.DeleteUserSQLi construye la consulta concatenando) toma el id directamente del URL sin sanearlo ni parametrizarlo.



Diferentes comportamientos

definicion de la base de datos

- ID como texto
- Id como varchar
- Id como int  (serial)
- Id como UUID  uuid_generate_v4()


sql-injection-dummy id serial


