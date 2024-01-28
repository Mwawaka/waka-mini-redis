tidy:
	@ go mod tidy
vendor:
	@ go mod vendor
build:
	@ go build -o ./bin/resp main.go
run:
	@ ./bin/resp
