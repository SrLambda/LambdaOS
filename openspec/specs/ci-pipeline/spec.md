# CI Pipeline Specification

## Purpose

Automated linting and unit test execution on every push and pull request to catch regressions before merge.

## Requirements

### Requirement: CI Trigger on Push and PR

The CI pipeline SHALL execute on every push to `main` and on every pull request targeting `main`.

#### Scenario: Push to main branch

- GIVEN a commit is pushed to the `main` branch
- WHEN the push event occurs
- THEN the CI workflow SHALL start automatically
- AND all lint and test jobs SHALL execute

#### Scenario: Pull request opened or updated

- GIVEN a pull request targets the `main` branch
- WHEN the PR is opened, synchronized, or reopened
- THEN the CI workflow SHALL start automatically

#### Scenario: Push to non-main branch without PR

- GIVEN a commit is pushed to a feature branch with no open PR
- WHEN the push event occurs
- THEN the CI workflow SHALL NOT trigger

### Requirement: Python Linting

The CI pipeline SHALL enforce Python code formatting with `black --check` and import ordering with `isort --check`.

#### Scenario: Python files pass linting

- GIVEN Python files are formatted with black and isort
- WHEN the CI pipeline runs `black --check` and `isort --check`
- THEN both commands SHALL exit with code 0
- AND the CI job SHALL pass

#### Scenario: Python files fail linting

- GIVEN at least one Python file violates black or isort rules
- WHEN the CI pipeline runs `black --check` or `isort --check`
- THEN the failing command SHALL exit with non-zero code
- AND the CI job SHALL fail with a descriptive error

### Requirement: Shell Linting

The CI pipeline SHALL validate shell scripts with `shellcheck` and `shfmt`.

#### Scenario: Shell scripts pass linting

- GIVEN all shell scripts comply with shellcheck and shfmt rules
- WHEN the CI pipeline runs `shellcheck` and `shfmt --diff`
- THEN both commands SHALL exit with code 0
- AND the CI job SHALL pass

#### Scenario: Shell scripts fail linting

- GIVEN at least one shell script violates shellcheck or shfmt rules
- WHEN the CI pipeline runs `shellcheck` or `shfmt --diff`
- THEN the failing command SHALL exit with non-zero code
- AND the CI job SHALL fail

### Requirement: Unit Test Execution

The CI pipeline SHALL execute all unit tests in `tests/unit/` using pytest.

#### Scenario: All unit tests pass

- GIVEN the codebase is in a valid state
- WHEN the CI pipeline runs `pytest tests/unit/ -v`
- THEN all tests SHALL pass (exit code 0)
- AND the test output SHALL include per-test results

#### Scenario: Unit test regression

- GIVEN a code change breaks one or more unit tests
- WHEN the CI pipeline runs `pytest tests/unit/ -v`
- THEN pytest SHALL exit with non-zero code
- AND the CI job SHALL fail
- AND the failing test names SHALL be visible in the output

### Requirement: CI Job Dependency Order

The CI pipeline SHALL run linting and unit tests as independent jobs that MAY run in parallel.

#### Scenario: Parallel execution

- GIVEN a push or PR event triggers CI
- WHEN the workflow starts
- THEN lint and test jobs SHALL start concurrently
- AND the workflow status SHALL reflect the aggregate result

#### Scenario: Partial failure

- GIVEN lint passes but unit tests fail
- WHEN both jobs complete
- THEN the overall workflow status SHALL be failure
