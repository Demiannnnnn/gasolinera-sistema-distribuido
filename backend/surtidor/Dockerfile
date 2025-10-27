# --- BUILDER- --
FROM golang:1.25-alpine AS builder

WORKDIR /app

#descargar paquetes repositorio git
RUN apk --no-cache add git

#se define primero para que no vuelva a descargar dependencias si no hay cambios
# archivos de módulos de Go
COPY go.mod go.sum ./

# dependencias
RUN go mod download

#código fuente del backend
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main .

# --- Runner ---
FROM alpine:latest

# usuario no-root
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

WORKDIR /app

# Copia el binario compilado desde builder
COPY --from=builder /app/main .


EXPOSE 8082

CMD ["./main"]
