VERSION := $(shell git tag)
BUILD := $(shell git rev-parse --short HEAD) 
PROJECTNAME := $(shell basename "$(PWD)")
EXT:=.exe
LDFLAGS=-ldflags "-s -w -X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

.PHONY : clean win linux all

all: clean win linux
clean:
	@rm -f $(PROJECTNAME)*
	@rm -rf extracted
win:
	@env GOOS=windows GOARCH=386 CGO_ENABLED=0 go build $(LDFLAGS) -o $(PROJECTNAME)$(EXT)
linux:
	@env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(PROJECTNAME)

