GC = go

TARGET = nacdlow-server

all: clean format test $(TARGET)

$(TARGET): main.go
	$(GC) build

format:
	$(GC) fmt ./...

test:
	$(GC) test -race ./...
	$(GC) vet ./...

clean:
	$(RM) $(TARGET)

sat:
	gocyclo -over 15 .
	golint ./...
	ineffassign .

