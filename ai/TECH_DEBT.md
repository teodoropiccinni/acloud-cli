# TECH_DEBT.md — Technical Debt & Refactoring Backlog

Issues are grouped by severity. Address Critical items before new features ship; High items before any public release.

## Resolved

| ID | Summary |
|----|---------|
| TD-001 | All `Run` handlers converted to `RunE`; `SilenceUsage: true`; `fmtAPIError` helper in `root.go` |
| TD-002 | Nil guards added for `LocationResponse`, `Metadata.ID`, `Metadata.Name` in all list/create/update responses |
| TD-003 | Errors propagated via `return fmt.Errorf(...)` instead of printing to stdout |
| TD-004 | Flag read errors checked via `RunE` return paths |
| TD-005 | `confirmDelete()` helper in `root.go` detects non-interactive stdin before prompting |
| TD-006 | `newCtx()` helper in `root.go` applies 30-second timeout to all SDK calls |
| TD-007 | `getContextFilePath` returns `(string, error)` instead of silently falling back to CWD |
| TD-008 | YAML unmarshal errors wrapped with user-friendly messages in `LoadConfig` and `LoadContext` |
| TD-013 | `Args: cobra.NoArgs` added to all `create` and `list` commands that take no positional arguments |
| TD-014 | `cmd/constants.go` created with `StateInCreation`, `DateLayout`, `FilePermConfig`, `FilePermDirAll`; all magic strings replaced |
| TD-009 | `MarkFlagRequired` used as the single mechanism for all required flags; redundant `if flag == ""` manual checks removed from all 19 affected commands |
| TD-011 | `readSecret()` helper added to `root.go` using `golang.org/x/term.ReadPassword`; `config set` now prompts interactively when `--client-secret` is not passed and no secret exists in config |

---

---

## Critical

### TD-001 · `Run` instead of `RunE` — CLI always exits 0 on error
**All 50+ command handlers use `Run` instead of `RunE`.** When an error occurs the handler prints to stdout and returns, so the process exits with code 0. Scripts, CI pipelines, and automation cannot detect failures.

**Fix:** Convert every `Run` to `RunE`. Return `fmt.Errorf(...)` instead of printing and returning. Let Cobra propagate exit codes.

```go
// Before
Run: func(cmd *cobra.Command, args []string) {
    if err != nil { fmt.Println("Error:", err); return }
}
// After
RunE: func(cmd *cobra.Command, args []string) error {
    if err != nil { return fmt.Errorf("...: %w", err) }
    return nil
}
```

**Files:** every file in `cmd/` that declares a command.

---

### TD-002 · Nil pointer dereferences on SDK response fields
Several list and get commands dereference `LocationResponse`, `Flavor`, and similar nested structs without nil-guarding. These panic at runtime if the API returns a partial response.

**Known unsafe patterns:**
- `response.Data.Metadata.LocationResponse.Value` — `cmd/network.securitygroup.go`, `cmd/network.vpc.go`, `cmd/storage.snapshot.go`
- `response.Data.Properties.Flavor.Name/CPU/RAM/HD` — `cmd/compute.cloudserver.go`
- Various table-building loops that dereference `*string` response fields without checking

**Fix:** Add nil guards for every nested struct and pointer field before use:
```go
if r.Metadata.LocationResponse != nil {
    region = r.Metadata.LocationResponse.Value
}
```

---

### TD-003 · Error messages go to stdout, not stderr
Every error is printed with `fmt.Printf`/`fmt.Println` to stdout. This corrupts stdout for any consumer trying to parse command output and prevents reliable stderr redirection.

**Fix:** Use `fmt.Fprintln(os.Stderr, ...)` (or `cmd.PrintErr`) for all error output. Pair with TD-001 (RunE) to propagate errors rather than printing them inline.

---

## High

### TD-004 · Swallowed errors on every flag read
Flag reads use `_, _` almost universally:
```go
name, _ := cmd.Flags().GetString("name")
```
If the flag is not registered or has a type mismatch the returned error is silently discarded and `name` becomes `""`. This causes confusing downstream failures.

**Fix:** Check the error from every `cmd.Flags().Get*()` call. With `RunE` this becomes easy — just return the error.

---

### TD-005 · `fmt.Scanln` blocks indefinitely in non-interactive environments
Delete confirmations use `fmt.Scanln(&response)` which blocks forever when stdin is a pipe or `/dev/null` (CI, containers, cron). The process hangs until the parent times out.

**Fix:** Detect non-interactive stdin before prompting. A minimal safe pattern:
```go
fi, _ := os.Stdin.Stat()
if (fi.Mode() & os.ModeCharDevice) == 0 {
    return fmt.Errorf("delete requires --yes in non-interactive mode")
}
```

**Files:** all `*DeleteCmd` handlers (`cmd/compute.cloudserver.go`, `cmd/storage.blockstorage.go`, `cmd/network.vpc.go`, `cmd/database.dbaas.go`, `cmd/security.kms.go`, etc.).

---

### TD-006 · No API call timeout — CLI can hang indefinitely
All SDK calls use `context.Background()` with no deadline. A slow or unresponsive API endpoint hangs the CLI forever.

**Fix:** Apply a configurable default timeout (e.g. 30 s) for all SDK calls:
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
```
Optionally expose `--timeout` as a global flag.

---

### TD-007 · `getContextFilePath` silently falls back to current working directory
```go
// cmd/context.go
func getContextFilePath() string {
    home, err := os.UserHomeDir()
    if err != nil {
        return ".acloud-context.yaml"  // writes context file to CWD
    }
    ...
}
```
If home-directory resolution fails, contexts are silently written to (and read from) the current directory. Any subsequent invocation from a different directory loses all contexts with no warning.

**Fix:** Return an error instead of a fallback path. Callers that ignore the error should be updated to propagate it.

---

### TD-008 · Corrupted config/context YAML gives raw parser error
When `~/.acloud.yaml` or `~/.acloud-context.yaml` is malformed, the user sees a raw `yaml:` parse error with no actionable guidance.

**Fix:** Wrap YAML unmarshal errors with user-friendly context:
```go
if err := yaml.Unmarshal(data, &cfg); err != nil {
    return nil, fmt.Errorf("config file %s is corrupted (%w). Delete it and run 'acloud config set' to reconfigure", configPath, err)
}
```

---

### TD-009 · Inconsistent use of `MarkFlagRequired`
Several `create` commands mark required flags properly; others rely on manual validation inside `Run`. The two approaches coexist in the same codebase:

- **Marked:** `cmd/storage.blockstorage.go`, `cmd/compute.cloudserver.go`, `cmd/storage.snapshot.go`
- **Not marked (validated manually in Run):** `cmd/network.vpc.go`, `cmd/network.securitygroup.go`, `cmd/config.go`
- **Both:** `cmd/storage.blockstorage.go` marks AND validates — redundant.

**Fix:** Use `MarkFlagRequired` as the single mechanism for all truly required flags. Remove manual flag-presence checks in `Run`.

---

### TD-010 · No test coverage for command `Run`/`RunE` functions
Only helper functions (`LoadConfig`, `SaveConfig`, `GetProjectID`, `PrintTable`, etc.) have tests. Zero of the 42 command files have tests that invoke an actual command handler. Every `Run` body is entirely untested.

**Affected files:** All `cmd/<category>.<resource>.go` files.

**Fix:** Add table-driven tests using `cobra.Command.Execute()` with flag injection. Mock the SDK client via an interface so tests don't require live credentials.

---

### TD-011 · Credentials passed as CLI flags — visible in shell history and `ps`
`acloud config set --client-id X --client-secret Y` exposes the secret in:
- Shell history files (`.bash_history`, `.zsh_history`)
- Process listings (`ps aux`)
- CI/CD log output

**Fix:** When `--client-secret` is not provided on the command line, prompt for it interactively with echo disabled (use `golang.org/x/term.ReadPassword`).

---

## Medium

### TD-012 · `--debug` may log credentials and tokens from HTTP traffic
`WithNativeLogger()` enables full SDK HTTP logging to stderr. Authorization headers and request bodies potentially containing tokens or secrets are logged without sanitization.

**Fix:** Add a warning to the `--debug` flag description. Investigate whether the SDK logger exposes authorization headers; if so, add a sanitizing log interceptor.

---

### TD-013 · Missing `Args: cobra.NoArgs` on create and list commands
`create` and `list` commands silently ignore any positional arguments the user mistypes:
```
acloud storage blockstorage create accidental-extra-arg --name foo ...
```
The extra arg is accepted without error.

**Fix:** Add `Args: cobra.NoArgs` to all `create` and `list` commands.

---

### TD-014 · Magic strings scattered across the codebase
Repeated literals with no named constant:
- State values: `"InCreation"`, `"Used"`, `"NotUsed"` — referenced in multiple network and storage files
- Default region: `"ITBG-Bergamo"` — `cmd/storage.blockstorage.go:27`
- Error prefix: `"Error: "` — 100+ occurrences
- File permission modes: `0600`, `0755` — no constants defined
- Date format: `"02-01-2006 15:04:05"` — multiple files

**Fix:** Define a `constants.go` (or equivalent) in `cmd/` for shared values.

---

### TD-015 · Fragile raw-JSON workaround for CloudServer ID extraction
`cmd/compute.cloudserver.go` manually unmarshals `response.RawBody` into `map[string]interface{}` to extract the resource ID because the SDK's typed response struct does not expose it. This breaks silently if the API response shape changes.

**Fix:** Either update the SDK to expose the ID in its typed response, or — if the SDK cannot be changed — add a test that fails if the JSON structure changes unexpectedly.

---

### TD-016 · No JSON/machine-readable output format
All output is human-readable tables and formatted strings. There is no `--output json` or `--output yaml` flag. Consumers integrating acloud-cli into scripts must parse table output with fragile text processing.

**Fix:** Add a global `--output` flag supporting `table` (default) and `json`. Implement a thin output-format layer that receives structured data and serializes it in the requested format.

---

### TD-017 · No pagination for list commands
List commands pass `nil` as the options parameter, fetching all resources in a single call:
```go
client.FromCompute().CloudServers().List(ctx, projectID, nil)
```
Large accounts could return thousands of resources, exhausting memory and producing slow, unusable table output.

**Fix:** Add `--limit` and `--offset` flags; pass them through the SDK options struct if the API supports pagination.

---

### TD-018 · Global mutable state breaks parallel tests and requires manual cache reset
Package-level variables in `cmd/root.go` (`clientCache`, `cachedClientID`, etc.) are reset manually in tests. Running tests with `-parallel` causes race conditions. Any future test that forgets the cleanup will inherit a stale client.

**Fix:** Encapsulate the client and its cached state in a struct. Inject it via a package-level variable that tests can replace, or use `t.Cleanup` to reset it reliably.

---

## Low

### TD-019 · No `--dry-run` flag on destructive operations
Users cannot validate that a delete command would succeed (permissions, resource existence) without actually deleting the resource.

**Fix:** Add `--dry-run` to all delete commands. In dry-run mode, perform a `get` to validate the resource exists and the client has access, then print what would happen without calling `delete`.

---

### TD-020 · Inconsistent success messages across resources
Success messages vary: `"...created successfully!"`, `"...initiated. Use 'get' to check status."`, `"Resource created, but no details returned."`. The last form is especially confusing — it implies partial success.

**Fix:** Standardise to two patterns: a definitive success message when the API confirms completion, and an async-operation message when the API returns accepted-but-pending.

---

### TD-021 · Many commands missing `Long` descriptions
Most leaf commands (get, create, update, delete) have only a `Short` description. `acloud <resource> create --help` provides minimal guidance on valid flag values or usage examples.

**Fix:** Add `Long` and `Example` fields to at minimum all `create` commands, documenting flag valid values and a copy-paste example invocation.

---

### TD-022 · Pre-release SDK version (v0.1.x)
`go.mod` depends on `github.com/Arubacloud/sdk-go v0.1.21`. The `0.x` major version provides no semantic versioning stability guarantee — a minor-version bump may introduce breaking changes.

**Fix:** Track the SDK release roadmap. When a `v1.0.0` is released, migrate and pin to it. Until then, pin to a specific minor version and treat any upgrade as potentially breaking.
