
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 packr build -a -installsuffix cgo -o main .
docker:
	docker build -t micro-go-chat -f Dockerfile.scratch .
