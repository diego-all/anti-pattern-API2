# Usa la misma imagen base de PostgreSQL que en tu docker-compose
FROM postgres:15-alpine

# Copia el script de inicialización SQL dentro del directorio de entrada de PostgreSQL
# Este script se ejecutará automáticamente cuando el contenedor de la base de datos inicie por primera vez
COPY init.sql /docker-entrypoint-initdb.d/init.sql

# El comando por defecto de la imagen de postgres ya inicia el servidor, no necesitamos un CMD explícito aquí