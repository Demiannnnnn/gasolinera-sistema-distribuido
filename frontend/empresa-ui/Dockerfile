FROM node:18-alpine
WORKDIR /app
# Instalar dependencias del sistema
RUN apk add --no-cache git

# Copiar archivos de dependencias
COPY package*.json ./

# Instalar dependencias
RUN npm ci
COPY . .
EXPOSE 3001

# Variable de entorno para hot-reload en Docker
ENV CHOKIDAR_USEPOLLING=true
ENV WATCHPACK_POLLING=true

CMD ["npm", "start"]

