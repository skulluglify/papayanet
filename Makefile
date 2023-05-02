build:
	mkdir -p build/bin
	go build -buildmode=pie -o build/bin/app.exe main.go

clean:
	rm -rf build

test:
	go test -v test/*.go

tidy:
	go mod tidy

push:
	bash beauty.sh record.sh
	git add .
	git commit -am 'heat/feat/update'
	git push

bpack-build:
	mkdir -p build/bin
	mv papaya/ant/bpack/data.go.bk papaya/ant/bpack/data.go
	go build -buildmode=pie -o build/bin/bpack.exe  bpack/main.go
	chmod a+x build/bin/bpack.exe

bpack-run: bpack-build
	./build/bin/bpack.exe
	mv papaya/ant/bpack/data.go papaya/ant/bpack/data.go.bk
	mv data.go papaya/ant/bpack/data.go