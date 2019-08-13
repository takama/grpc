PROTOC_IMAGE ?= takama/protoc:v0.1.0 

PROTOC_RUN = docker run --rm -v $(shell pwd):/contracts -w /contracts $(PROTOC_IMAGE) -I.

all: generate

generate: info

info:
	@echo "+ $@"
	@$(PROTOC_RUN) --go_out=plugins=grpc,paths=source_relative:. info/info.proto
	@$(PROTOC_RUN) --grpc-gateway_out=logtostderr=true,paths=source_relative:. info/info.proto
	@$(PROTOC_RUN) --swagger_out=logtostderr=true:. info/info.proto


.PHONY: all \
	info \
	generate
