SRC_CODE = telegraf/src
unit-test:
	cd telegraf/src && \
	go build . && \
	go test -v -coverprofile cover.out ./ && \

code-coverage: 
	cd telegraf/src && \
	go build . && \
	go test -v -coverprofile cover.out ./ && \
	go tool cover -html=cover.out -o cover.html && \
	open cover.html

ifeq ($(OS),Windows_NT)
    OS_String = windows
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        OS_String = linux
    endif
    ifeq ($(UNAME_S),Darwin)
        OS_String = macos
    endif
endif

lint: 
 ifeq ($(OS_String),macos)
	wget -c https://golang.org/dl/go1.15.2.darwin-amd64.tar.gz && \
	shasum -a 256 go1.15.2.darwin-amd64.tar.gz && \
	sudo tar -C /usr/local -xvzf go1.15.2.darwin-amd64.tar.gz && \
	export  PATH=$PATH:/usr/local/go/bin
	export GOPATH="$HOME/go_projects"
	export GOBIN="$GOPATH/bin"
	export GOROOT=$HOME/go
	export PATH=$PATH:$GOROOT/bin	
 endif
