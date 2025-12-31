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
echo "Note: KaaS tests require the following environment variables:"
echo "  - ACLOUD_VPC_URI (required)"
echo "  - ACLOUD_SUBNET_URI (required)"
echo "  - ACLOUD_NODE_POOL_INSTANCE (required)"
echo "  - ACLOUD_NODE_POOL_ZONE (required)"
echo "  - ACLOUD_NODE_CIDR (optional, default: 10.0.0.0/16)"
echo "  - ACLOUD_NODE_CIDR_NAME (optional, default: node-cidr)"
echo "  - ACLOUD_SECURITY_GROUP_NAME (optional, default: kaas-sg)"
echo "  - ACLOUD_NODE_POOL_NAME (optional, default: default-pool)"
echo "  - ACLOUD_NODE_POOL_NODES (optional, default: 1)"
echo "  - ACLOUD_K8S_VERSION (optional, default: 1.28.0)"
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
    
    # Required environment variables for KaaS creation
    local vpc_uri="${ACLOUD_VPC_URI}"
    local subnet_uri="${ACLOUD_SUBNET_URI}"
    local node_cidr_address="${ACLOUD_NODE_CIDR:-10.0.0.0/16}"
    local node_cidr_name="${ACLOUD_NODE_CIDR_NAME:-node-cidr}"
    local security_group_name="${ACLOUD_SECURITY_GROUP_NAME:-kaas-sg}"
    local node_pool_name="${ACLOUD_NODE_POOL_NAME:-default-pool}"
    local node_pool_nodes="${ACLOUD_NODE_POOL_NODES:-1}"
    local node_pool_instance="${ACLOUD_NODE_POOL_INSTANCE}"
    local node_pool_zone="${ACLOUD_NODE_POOL_ZONE}"
    
    echo -e "${BLUE}--- Testing KaaS Operations ---${NC}"
    
    # Check if required variables are set
    if [ -z "$vpc_uri" ] || [ -z "$subnet_uri" ] || [ -z "$node_pool_instance" ] || [ -z "$node_pool_zone" ]; then
        echo -e "${YELLOW}⚠ Skipping KaaS create test - required environment variables not set${NC}"
        echo "Required: ACLOUD_VPC_URI, ACLOUD_SUBNET_URI, ACLOUD_NODE_POOL_INSTANCE, ACLOUD_NODE_POOL_ZONE"
        echo "Optional: ACLOUD_NODE_CIDR, ACLOUD_NODE_CIDR_NAME, ACLOUD_SECURITY_GROUP_NAME, ACLOUD_NODE_POOL_NAME, ACLOUD_NODE_POOL_NODES"
        return
    fi
    
    # Create
    echo -e "${YELLOW}Creating KaaS cluster...${NC}"
    CREATE_OUTPUT=$($ACLOUD_CMD container kaas create \
        --project-id "$PROJECT_ID" \
        --name "$cluster_name" \
        --region "$REGION" \
        --vpc-uri "$vpc_uri" \
        --subnet-uri "$subnet_uri" \
        --node-cidr-address "$node_cidr_address" \
        --node-cidr-name "$node_cidr_name" \
        --security-group-name "$security_group_name" \
        --kubernetes-version "$version" \
        --node-pool-name "$node_pool_name" \
        --node-pool-nodes "$node_pool_nodes" \
        --node-pool-instance "$node_pool_instance" \
        --node-pool-zone "$node_pool_zone" \
        --tags "e2e-test,kaas" 2>&1)
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
            --name "${cluster_name}-updated" \
            --tags "e2e-test,kaas,updated" 2>&1)
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}✓ KaaS cluster update successful${NC}"
        else
            echo -e "${RED}✗ Failed to update KaaS cluster${NC}"
            echo "$UPDATE_OUTPUT"
        fi
        
        # Connect (optional - requires kubectl)
        if command -v kubectl >/dev/null 2>&1; then
            echo -e "${YELLOW}Testing KaaS connect...${NC}"
            CONNECT_OUTPUT=$($ACLOUD_CMD container kaas connect "$CLUSTER_ID" \
                --project-id "$PROJECT_ID" 2>&1)
            if [ $? -eq 0 ]; then
                echo -e "${GREEN}✓ KaaS connect successful${NC}"
            else
                echo -e "${YELLOW}⚠ KaaS connect failed (may be expected if cluster is not ready)${NC}"
                echo "$CONNECT_OUTPUT"
            fi
        else
            echo -e "${YELLOW}⚠ Skipping KaaS connect test - kubectl not found${NC}"
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

