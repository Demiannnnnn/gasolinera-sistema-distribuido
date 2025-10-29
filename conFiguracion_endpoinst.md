

Comunicacion de distribuidor con empresa 

| Nivel | Rol de Comunicación | Servicio Docker | Host/Puerto (Interno a Docker) | Host/Puerto (Acceso Externo / Host) | Propósito de la Conexión |
|-------|---------------------|-----------------|--------------------------------|-------------------------------------|--------------------------|
| N3 (Empresa) | Servidor TCP | empresa | 0.0.0.0:8081 | localhost:8081 | Escucha la recepción de transacciones y envía la actualización de precios base a los N2. |
| N2 (Distribuidor) | Cliente TCP | distribuidor1 | Cliente que inicia la conexión a empresa:8081. | N/A | Inicia la conexión periódicamente para enviar datos y consultar precios. |
| N2 (Distribuidor) | Cliente TCP | distribuidor2 | Cliente que inicia la conexión a empresa:8081. | N/A | Lo mismo que el distribuidor 1. |
| N2 (Distribuidor) | Cliente TCP | distribuidor3 | Cliente que inicia la conexión a empresa:8081. | N/A | Lo mismo que el distribuidor 1. |