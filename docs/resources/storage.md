# Storage Resources

This guide covers managing storage resources in Aruba Cloud, including block storage volumes and snapshots.

## Overview

Storage resources in Aruba Cloud include:
- **Block Storage**: Persistent volumes that can be attached to cloud servers
- **Snapshots**: Point-in-time copies of block storage volumes

## Prerequisites

- Configured CLI credentials (see [Getting Started](../getting-started.md))
- A project ID (or set up a [context](../getting-started.md#context-management))

## Block Storage

Block storage volumes provide persistent storage for your cloud infrastructure.

### List Block Storage Volumes

```bash
# Using context
acloud storage blockstorage list

# With explicit project ID
acloud storage blockstorage list --project-id "66a10244f62b99c686572a9f"

# With verbose debug output
acloud storage blockstorage list --verbose
```

**Output:**
```
NAME                           ID                         SIZE(GB)     REGION          ZONE            TYPE            STATUS          
test-volume                    69442fe38f4a09c12b5ded74   10           IT BG           ITBG-3          Standard        NotUsed         
my-volume                      69087e12cc4c7793b9e4d2eb   100          IT BG           ITBG-3          Standard        Available 
data-volume                    6901cafe42360f3845c1324c   500          IT BG           ITBG-3          Standard        In-use    
backup-volume                  6901c6b98b5ab53f0516c397   250          IT BG           ITBG-3          Standard        Available 
```

**Flags:**
- `--project-id` - Project ID (optional if context is set)
- `-v, --verbose` - Show detailed debug information including raw API response

### Get Block Storage Details

```bash
acloud storage blockstorage get <volume-id>
```

**Example:**
```bash
acloud storage blockstorage get vol-123456789
```

**Output:**
```
Block Storage Details:
ID:              vol-123456789
Name:            my-volume
Size (GB):       100
Type:            Standard
Zone:            it-mil1-1
Region:          it-mil1
Status:          Available
Billing Period:  Hour
Created:         2025-12-18T10:30:00Z
Tags:            environment=production, owner=devops
```

### Create Block Storage

```bash
acloud storage blockstorage create \
  --name "my-volume" \
  --zone "ITBG-3" \
  --size 100
```

**With all options:**
```bash
acloud storage blockstorage create \
  --name "my-volume" \
  --region "ITBG-Bergamo" \
  --zone "ITBG-3" \
  --size 100 \
  --type "Standard" \
  --billing-period "Hour" \
  --tags "environment=production,owner=devops"
```

**Required Flags:**
- `--name`: Name for the block storage volume
- `--zone`: Zone/datacenter (e.g., `ITBG-3`)
- `--size`: Size in GB

**Optional Flags:**
- `--region`: Region code (default: `ITBG-Bergamo`)
- `--type`: Volume type (`Standard` or `Performance`, default: `Standard`)
- `--billing-period`: Billing period (`Hour`, `Month`, `Year`, default: `Hour`)
- `--tags`: Comma-separated tags
- `--project-id`: Project ID (uses context if not specified)

**Example Output:**
```
Creating block storage with:
  Name: my-volume
  Region: ITBG-Bergamo
  Zone: ITBG-3
  Size: 100 GB
  Type: Standard
  Billing Period: Hour
  Project ID: 68398923fb2cb026400d4d31

Block storage created successfully!
ID:              vol-987654321
Name:            my-volume
Size (GB):       100
Type:            Standard
Zone:            ITBG-3
Region:          IT BG
Status:          InCreation
Creation Date:   18-12-2025 16:46:27
```

**Common Region Codes:**
- `ITBG-Bergamo` - Italy, Bergamo (default)
- `ITMIL-Milano` - Italy, Milan
- `CZPRG-Prague` - Czech Republic, Prague

### Update Block Storage

**Note:** The update functionality for block storage is currently limited by the SDK. Size updates are not supported at this time.

You can update the name and/or tags of a block storage volume:

```bash
# Update name
acloud storage blockstorage update vol-123456789 --name "new-volume-name"

# Update tags
acloud storage blockstorage update vol-123456789 --tags "env=prod,team=backend"

# Update both
acloud storage blockstorage update vol-123456789 \
  --name "new-name" \
  --tags "env=prod,team=backend"
```

**Note:** Currently, only name and tags can be updated. Size and type updates are not supported by the API at this time.

### Delete Block Storage

```bash
# With confirmation prompt
acloud storage blockstorage delete vol-123456789

# Skip confirmation
acloud storage blockstorage delete vol-123456789 --yes
```

**Example with confirmation:**
```bash
$ acloud storage blockstorage delete 69442fe38f4a09c12b5ded74
Are you sure you want to delete block storage 69442fe38f4a09c12b5ded74? (yes/no): yes

Block storage 69442fe38f4a09c12b5ded74 deleted successfully!
```

**Flags:**
- `-y, --yes` - Skip confirmation prompt (useful for scripts)

**Warning:** Deleting a block storage volume is permanent and cannot be undone.

## Snapshots

Snapshots provide point-in-time copies of block storage volumes for backup and recovery.

### List Snapshots

```bash
# Using context
acloud storage snapshot list

# With explicit project ID
acloud storage snapshot list --project-id "66a10244f62b99c686572a9f"
```

**Output:**
```
NAME                SIZE(GB)  SOURCE                  STATUS    
daily-backup        100       my-volume               Available 
pre-upgrade         500       data-volume             Available 
weekly-snap         250       backup-volume           Creating  
```

### Get Snapshot Details

```bash
acloud storage snapshot get <snapshot-id>
```

**Example:**
```bash
acloud storage snapshot get snap-123456789
```

**Output:**
```
Snapshot Details:
ID:              snap-123456789
Name:            daily-backup
Size (GB):       100
Region:          it-mil1
Status:          Available
Source Volume:   /storage/volumes/vol-123456789
Created:         2025-12-18T02:00:00Z
Tags:            type=backup, schedule=daily
```

### Create Snapshot

```bash
acloud storage snapshot create \
  --name "my-snapshot" \
  --region "it-mil1" \
  --volume-uri "/storage/volumes/vol-123456789" \
  --tags "type=backup,schedule=daily"
```

**Required Flags:**
- `--name`: Name for the snapshot
- `--region`: Region code
- `--volume-uri`: URI of the source volume (e.g., `/storage/volumes/vol-123456789`)

**Optional Flags:**
- `--tags`: Comma-separated tags
- `--project-id`: Project ID (uses context if not specified)

**Example Output:**
```
Snapshot created successfully!
ID:              snap-987654321
Name:            my-snapshot
Size (GB):       100
Region:          it-mil1
Status:          Creating
Source Volume:   /storage/volumes/vol-123456789
```

### Update Snapshot

You can update the name and/or tags of a snapshot:

```bash
# Update name
acloud storage snapshot update snap-123456789 --name "new-snapshot-name"

# Update tags
acloud storage snapshot update snap-123456789 --tags "type=backup,retention=30days"

# Update both
acloud storage snapshot update snap-123456789 \
  --name "new-name" \
  --tags "type=backup,retention=30days"
```

**Note:** Only name and tags can be updated.

### Delete Snapshot

```bash
# With confirmation prompt
acloud storage snapshot delete snap-123456789

# Skip confirmation
acloud storage snapshot delete snap-123456789 --yes
```

**Warning:** Deleting a snapshot is permanent and cannot be undone.

## Best Practices

### Block Storage

1. **Naming Convention**: Use descriptive names that indicate purpose and environment
   ```bash
   --name "prod-db-data-volume"
   --name "dev-app-logs"
   ```

2. **Tagging Strategy**: Use consistent tags for organization and cost tracking
   ```bash
   --tags "environment=production,application=database,cost-center=engineering"
   ```

3. **Size Planning**: Choose appropriate size based on growth projections
   - Start with minimum required size
   - Monitor usage regularly
   - Note: Size cannot be increased after creation (API limitation)
   - Consider costs vs. performance

4. **Type Selection**:
   - **Standard**: General-purpose workloads, cost-effective
   - **Performance**: High I/O workloads, databases, critical applications

5. **Regional Placement**: Create volumes in the same region as your compute resources
   - Default region: `ITBG-Bergamo`
   - Common zones: `ITBG-3` (Bergamo)

6. **Billing Period**: Choose based on usage duration
   - **Hour**: Short-term testing, development
   - **Month**: Production workloads
   - **Year**: Long-term stable infrastructure (best cost savings)

### Snapshots

1. **Regular Backups**: Create snapshots on a regular schedule
   ```bash
   --tags "schedule=daily,retention=7days"
   --tags "schedule=weekly,retention=4weeks"
   ```

2. **Pre-Change Snapshots**: Always create a snapshot before major changes
   ```bash
   acloud storage snapshot create \
     --name "pre-upgrade-$(date +%Y%m%d)" \
     --region "it-mil1" \
     --volume-uri "/storage/volumes/vol-123"
   ```

3. **Naming Convention**: Include date and purpose in snapshot names
   ```bash
   --name "prod-db-daily-20251218"
   --name "pre-migration-backup-20251218"
   ```

4. **Retention Policy**: Implement a clear retention policy
   - Daily snapshots: 7 days
   - Weekly snapshots: 4 weeks
   - Monthly snapshots: 12 months

5. **Cleanup**: Delete old snapshots to manage costs
   ```bash
   # List and identify old snapshots
   acloud storage snapshot list
   
   # Delete outdated snapshots
   acloud storage snapshot delete snap-old-123 --yes
   ```

## Common Workflows

### Creating a Volume with Backup

```bash
# 1. Create the volume
acloud storage blockstorage create \
  --name "my-app-data" \
  --zone "ITBG-3" \
  --size 100 \
  --type "Standard" \
  --billing-period "Month" \
  --tags "app=myapp,environment=production"

# 2. Create initial snapshot
acloud storage snapshot create \
  --name "my-app-data-initial" \
  --region "ITBG-Bergamo" \
  --volume-uri "/storage/volumes/vol-123456789" \
  --tags "type=initial,app=myapp"
```

### Migrating Data Between Regions

```bash
# 1. Create snapshot of source volume
acloud storage snapshot create \
  --name "migration-source-$(date +%Y%m%d)" \
  --region "it-mil1" \
  --volume-uri "/storage/volumes/vol-source-123"

# 2. Create new volume in target region from snapshot
# (Note: This would typically be done through the API or console)

# 3. Verify and cleanup old resources after migration
```

### Disaster Recovery Setup

```bash
# Set up context for production project
acloud context set prod --project-id "prod-project-id"
acloud context use prod

# Create daily snapshot (run via cron/scheduled task)
acloud storage snapshot create \
  --name "dr-backup-$(date +%Y%m%d-%H%M)" \
  --region "it-mil1" \
  --volume-uri "/storage/volumes/vol-critical-123" \
  --tags "type=dr,schedule=daily,retention=7days"
```

## Using Context for Efficiency

Set up contexts to avoid specifying `--project-id` repeatedly:

```bash
# Configure contexts for different environments
acloud context set prod --project-id "prod-project-id"
acloud context set dev --project-id "dev-project-id"

# Switch to production
acloud context use prod

# Now all commands use the prod project ID
acloud storage blockstorage list
acloud storage snapshot list

# Switch to development
acloud context use dev

# Commands now use the dev project ID
acloud storage blockstorage list
```

## Troubleshooting

### "Error: project ID not specified"

**Solution:** Set a context or provide `--project-id`:
```bash
# Option 1: Set context
acloud context use my-prod

# Option 2: Use explicit project ID
acloud storage blockstorage list --project-id "your-project-id"
```

### "Error: Location Value: IT-BG not found"

**Solution:** Use the correct region code format:
```bash
# Correct format
--region "ITBG-Bergamo"

# Not: IT-BG, it-bg1, IT BG
```

Common region codes:
- `ITBG-Bergamo` (Italy - Bergamo)
- `ITMIL-Milano` (Italy - Milan)
- `CZPRG-Prague` (Czech Republic - Prague)

### "Error: validation errors occurred"

**Solution:** Check all required parameters:
```bash
acloud storage blockstorage create \
  --name "test-volume" \
  --zone "ITBG-3" \
  --size 10
```

Required: `--name`, `--zone`, `--size`

### Volume shows empty zone

**Issue:** Some older volumes may not have a zone assigned.

**Solution:** This is expected for volumes created before zone assignment was enforced. New volumes will always have a zone.

### "panic: unimplemented" when updating

**Issue:** The SDK's update method is not yet implemented for block storage.

**Solution:** Size and type updates are not currently supported. Only name and tags can be updated through the console or API directly.

### Using verbose mode for debugging

If you encounter issues, use the `--verbose` flag to see the full API response:

```bash
acloud storage blockstorage list --verbose
```

This will show:
- HTTP status codes
- Full volume details
- All metadata and properties
- Exact field values from the API

## Next Steps

- Explore [Compute Resources](../compute.md)
- Learn about [Network Resources](../network.md)
- Review [Project Management](management/projects.md)
