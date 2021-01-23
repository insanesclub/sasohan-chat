GOCMD = go
GOBUILD = $(GOCMD) build
GORUN = $(GOCMD) run
GOCLEAN = $(GOCMD) clean
BINARY_FILE = bin/chat
RM = rm -f

all: run

build:
	$(GOBUILD) -o $(BINARY_FILE) -v .

run:
	$(GORUN) -v .

clean:
	$(RM) $(BINARY_FILE)
