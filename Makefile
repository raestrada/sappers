build:
	go build -o bin/sappers main.go

run:
	go run main.go

compile:
	@echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=386 go build -o bin/sappers-linux-386 main.go

run-cluster: compile
	@echo "Compiling for every OS and Platform"
	curl https://ops.city/get.sh -sSfL | sh
	cp bin/sappers-linux-386 deploy/sappers

all: build compile
