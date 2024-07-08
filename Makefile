build-windows:
	set GOOS= 
	set GOARCH= 
	go build -ldflags "-H=windowsgui" -o build/rayGoGame.exe

build-linux:
	set GOOS=linux
	set GOARCH=amd64
	go build -o build/rayGoGame

.PHONY: build build-windows build-linux