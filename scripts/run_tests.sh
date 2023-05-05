# Local .env
cd ..
export $(xargs < .env.test)
go test ./...