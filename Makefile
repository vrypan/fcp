HUBBLE_VER := "1.19.1"
FCP_VER := $(shell git describe --tags 2>/dev/null || echo "v0.0.0")

all:

proto:
	curl -s -L "https://github.com/farcasterxyz/hub-monorepo/archive/refs/tags/@farcaster/hubble@"${HUBBLE_VER}".tar.gz" \
	| tar -zxvf - -C . --strip-components 2 hub-monorepo--farcaster-hubble-${HUBBLE_VER}/protobufs/schemas/

farcaster-go: $(wildcard schemas/*.proto)
	@echo "Compiling .proto files"
	protoc --proto_path=schemas --go_out=. \
	$(shell cd schemas; ls | xargs -I \{\} echo -n '--go_opt=M'{}=farcaster/" " '--go-grpc_opt=M'{}=farcaster/" " ) \
	--go-grpc_out=. \
	schemas/*.proto

local:
	@echo Building tcp v${FCP_VER}
	go build \
	-ldflags "-w -s" \
	-ldflags "-X github.com/vrypan/fcp/main.FCP_VERSION=${FCP_VER}" \
	-o fcp

release-notes:
	# Autmatically generate release_notes.md
	./bin/generate_release_notes.sh

tag:
	./bin/auto_increment_tag.sh patch

tag-minor:
	./bin/auto_increment_tag.sh minor

tag-major:
	./bin/auto_increment_tag.sh major

releases:
	goreleaser release --clean
