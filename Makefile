# Makefile — LambdaOS CI/CD local parity
# Mirrors CI commands for lint, test, and build

.PHONY: lint lint-go test test-go build build-go release clean clean-go validate-specs

lint:
	black --check . && isort --check . && shellcheck **/*.sh && shfmt -d **/*.sh && luacheck .

lint-go:
	cd src/lambda-env && go vet ./...

validate-specs:
	./scripts/validate-specs.sh

test:
	python -m pytest tests/unit/ -v

test-go:
	cd src/lambda-env && go test ./... -v

build:
	@if [ "$$(id -u)" -ne 0 ]; then \
		echo "Warning: mkarchiso requires root privileges. Using sudo..."; \
	fi
	sudo mkarchiso -v -w work/ -o out/ .

release: build
	@echo "Release build complete. Artifacts in out/"

build-go:
	cd src/lambda-env && go build -o bin/lambda-env ./cmd/lambda-env

clean:
	sudo rm -rf work/ out/ .venv/

clean-go:
	rm -rf src/lambda-env/bin/
