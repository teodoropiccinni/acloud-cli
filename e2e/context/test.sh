#!/bin/bash

# E2E Test Script for context commands
# Validates --format json and --format yaml output structure
# against the fixture files list.json, list.yaml, use.json, use.yaml

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
if [ -f "$SCRIPT_DIR/../../acloud" ]; then
    ACLOUD_CMD="$SCRIPT_DIR/../../acloud"
elif [ -f "./acloud" ]; then
    ACLOUD_CMD="./acloud"
elif command -v acloud >/dev/null 2>&1; then
    ACLOUD_CMD="acloud"
else
    ACLOUD_CMD="${ACLOUD_CMD:-./acloud}"
fi

echo -e "${BLUE}=== Context Command E2E Test ===${NC}\n"
echo "ACLOUD command: $ACLOUD_CMD"
echo "Fixtures dir:   $SCRIPT_DIR"
echo ""

PASS=0
FAIL=0

pass() { echo -e "${GREEN}✓ $1${NC}"; PASS=$((PASS + 1)); }
fail() { echo -e "${RED}✗ $1${NC}"; FAIL=$((FAIL + 1)); }

# ---------------------------------------------------------------------------
# JSON key extraction (no external deps beyond python3/python/jq)
# ---------------------------------------------------------------------------
json_keys() {
    local input="$1"
    if command -v jq >/dev/null 2>&1; then
        echo "$input" | jq -r 'keys[]' 2>/dev/null
    elif command -v python3 >/dev/null 2>&1; then
        echo "$input" | python3 -c "import sys,json; [print(k) for k in json.load(sys.stdin).keys()]" 2>/dev/null
    elif command -v python >/dev/null 2>&1; then
        echo "$input" | python -c "import sys,json; [print(k) for k in json.load(sys.stdin).keys()]" 2>/dev/null
    fi
}

is_valid_json() {
    local input="$1"
    if command -v python3 >/dev/null 2>&1; then
        echo "$input" | python3 -c "import sys,json; json.load(sys.stdin)" 2>/dev/null
    elif command -v python >/dev/null 2>&1; then
        echo "$input" | python -c "import sys,json; json.load(sys.stdin)" 2>/dev/null
    else
        echo "$input" | grep -qE '^\{|^\[' 2>/dev/null
    fi
}

# Returns the top-level keys of the fixture JSON file
fixture_json_keys() {
    local fixture="$1"
    json_keys "$(cat "$fixture")"
}

# ---------------------------------------------------------------------------
# test_context_list_json
# ---------------------------------------------------------------------------
test_context_list_json() {
    echo -e "${YELLOW}--- context list --format json ---${NC}"
    local fixture="$SCRIPT_DIR/list.json"

    OUTPUT=$($ACLOUD_CMD context list --format json 2>&1)
    EXIT_CODE=$?
    if [ $EXIT_CODE -ne 0 ]; then
        fail "context list --format json: command failed (exit $EXIT_CODE)"
        echo "    Output: $OUTPUT"
        echo ""
        return
    fi

    # Must be valid JSON
    if ! is_valid_json "$OUTPUT"; then
        fail "context list --format json: output is not valid JSON"
        echo "    Output: $OUTPUT"
        echo ""
        return
    fi
    pass "context list --format json: output is valid JSON"

    # Validate that all top-level keys from the fixture are present in the output
    local missing=0
    while IFS= read -r key; do
        [ -z "$key" ] && continue
        if echo "$OUTPUT" | grep -q "\"$key\""; then
            pass "context list --format json: key '$key' present"
        else
            fail "context list --format json: key '$key' missing (expected from fixture)"
            missing=$((missing + 1))
        fi
    done < <(fixture_json_keys "$fixture")

    echo ""
}

# ---------------------------------------------------------------------------
# test_context_list_yaml
# ---------------------------------------------------------------------------
test_context_list_yaml() {
    echo -e "${YELLOW}--- context list --format yaml ---${NC}"
    local fixture="$SCRIPT_DIR/list.yaml"

    OUTPUT=$($ACLOUD_CMD context list --format yaml 2>&1)
    EXIT_CODE=$?
    if [ $EXIT_CODE -ne 0 ]; then
        fail "context list --format yaml: command failed (exit $EXIT_CODE)"
        echo "    Output: $OUTPUT"
        echo ""
        return
    fi

    if [ -z "$OUTPUT" ]; then
        fail "context list --format yaml: output is empty"
        echo ""
        return
    fi
    pass "context list --format yaml: output is non-empty"

    # Must NOT look like JSON
    if echo "$OUTPUT" | head -1 | grep -qE '^\{|^\['; then
        fail "context list --format yaml: output looks like JSON, not YAML"
        echo ""
        return
    fi
    pass "context list --format yaml: output does not start with JSON braces"

    # Validate that all top-level keys from the fixture are present as YAML keys
    while IFS= read -r line; do
        # Extract "key:" lines from the fixture (top-level, no leading spaces)
        key=$(echo "$line" | grep -E '^[a-zA-Z_-][a-zA-Z0-9_-]*:' | sed 's/:.*//')
        [ -z "$key" ] && continue
        if echo "$OUTPUT" | grep -qE "^${key}:"; then
            pass "context list --format yaml: key '$key' present"
        else
            fail "context list --format yaml: key '$key' missing (expected from fixture)"
        fi
    done < "$fixture"

    echo ""
}

# ---------------------------------------------------------------------------
# test_context_use_json
# ---------------------------------------------------------------------------
test_context_use_json() {
    echo -e "${YELLOW}--- context use --format json ---${NC}"
    local fixture="$SCRIPT_DIR/use.json"

    # Determine a context name to use from list output
    CONTEXT_NAME=$($ACLOUD_CMD context list --format json 2>/dev/null \
        | python3 -c "import sys,json; d=json.load(sys.stdin); print(list(d.get('contexts',{}).keys())[0])" 2>/dev/null \
        || $ACLOUD_CMD context list --format json 2>/dev/null \
        | python -c "import sys,json; d=json.load(sys.stdin); print(list(d.get('contexts',{}).keys())[0])" 2>/dev/null)

    if [ -z "$CONTEXT_NAME" ]; then
        echo -e "  ${YELLOW}Skipping: could not determine a context name to use${NC}"
        echo ""
        return
    fi

    OUTPUT=$($ACLOUD_CMD context use "$CONTEXT_NAME" --format json 2>&1)
    EXIT_CODE=$?
    if [ $EXIT_CODE -ne 0 ]; then
        fail "context use --format json: command failed (exit $EXIT_CODE)"
        echo "    Output: $OUTPUT"
        echo ""
        return
    fi

    if ! is_valid_json "$OUTPUT"; then
        fail "context use --format json: output is not valid JSON"
        echo "    Output: $OUTPUT"
        echo ""
        return
    fi
    pass "context use --format json: output is valid JSON"

    # Validate that all top-level keys from the fixture are present
    while IFS= read -r key; do
        [ -z "$key" ] && continue
        if echo "$OUTPUT" | grep -q "\"$key\""; then
            pass "context use --format json: key '$key' present"
        else
            fail "context use --format json: key '$key' missing (expected from fixture)"
        fi
    done < <(fixture_json_keys "$fixture")

    echo ""
}

# ---------------------------------------------------------------------------
# test_context_use_yaml
# ---------------------------------------------------------------------------
test_context_use_yaml() {
    echo -e "${YELLOW}--- context use --format yaml ---${NC}"
    local fixture="$SCRIPT_DIR/use.yaml"

    CONTEXT_NAME=$($ACLOUD_CMD context list --format json 2>/dev/null \
        | python3 -c "import sys,json; d=json.load(sys.stdin); print(list(d.get('contexts',{}).keys())[0])" 2>/dev/null \
        || $ACLOUD_CMD context list --format json 2>/dev/null \
        | python -c "import sys,json; d=json.load(sys.stdin); print(list(d.get('contexts',{}).keys())[0])" 2>/dev/null)

    if [ -z "$CONTEXT_NAME" ]; then
        echo -e "  ${YELLOW}Skipping: could not determine a context name to use${NC}"
        echo ""
        return
    fi

    OUTPUT=$($ACLOUD_CMD context use "$CONTEXT_NAME" --format yaml 2>&1)
    EXIT_CODE=$?
    if [ $EXIT_CODE -ne 0 ]; then
        fail "context use --format yaml: command failed (exit $EXIT_CODE)"
        echo "    Output: $OUTPUT"
        echo ""
        return
    fi

    if [ -z "$OUTPUT" ]; then
        fail "context use --format yaml: output is empty"
        echo ""
        return
    fi
    pass "context use --format yaml: output is non-empty"

    if echo "$OUTPUT" | head -1 | grep -qE '^\{|^\['; then
        fail "context use --format yaml: output looks like JSON, not YAML"
        echo ""
        return
    fi
    pass "context use --format yaml: output does not start with JSON braces"

    # Validate top-level keys from fixture are present
    while IFS= read -r line; do
        key=$(echo "$line" | grep -E '^[a-zA-Z_-][a-zA-Z0-9_-]*:' | sed 's/:.*//')
        [ -z "$key" ] && continue
        if echo "$OUTPUT" | grep -qE "^${key}:"; then
            pass "context use --format yaml: key '$key' present"
        else
            fail "context use --format yaml: key '$key' missing (expected from fixture)"
        fi
    done < "$fixture"

    echo ""
}

# ---------------------------------------------------------------------------
# Run
# ---------------------------------------------------------------------------
echo -e "${BLUE}Starting Context E2E Tests...${NC}\n"

test_context_list_json
test_context_list_yaml
test_context_use_json
test_context_use_yaml

echo -e "${BLUE}=== Test Summary ===${NC}"
echo -e "${GREEN}Passed: $PASS${NC}"
if [ $FAIL -gt 0 ]; then
    echo -e "${RED}Failed: $FAIL${NC}"
else
    echo -e "${GREEN}Failed: $FAIL${NC}"
fi
echo ""

if [ $FAIL -gt 0 ]; then
    exit 1
else
    exit 0
fi
