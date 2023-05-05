# Local .env
cd ..
export $(xargs < .env.dev)
go run .