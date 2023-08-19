all: run

build:
	go build -o bin/bifrost-env-manager

run: build
	./bin/bifrost-env-manager