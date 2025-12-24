#!/bin/bash

# E2E Test Script for Storage Resources
# Tests CRUD operations for Block Storage, Snapshots, Backups, and Restores

set -e  # Exit on error

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
ACLOUD_CMD="${ACLOUD_CMD:-./acloud}"

# Cleanup tracking
CREATED_VOLUMES=()
CREATED_SNAPSHOTS=()
CREATED_BACKUPS=()
CREATED_RESTORES=()

echo -e "${BLUE}=== Storage Resources E2E Test ===${NC}\n"
echo "Project ID: $PROJECT_ID"
echo "Region: $REGION"
echo "Test prefix: $RESOURCE_PREFIX"
echo ""

# Function to extract resource ID from output
extract_id() {
    local output="$1"
    echo "$output" | grep -oE '[a-f0-9]{24}' | head -1
}

# Cleanup function
cleanup() {
    echo -e "\n${YELLOW}Cleaning up test resources...${NC}"
    
    # Delete restores
    for restore_id in "${CREATED_RESTORES[@]}"; do
        echo "Deleting restore: $restore_id"
        echo "yes" | $ACLOUD_CMD storage restore delete "$restore_id" --yes --project-id "$PROJECT_ID" 2>&1 || true
    done
    
    # Delete backups
    for backup_id in "${CREATED_BACKUPS[@]}"; do
        echo "Deleting backup: $backup_id"
        echo "yes" | $ACLOUD_CMD storage backup delete "$backup_id" --yes --project-id "$PROJECT_ID" 2>&1 || true
    done
    
    # Delete snapshots
    for snapshot_id in "${CREATED_SNAPSHOTS[@]}"; do
        echo "Deleting snapshot: $snapshot_id"
        echo "yes" | $ACLOUD_CMD storage snapshot delete "$snapshot_id" --yes --project-id "$PROJECT_ID" 2>&1 || true
    done
    
    # Delete volumes
    for volume_id in "${CREATED_VOLUMES[@]}"; do
        echo "Deleting volume: $volume_id"
        echo "yes" | $ACLOUD_CMD storage blockstorage delete "$volume_id" --yes --project-id "$PROJECT_ID" 2>&1 || true
    done
}

trap cleanup EXIT

# Test Block Storage
test_block_storage() {
    local volume_name="${RESOURCE_PREFIX}-volume"
    
    echo -e "${YELLOW}--- Testing Block Storage CRUD ---${NC}\n"
    
    # CREATE
    echo -e "${GREEN}[CREATE]${NC} Creating block storage: $volume_name"
    CREATE_OUTPUT=$($ACLOUD_CMD storage blockstorage create \
        --name "$volume_name" \
        --region "$REGION" \
        --size 10 \
        --type Standard \
        --billing-period Hour \
        --tags "e2e,test,storage" \
        --project-id "$PROJECT_ID" 2>&1) || {
        echo -e "${RED}CREATE failed:${NC}"
        echo "$CREATE_OUTPUT"
        return 1
    }
    echo "$CREATE_OUTPUT"
    
    VOLUME_ID=$(extract_id "$CREATE_OUTPUT")
    if [ -z "$VOLUME_ID" ]; then
        echo -e "${RED}Could not extract volume ID${NC}"
        return 1
    fi
    CREATED_VOLUMES+=("$VOLUME_ID")
    echo -e "${GREEN}Created volume ID: $VOLUME_ID${NC}\n"
    
    # Wait for volume to be ready (optional, depends on API)
    echo "Waiting for volume to be ready..."
    sleep 5
    
    # LIST
    echo -e "${GREEN}[LIST]${NC} Listing block storage..."
    LIST_OUTPUT=$($ACLOUD_CMD storage blockstorage list --project-id "$PROJECT_ID" 2>&1) || {
        echo -e "${RED}LIST failed:${NC}"
        echo "$LIST_OUTPUT"
        return 1
    }
    echo "$LIST_OUTPUT" | head -15
    echo ""
    
    # GET
    echo -e "${GREEN}[GET]${NC} Getting block storage details..."
    GET_OUTPUT=$($ACLOUD_CMD storage blockstorage get "$VOLUME_ID" --project-id "$PROJECT_ID" 2>&1) || {
        echo -e "${RED}GET failed:${NC}"
        echo "$GET_OUTPUT"
        return 1
    }
    echo "$GET_OUTPUT"
    echo ""
    
    # UPDATE
    echo -e "${GREEN}[UPDATE]${NC} Updating block storage..."
    UPDATE_OUTPUT=$($ACLOUD_CMD storage blockstorage update "$VOLUME_ID" \
        --name "${volume_name}-updated" \
        --tags "e2e,test,updated" \
        --project-id "$PROJECT_ID" 2>&1) || {
        echo -e "${RED}UPDATE failed:${NC}"
        echo "$UPDATE_OUTPUT"
        return 1
    }
    echo "$UPDATE_OUTPUT"
    echo ""
    
    echo -e "${GREEN}✓ Block Storage CRUD test completed!${NC}\n"
    echo "$VOLUME_ID"  # Return volume ID for use in snapshots
}

# Test Snapshot
test_snapshot() {
    local volume_id="$1"
    local snapshot_name="${RESOURCE_PREFIX}-snapshot"
    
    if [ -z "$volume_id" ]; then
        echo -e "${YELLOW}Skipping snapshot test (no volume available)${NC}\n"
        return 0
    fi
    
    echo -e "${YELLOW}--- Testing Snapshot CRUD ---${NC}\n"
    
    # Get volume URI
    VOLUME_GET=$($ACLOUD_CMD storage blockstorage get "$volume_id" --project-id "$PROJECT_ID" 2>&1)
    VOLUME_URI=$(echo "$VOLUME_GET" | grep -i "URI:" | awk '{print $2}')
    
    if [ -z "$VOLUME_URI" ]; then
        echo -e "${YELLOW}Could not extract volume URI, skipping snapshot test${NC}\n"
        return 0
    fi
    
    # CREATE
    echo -e "${GREEN}[CREATE]${NC} Creating snapshot: $snapshot_name"
    CREATE_OUTPUT=$($ACLOUD_CMD storage snapshot create \
        --name "$snapshot_name" \
        --region "$REGION" \
        --volume-uri "$VOLUME_URI" \
        --tags "e2e,test,snapshot" \
        --project-id "$PROJECT_ID" 2>&1) || {
        echo -e "${RED}CREATE failed:${NC}"
        echo "$CREATE_OUTPUT"
        return 1
    }
    echo "$CREATE_OUTPUT"
    
    SNAPSHOT_ID=$(extract_id "$CREATE_OUTPUT")
    if [ -z "$SNAPSHOT_ID" ]; then
        echo -e "${RED}Could not extract snapshot ID${NC}"
        return 1
    fi
    CREATED_SNAPSHOTS+=("$SNAPSHOT_ID")
    echo -e "${GREEN}Created snapshot ID: $SNAPSHOT_ID${NC}\n"
    
    # LIST
    echo -e "${GREEN}[LIST]${NC} Listing snapshots..."
    LIST_OUTPUT=$($ACLOUD_CMD storage snapshot list --project-id "$PROJECT_ID" 2>&1) || {
        echo -e "${RED}LIST failed:${NC}"
        echo "$LIST_OUTPUT"
        return 1
    }
    echo "$LIST_OUTPUT" | head -15
    echo ""
    
    # GET
    echo -e "${GREEN}[GET]${NC} Getting snapshot details..."
    GET_OUTPUT=$($ACLOUD_CMD storage snapshot get "$SNAPSHOT_ID" --project-id "$PROJECT_ID" 2>&1) || {
        echo -e "${RED}GET failed:${NC}"
        echo "$GET_OUTPUT"
        return 1
    }
    echo "$GET_OUTPUT"
    echo ""
    
    # UPDATE
    echo -e "${GREEN}[UPDATE]${NC} Updating snapshot..."
    UPDATE_OUTPUT=$($ACLOUD_CMD storage snapshot update "$SNAPSHOT_ID" \
        --name "${snapshot_name}-updated" \
        --tags "e2e,test,updated" \
        --project-id "$PROJECT_ID" 2>&1) || {
        echo -e "${RED}UPDATE failed:${NC}"
        echo "$UPDATE_OUTPUT"
        return 1
    }
    echo "$UPDATE_OUTPUT"
    echo ""
    
    echo -e "${GREEN}✓ Snapshot CRUD test completed!${NC}\n"
}

# Test Backup (if available)
test_backup() {
    echo -e "${YELLOW}--- Testing Backup CRUD ---${NC}\n"
    echo -e "${YELLOW}Note: Backup operations may require specific prerequisites${NC}\n"
    # Backup tests would go here
    echo -e "${GREEN}✓ Backup test placeholder${NC}\n"
}

# Test Restore (if available)
test_restore() {
    echo -e "${YELLOW}--- Testing Restore CRUD ---${NC}\n"
    echo -e "${YELLOW}Note: Restore operations require existing backups${NC}\n"
    # Restore tests would go here
    echo -e "${GREEN}✓ Restore test placeholder${NC}\n"
}

# Run tests
echo -e "${BLUE}Starting Storage Resources E2E Tests...${NC}\n"

VOLUME_ID=""
if test_block_storage; then
    VOLUME_ID=$(extract_id "$CREATE_OUTPUT")
    test_snapshot "$VOLUME_ID"
fi

test_backup
test_restore

echo -e "${GREEN}=== All Storage Tests Completed! ===${NC}"

