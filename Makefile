build:
	mkdir -p build/bin
	go build -buildmode=pie -o build/bin/app.exe main.go

clean:
	rm -rf build

test:
	go test -v test/*.go
