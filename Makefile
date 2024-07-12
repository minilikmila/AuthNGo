url=http://localhost:9090/api/v1/auth/health

Run:
	go run ./main.go
Run-dev:
	go run ./main.go -dev
Health:
	curl -X GET $(url) -w "\n"
# Build:
# 	go build -o auth