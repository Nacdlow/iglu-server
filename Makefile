GC = go

TARGET = nacdlow-server

all: clean format test $(TARGET)
full: clean bindata format test $(TARGET)

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

bindata:
	go-bindata -pkg templates -prefix "templates/" -o templates/bindata.go templates/...
	go-bindata -pkg public -prefix "public/" -o public/bindata.go public/...

bindata-dev:
	go-bindata -debug -pkg templates -prefix "templates/" -o templates/bindata.go templates/...
	go-bindata -debug -pkg public -prefix "public/" -o public/bindata.go public/...
