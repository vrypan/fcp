SNAPCHAIN_VER := 0.2.20
FCP_VER := $(shell git describe --tags 2>/dev/null || echo "v0.0.0")

BINS = fcp fcp-inspect
COMMON_SOURCES := $(wildcard utils/*.go fctools/*.go)
FCP_SOURCES := $(wildcard cmd/fcp/*.go)
FCP_INSPECT_SOURCES := $(wildcard cmd/fcp-inspect/*.go)

# Colors for output
GREEN = \033[0;32m
NC = \033[0m

all: $(BINS)

clean:
	@echo -e "$(GREEN)Cleaning up...$(NC)"
	rm -f $(BINS)

.PHONY: all clean local release-notes tag tag-minor tag-major releases

fcp: $(FCP_SOURCES) $(COMMON_SOURCES)
	@echo -e "$(GREEN)Building fcp ${FCP_VER} $(NC)"
	go build -o $@ -ldflags "-w -s -X main.FCP_VERSION=${FCP_VER}" ./cmd/$@

fcp-inspect: $(FCP_INSPECT_SOURCES) $(COMMON_SOURCES)
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
