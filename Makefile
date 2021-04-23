default: test

GO=go
GOTEST=$(GO) test
GOCOVER=$(GO) tool cover

.PHONY: test
test: test/cover test/report

.PHONY: test/cover
test/cover:
	$(GOTEST) -v -race -failfast -parallel 4 -cpu 4 -coverprofile main.cover.out ./...

.PHONY: test/report
test/report:
	$(GOCOVER) -html=main.cover.out
