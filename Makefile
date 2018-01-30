PROJ:=sigridr-metrics
ORG_PATH:=git.nb.no/nettarkiv
REPO_PATH:=$(ORG_PATH)/$(PROJ)

## https://golang.org/cmd/link/
## -w Omit the DWARF symbol table.
LD_FLAGS="-w"

.PHONY: release-binary install install-dep

install:
	@CGO_ENABLED=0 go build -a -tags netgo -v -ldflags $(LD_FLAGS) $(REPO_PATH)

install-dep:
	@go get ./...

release-binary: install-dep install
