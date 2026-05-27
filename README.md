# IncogniCloud
A personal, customizable Cloud that works with any kind of Server. Uses Tailscale to make sure its not publically visible in the internet. Designed for NAS Systems to make backups and cloud saves easy as cake. Very lightweight to be able to work with any Hardware. One or two Harddrives are required to make the Cloud work effectively.

# Dev

## Starten der Datenbank

`docker run --name cloud-db -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=geheim -e POSTGRES_DB=cloud_db -p 5432:5432 -d postgres`

## Starten des backends

`go run cmd/server/main.go`

## Starten des Frontends

`npm run dev`

## Database Connect

`docker exec -it cloud_postgres psql -U postgres -d cloud_db`

## Prod Database Connect

`docker exec -it incognicloud_database psql -U postgres -d cloud_db`

## Docker Compose Starten

`docker compose up -d`

## Docker Compose Build Starten

`docker compose up --build -d`

## Docker Compose Down

`docker compose down -v `

## Neue Migrations Erstellen

`migrate create -ext sql -dir internal/database/migrations -seq [Name]`