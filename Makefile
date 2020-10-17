all: build

build:
	go build ./...

build-win64:
	rm -rf dist/win64
	GOTRACEBACK=none GOOS=windows GOARCH=amd64 go build -o dist/win64/vvfst.exe main.go

build-win32:
	rm -rf dist/win32
	GOTRACEBACK=none GOOS=windows GOARCH=386 go build -o dist/win32/vvfst.exe main.go

build-linux64:
	rm -rf dist/linux64
	GOTRACEBACK=none GOOS=linux GOARCH=amd64 go build -o dist/linux64/vvfst.exe main.go

build-linux32:
	rm -rf dist/linux32
	GOTRACEBACK=none GOOS=linux GOARCH=386 go build -o dist/linux32/vvfst.exe main.go

build-osx:
	echo "Building for osx"
	rm -rf dist/osx
	GOTRACEBACK=none GOOS=darwin GOARCH=amd64 go build -o dist/osx/vvfst main.go

install:
	brew install golangci/tap/golangci-lint

lint:
	golangci-lint run --exclude=vendor --exclude=repos --disable-all --enable=golint --enable=vet --enable=gofmt ./...
	find . -name '*.go' | xargs gofmt -w -s

run:
	go run main.go
