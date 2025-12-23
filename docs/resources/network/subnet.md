# Subnet

Subnets allow you to segment your VPC network into smaller, isolated sections. Each subnet can have its own CIDR block and can be used to organize resources, control routing, and apply security policies within a VPC.

## Commands

### List Subnets
List all subnets in a VPC.

```bash
acloud network subnet list <vpc-id>
```

**Arguments:**
- `vpc-id` - The ID of the VPC

**Example:**
```bash
acloud network subnet list 689307f4745108d3c6343b5a
```

**Output:**
```
NAME         ID                        CIDR           STATUS
subnet-1     1234567890abcdef          10.0.1.0/24    Active
subnet-2     0987654321fedcba          10.0.2.0/24    Active
```

### Get Subnet Details
Get details about a specific subnet.

```bash
acloud network subnet get <vpc-id> <subnet-id>
```

**Arguments:**
- `vpc-id` - The ID of the VPC
- `subnet-id` - The ID of the subnet

**Example:**
```bash
acloud network subnet get 689307f4745108d3c6343b5a 1234567890abcdef
```

**Output:**
```
Subnet Details:
===============
ID:            1234567890abcdef
Name:          subnet-1
CIDR:          10.0.1.0/24
Status:        Active
```

### Create Subnet
Create a new subnet in a VPC.

```bash
acloud network subnet create <vpc-id> --cidr <cidr> --name <name>
```

**Required Flags:**
- `--cidr string` - CIDR block for the subnet
- `--name string` - Name for the subnet

**Example:**
```bash
acloud network subnet create 689307f4745108d3c6343b5a --cidr 10.0.1.0/24 --name subnet-1
```

**Output:**
```
Subnet created successfully!
ID:      1234567890abcdef
Name:    subnet-1
CIDR:    10.0.1.0/24
```

### Update Subnet
Update an existing subnet's name or CIDR.

```bash
acloud network subnet update <vpc-id> <subnet-id> [flags]
```

**Arguments:**
- `vpc-id` - The ID of the VPC
- `subnet-id` - The ID of the subnet

**Flags:**
- `--name string` - New name for the subnet
- `--cidr string` - New CIDR block

**Example:**
```bash
acloud network subnet update 689307f4745108d3c6343b5a 1234567890abcdef --name new-subnet-name
```

**Output:**
```
Subnet updated successfully!
ID:      1234567890abcdef
Name:    new-subnet-name
```

### Delete Subnet
Delete a subnet from a VPC.

```bash
acloud network subnet delete <vpc-id> <subnet-id>
```

**Arguments:**
- `vpc-id` - The ID of the VPC
- `subnet-id` - The ID of the subnet

**Example:**
```bash
acloud network subnet delete 689307f4745108d3c6343b5a 1234567890abcdef
```

**Output:**
```
Subnet 1234567890abcdef deleted successfully!
```

## Shell Auto-completion

The subnet commands support auto-completion for VPC IDs and subnet IDs.

## Best Practices
- Use descriptive names for subnets based on their purpose.
- Avoid overlapping CIDR blocks between subnets.

## Troubleshooting
- Ensure the VPC is **Active** before creating subnets.
- Check for CIDR conflicts when adding or updating subnets.

## Related Commands
- [VPC](vpc.md)
- [Security Group](securitygroup.md)
