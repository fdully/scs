run:
	go run cmd/scs/main.go

deploy:
	gcloud app deploy ./deploy/app.yaml

lint:
	golangci-lint run ./...

test:
	go test -v -race ./... -count=1


.PHONY: lint deploy test