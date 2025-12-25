#!/bin/bash

# E2E Test Script for Container Resources
# Tests CRUD operations for KaaS (Kubernetes as a Service) clusters

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
CREATED_CLUSTERS=()

echo -e "${BLUE}=== Container Resources E2E Test ===${NC}\n"
echo "Project ID: $PROJECT_ID"
echo "Region: $REGION"
echo "Test prefix: $RESOURCE_PREFIX"
echo "ACLOUD command: $ACLOUD_CMD"
echo ""

# Function to extract resource ID from output
extract_id() {
    local output="$1"
    # Try to extract ID from table output or JSON
    echo "$output" | grep -oE '[a-f0-9]{24}' | tail -1 || echo ""
}

# Function to test KaaS operations
test_kaas() {
    local cluster_name="${RESOURCE_PREFIX}-kaas"
    local version="${ACLOUD_K8S_VERSION:-1.28.0}"
    
    echo -e "${BLUE}--- Testing KaaS Operations ---${NC}"
    
    # Create
    echo -e "${YELLOW}Creating KaaS cluster...${NC}"
    CREATE_OUTPUT=$($ACLOUD_CMD container kaas create \
        --project-id "$PROJECT_ID" \
        --name "$cluster_name" \
        --region "$REGION" \
        --version "$version" 2>&1)
    CREATE_EXIT=$?
    
    if [ $CREATE_EXIT -eq 0 ]; then
        CLUSTER_ID=$(extract_id "$CREATE_OUTPUT")
        if [ -n "$CLUSTER_ID" ]; then
            CREATED_CLUSTERS+=("$CLUSTER_ID")
            echo -e "${GREEN}✓ KaaS cluster created: $CLUSTER_ID${NC}"
        else
            echo -e "${YELLOW}⚠ KaaS cluster creation may have succeeded but ID not found${NC}"
        fi
    else
        echo -e "${RED}✗ Failed to create KaaS cluster${NC}"
        echo "$CREATE_OUTPUT"
    fi
    
    # List
    echo -e "${YELLOW}Listing KaaS clusters...${NC}"
    LIST_OUTPUT=$($ACLOUD_CMD container kaas list --project-id "$PROJECT_ID" 2>&1)
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ KaaS cluster list successful${NC}"
    else
        echo -e "${RED}✗ Failed to list KaaS clusters${NC}"
        echo "$LIST_OUTPUT"
    fi
    
    # Get (if we have an ID)
    if [ -n "$CLUSTER_ID" ]; then
        echo -e "${YELLOW}Getting KaaS cluster details...${NC}"
        GET_OUTPUT=$($ACLOUD_CMD container kaas get "$CLUSTER_ID" --project-id "$PROJECT_ID" 2>&1)
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}✓ KaaS cluster get successful${NC}"
        else
            echo -e "${RED}✗ Failed to get KaaS cluster${NC}"
            echo "$GET_OUTPUT"
        fi
        
        # Update
        echo -e "${YELLOW}Updating KaaS cluster...${NC}"
        UPDATE_OUTPUT=$($ACLOUD_CMD container kaas update "$CLUSTER_ID" \
            --project-id "$PROJECT_ID" \
            --name "${cluster_name}-updated" 2>&1)
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}✓ KaaS cluster update successful${NC}"
        else
            echo -e "${RED}✗ Failed to update KaaS cluster${NC}"
            echo "$UPDATE_OUTPUT"
        fi
    fi
    
    echo ""
}

# Cleanup function
cleanup() {
    echo -e "${BLUE}--- Cleanup ---${NC}"
    
    # Delete KaaS clusters
    for cluster_id in "${CREATED_CLUSTERS[@]}"; do
        echo -e "${YELLOW}Deleting KaaS cluster: $cluster_id${NC}"
        $ACLOUD_CMD container kaas delete "$cluster_id" --project-id "$PROJECT_ID" --yes 2>&1 >/dev/null
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}✓ KaaS cluster deleted: $cluster_id${NC}"
        else
            echo -e "${RED}✗ Failed to delete KaaS cluster: $cluster_id${NC}"
        fi
    done
    
    echo ""
}

# Trap to ensure cleanup on exit
trap cleanup EXIT

# Run tests
test_kaas

# Summary
echo -e "${BLUE}=== Test Summary ===${NC}"
echo "Created KaaS clusters: ${#CREATED_CLUSTERS[@]}"
echo ""
echo -e "${GREEN}E2E tests completed!${NC}"

