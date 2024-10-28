NAME=otlh

build:
	GOOS=windows GOARCH=amd64 go build -o bin/${NAME}-win64.exe .
	GOOS=linux GOARCH=amd64 go build -o bin/${NAME}-linux-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o bin/${NAME}-darwin-arm64 .

run:
	go run .

clean:
	go clean
	rm bin/${NAME}-darwin-arm64
	rm bin/${NAME}-linux
	rm bin/${NAME}-win64.exe
