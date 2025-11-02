# Etapa de construcción
FROM golang:1.18-alpine AS builder

WORKDIR /app

# Copiar los archivos de módulos y descargar dependencias
COPY go.mod ./
# Asegurarse de que las dependencias estén descargadas antes de copiar el resto del código
RUN go mod download

# Copiar el resto del código fuente
COPY . .

# Compilar la aplicación
# -o /app/main crea el binario en la raíz del directorio de la app
# CGO_ENABLED=0 deshabilita CGO para crear un binario estático
RUN CGO_ENABLED=0 go build -o /app/main .

# Etapa final
FROM alpine:latest

WORKDIR /app

# Copiar solo el binario compilado desde la etapa de construcción
COPY --from=builder /app/main .

# Exponer el puerto si fuera necesario (no es necesario para este cliente)
# EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./main"]
