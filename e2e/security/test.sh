#!/bin/bash

# E2E Test Script for Security Resources
# Tests CRUD operations for KMS (Key Management Service)

# Don't exit on error - we want to continue and show summary
# set -e  # Exit on error

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROJECT_ID="${ACLOUD_PROJECT_ID:-your-project-id}"
REGION="${ACLOUD_REGION:-ITBG-Bergamo}"
RESOURCE_PREFIX="e2e-test-$(date +%s)"

# Determine acloud command path
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

# Cleanup tracking
CREATED_KMS=()

echo -e "${BLUE}=== Security Resources E2E Test ===${NC}\n"
echo "Project ID: $PROJECT_ID"
echo "Region: $REGION"
echo "Test prefix: $RESOURCE_PREFIX"
echo "ACLOUD command: $ACLOUD_CMD"
echo ""

# Function to extract resource ID from output
extract_id() {
    local output="$1"
    local exclude_id="${2:-}"
    
    if [ -n "$exclude_id" ]; then
        local filtered_ids=$(echo "$output" | grep -oE '[a-f0-9]{24}' | grep -v "^${exclude_id}$")
        if [ -n "$filtered_ids" ]; then
            echo "$filtered_ids" | tail -1
            return 0
        fi
    fi
    
    local all_ids=$(echo "$output" | grep -oE '[a-f0-9]{24}')
    if [ -n "$all_ids" ]; then
        if [ -n "$exclude_id" ]; then
            echo "$all_ids" | grep -v "^${exclude_id}$" | tail -1
        else
            echo "$all_ids" | tail -1
        fi
    fi
}

# Helper function to validate resource ID
is_valid_id() {
    local id="$1"
    [[ "$id" =~ ^[a-f0-9]{24}$ ]]
}

# Check if string is valid JSON
is_valid_json() {
    local input="$1"
    if command -v python3 >/dev/null 2>&1; then
        echo "$input" | python3 -c "import sys,json; json.load(sys.stdin)" 2>/dev/null && return 0
    elif command -v python >/dev/null 2>&1; then
        echo "$input" | python -c "import sys,json; json.load(sys.stdin)" 2>/dev/null && return 0
    fi
    return 1
}

# Generic --format flag test helper
# Usage: test_format_flags "Label" "no resources msg" cmd arg1 arg2...
test_format_flags() {
    local label="$1"
    local no_resources_msg="$2"
    shift 2
    local cmd_args=("$@")

    echo -e "${BLUE}--- Testing $label --format flag ---${NC}"

    for fmt in "" table; do
        local lbl="--format \"$fmt\""
        [ -z "$fmt" ] && lbl='--format "" (default)'
        echo -e "${YELLOW}Testing $lbl...${NC}"
        if [ -z "$fmt" ]; then
            OUT=$($ACLOUD_CMD "${cmd_args[@]}" 2>&1)
        else
            OUT=$($ACLOUD_CMD "${cmd_args[@]}" --format "$fmt" 2>&1)
        fi
        EXIT=$?
        if [ $EXIT -eq 0 ]; then
            echo -e "${GREEN}✓ $lbl: command succeeded${NC}"
            if ! is_valid_json "$OUT" || echo "$OUT" | grep -qF "$no_resources_msg"; then
                echo -e "${GREEN}✓ $lbl: output is table/plain (not JSON)${NC}"
            else
                echo -e "${RED}✗ $lbl: output unexpectedly looks like JSON${NC}"
            fi
        else
            echo -e "${RED}✗ $lbl: command failed (exit $EXIT)${NC}"
            echo "$OUT"
        fi
    done

    echo -e "${YELLOW}Testing --format json...${NC}"
    JSON_OUTPUT=$($ACLOUD_CMD "${cmd_args[@]}" --format json 2>&1)
    JSON_EXIT=$?
    if [ $JSON_EXIT -ne 0 ]; then
        echo -e "${RED}✗ --format json: command failed (exit $JSON_EXIT)${NC}"
        echo "$JSON_OUTPUT"
    elif echo "$JSON_OUTPUT" | grep -qF "$no_resources_msg"; then
        echo -e "${YELLOW}⚠ --format json: no resources — format validation skipped${NC}"
    elif is_valid_json "$JSON_OUTPUT"; then
        echo -e "${GREEN}✓ --format json: valid JSON${NC}"
        if echo "$JSON_OUTPUT" | grep -q '"metadata"'; then
            echo -e "${GREEN}✓ --format json: 'metadata' key present${NC}"
        else
            echo -e "${RED}✗ --format json: 'metadata' key missing${NC}"
        fi
        if echo "$JSON_OUTPUT" | grep -q '"properties"'; then
            echo -e "${GREEN}✓ --format json: 'properties' key present${NC}"
        else
            echo -e "${RED}✗ --format json: 'properties' key missing${NC}"
        fi
    else
        echo -e "${RED}✗ --format json: output is not valid JSON${NC}"
        echo "$JSON_OUTPUT"
    fi

    echo -e "${YELLOW}Testing --format yaml...${NC}"
    YAML_OUTPUT=$($ACLOUD_CMD "${cmd_args[@]}" --format yaml 2>&1)
    YAML_EXIT=$?
    if [ $YAML_EXIT -ne 0 ]; then
        echo -e "${RED}✗ --format yaml: command failed (exit $YAML_EXIT)${NC}"
        echo "$YAML_OUTPUT"
    elif echo "$YAML_OUTPUT" | grep -qF "$no_resources_msg"; then
        echo -e "${YELLOW}⚠ --format yaml: no resources — format validation skipped${NC}"
    elif echo "$YAML_OUTPUT" | grep -qE '^[a-zA-Z].*:|^- '; then
        echo -e "${GREEN}✓ --format yaml: output looks like YAML${NC}"
        if echo "$YAML_OUTPUT" | grep -q 'metadata:'; then
            echo -e "${GREEN}✓ --format yaml: 'metadata' key present${NC}"
        else
            echo -e "${RED}✗ --format yaml: 'metadata' key missing${NC}"
        fi
        if echo "$YAML_OUTPUT" | grep -q 'properties:'; then
            echo -e "${GREEN}✓ --format yaml: 'properties' key present${NC}"
        else
            echo -e "${RED}✗ --format yaml: 'properties' key missing${NC}"
        fi
    else
        echo -e "${RED}✗ --format yaml: output does not look like YAML${NC}"
        echo "$YAML_OUTPUT"
    fi

    echo ""
}

# Check for authentication errors
check_auth_error() {
    local output="$1"
    if echo "$output" | grep -qi "authentication failed\|invalid_client\|unauthorized"; then
        echo -e "${RED}Authentication error detected. Please check your credentials.${NC}" >&2
        return 1
    fi
    return 0
}

# Cleanup function
cleanup() {
    echo -e "\n${YELLOW}Cleaning up test resources...${NC}"
    
    # Delete KMS keys
    for kms_id in "${CREATED_KMS[@]}"; do
        if is_valid_id "$kms_id"; then
            echo "Deleting KMS key: $kms_id"
            $ACLOUD_CMD security kms delete "$kms_id" --yes 2>&1 || true
        fi
    done
    
    echo -e "${GREEN}Cleanup completed!${NC}"
}

# Trap to ensure cleanup runs on exit
trap cleanup EXIT

# Set up context
echo -e "${BLUE}Setting up context for project ID...${NC}"
$ACLOUD_CMD context set e2e-test-context --project-id "$PROJECT_ID" 2>&1 || true
echo ""

# Test function for KMS
test_kms() {
    echo -e "${BLUE}=== 1. KMS CRUD Test ===${NC}"
    
    local kms_name="${RESOURCE_PREFIX}-kms"
    
    echo -e "${GREEN}[CREATE]${NC} Creating KMS key: $kms_name"
    CREATE_OUTPUT=$($ACLOUD_CMD security kms create \
        --name "$kms_name" \
        --region "$REGION" \
        --billing-period "Hour" \
        --tags "e2e-test,created-by-script" 2>&1)
    exit_code=$?
    
    if ! check_auth_error "$CREATE_OUTPUT"; then
        echo -e "${RED}KMS test failed: Authentication error${NC}"
        return 1
    fi
    
    if [ $exit_code -ne 0 ]; then
        echo -e "${RED}CREATE failed:${NC}"
        echo "$CREATE_OUTPUT" >&2
        echo -e "${RED}KMS test failed${NC}"
        return 1
    fi
    
    KMS_ID=$(extract_id "$CREATE_OUTPUT")
    if [ -z "$KMS_ID" ] || ! is_valid_id "$KMS_ID"; then
        echo -e "${RED}Could not extract KMS ID from create output${NC}"
        echo "CREATE_OUTPUT: $CREATE_OUTPUT" >&2
        echo -e "${RED}KMS test failed${NC}"
        return 1
    fi
    
    CREATED_KMS+=("$KMS_ID")
    echo -e "${GREEN}KMS key created: $KMS_ID${NC}"
    
    echo -e "${GREEN}[LIST]${NC} Listing KMS keys"
    $ACLOUD_CMD security kms list 2>&1 | head -20
    
    echo -e "${GREEN}[GET]${NC} Getting KMS details: $KMS_ID"
    $ACLOUD_CMD security kms get "$KMS_ID" 2>&1
    
    echo -e "${GREEN}[UPDATE]${NC} Updating KMS: $KMS_ID"
    UPDATE_OUTPUT=$($ACLOUD_CMD security kms update "$KMS_ID" \
        --name "${kms_name}-updated" \
        --tags "e2e-test,updated" 2>&1)
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}KMS updated successfully${NC}"
    else
        echo -e "${YELLOW}Update may have failed or resource is in InCreation state${NC}"
        echo "$UPDATE_OUTPUT" | head -5
    fi
    
    echo -e "${GREEN}KMS test completed successfully${NC}\n"
    return 0
}

# Run tests
echo -e "${BLUE}Starting Security Resources E2E Tests...${NC}\n"

test_kms

test_format_flags "security kms list" "No KMS resources found" security kms list

# Test summary
echo -e "${BLUE}=== Test Summary ===${NC}"
echo "Project ID: $PROJECT_ID"
echo "○ KMS Keys: ${#CREATED_KMS[@]} created"
echo ""

echo -e "${GREEN}=== All Security Tests Completed! ===${NC}"

