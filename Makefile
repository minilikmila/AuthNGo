url=http://localhost:5001/api/v1/auth/health

Run:
	go run ./main.go
Run-dev:
	go run ./main.go -mode=dev
Health:
	curl -X GET $(url) -w "\n"
# Build:
# 	go build -o auth