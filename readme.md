# Readme

##### Made by
- Tinius J. Presterud
- Sindre Olsen
- Jonas Tomren

##### Main repo URL
[https://github.com/tiniuspre/prog2005-2026_assignment_2](https://github.com/tiniuspre/prog2005-2026_assignment_2)

## How to setup and run

### Prerequisites
- Docker


### First time setup with Docker

1. Copy .env.example to .env and fill in the values:
```
cp .env.example .env
```

2. Download the firebase key and save it as `secrets/fire-key.json`:
3. Download openaq api key and set it in the .env file.
4. Then build and run the docker container:
```bash
docker compose up --build
```
or if you want to run it in the background:
```bash
docker compose up --build -d
```

### First time setup without Docker
1. Download the firebase key and save it as `secrets/fire-key.json`:
2. Download openaq api key and set it in the .env file.
3. Then set environment variable:
```shell
export GOOGLE_APPLICATION_CREDENTIALS="secrets/fire-key.json"
export OPENAQ_API_KEY="your openaq key"
export PORT=8080
```
4. Run the project:
```bash
go run cmd/main.go
```

---

## Firebase setup

1. Download key from:
```
https://console.firebase.google.com/u/1/project/.../settings/serviceaccounts/adminsdk
```

2. Save the file as `secrets/fire-key.json` in the project root.


3. Then set environment variable:
```
export GOOGLE_APPLICATION_CREDENTIALS="secrets/fire-key.json"
```

4. Done

---
## Docker
Building and running:
```bash
docker compose up --build
```

If in need of detached / running in the background:
```bash
docker compose up --build -d
```

---
# How to use

1. Register a user:
```bash
curl -s -X POST http://localhost:8080/auth/ \
  -H "Content-Type: application/json" \
  -d '{"name":"your name","email":"test@stud.ntnu.no"}'
```

You will get something like:
```json
{"key":"sk-envdash-ba2c...0d8","createdAt":"20260412 18:49"}
```

Use the key like:
```bash
curl -s -X GET http://localhost:8080/somepath \
    -H "X-API-Key: sk-envdash-ba2c...0d8"
```

# About the project

### Project Structure

```
.
├── Dockerfile
├── cmd
│   └── main.go
├── docker-compose.yml
├── files.txt
├── go.mod
├── go.sum
├── internal
│   ├── clients
│   │   ├── countries.go
│   │   ├── countries_test.go
│   │   ├── currencies.go
│   │   ├── currencies_test.go
│   │   ├── meteo.go
│   │   ├── meteo_test.go
│   │   ├── nominatim.go
│   │   ├── nominatim_test.go
│   │   ├── openaq.go
│   │   └── openaq_test.go
│   ├── firebase
│   │   ├── api-keys.go
│   │   ├── cache.go
│   │   ├── client.go
│   │   ├── notifications.go
│   │   └── registrations.go
│   ├── handlers
│   │   ├── auth.go
│   │   ├── auth_test.go
│   │   ├── dashboard_success_test.go
│   │   ├── dashboards.go
│   │   ├── dashboards_test.go
│   │   ├── deps.go
│   │   ├── dispatch.go
│   │   ├── dispatch_test.go
│   │   ├── helpers.go
│   │   ├── notifications.go
│   │   ├── notifications_test.go
│   │   ├── registrations.go
│   │   ├── registrations_success_test.go
│   │   ├── registrations_test.go
│   │   ├── status.go
│   │   ├── status_test.go
│   │   ├── store.go
│   │   ├── store_firestore.go
│   │   └── store_memory.go
│   ├── middleware
│   │   ├── auth.go
│   │   ├── auth_test.go
│   │   └── deps.go
│   └── models
│       ├── country.go
│       └── models.go
├── readme.md
└── secrets
    └── fire-key.json
```

### API Endpoints

```
### API Endpoints

POST   /auth/                             # Register a new user and get an API key
DELETE /auth/{key}                        # Revoke an API key

POST   /envdash/v1/registrations/         # Create a new dashboard configuration
GET    /envdash/v1/registrations/{id}     # Get one registration by ID
GET    /envdash/v1/registrations/         # List all registrations
PUT    /envdash/v1/registrations/{id}     # Update an existing registration
DELETE /envdash/v1/registrations/{id}     # Delete a registration

GET    /envdash/v1/dashboards/{id}        # Get a populated dashboard for a registration

POST   /envdash/v1/notifications/         # Register a new webhook notification
GET    /envdash/v1/notifications/{id}     # Get one notification by ID
GET    /envdash/v1/notifications/         # List all notifications
DELETE /envdash/v1/notifications/{id}     # Delete a notification

GET    /envdash/v1/status/                # Get service and dependency health status
```

---

## Caching Strategy

Country data from the REST Countries API is cached in Firestore with a 24-hour TTL.
This is the only upstream response that is cached, for the following reasons:

- **Country data is static** — population, area, capital, coordinates and currency codes
  change rarely if ever. Caching for 24 hours introduces no meaningful data staleness.
- **Weather data is not cached** — temperature and precipitation forecasts from Open-Meteo
  are time-sensitive and must reflect current conditions on every dashboard retrieval.
- **Air quality data is not cached** — PM2.5 and PM10 readings from OpenAQ represent live
  sensor measurements. Caching these would defeat the purpose of the feature.
- **Currency rates are not cached** — exchange rates fluctuate constantly and should always
  reflect the latest available values.

This approach maximises the reduction in outbound API traffic where it is safe to do so,
while ensuring all time-sensitive data remains accurate on every request.

Nominatim coordinate lookups are not cached. While capital city coordinates are static,
Nominatim responses are fast and the service has no rate limit concerns at the scale
of this application. Adding a separate cache collection for coordinates would introduce
additional complexity without meaningful reduction in outbound traffic, since Nominatim
is only called when air quality data is requested.
