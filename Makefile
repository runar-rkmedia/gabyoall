bob="Gjerne det"
repo=github.com/runar-rkmedia/gabyoall
version := $(shell git describe --tags)
gitHash := $(shell git rev-parse --short HEAD)
buildDate := $(shell TZ=UTC date +"%Y-%m-%dT%H:%M:%S")
ldflags=-X 'main.Version=$(version)' -X 'main.BuildDateStr=$(buildDate)' -X 'main.GitHash=$(gitHash)' -X 'main.IsDevStr=0'
watch:
	cd frontend && yarn watch &
	find -E . -regex ".*\.(go)"  | entr -r  sh -c "go generate ./... & go run api/apiMain.go"
gen:
	go generate ./...
build-server:
	go build -ldflags="${ldflags}" -o dist/gobyoall-server${SUFFIX} api/apiMain.go
build-cli:
	go build -ldflags="${ldflags}" -o dist/gobyoall${SUFFIX} main.go
clean:
	rm -rf dist
build:
	${MAKE} clean
	@GOOS=linux   GOARCH=amd64    SUFFIX="-linux-amd64"  ${MAKE} build-server
	@GOOS=darwin                  SUFFIX="-darwin"       ${MAKE} build-server
	@GOOS=windows                 SUFFIX=".exe"         ${MAKE} build-server

	@GOOS=linux   GOARCH=amd64    SUFFIX="-linux-amd64"  ${MAKE} build-cli
	@GOOS=darwin                  SUFFIX="-darwin"       ${MAKE} build-cli
	@GOOS=windows                 SUFFIX=".exe"         ${MAKE} build-cli
	ls -lah dist