default: test

SHELL=/usr/bin/env bash
GO=go
GOTEST=$(GO) test
GOCOVER=$(GO) tool cover

.PHONY: test
test: test/cover test/report

.PHONY: test/cover
test/cover: test/httpin test/echo

.PHONY: test/httpin
test/httpin:
	$(GOTEST) -v -race -failfast -parallel 4 -cpu 4 -coverprofile httpin.cover.out ./...

.PHONY: test/echo
test/echo:
	cd integration/echo \
		&& $(GOTEST) -v -race -failfast -parallel 4 -cpu 4 -coverprofile echo.cover.out ./...

.PHONY: test/report
test/report:
	if [[ "$$HOSTNAME" =~ "codespaces-"* ]]; then \
		mkdir -p /tmp/httpin_test; \
		$(GOCOVER) -html=httpin.cover.out -o /tmp/httpin_test/coverage-httpin.html; \
		$(GOCOVER) -html=integration/echo/echo.cover.out -o /tmp/httpin_test/coverage-echo.html; \
		sudo python -m http.server -d /tmp/httpin_test -b localhost 80; \
	else \
		$(GOCOVER) -html=httpin.cover.out; \
		$(GOCOVER) -html=integration/echo/echo.cover.out; \
	fi
