MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --no-builtin-variables

.PHONY: build
build:
	mkdir -p build
	go build -mod=readonly -v -o build/ .

.PHONY: generate
generate:
	buf generate protos

.PHONY: test
test:
	go test -v ./... ./example/...

.PHONY: example
example: build
	buf --debug generate --template buf.example.gen.yaml --path example/bookstore

.PHONY: fmt
fmt:
	buf format -w 

.PHONY: lint
lint:
	buf lint ./protos

GOMOD_DIRS := .
ADDDEP_TARGET := $(addprefix adddep_, $(GOMOD_DIRS))
UPDATEDEP_TARGET := $(addprefix updatedep_, $(GOMOD_DIRS))

${ADDDEP_TARGET}:
	cd $(@:adddep_%=%) && go mod tidy -v

${UPDATEDEP_TARGET}:
	cd $(@:updatedep_%=%) && go get -d -u ./...
	cd $(@:updatedep_%=%) && go mod tidy -v


.PHONY: $(ADDDEP_TARGET)
adddep: $(ADDDEP_TARGET)
	go mod vendor


# .PHONY: adddep
# adddep:
# 	go mod tidy -v
# 	go mod vendor

.PHONY: ${UPDATEDEP_TARGET}
updatedeps: ${UPDATEDEP_TARGET}
	go mod vendor

