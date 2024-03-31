.PHONY: lint
lint:
	golangci-lint run

.PHONY: generate-mocks
generate-mocks:
	mockery

.PHONY: test
test:
	go clean -testcache
	go test ./... -race

.PHONY: docker-compose
docker-compose:
	docker-compose -f deployments/docker-compose/docker-compose.yaml up

.PHONY: docker-build
docker-build:
	docker build --tag request-counter:1.0.0 -f build/package/Dockerfile .

.PHONY: helm
helm: docker-build
	 helm install --replace requestcounter ./deployments/helm/requestcounter