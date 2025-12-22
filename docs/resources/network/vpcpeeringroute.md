# VPC Peering Route

VPC Peering Routes allow you to define custom routing rules for traffic between peered VPCs in Aruba Cloud. These routes control how network traffic is directed between VPCs connected via a VPC Peering connection.

## Commands

### List VPC Peering Routes
List all routes for a specific VPC peering connection.

```bash
acloud network vpcpeeringroute list <vpc-id> <peering-id>
```

**Arguments:**
- `vpc-id` - The ID of the VPC
- `peering-id` - The ID of the VPC peering connection

**Example:**
```bash
acloud network vpcpeeringroute list 689307f4745108d3c6343b5a 6949666e4d0cdc87949b7204
```

**Output:**
```
DESTINATION      NEXT HOP      STATUS
10.0.2.0/24      10.0.1.1      Active
10.0.3.0/24      10.0.1.2      Active
```

### Get VPC Peering Route Details
Get details about a specific VPC peering route.

```bash
acloud network vpcpeeringroute get <vpc-id> <peering-id> <route-id>
```

**Arguments:**
- `vpc-id` - The ID of the VPC
- `peering-id` - The ID of the VPC peering connection
- `route-id` - The ID of the route

**Example:**
```bash
acloud network vpcpeeringroute get 689307f4745108d3c6343b5a 6949666e4d0cdc87949b7204 123456
```

**Output:**
```
Route Details:
==============
ID:            123456
Destination:   10.0.2.0/24
Next Hop:      10.0.1.1
Status:        Active
```

### Create VPC Peering Route
Create a new route for a VPC peering connection.

```bash
acloud network vpcpeeringroute create <vpc-id> <peering-id> --destination-cidr <cidr> --next-hop <ip>
```

**Required Flags:**
- `--destination-cidr string` - Destination CIDR for the route
- `--next-hop string` - Next hop IP address

**Example:**
```bash
acloud network vpcpeeringroute create 689307f4745108d3c6343b5a 6949666e4d0cdc87949b7204 --destination-cidr 10.0.2.0/24 --next-hop 10.0.1.1
```

**Output:**
```
VPC Peering Route created successfully!
ID:          123456
Destination: 10.0.2.0/24
Next Hop:    10.0.1.1
```

### Update VPC Peering Route
Update an existing VPC peering route.

```bash
acloud network vpcpeeringroute update <vpc-id> <peering-id> <route-id> [flags]
```

**Arguments:**
- `vpc-id` - The ID of the VPC
- `peering-id` - The ID of the VPC peering connection
- `route-id` - The ID of the route

**Flags:**
- `--destination-cidr string` - New destination CIDR
- `--next-hop string` - New next hop IP address

**Example:**
```bash
acloud network vpcpeeringroute update 689307f4745108d3c6343b5a 6949666e4d0cdc87949b7204 123456 --next-hop 10.0.1.2
```

**Output:**
```
VPC Peering Route updated successfully!
ID:          123456
Next Hop:    10.0.1.2
```

### Delete VPC Peering Route
Delete a VPC peering route.

```bash
acloud network vpcpeeringroute delete <vpc-id> <peering-id> <route-id>
```

**Arguments:**
- `vpc-id` - The ID of the VPC
- `peering-id` - The ID of the VPC peering connection
- `route-id` - The ID of the route

**Example:**
```bash
acloud network vpcpeeringroute delete 689307f4745108d3c6343b5a 6949666e4d0cdc87949b7204 123456
```

**Output:**
```
VPC Peering Route 123456 deleted successfully!
```

## Shell Auto-completion

The VPC Peering Route commands support auto-completion for VPC IDs, peering IDs, and route IDs.

## Best Practices
- Use descriptive destination CIDRs and next hop IPs.
- Regularly review and clean up unused routes.

## Troubleshooting
- Ensure the VPC peering connection is **Active** before adding routes.
- Check for overlapping CIDRs that may cause routing conflicts.

## Related Commands
- [VPC Peering](vpcpeering.md)
- [VPC](vpc.md)
