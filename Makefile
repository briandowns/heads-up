GO ?= go
GLIDE ?= glide
GOFLAGS = CGO_ENABLED=0
LINTER ?= golint

BINDIR := bin
BINARY := heads-up

VERSION := 0.1.0
LDFLAGS = -ldflags "-X main.gitSHA=$(shell git rev-parse HEAD) -X main.version=$(VERSION) -X main.name=$(BINARY)"

OS := $(shell uname)

.PHONY:
$(BINDIR)/$(BINARY): clean
	if [ ! -d $(BINDIR) ]; then mkdir $(BINDIR); fi
ifeq ($(OS),Darwin)
	GOOS=darwin $(GOFLAGS) $(GO) build -v -o $(BINDIR)/$(BINARY) $(LDFLAGS)
endif
ifeq ($(OS),Linux)
	GOOS=linux $(GOFLAGS) $(GO) build -v -o $(BINDIR)/$(BINARY) $(LDFLAGS)
endif

.PHONY:
test:
	$(GO) test -v -cover ./...

.PHONY:
deps:
ifeq (,$(wildcard glide.yaml))
	$(GLIDE) init
else
	$(GLIDE) update
endif

.PHONY:
clean:
	$(GO) clean
	rm -f $(BINDIR)/$(BINARY)

.PHONY:
docs:
	@godoc -http=:6060 2>/dev/null &
	@printf "To view heads-up docs, point your browser to:\n"
	@printf "\n\thttp://127.0.0.1:6060/pkg/github.com/briandowns/$(BINARY)/$(pkg)\n\n"
	@sleep 1
	@open "http://127.0.0.1:6060/pkg/github.com/briandowns/$(BINARY)/$(pkg)"

.PHONY:
lint:
	$(LINTER) `$(GO) list ./... | grep -v /vendor/`
