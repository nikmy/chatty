# RPC Chat
## Requirements
- `make`
- `go 1.17` (if you have no executables)
- `docker-compose 3.17`

## Usage
```shell
# Launch RPC Server
# You can use executable file instead of
# go run, but make is necessary
make env
go run ./cmd/rpc_server # press Ctrl+C to gracefully shutdown
make down
```
```shell
# Launch CLI client (you can also use executable file)
go run ./cmd/rpc_cli_client
```

## RPC Server command line options
- `-help` Show all available flags with description
- `-port 80` Listening port, default is `38120`
- `-logfile session.log` File for logging, default stdout (no file)

## RPC CLI Client
- History will be saved in file `.chatty_history.txt`

## Executables
Now there are executables for `windows amd64` in exe directory