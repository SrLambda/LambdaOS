# Makefile — LambdaOS CI/CD local parity
# Mirrors CI commands for lint, test, and build

.PHONY: lint test build release clean validate-specs

lint:
	black --check . && isort --check . && shellcheck **/*.sh && shfmt -d **/*.sh && luacheck .

validate-specs:
	./scripts/validate-specs.sh

test:
	python -m pytest tests/unit/ -v

build:
	@if [ "$$(id -u)" -ne 0 ]; then \
		echo "Warning: mkarchiso requires root privileges. Using sudo..."; \
	fi
	sudo mkarchiso -v -w work/ -o out/ .

release: build
	@echo "Release build complete. Artifacts in out/"

clean:
	sudo rm -rf work/ out/ .venv/
