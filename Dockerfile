# Usa una imagen base de Go directamente
FROM golang:1.23.2-alpine

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia los archivos go.mod y go.sum para descargar las dependencias
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copia todo el código fuente de la aplicación
COPY main.go .
COPY routes.go .
COPY handlers/ ./handlers/
COPY models/ ./models/
COPY db/ ./db/

COPY cert.pem .
COPY key.pem .

# Construye la aplicación Go
# NOTA: Sin CGO_ENABLED=0, la compilación usará el valor por defecto
# Si tu aplicación o sus dependencias requieren bibliotecas C,
# y estas no están presentes en la imagen alpine, el binario podría fallar al ejecutarse.
RUN GOOS=linux go build -o main .

# Expone el puerto en el que la aplicación Go escuchará
EXPOSE 8080

# Comando para iniciar la aplicación cuando el contenedor se inicie
CMD ["./main"]