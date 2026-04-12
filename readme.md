# Readme

## Running Locally
1. Enter cmd directory:
```shell
cd cmd
```
Then run
```bash
go run .
# Service starts on http://localhost:8080
```


## Firebase setup

1. Download key from:
```
https://console.firebase.google.com/u/1/project/.../settings/serviceaccounts/adminsdk
```

2. Save the file as `secrets/fire-key.json` in the project root.


3. Then set environment variable:
```
export GOOGLE_APPLICATION_CREDENTIALS="../secrets/fire-key.json"
```

4. Done


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
