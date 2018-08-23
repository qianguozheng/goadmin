SOURCES = main.go pprof.go

all : goadmin

goadmin: $(SOURCES)
	go build -x -o goadmin $(SOURCES)

clean:
	go clean -x
	rm -vf grab

check:
	go test -v .

install:
	go install -v .

.PHONY: all clean check install
