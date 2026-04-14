# Readme

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
export OPENAQ_KEY="your openaq key"
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
