SNAPCHAIN_VER := 0.2.20
FCP_VER := $(shell git describe --tags 2>/dev/null || echo "v0.0.0")

BINS = fcp fcp-inspect
PROTO_FILES := $(wildcard schemas/*.proto)
COMMON_SOURCES := $(wildcard utils/*.go fctools/*.go)
FCP_SOURCES := $(wildcard cmd/fcp/*.go)
FCP_INSPECT_SOURCES := $(wildcard cmd/fcp-inspect/*.go)

# Colors for output
GREEN = \033[0;32m
NC = \033[0m

all: $(BINS)

# Build binaries, depends on compiled protos
# $(BINS): .farcaster-built
#	@echo -e "$(GREEN)Building $@...$(NC)"
#	go build -o $@ ./cmd/$@

# Compile .proto files, touch stamp file
.farcaster-built: $(PROTO_FILES) .protobufs-downloaded
	@echo -e "$(GREEN)Compiling .proto files...$(NC)"
	protoc --proto_path=proto --go_out=. \
	$(shell cd proto; ls | xargs -I \{\} echo -n '--go_opt=M'{}=farcaster/" " '--go-grpc_opt=M'{}=farcaster/" " ) \
	--go-grpc_out=. \
	proto/*.proto
	@touch .farcaster-built

.protobufs-downloaded:
	@echo -e "$(GREEN)Downloading proto files (Hubble v$(SNAPCHAIN_VER))...$(NC)"
	curl -s -L "https://codeload.github.com/farcasterxyz/snapchain/tar.gz/refs/tags/v$(SNAPCHAIN_VER)" \
	| tar -zxvf - -C . --strip-components 2 "snapchain-$(SNAPCHAIN_VER)/src/proto/"
	@touch .protobufs-downloaded

clean-proto:
	@echo -e "$(GREEN)Cleaning up protobuf definitions...$(NC)"
	rm -f $(BINS) proto/*.proto .protobufs-downloaded

clean:
	@echo -e "$(GREEN)Cleaning up...$(NC)"
	rm -f $(BINS) farcaster/*.pb.go farcaster/*.pb.gw.go .farcaster-built

.PHONY: all proto clean local release-notes tag tag-minor tag-major releases

fcp: .farcaster-built $(FCP_SOURCES) $(COMMON_SOURCES)
	@echo -e "$(GREEN)Building fcp ${FCP_VER} $(NC)"
	go build -o $@ -ldflags "-w -s -X main.FCP_VERSION=${FCP_VER}" ./cmd/$@

fcp-inspect: .farcaster-built $(FCP_INSPECT_SOURCES) $(COMMON_SOURCES)
	@echo -e "$(GREEN)Building fcp-inspect ${FCP_VER} $(NC)"
	go build -o $@ -ldflags "-w -s -X main.FCP_VERSION=${FCP_VER}" ./cmd/$@

release-notes:
	# Automatically generate release_notes.md
	./bin/generate_release_notes.sh

tag:
	./bin/auto_increment_tag.sh patch

tag-minor:
	./bin/auto_increment_tag.sh minor

tag-major:
	./bin/auto_increment_tag.sh major

releases:
	goreleaser release --clean
