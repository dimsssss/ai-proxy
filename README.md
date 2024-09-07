# ai-proxy
please Create a .env file in your project root directory. Then, enter the environment items in docker-compose.yml into .env.

## Required
go 1.23

docker 27.2.0

## How to test
```sh
go test ./internal/ratelimitratelimit_testgo
```

## How to run
```sh
docker compose -f ./deployments/docker-compose.yml up -d

go run main.go
```


