# Local .env
export $(xargs < .env.dev)
go run .