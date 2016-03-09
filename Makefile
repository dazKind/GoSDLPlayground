.PHONY: all build dep

build:
	export GOPATH=$(shell pwd) && go build src/main.go

run:
	export GOPATH=$(shell pwd) && go run src/main.go

deps:
	export GOPATH=$(shell pwd) && go get -v github.com/veandco/go-sdl2/sdl
	export GOPATH=$(shell pwd) && go get -v github.com/veandco/go-sdl2/sdl_mixer
	export GOPATH=$(shell pwd) && go get -v github.com/veandco/go-sdl2/sdl_image
	export GOPATH=$(shell pwd) && go get -v github.com/veandco/go-sdl2/sdl_ttf
	export GOPATH=$(shell pwd) && go get -v github.com/chsc/gogl/gl33