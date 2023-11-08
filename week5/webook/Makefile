.PHONY: docker
docker:
	@rm webook || true
	@go mod tidy
	@GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -tags=k8s -o webook .
	@docker rmi -f flycash/webook:v0.0.1
	@docker build -t flycash/webook:v0.0.1 .
	@docker tag flycash/webook:v0.0.1 127.0.0.1:5000/flycash/webook:v0.0.1
	@docker push 127.0.0.1:5000/flycash/webook:v0.0.1
