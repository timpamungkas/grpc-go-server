ifeq ($(OS), Windows_NT)
	BIN_FILENAME  := grpc-server.exe
else
	BIN_FILENAME  := grpc-server
endif

.PHONY: tidy
tidy:
	go mod tidy


.PHONY: clean
clean:
ifeq ($(OS), Windows_NT)
	if exist "bin" rd /s /q bin	
else
	rm -fR ./bin
endif


.PHONY: build
build: clean
	go build -o ./bin/${BIN_FILENAME} ./main


.PHONY: execute
execute: clean build
	./bin/${BIN_FILENAME}