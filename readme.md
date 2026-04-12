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