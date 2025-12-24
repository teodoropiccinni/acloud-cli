#!/bin/bash

# E2E Test Script for Network Resources
# Tests CRUD operations for VPC, Subnet, Security Group, Security Rule,
# Elastic IP, VPC Peering, VPC Peering Route, VPN Tunnel, and VPN Route

set -e  # Exit on error

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration - UPDATE THESE VALUES
PROJECT_ID="${ACLOUD_PROJECT_ID:-your-project-id}"
VPC_ID="${ACLOUD_VPC_ID:-your-vpc-id}"
PEER_VPC_ID="${ACLOUD_PEER_VPC_ID:-your-peer-vpc-id}"
REGION="${ACLOUD_REGION:-ITBG-Bergamo}"
ELASTIC_IP_URI="${ACLOUD_ELASTIC_IP_URI:-your-elastic-ip-uri}"
ACLOUD_CMD="${ACLOUD_CMD:-./acloud}"

# Derived values
PEER_VPC_URI="${ACLOUD_PEER_VPC_URI:-/projects/${PROJECT_ID}/providers/Aruba.Network/vpcs/${PEER_VPC_ID}}"
RESOURCE_PREFIX="e2e-test-$(date +%s)"

# Cleanup tracking
CREATED_VPCS=()
CREATED_SUBNETS=()
CREATED_SECURITY_GROUPS=()
CREATED_SECURITY_RULES=()
CREATED_ELASTIC_IPS=()
CREATED_PEERINGS=()
CREATED_PEERING_ROUTES=()
CREATED_VPN_TUNNELS=()
CREATED_VPN_ROUTES=()

echo -e "${BLUE}=== Network Resources E2E Test ===${NC}\n"
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
    
    # Delete VPN routes
    for route_id in "${CREATED_VPN_ROUTES[@]}"; do
        echo "Deleting VPN route: $route_id"
        echo "yes" | $ACLOUD_CMD network vpnroute delete "$VPN_TUNNEL_ID" "$route_id" --yes --project-id "$PROJECT_ID" 2>&1 || true
    done
    
    # Delete VPN tunnels
    for tunnel_id in "${CREATED_VPN_TUNNELS[@]}"; do
        echo "Deleting VPN tunnel: $tunnel_id"
        echo "yes" | $ACLOUD_CMD network vpntunnel delete "$tunnel_id" --yes --project-id "$PROJECT_ID" 2>&1 || true
    done
    
    # Delete peering routes
    for route_id in "${CREATED_PEERING_ROUTES[@]}"; do
        echo "Deleting peering route: $route_id"
        echo "yes" | $ACLOUD_CMD network vpcpeeringroute delete "$VPC_ID" "$PEERING_ID" "$route_id" --yes --project-id "$PROJECT_ID" 2>&1 || true
    done
    
    # Delete peerings
    for peering_id in "${CREATED_PEERINGS[@]}"; do
        echo "Deleting VPC peering: $peering_id"
        echo "yes" | $ACLOUD_CMD network vpcpeering delete "$VPC_ID" "$peering_id" --yes --project-id "$PROJECT_ID" 2>&1 || true
    done
    
    # Delete security rules
    for rule_id in "${CREATED_SECURITY_RULES[@]}"; do
        echo "Deleting security rule: $rule_id"
        echo "yes" | $ACLOUD_CMD network securityrule delete "$VPC_ID" "$SECURITY_GROUP_ID" "$rule_id" --yes --project-id "$PROJECT_ID" 2>&1 || true
    done
    
    # Delete security groups
    for sg_id in "${CREATED_SECURITY_GROUPS[@]}"; do
        echo "Deleting security group: $sg_id"
        echo "yes" | $ACLOUD_CMD network securitygroup delete "$VPC_ID" "$sg_id" --yes --project-id "$PROJECT_ID" 2>&1 || true
    done
    
    # Delete elastic IPs
    for eip_id in "${CREATED_ELASTIC_IPS[@]}"; do
        echo "Deleting elastic IP: $eip_id"
        echo "yes" | $ACLOUD_CMD network elasticip delete "$eip_id" --yes --project-id "$PROJECT_ID" 2>&1 || true
    done
    
    # Delete subnets
    for subnet_id in "${CREATED_SUBNETS[@]}"; do
        echo "Deleting subnet: $subnet_id"
        echo "yes" | $ACLOUD_CMD network subnet delete "$VPC_ID" "$subnet_id" --yes --project-id "$PROJECT_ID" 2>&1 || true
    done
    
    # Delete VPCs (only if we created them)
    for vpc_id in "${CREATED_VPCS[@]}"; do
        echo "Deleting VPC: $vpc_id"
        echo "yes" | $ACLOUD_CMD network vpc delete "$vpc_id" --yes --project-id "$PROJECT_ID" 2>&1 || true
    done
}

trap cleanup EXIT

# Function to test a resource
test_resource() {
    local resource_name=$1
    local create_cmd=$2
    local list_cmd=$3
    local get_cmd=$4
    local update_cmd=$5
    local delete_cmd=$6
    
    echo -e "${YELLOW}--- Testing $resource_name ---${NC}"
    
    # CREATE
    echo -e "${GREEN}[CREATE]${NC} Creating $resource_name..."
    CREATE_OUTPUT=$(eval "$create_cmd" 2>&1) || {
        echo -e "${RED}CREATE failed:${NC}"
        echo "$CREATE_OUTPUT"
        return 1
    }
    echo "$CREATE_OUTPUT"
    
    # Extract resource ID from output
    RESOURCE_ID=$(extract_id "$CREATE_OUTPUT")
    if [ -z "$RESOURCE_ID" ]; then
        echo -e "${RED}Could not extract resource ID from create output${NC}"
        return 1
    fi
    echo -e "${GREEN}Created resource ID: $RESOURCE_ID${NC}\n"
    
    # LIST
    echo -e "${GREEN}[LIST]${NC} Listing $resource_name..."
    LIST_OUTPUT=$(eval "$list_cmd" 2>&1) || {
        echo -e "${RED}LIST failed:${NC}"
        echo "$LIST_OUTPUT"
        return 1
    }
    echo "$LIST_OUTPUT" | head -10
    echo ""
    
    # GET
    echo -e "${GREEN}[GET]${NC} Getting $resource_name details..."
    GET_OUTPUT=$(eval "$get_cmd" 2>&1) || {
        echo -e "${RED}GET failed:${NC}"
        echo "$GET_OUTPUT"
        return 1
    }
    echo "$GET_OUTPUT"
    echo ""
    
    # UPDATE
    echo -e "${GREEN}[UPDATE]${NC} Updating $resource_name..."
    UPDATE_OUTPUT=$(eval "$update_cmd" 2>&1) || {
        echo -e "${RED}UPDATE failed:${NC}"
        echo "$UPDATE_OUTPUT"
        return 1
    }
    echo "$UPDATE_OUTPUT"
    echo ""
    
    # DELETE
    echo -e "${GREEN}[DELETE]${NC} Deleting $resource_name..."
    DELETE_OUTPUT=$(eval "$delete_cmd" 2>&1) || {
        echo -e "${RED}DELETE failed:${NC}"
        echo "$DELETE_OUTPUT"
        return 1
    }
    echo "$DELETE_OUTPUT"
    echo ""
    
    echo -e "${GREEN}✓ $resource_name CRUD test completed successfully!${NC}\n"
    echo "$RESOURCE_ID"  # Return resource ID
}

# Test VPC
test_vpc() {
    echo -e "${YELLOW}=== 1. VPC CRUD Test ===${NC}\n"
    VPC_ID_OUTPUT=$(test_resource "VPC" \
        "$ACLOUD_CMD network vpc create --name ${RESOURCE_PREFIX}-vpc --region $REGION --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network vpc list --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network vpc get \$RESOURCE_ID --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network vpc update \$RESOURCE_ID --name ${RESOURCE_PREFIX}-vpc-updated --tags updated --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network vpc delete \$RESOURCE_ID --yes --project-id $PROJECT_ID")
    
    if [ -n "$VPC_ID_OUTPUT" ]; then
        CREATED_VPCS+=("$VPC_ID_OUTPUT")
        VPC_ID="$VPC_ID_OUTPUT"
    fi
}

# Test Subnet
test_subnet() {
    if [ -z "$VPC_ID" ] || [ "$VPC_ID" = "your-vpc-id" ]; then
        echo -e "${YELLOW}Skipping subnet test (no VPC available)${NC}\n"
        return 0
    fi
    
    echo -e "${YELLOW}=== 2. Subnet CRUD Test ===${NC}\n"
    SUBNET_ID=$(test_resource "Subnet" \
        "$ACLOUD_CMD network subnet create $VPC_ID --name ${RESOURCE_PREFIX}-subnet --cidr 10.150.0.0/24 --region $REGION --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network subnet list $VPC_ID --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network subnet get $VPC_ID \$RESOURCE_ID --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network subnet update $VPC_ID \$RESOURCE_ID --name ${RESOURCE_PREFIX}-subnet-updated --tags updated --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network subnet delete $VPC_ID \$RESOURCE_ID --yes --project-id $PROJECT_ID")
    
    if [ -n "$SUBNET_ID" ]; then
        CREATED_SUBNETS+=("$SUBNET_ID")
    fi
}

# Test Security Group
test_security_group() {
    if [ -z "$VPC_ID" ] || [ "$VPC_ID" = "your-vpc-id" ]; then
        echo -e "${YELLOW}Skipping security group test (no VPC available)${NC}\n"
        return 0
    fi
    
    echo -e "${YELLOW}=== 3. Security Group CRUD Test ===${NC}\n"
    SECURITY_GROUP_ID=$(test_resource "Security Group" \
        "$ACLOUD_CMD network securitygroup create $VPC_ID --name ${RESOURCE_PREFIX}-sg --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network securitygroup list $VPC_ID --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network securitygroup get $VPC_ID \$RESOURCE_ID --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network securitygroup update $VPC_ID \$RESOURCE_ID --name ${RESOURCE_PREFIX}-sg-updated --tags updated --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network securitygroup delete $VPC_ID \$RESOURCE_ID --yes --project-id $PROJECT_ID")
    
    if [ -n "$SECURITY_GROUP_ID" ]; then
        CREATED_SECURITY_GROUPS+=("$SECURITY_GROUP_ID")
    fi
}

# Test Security Rule
test_security_rule() {
    if [ -z "$VPC_ID" ] || [ "$VPC_ID" = "your-vpc-id" ] || [ -z "$SECURITY_GROUP_ID" ]; then
        echo -e "${YELLOW}Skipping security rule test (no VPC or security group available)${NC}\n"
        return 0
    fi
    
    echo -e "${YELLOW}=== 4. Security Rule CRUD Test ===${NC}\n"
    SECURITY_RULE_ID=$(test_resource "Security Rule" \
        "$ACLOUD_CMD network securityrule create $VPC_ID $SECURITY_GROUP_ID --name ${RESOURCE_PREFIX}-rule --region $REGION --direction Ingress --protocol TCP --port 80 --target-kind Ip --target-value 0.0.0.0/0 --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network securityrule list $VPC_ID $SECURITY_GROUP_ID --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network securityrule get $VPC_ID $SECURITY_GROUP_ID \$RESOURCE_ID --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network securityrule update $VPC_ID $SECURITY_GROUP_ID \$RESOURCE_ID --name ${RESOURCE_PREFIX}-rule-updated --tags updated --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network securityrule delete $VPC_ID $SECURITY_GROUP_ID \$RESOURCE_ID --yes --project-id $PROJECT_ID")
    
    if [ -n "$SECURITY_RULE_ID" ]; then
        CREATED_SECURITY_RULES+=("$SECURITY_RULE_ID")
    fi
}

# Test Elastic IP
test_elastic_ip() {
    echo -e "${YELLOW}=== 5. Elastic IP CRUD Test ===${NC}\n"
    EIP_ID=$(test_resource "Elastic IP" \
        "$ACLOUD_CMD network elasticip create --name ${RESOURCE_PREFIX}-eip --region $REGION --billing-period Hour --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network elasticip list --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network elasticip get \$RESOURCE_ID --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network elasticip update \$RESOURCE_ID --name ${RESOURCE_PREFIX}-eip-updated --tags updated --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network elasticip delete \$RESOURCE_ID --yes --project-id $PROJECT_ID")
    
    if [ -n "$EIP_ID" ]; then
        CREATED_ELASTIC_IPS+=("$EIP_ID")
    fi
}

# Test VPC Peering
test_vpc_peering() {
    if [ -z "$VPC_ID" ] || [ "$VPC_ID" = "your-vpc-id" ]; then
        echo -e "${YELLOW}Skipping VPC peering test (no VPC available)${NC}\n"
        return 0
    fi
    
    echo -e "${YELLOW}=== 6. VPC Peering CRUD Test ===${NC}\n"
    PEERING_ID=$(test_resource "VPC Peering" \
        "$ACLOUD_CMD network vpcpeering create $VPC_ID --name ${RESOURCE_PREFIX}-peering --peer-vpc-id $PEER_VPC_URI --region $REGION --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network vpcpeering list $VPC_ID --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network vpcpeering get $VPC_ID \$RESOURCE_ID --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network vpcpeering update $VPC_ID \$RESOURCE_ID --name ${RESOURCE_PREFIX}-peering-updated --tags updated --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network vpcpeering delete $VPC_ID \$RESOURCE_ID --yes --project-id $PROJECT_ID")
    
    if [ -n "$PEERING_ID" ]; then
        CREATED_PEERINGS+=("$PEERING_ID")
    fi
}

# Test VPC Peering Route
test_vpc_peering_route() {
    if [ -z "$VPC_ID" ] || [ "$VPC_ID" = "your-vpc-id" ] || [ -z "$PEERING_ID" ]; then
        echo -e "${YELLOW}Skipping VPC peering route test (no VPC or peering available)${NC}\n"
        return 0
    fi
    
    echo -e "${YELLOW}=== 7. VPC Peering Route CRUD Test ===${NC}\n"
    ROUTE_ID=$(test_resource "VPC Peering Route" \
        "$ACLOUD_CMD network vpcpeeringroute create $VPC_ID $PEERING_ID --name ${RESOURCE_PREFIX}-route --local-network 10.0.1.0/24 --remote-network 10.0.2.0/24 --billing-period Hour --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network vpcpeeringroute list $VPC_ID $PEERING_ID --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network vpcpeeringroute get $VPC_ID $PEERING_ID \$RESOURCE_ID --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network vpcpeeringroute update $VPC_ID $PEERING_ID \$RESOURCE_ID --name ${RESOURCE_PREFIX}-route-updated --tags updated --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network vpcpeeringroute delete $VPC_ID $PEERING_ID \$RESOURCE_ID --yes --project-id $PROJECT_ID")
    
    if [ -n "$ROUTE_ID" ]; then
        CREATED_PEERING_ROUTES+=("$ROUTE_ID")
    fi
}

# Test VPN Tunnel
test_vpn_tunnel() {
    if [ -z "$VPC_ID" ] || [ "$VPC_ID" = "your-vpc-id" ] || [ -z "$ELASTIC_IP_URI" ] || [ "$ELASTIC_IP_URI" = "your-elastic-ip-uri" ]; then
        echo -e "${YELLOW}Skipping VPN tunnel test (missing VPC or Elastic IP)${NC}\n"
        return 0
    fi
    
    echo -e "${YELLOW}=== 8. VPN Tunnel CRUD Test ===${NC}\n"
    VPN_TUNNEL_ID=$(test_resource "VPN Tunnel" \
        "$ACLOUD_CMD network vpntunnel create --name ${RESOURCE_PREFIX}-vpn --region $REGION --peer-ip 203.0.113.1 --vpc-uri /projects/$PROJECT_ID/providers/Aruba.Network/vpcs/$VPC_ID --subnet-cidr 10.0.1.0/24 --elastic-ip-uri $ELASTIC_IP_URI --billing-period Hour --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network vpntunnel list --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network vpntunnel get \$RESOURCE_ID --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network vpntunnel update \$RESOURCE_ID --name ${RESOURCE_PREFIX}-vpn-updated --tags updated --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network vpntunnel delete \$RESOURCE_ID --yes --project-id $PROJECT_ID")
    
    if [ -n "$VPN_TUNNEL_ID" ]; then
        CREATED_VPN_TUNNELS+=("$VPN_TUNNEL_ID")
    fi
}

# Test VPN Route
test_vpn_route() {
    if [ -z "$VPN_TUNNEL_ID" ]; then
        echo -e "${YELLOW}Skipping VPN route test (no VPN tunnel available)${NC}\n"
        return 0
    fi
    
    echo -e "${YELLOW}=== 9. VPN Route CRUD Test ===${NC}\n"
    ROUTE_ID=$(test_resource "VPN Route" \
        "$ACLOUD_CMD network vpnroute create $VPN_TUNNEL_ID --name ${RESOURCE_PREFIX}-vpn-route --region $REGION --cloud-subnet 10.0.1.0/24 --onprem-subnet 192.168.1.0/24 --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network vpnroute list $VPN_TUNNEL_ID --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network vpnroute get $VPN_TUNNEL_ID \$RESOURCE_ID --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network vpnroute update $VPN_TUNNEL_ID \$RESOURCE_ID --name ${RESOURCE_PREFIX}-vpn-route-updated --tags updated --project-id $PROJECT_ID" \
        "$ACLOUD_CMD network vpnroute delete $VPN_TUNNEL_ID \$RESOURCE_ID --yes --project-id $PROJECT_ID")
    
    if [ -n "$ROUTE_ID" ]; then
        CREATED_VPN_ROUTES+=("$ROUTE_ID")
    fi
}

# Run tests
echo -e "${BLUE}Starting Network Resources E2E Tests...${NC}\n"

# Only test VPC if not provided
if [ "$VPC_ID" = "your-vpc-id" ]; then
    test_vpc
fi

test_subnet
test_security_group
test_security_rule
test_elastic_ip

# VPC Peering tests (require peer VPC)
if [ "$PEER_VPC_ID" != "your-peer-vpc-id" ]; then
    test_vpc_peering
    test_vpc_peering_route
fi

# VPN tests (require Elastic IP)
if [ "$ELASTIC_IP_URI" != "your-elastic-ip-uri" ]; then
    test_vpn_tunnel
    test_vpn_route
fi

echo -e "${GREEN}=== All Network Tests Completed! ===${NC}"

