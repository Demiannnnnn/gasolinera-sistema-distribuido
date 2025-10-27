# --- BUILDER- --
# Usar una versión de Go compatible con tu go.mod
FROM golang:1.25-alpine AS builder

WORKDIR /app

#descargar paquetes repositorio git
# Es necesario para 'go mod download'
RUN apk --no-cache add git

#se define primero para que no vuelva a descargar dependencias si no hay cambios
# archivos de módulos de Go
COPY go.mod go.sum ./

# dependencias
RUN go mod download

#código fuente del backend
COPY . .

# Compilar para Alpine
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main main.go


# --- Runner ---
FROM alpine:latest

# usuario no-root
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copia el binario compilado desde builder
COPY --from=builder /app/main .

# permisos de ejecución al binario (y asignar propietario correcto)
RUN chmod +x /app/main && \
    chown appuser:appgroup /app/main

# todo está configurado, se cambia el usuario no-root
USER appuser

EXPOSE 8080

CMD ["/app/main"]

