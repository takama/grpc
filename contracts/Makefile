PROTOC_IMAGE ?= takama/protoc:v0.1.8

PROTOC_RUN = docker run --rm -v $(shell pwd):/contracts -w /contracts $(PROTOC_IMAGE) -I.

all: generate

generate: info echo

info:
	@echo "+ $@"
	@$(PROTOC_RUN) --go_out=plugins=grpc,paths=source_relative:. info/info.proto
	@$(PROTOC_RUN) --grpc-gateway_out=logtostderr=true,paths=source_relative:. info/info.proto
	@$(PROTOC_RUN) --swagger_out=logtostderr=true:. info/info.proto

echo:
	@echo "+ $@"
	@$(PROTOC_RUN) --go_out=plugins=grpc,paths=source_relative:. echo/echo.proto
	@$(PROTOC_RUN) --grpc-gateway_out=logtostderr=true,paths=source_relative:. echo/echo.proto
	@$(PROTOC_RUN) --swagger_out=logtostderr=true:. echo/echo.proto

.PHONY: all \
	info \
	echo \
	generate
