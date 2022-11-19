lint:
	@golangci-lint run

tools:
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

update-pkgs:
	go mod tidy
	go mod vendor

openapi:
	@oapi-codegen --config .oapi-codegen.yaml api/spec.yaml
