# Multi-stage build para una imagen final pequeña

# --- Etapa de compilación ---
FROM golang:1.23.2-alpine AS builder

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia los archivos go.mod y go.sum para descargar las dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copia todo el código fuente de la aplicación
COPY . .

# Construye la aplicación Go
# CGO_ENABLED=0 deshabilita la vinculación de CGO para una compilación estática
# GOOS=linux para asegurar que se compile para Linux (el SO base del contenedor)
# -o /go-app especifica el nombre del ejecutable y su ubicación
# main.go es el archivo de entrada principal
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go-app main.go

# --- Etapa de ejecución (imagen final) ---
FROM alpine:latest

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /root/

# Copia el ejecutable compilado desde la etapa de construcción
COPY --from=builder /go-app .

# Expone el puerto en el que la aplicación Go escuchará
EXPOSE 8080

# Comando para iniciar la aplicación
CMD ["./go-app"]