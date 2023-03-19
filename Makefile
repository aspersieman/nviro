dist:
	npm install
	npm run build

build: dist
	go build -o bin/nviro main.go

serve: dist
	go run main.go serve

compile: dist
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=386 go build -o bin/nviro-linux-386 main.go
	GOOS=linux GOARCH=amd64 go build -o bin/nviro-linux-amd64 main.go
	GOOS=freebsd GOARCH=386 go build -o bin/nviro-freebsd-386 main.go
	GOOS=windows GOARCH=386 go build -o bin/nviro-windows-386.exe main.go
	#GOOS=js GOARCH=amd64 go build -o bin/nviro-wasm.wasm main.go

all: build

clean:
	rm -rf cmd/static/dist
