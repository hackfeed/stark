GOCMD=go
GOBUILD=$(GOCMD) build

install:
	$(GOBUILD) -o build/stark.exe cmd/stark/main.go

clean:
	rm -rf build/*.exe