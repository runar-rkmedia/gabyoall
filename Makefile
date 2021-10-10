watch:
	cd frontend && yarn watch &
	find -E . -regex ".*\.(go)"  | entr -r  sh -c "go generate ./... & go run api/apiMain.go"
gen:
	go generate ./...
build:
	GOOS=linux GOARCH=amd64 go build -o dist/gobyoall-linux main.go
	GOOS=darwin go build -o dist/gobyoall-darwin main.go
	GOOS=windows go build -o dist/gobyoall-windows main.go
	
	GOOS=linux GOARCH=amd64 go build -o dist/gobyoall-server-linux api/apiMain.go
	GOOS=darwin go build -o dist/gobyoall-server-darwin api/apiMain.go
	GOOS=windows go build -o dist/gobyoall-server-windows api/apiMain.go