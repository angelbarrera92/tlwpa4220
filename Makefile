build: clean
	@docker run --rm -v $(shell pwd):/home/tlwpa4220 -u $(shell id -u):$(shell id -g) -w /home/tlwpa4220 -e GOOS=linux -e GOARCH=amd64 -e GOCACHE=/home/tlwpa4220/.cache golang:1.17 go build -o tlwpa4220 cmd/main.go

lint:
	@docker run --rm -e RUN_LOCAL=true -v $(shell pwd):/tmp/lint github/super-linter:v4

clean:
	@rm -rf .cache
	@rm -rf tlwpa4220
	@rm -rf super-linter.log
