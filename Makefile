VERSION = 	v0.0.1-snapshot

all:
	@env GCO_ENABLED=0 go build -ldflags "-s -X main.version=$(VERSION)"

version:
	@echo $(VERSION)
