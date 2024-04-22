copy-config:
	cp application.yaml.sample application.yaml

run-dev:
	GIN_MODE=debug go run main.go

run-prod:
	GIN_MODE=release go run main.go

test:
	go clean -testcache && go test ./... -cover -coverprofile=coverage.out

cover-html:
	go tool cover -html=coverage.out

cover-total:
	go tool cover -func coverage.out
