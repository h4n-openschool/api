tools:
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

update-pkgs:
	go mod tidy
	go mod vendor

