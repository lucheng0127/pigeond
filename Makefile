GOCMD=go

all: build

build:
	$(GOCMD) build -o pigeond pigeond.go

gcflag:
	$(GOCMD) build -gcflags "-N -l" -o pigeond pigeond.go

clean:
	rm -f pigeond
