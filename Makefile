default: test

GO=go
GOTEST=$(GO) test
GOCOVER=$(GO) tool cover

.PHONY:
test:
	$(GOTEST) -v -race -failfast -parallel 4 -cpu 4 -coverprofile main.cover.out ./...
	$(GOCOVER) -html=main.cover.out
