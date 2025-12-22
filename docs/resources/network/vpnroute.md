# VPN Tunnel Route

VPN Tunnel Routes allow you to define custom routing rules for traffic flowing through VPN tunnels in Aruba Cloud. These routes control how network traffic is directed between your on-premises network and your VPC via a VPN tunnel.

## Commands

### List VPN Tunnel Routes
List all routes for a specific VPN tunnel.

```bash
acloud network vpnroute list <vpn-tunnel-id>
```

**Arguments:**
- `vpn-tunnel-id` - The ID of the VPN tunnel

**Example:**
```bash
acloud network vpnroute list 1234567890abcdef
```

**Output:**
```
DESTINATION      NEXT HOP      STATUS
10.0.3.0/24      10.0.2.1      Active
10.0.4.0/24      10.0.2.2      Active
```

### Get VPN Tunnel Route Details
Get details about a specific VPN tunnel route.

```bash
acloud network vpnroute get <vpn-tunnel-id> <route-id>
```

**Arguments:**
- `vpn-tunnel-id` - The ID of the VPN tunnel
- `route-id` - The ID of the route

**Example:**
```bash
acloud network vpnroute get 1234567890abcdef 654321
```

**Output:**
```
Route Details:
==============
ID:            654321
Destination:   10.0.3.0/24
Next Hop:      10.0.2.1
Status:        Active
```

### Create VPN Tunnel Route
Create a new route for a VPN tunnel.

```bash
acloud network vpnroute create <vpn-tunnel-id> --destination-cidr <cidr> --next-hop <ip>
```

**Required Flags:**
- `--destination-cidr string` - Destination CIDR for the route
- `--next-hop string` - Next hop IP address

**Example:**
```bash
acloud network vpnroute create 1234567890abcdef --destination-cidr 10.0.3.0/24 --next-hop 10.0.2.1
```

**Output:**
```
VPN Tunnel Route created successfully!
ID:          654321
Destination: 10.0.3.0/24
Next Hop:    10.0.2.1
```

### Update VPN Tunnel Route
Update an existing VPN tunnel route.

```bash
acloud network vpnroute update <vpn-tunnel-id> <route-id> [flags]
```

**Arguments:**
- `vpn-tunnel-id` - The ID of the VPN tunnel
- `route-id` - The ID of the route

**Flags:**
- `--destination-cidr string` - New destination CIDR
- `--next-hop string` - New next hop IP address

**Example:**
```bash
acloud network vpnroute update 1234567890abcdef 654321 --next-hop 10.0.2.2
```

**Output:**
```
VPN Tunnel Route updated successfully!
ID:          654321
Next Hop:    10.0.2.2
```

### Delete VPN Tunnel Route
Delete a VPN tunnel route.

```bash
acloud network vpnroute delete <vpn-tunnel-id> <route-id>
```

**Arguments:**
- `vpn-tunnel-id` - The ID of the VPN tunnel
- `route-id` - The ID of the route

**Example:**
```bash
acloud network vpnroute delete 1234567890abcdef 654321
```

**Output:**
```
VPN Tunnel Route 654321 deleted successfully!
```

## Shell Auto-completion

The VPN Tunnel Route commands support auto-completion for VPN tunnel IDs and route IDs.

## Best Practices
- Use descriptive destination CIDRs and next hop IPs.
- Regularly review and clean up unused routes.

## Troubleshooting
- Ensure the VPN tunnel is **Active** before adding routes.
- Check for overlapping CIDRs that may cause routing conflicts.

## Related Commands
- [VPN Tunnel](vpntunnel.md)
- [VPC](vpc.md)
