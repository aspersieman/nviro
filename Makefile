build:
	go build -o bin/nviro main.go

run:
	go run main.go

compile:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=arm go build -o bin/nviro-linux-arm main.go
	GOOS=linux GOARCH=arm64 go build -o bin/nviro-linux-arm64 main.go
	GOOS=freebsd GOARCH=386 go build -o bin/nviro-freebsd-386 main.go
	GOOS=windows GOARCH=386 go build -o bin/nviro-windows-386 main.go

all: build
