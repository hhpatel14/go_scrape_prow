### Important, the LANG_DIR must match one from Makefile
LANG_DIR=go

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

$(LANG_DIR)-dev-env:
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


$(LANG_DIR)-clean-dev-env:
	sudo rm -rvf /usr/local/go/