#!/usr/bin/env bash
# tests/unit/test_versioning.sh
# Tests profiledef.sh version resolution: env var → tag → describe → hash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
PROFILEDEF="${PROJECT_ROOT}/profiledef.sh"

PASSED=0
FAILED=0

# Extract the version resolution logic from profiledef.sh into a temp script
# and test it in isolation.
extract_version_logic() {
    local tmpfile
    tmpfile=$(mktemp)

    # We look for lines that set iso_version.
    # The design expects either a one-liner or an if-elif-else block.
    # We grep the relevant lines and write them to a temp script.
    grep -n "iso_version" "${PROFILEDEF}" | head -20 > "${tmpfile}.debug" || true

    # Strategy: write a minimal script that sources just the variable assignments
    # we care about. Since the full file has a problematic array, we isolate.
    cat > "${tmpfile}" << 'EOF'
#!/usr/bin/env bash
set -euo pipefail

# The version logic will be injected here by the test
EOF

    # Extract lines between the first occurrence of iso_version logic
    # and the end of that block. This handles both one-liner and multi-line.
    local in_block=false
    while IFS= read -r line; do
        if [[ "${line}" == *iso_version*"="* ]]; then
            in_block=true
        fi
        if [[ "${in_block}" == true ]]; then
            echo "${line}" >> "${tmpfile}"
            # Heuristic: stop after a line that doesn't end with continuation
            # and isn't an if/elif/else keyword
            if [[ ! "${line}" =~ \$ ]] && [[ ! "${line}" =~ \|\| ]] && [[ ! "${line}" == *"fi"* ]]; then
                break
            fi
            # Also stop at 'fi' for if-blocks
            if [[ "${line}" == *"fi"* ]]; then
                break
            fi
        fi
    done < "${PROFILEDEF}"

    echo "${tmpfile}"
}

# Actually, a simpler and more robust approach:
# We grep iso_version lines and eval them in a controlled subshell.
get_iso_version() {
    # We source a dynamically generated script that contains only the version logic
    # extracted from profiledef.sh.
    local tmpfile
    tmpfile=$(mktemp)

    # Write a wrapper that only defines the version logic
    {
        echo '#!/usr/bin/env bash'
        echo 'set -euo pipefail'
        # Extract all lines containing iso_version from profiledef.sh
        grep -n "iso_version" "${PROFILEDEF}" | while IFS= read -r line; do
            echo "# ${line}"
        done
        # For now, we just check the raw assignment line
        grep "^iso_version=" "${PROFILEDEF}" || true
        grep "^iso_version " "${PROFILEDEF}" || true
    } > "${tmpfile}"

    (
        # shellcheck source=/dev/null
        source "${tmpfile}"
        echo "${iso_version:-UNDEFINED}"
    )

    rm -f "${tmpfile}"
}

# A better approach: create a synthetic script that mirrors the expected logic
# and compare what profiledef.sh actually contains.
# But for TDD, we want to test the ACTUAL file.

# Extract the version resolution block from profiledef.sh and source it in a subshell.
get_iso_version_direct() {
    local tmpfile
    tmpfile=$(mktemp)
    {
        echo '#!/usr/bin/env bash'
        echo 'set -euo pipefail'
        sed -n '/# Version resolution/,/^fi$/p' "${PROFILEDEF}"
    } > "${tmpfile}"
    (
        # shellcheck source=/dev/null
        source "${tmpfile}"
        echo "${iso_version}"
    )
    rm -f "${tmpfile}"
}

test_env_var_priority() {
    local name="env var priority"
    local result
    result="$(LAMBDAOS_VERSION="1.5.0" get_iso_version_direct)"
    if [[ "${result}" == "1.5.0" ]]; then
        echo "PASS: ${name}"
        PASSED=$((PASSED + 1))
    else
        echo "FAIL: ${name} — expected '1.5.0', got '${result}'"
        FAILED=$((FAILED + 1))
    fi
}

test_no_date_in_version() {
    local name="no date-based version"
    local result
    result="$(get_iso_version_direct)"
    # Date-based versions look like 2026.05.24 or 202605
    if [[ "${result}" =~ ^[0-9]{4}\.[0-9]{2}\.[0-9]{2}$ ]] || [[ "${result}" =~ ^[0-9]{6}$ ]]; then
        echo "FAIL: ${name} — looks date-based: '${result}'"
        FAILED=$((FAILED + 1))
    else
        echo "PASS: ${name} (not date-based: '${result}')"
        PASSED=$((PASSED + 1))
    fi
}

test_not_empty() {
    local name="version is not empty"
    local result
    result="$(get_iso_version_direct)"
    if [[ -n "${result}" && "${result}" != "UNDEFINED" ]]; then
        echo "PASS: ${name} (value: '${result}')"
        PASSED=$((PASSED + 1))
    else
        echo "FAIL: ${name} — version is empty or undefined"
        FAILED=$((FAILED + 1))
    fi
}

test_v_prefix_stripped() {
    local name="v prefix stripped from tag"
    # This test checks the logic itself, not the current git state.
    # We verify that if a tag v1.0.0 is present, the version logic strips the v.
    local tag
    tag=$(git describe --tags --exact-match 2> /dev/null || echo "")
    if [[ -n "${tag}" && "${tag}" == v* ]]; then
        local result
        result="$(get_iso_version_direct)"
        if [[ "${result}" == "${tag#v}" ]]; then
            echo "PASS: ${name} (tag '${tag}' → version '${result}')"
            PASSED=$((PASSED + 1))
        else
            echo "FAIL: ${name} — tag '${tag}' should produce '${tag#v}', got '${result}'"
            FAILED=$((FAILED + 1))
        fi
    else
        echo "SKIP: ${name} (no v-prefixed tag on current commit)"
    fi
}

echo "=== Version Resolution Tests ==="
echo "Testing: ${PROFILEDEF}"
echo ""

test_env_var_priority
test_no_date_in_version
test_not_empty
test_v_prefix_stripped

echo ""
echo "=== Results ==="
echo "Passed: ${PASSED}"
echo "Failed: ${FAILED}"

if [[ ${FAILED} -gt 0 ]]; then
    exit 1
fi
