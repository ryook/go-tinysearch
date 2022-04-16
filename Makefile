ifdef update
	u=-u
endif

export GO11MODULE=on

.PHONY: deps
deps:
	go get ${u} -d
	go mod tidy

.PHONY: devel-deps
devel-deps:
	GO11MODULE=off go get ${u} \
	golang.org/x/tools/cmd/goiimports

.PHONY: test
test: deps devel-deps
		docker-compose up -d
		goiimports -l -w .
