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
acloud storage blockstorage get 69442fe38f4a09c12b5ded74
```

**Output:**
```
Block Storage Details:
======================
ID:              69442fe38f4a09c12b5ded74
URI:             /projects/68398923fb2cb026400d4d31/providers/Aruba.Storage/blockStorages/69442fe38f4a09c12b5ded74
Name:            my-volume
Size (GB):       100
Type:            Standard
Zone:            ITBG-3
Region:          IT BG
Status:          NotUsed
Creation Date:   18-12-2025 10:30:00
Created By:      aru-123456
Tags:            [environment=production owner=devops]
```

**Note:** The URI field shows the full resource path needed for creating snapshots.

### Create Block Storage

```bash
acloud storage blockstorage create \
  --name "my-volume" \
  --size 100
```

**With zone:**
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
- `--size`: Size in GB

**Optional Flags:**
- `--zone`: Zone/datacenter (e.g., `ITBG-3`) - Optional, can be omitted
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

**Important:** Block storage can only be updated when its status is `Used` or `NotUsed`. Updates are not allowed during provisioning states like `InCreation`.

You can update the name and/or tags of a block storage volume:

```bash
# Update name
acloud storage blockstorage update 69442fe38f4a09c12b5ded74 --name "new-volume-name"

# Update tags
acloud storage blockstorage update 69442fe38f4a09c12b5ded74 --tags "env=prod,team=backend"

# Update both
acloud storage blockstorage update 69442fe38f4a09c12b5ded74 \
  --name "new-name" \
  --tags "env=prod,team=backend"
```

**Example Output:**
```bash
$ acloud storage blockstorage update 694557430d0972656501d43c --name "my-updated-volume"

Block storage updated successfully!
ID:              694557430d0972656501d43c
Name:            my-updated-volume
Size (GB):       5
Type:            Standard
```

**Status Validation:**
```bash
$ acloud storage blockstorage update 69455aa70d0972656501d45d --name "should-fail"
Error: Cannot update block storage with status 'InCreation'
Block storage can only be updated when status is 'Used' or 'NotUsed'
```

**Note:** Size and type updates are not supported by the API at this time. Only name and tags can be modified.

### Delete Block Storage

```bash
# With confirmation prompt
acloud storage blockstorage delete 69442fe38f4a09c12b5ded74

# Skip confirmation
acloud storage blockstorage delete 69442fe38f4a09c12b5ded74 --yes
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

Snapshots provide point-in-time copies of block storage volumes for backup and recovery purposes.

### List Snapshots

List snapshots for a specific block storage volume:

```bash
# List snapshots for a volume (requires volume URI)
acloud storage snapshot list --volume-uri "/projects/68398923fb2cb026400d4d31/providers/Aruba.Storage/blockStorages/69442fe38f4a09c12b5ded74"

# With explicit project ID
acloud storage snapshot list \
  --volume-uri "/projects/68398923fb2cb026400d4d31/providers/Aruba.Storage/blockStorages/69442fe38f4a09c12b5ded74" \
  --project-id "66a10244f62b99c686572a9f"
```

**Output:**
```
NAME                           ID                         SIZE(GB)     STATUS          
backup-snapshot-01             69442abc8f4a09c12b5def12   100          Active       
daily-backup                   69442def8f4a09c12b5def34   250          Active       
updated-test-snapshot          69455d620d0972656501d477   824634892352 Active
```

**Required Flags:**
- `--volume-uri` - Full URI of the source block storage volume (get this from `blockstorage get` command)

**Optional Flags:**
- `--project-id` - Project ID (uses context if not specified)

**Note:** The `--volume-uri` flag is required to filter snapshots for a specific volume. Use `blockstorage get <volume-id>` to get the volume URI.

### Get Snapshot Details

```bash
acloud storage snapshot get <snapshot-id>
```

**Example:**
```bash
acloud storage snapshot get 69442abc8f4a09c12b5def12
```

**Output:**
```
Snapshot Details:
=================
ID:              69442abc8f4a09c12b5def12
URI:             /projects/68398923fb2cb026400d4d31/providers/Aruba.Storage/snapshots/69442abc8f4a09c12b5def12
Name:            daily-backup
Size (GB):       100
Source Volume:   /projects/68398923fb2cb026400d4d31/providers/Aruba.Storage/blockStorages/69442fe38f4a09c12b5ded74
Region:          IT BG
Status:          Active
Creation Date:   19-12-2025 14:12:50
Created By:      aru-297647
```

**Note:** The Region field shows the format returned by the API (`IT BG`). When updating snapshots, the CLI automatically converts this to the required format (`ITBG-Bergamo`).
### Create Snapshot

To create a snapshot, you need the full URI of the source block storage volume. You can get this from the `blockstorage get` command.

```bash
# First, get the volume URI
acloud storage blockstorage get 69442fe38f4a09c12b5ded74

# Then create the snapshot using the URI
acloud storage snapshot create \
  --name "my-snapshot" \
  --region "ITBG-Bergamo" \
  --volume-uri "/projects/68398923fb2cb026400d4d31/providers/Aruba.Storage/blockStorages/69442fe38f4a09c12b5ded74"
```

**With tags:**
```bash
acloud storage snapshot create \
  --name "daily-backup" \
  --region "ITBG-Bergamo" \
  --volume-uri "/projects/68398923fb2cb026400d4d31/providers/Aruba.Storage/blockStorages/69442fe38f4a09c12b5ded74" \
  --tags "type=backup,schedule=daily"
```

**Required Flags:**
- `--name`: Name for the snapshot
- `--region`: Region code (e.g., `ITBG-Bergamo`)
- `--volume-uri`: Full URI of the source volume (get this from `blockstorage get` command)

**Optional Flags:**
- `--tags`: Comma-separated tags
- `--project-id`: Project ID (uses context if not specified)

**Example Output:**
```
Creating snapshot with the following parameters:
  Name:       daily-backup
  Region:     ITBG-Bergamo
  Volume URI: /projects/68398923fb2cb026400d4d31/providers/Aruba.Storage/blockStorages/69442fe38f4a09c12b5ded74
  Tags:       [type=backup schedule=daily]

Snapshot created successfully!
ID:              69442abc8f4a09c12b5def12
Name:            daily-backup
Creation Date:   19-12-2025 14:12:50
```

**Important Notes:**
- The snapshot will initially be in `InCreation` status and will transition to `Active` when ready
- SDK v0.1.2+ properly handles block storage state validation during snapshot creation
- Creation may take several minutes depending on volume size
Name:            daily-backup
Size (GB):       100
Creation Date:   19-12-2025 14:30:00
```

### Update Snapshot

You can update the name and/or tags of a snapshot:

```bash
# Update name
acloud storage snapshot update 69442abc8f4a09c12b5def12 --name "new-snapshot-name"

# Update tags
acloud storage snapshot update 69442abc8f4a09c12b5def12 --tags "type=backup,retention=30days"

# Update both
acloud storage snapshot update 69442abc8f4a09c12b5def12 \
  --name "new-name" \
  --tags "type=backup,retention=30days"
```

**Example:**
```bash
acloud storage snapshot update 69455d620d0972656501d477 --name "updated-test-snapshot"
```

**Output:**
```
Snapshot updated successfully!
```

**Important Notes:**
- Only name and tags can be updated
- The CLI automatically converts region codes from API format (`IT BG`) to the required update format (`ITBG-Bergamo`)
- This region conversion happens automatically - you don't need to specify the region during updates

### Delete Snapshot

```bash
# With confirmation prompt
acloud storage snapshot delete 69442abc8f4a09c12b5def12

# Skip confirmation (useful for automation)
acloud storage snapshot delete 69442abc8f4a09c12b5def12 --yes
```

**Example with confirmation:**
```bash
$ acloud storage snapshot delete 69455d620d0972656501d477
Are you sure you want to delete snapshot 69455d620d0972656501d477? (yes/no): yes

Snapshot 69455d620d0972656501d477 deleted successfully!
```

**Flags:**
- `-y, --yes` - Skip confirmation prompt (useful for scripts and automation)

**Warning:** Deleting a snapshot is permanent and cannot be undone. Always verify the snapshot ID before deletion.

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
   - Note: Size cannot be increased after creation via update command
   - Consider costs vs. performance

4. **Type Selection**:
   - **Standard**: General-purpose workloads, cost-effective
   - **Performance**: High I/O workloads, databases, critical applications

5. **Zone Selection**:
   - Zone parameter is **optional** when creating block storage
   - Can be omitted for region-wide deployment
   - Specify zone (e.g., `ITBG-3`) for specific datacenter placement

6. **Regional Placement**: Create volumes in the same region as your compute resources
   - Default region: `ITBG-Bergamo`
   - Use proper region format: `ITBG-Bergamo` (not `IT BG`)

7. **Billing Period**: Choose based on usage duration
   - **Hour**: Short-term testing, development (default)
   - **Month**: Production workloads
   - **Year**: Long-term stable infrastructure (best cost savings)

8. **Update Constraints**: Block storage can only be updated when status is `Used` or `NotUsed`
   - Cannot update during provisioning (`InCreation` state)
   - Only name and tags can be modified
   - Size and type changes not supported

### Snapshots

1. **Getting Volume URI**: Always use `blockstorage get` to get the full URI for snapshot operations
   ```bash
   # Get the URI from the block storage details
   acloud storage blockstorage get 69442fe38f4a09c12b5ded74
   # Copy the URI field for snapshot creation and listing
   ```

2. **Listing Snapshots**: The volume URI is required to list snapshots
   ```bash
   acloud storage snapshot list --volume-uri "/projects/PROJECT_ID/providers/Aruba.Storage/blockStorages/VOLUME_ID"
   ```

3. **Regular Backups**: Create snapshots on a regular schedule
   ```bash
   --tags "schedule=daily,retention=7days"
   --tags "schedule=weekly,retention=4weeks"
   ```

4. **Pre-Change Snapshots**: Always create a snapshot before major changes
   ```bash
   acloud storage snapshot create \
     --name "pre-upgrade-$(date +%Y%m%d)" \
     --region "ITBG-Bergamo" \
     --volume-uri "/projects/PROJECT_ID/providers/Aruba.Storage/blockStorages/VOLUME_ID"
   ```

5. **Naming Convention**: Include date and purpose in snapshot names
   ```bash
   --name "prod-db-daily-20251218"
   --name "pre-migration-backup-20251218"
   ```

6. **Retention Policy**: Implement a clear retention policy
   - Daily snapshots: 7 days
   - Weekly snapshots: 4 weeks
   - Monthly snapshots: 12 months

7. **Cleanup**: Delete old snapshots to manage costs
   ```bash
   # List snapshots for a volume
   acloud storage snapshot list --volume-uri "/projects/PROJECT_ID/providers/Aruba.Storage/blockStorages/VOLUME_ID"
   
   # Delete outdated snapshots
   acloud storage snapshot delete snap-old-123 --yes
   ```

8. **Region Code Handling**: The CLI automatically handles region format conversions
   - API returns region as `IT BG` in GET operations
   - API requires region as `ITBG-Bergamo` in POST/PUT operations
   - The CLI automatically converts formats during updates - no manual intervention needed

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
  --region "ITBG-Bergamo" \
  --volume-uri "/projects/PROJECT_ID/providers/Aruba.Storage/blockStorages/VOLUME_ID"

# 2. Create new volume in target region from snapshot
# (Note: This would typically be done through the API or console)

# 3. Verify and cleanup old resources after migration
```

### Disaster Recovery Setup

```bash
# Set up context for production project
acloud context set prod --project-id "prod-project-id"
acloud context use prod

# Get the volume URI
acloud storage blockstorage get VOLUME_ID

# Create daily snapshot (run via cron/scheduled task)
acloud storage snapshot create \
  --name "dr-backup-$(date +%Y%m%d-%H%M)" \
  --region "ITBG-Bergamo" \
  --volume-uri "/projects/PROJECT_ID/providers/Aruba.Storage/blockStorages/VOLUME_ID" \
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
  --size 10
```

Required: `--name`, `--size`
Optional: `--zone`, `--region`, `--type`, `--billing-period`

### Volume shows empty zone

**Issue:** Volumes created without a zone parameter have an empty zone field.

**Solution:** This is expected behavior. Zone is optional when creating block storage. You can:
- Create new volumes without specifying a zone (region-wide)
- Create new volumes with a specific zone (e.g., `--zone "ITBG-3"`)

### "Cannot update block storage with status 'InCreation'"

**Issue:** Attempting to update a block storage while it's being provisioned.

**Solution:** Wait for the volume to reach `Used` or `NotUsed` status before updating:
```bash
# Check current status
acloud storage blockstorage get VOLUME_ID

# Wait until Status is "Used" or "NotUsed", then update
acloud storage blockstorage update VOLUME_ID --name "new-name"
```

Block storage can only be updated when status is `Used` or `NotUsed`.

### "Location Value: IT BG not found" during update

**Issue:** The API doesn't accept the space-separated region format.

**Solution:** The CLI automatically converts `IT BG` to `ITBG-Bergamo` during updates. If you still see this error, ensure you're using the latest version of the CLI.

### Getting the correct volume URI for snapshots

**Issue:** Not sure what format the volume-uri should be.

**Solution:** Use the `blockstorage get` command to get the full URI:
```bash
# Get volume details including URI
acloud storage blockstorage get 69442fe38f4a09c12b5ded74

# Look for the URI field in the output:
# URI: /projects/PROJECT_ID/providers/Aruba.Storage/blockStorages/VOLUME_ID

# Use this exact URI in snapshot create
acloud storage snapshot create \
  --name "my-snapshot" \
  --region "ITBG-Bergamo" \
  --volume-uri "/projects/PROJECT_ID/providers/Aruba.Storage/blockStorages/VOLUME_ID"
```

### Using verbose mode for debugging

If you encounter issues, use the `--verbose` flag to see the full API response:

```bash
acloud storage blockstorage list --verbose
acloud storage snapshot list --verbose
```

This will show:
- HTTP status codes
- Full API response
- Error details with field-level validation messages
- Raw response body for detailed debugging

### API Error Response Decoding

When an operation fails, the CLI automatically decodes and displays the error:
```bash
$ acloud storage blockstorage update VOLUME_ID --name "test"
API Error (Status 400):
  Title: One or more validation errors occurred.
  Extensions: map[errors:[map[errorMessage:invalid datacenter field:DataCenter]]]
  Raw Response: {"title":"...","status":400,"errors":[...]}
```

Use this information to identify and fix the issue.

## Next Steps

- Explore [Compute Resources](../compute.md)
- Learn about [Network Resources](../network.md)
- Review [Project Management](management/projects.md)
