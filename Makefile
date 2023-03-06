build:
	go build -o bin/nviro main.go

serve:
	go run main.go serve

compile:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=386 go build -o bin/nviro-linux-386 main.go
	GOOS=linux GOARCH=amd64 go build -o bin/nviro-linux-amd64 main.go
	GOOS=freebsd GOARCH=386 go build -o bin/nviro-freebsd-386 main.go
	GOOS=windows GOARCH=386 go build -o bin/nviro-windows-386.exe main.go
	#GOOS=js GOARCH=amd64 go build -o bin/nviro-wasm.wasm main.go

static:
	npx tailwindcss -i ./cmd/static/css/style.css -o ./cmd/static/css/output.css

static-watch:
	cd cmd/static/js && npx tailwindcss -i ./../css/style.css -o ./../css/output.css --watch

all: build
