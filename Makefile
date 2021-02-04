GOCMD = go
GOBUILD = $(GOCMD) build
GORUN = $(GOCMD) run
GOCLEAN = $(GOCMD) clean

all: build

build:
	$(GOBUILD) -gcflags -m -v .

run:
	$(GORUN) -gcflags -m -v .

clean:
	$(GOCLEAN)
