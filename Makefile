APP?=network_exporter

build: clean
	go build -o bin/${APP} cmd/${APP}/main.go

clean:
	@rm -f bin/${APP}
