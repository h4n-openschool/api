tools:
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

openapi:
	oapi-codegen --config .oapi-codegen.yaml api/openapi.yaml
