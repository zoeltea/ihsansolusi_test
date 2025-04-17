# ihsansolusi_test
It is a repository for assessment BE on ihsansolusi

## Requirements
- [Go]
- [Docker (Optional)]

## Getting Started
This section will guide you to get the project up and running

### Go Mod
This service need to install all of required dependencies
```
$ go mod tidy
```

### Running docker for postgres



### Application Properties or Environment
Copy and rename `.env.development` to `.env` and ensure all the propeties correct.

### Migration Database
Install goose
```
$ go install github.com/pressly/goose/v3/cmd/goose@latest

```
Migrate up using this command
```
$ goose -dir migrations postgres "user=postgres password=root dbname=postgres sslmode=disable" up
```

Migrate down using this command
```
$ goose -dir migrations postgres "user=postgres password=root dbname=postgres sslmode=disable" down
```

### Run service
using `go run`
```
$ go run main.go
```

### Visual studio debug
Create `launch.json` and apply with this
```
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "app debug",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}",
            "envFile": "${workspaceFolder}/cmd/github.com/lapakgaming/go-archetype/.env",
        }
    ]
}
```