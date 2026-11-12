# sistema-combustible-distribuido
Proyecto de la unidad 2 de la materia "Sistemas distribuidos". App para gasolinera la cual gestiona la distribucion de gasolina, precios, etc.ribucion de 


## Tests

- Backend distribuidor (Go): `cd backend/distribuidor && go test -race ./...`
- Monitor de distribuidores (React): `cd frontend/distribuidor-ui && CI=true npx react-scripts test --watchAll=false`

El workflow `.github/workflows/ci.yml` ejecuta vet, build y tests de los backends, tests y build de los tres frontends, y el build de las imágenes Docker en cada push y pull request.
