all: build

build:
	go build ./...

build-all: build-win64 build-win32 build-linux32 build-linux64 build-osx

build-win64:
	echo "Building for win64"
	rm -rf dist/win64
	GOOS=windows GOARCH=amd64 go build -o dist/win64/vvfst.exe main.go

build-win32:
	echo "Building for win32"
	rm -rf dist/win32
	GOOS=windows GOARCH=386 go build -o dist/win32/vvfst.exe main.go

build-linux64:
	echo "Building for linux64"
	rm -rf dist/linux64
	GOOS=linux GOARCH=amd64 go build -o dist/linux64/vvfst.exe main.go

build-linux32:
	echo "Building for linux32"
	rm -rf dist/linux32
	GOOS=linux GOARCH=386 go build -o dist/linux32/vvfst.exe main.go

build-osx:
	echo "Building for osx"
	rm -rf dist/osx
	GOOS=darwin GOARCH=amd64 go build -o dist/osx/vvfst main.go

install:
	brew install golangci/tap/golangci-lint

lint:
	golangci-lint run --exclude=vendor --exclude=repos --disable-all --enable=golint --enable=vet --enable=gofmt ./...
	find . -name '*.go' | xargs gofmt -w -s

run:
	go run main.go
