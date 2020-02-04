GO_CLEAN := go clean
GO_TEST := go test
GO_TOOL := go tool
GO_BUILD := GOOS=linux go build
TARGET := skeleton
TARGET_ZIP := $(TARGET).zip

clean:
	go clean
	rm -f $(TARGET)
	rm -f $(TARGET_ZIP)

test:
	$(GO_TEST) -v

build:
	$(GO_BUILD) -v -gcflags "-N -l" -o $(TARGET)

install: clean build
	zip $(TARGET_ZIP) $(TARGET)