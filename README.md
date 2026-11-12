# gasolinera-sistema-distribuido

[![CI](https://github.com/Demiannnnnn/gasolinera-sistema-distribuido/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/Demiannnnnn/gasolinera-sistema-distribuido/actions/workflows/ci.yml)

Sistema distribuido para gasolineras (Unidad 2 de Sistemas Distribuidos) que gestiona la distribución de combustible, surtidores concurrentes, registro de ventas y actualización periódica en bases de datos PostgreSQL (Nivel 2 - Distribuidores y Nivel 3 - Empresa).

## Stack

- **Backend:** Go (`backend/distribuidor` con Gorilla Mux y sqlmock, `backend/empresa`).
- **Frontend:** React (`frontend/distribuidor-ui`, `frontend/empresa-ui`, `frontend/surtidor-ui`).
- **Base de datos:** PostgreSQL 15 vía Docker Compose (esquemas en `db/`).
- **Contenedores:** Docker y Docker Compose (`docker-compose.yml`).
- **CI:** GitHub Actions (`.github/workflows/ci.yml`).

## Cómo correrlo

1. Configurar el archivo de entorno desde la plantilla:
   ```bash
   cp .env.example .env
   ```
   Definir valores para `POSTGRES_USER` y `POSTGRES_PASSWORD`.

2. Iniciar todos los contenedores con Docker Compose:
   ```bash
   docker-compose up --build
   ```

   O iniciar únicamente un distribuidor junto con su base de datos y su interfaz web:
   ```bash
   docker-compose up -d --no-deps db_distribuidor1 distribuidor1 distribuidor1-ui
   ```

### Puertos de los servicios

- **Backend Distribuidor 1:** `http://localhost:8080` (UI: `http://localhost:3002`)
- **Backend Distribuidor 2:** `http://localhost:8083` (UI: `http://localhost:3003`)
- **Backend Distribuidor 3:** `http://localhost:8084` (UI: `http://localhost:3004`)
- **Backend Empresa:** `http://localhost:8081` (UI: `http://localhost:3001`)

## Cómo correr los tests

### Backend (Go)

- **Distribuidor:**
  ```bash
  cd backend/distribuidor && go test -race -count=1 ./...
  ```
- **Empresa:**
  ```bash
  cd backend/empresa && go test -race -count=1 ./...
  ```
- **Vet en backend:**
  ```bash
  cd backend/distribuidor && go vet ./...
  cd backend/empresa && go vet ./...
  ```

### Frontend (React)

- **Distribuidor UI:**
  ```bash
  cd frontend/distribuidor-ui && CI=true npx react-scripts test --watchAll=false
  ```
- **Empresa UI:**
  ```bash
  cd frontend/empresa-ui && CI=true npx react-scripts test --watchAll=false
  ```
- **Surtidor UI:**
  ```bash
  cd frontend/surtidor-ui && CI=true npx react-scripts test --watchAll=false
  ```

### CI (GitHub Actions)

El workflow en `.github/workflows/ci.yml` ejecuta `go vet`, `go build` y `go test` en los servicios de Go, los tests y builds de los frontends en Node 20, y la construcción de imágenes Docker en push y pull requests.
