build:
	go build -o bin/nviro main.go

run:
	go run main.go

compile:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=386 go build -o bin/nviro-linux-386 main.go
	GOOS=linux GOARCH=amd64 go build -o bin/nviro-linux-amd64 main.go
	GOOS=freebsd GOARCH=386 go build -o bin/nviro-freebsd-386 main.go
	GOOS=windows GOARCH=386 go build -o bin/nviro-windows-386.exe main.go

all: build
