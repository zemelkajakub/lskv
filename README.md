# lskv

`lskv` is a CLI for discovering and retrieving Azure Key Vault secrets.

It is built for:
- fast secret lookup from local cache
- profile-based work across subscriptions/environments
- pipeline-friendly batch retrieval (`find | get -`)

## Why lskv first

If your daily flow is "find secret name quickly" and "fetch values safely across multiple vaults", `lskv` is faster and simpler than stitching many `az` commands together.

What you get immediately:
- **search once, reuse often** with local cache (`cache refresh` + `find`)
- **stable pipeline format** (`vault:secret`) for automation
- **environment context** via profiles (DEV/QA/PROD)
- **direct mode** for quick one-off calls when profile is unnecessary

> 💡 Tip: You can also use `lskv` directly without any profile for single-vault operations.

---

## 1) 🔐 First step: authenticate to Azure

```bash
az login
```

For a specific tenant:

```bash
az login --tenant <tenant-id>
```

---

## 2) 📦 Install

### Option A: one-line install (recommended)

Installs latest release to `/usr/local/bin`:

```bash
curl -fsSL https://raw.githubusercontent.com/zemelkajakub/lskv/develop/scripts/install.sh | sudo bash
```

Install to user path instead:

```bash
curl -fsSL https://raw.githubusercontent.com/zemelkajakub/lskv/develop/scripts/install.sh | bash -s -- --install-dir "$HOME/.local/bin"
```

### Option B: build locally

```bash
go build -o lskv
./lskv --help
```

---

## 3) 🚀 Quick start

### No profile needed (direct mode)

If you already know the vault, these commands work without configuring a profile:

```bash
lskv get <vault:secret>
lskv list secrets <vault>
```

Use profiles when you want cache/search workflows (`find`, `list secrets all`, `cache ...`).

### Create and activate profile

```bash
lskv profile init DEV --subscription-id <subscription-id>
```

Scoped profile (only selected vaults):

```bash
lskv profile init PROD --subscription-id <subscription-id> --description "Production" --vaults kv-prod1,kv-prod2
```

### See active profile

```bash
lskv profile show
```

### Refresh cache

```bash
lskv cache refresh
```

### Find secret names

```bash
lskv find traefik
```

### Get values

Single secret:

```bash
lskv get <vault:secret>
```

List one vault directly (cache-first, API fallback):

```bash
lskv list secrets <vault>
```

Batch from pipeline:

```bash
lskv find traefik | grep -E "dev-" | lskv get -
```

---

## 4) 👤 Profiles explained (important)

Profiles define the working context for cache/search/list workflows.

Think of a profile as:
- target Azure subscription
- optional list of allowed vaults
- optional description for team clarity

### How profile scope works

1. **Unscoped profile** (no `--vaults`):
	- `cache refresh` discovers all accessible vaults in subscription
	- best for exploration and broad discovery

2. **Scoped profile** (with `--vaults`):
	- only listed vaults are considered
	- best for production safety and team boundaries

### Typical team setup

```bash
# Developer broad discovery profile
lskv profile init DEV --subscription-id <subscription-id>

# Production restricted profile
lskv profile init PROD --subscription-id <subscription-id> --vaults kv-prod1,kv-prod2 --description "Production vaults only"

# Switch context explicitly
lskv profile switch DEV
lskv profile show
```

### Profile lifecycle commands

```bash
lskv profile list
lskv profile show
lskv profile show DEV
lskv profile switch PROD
lskv profile delete DEV
```

---

## 5) ⚖️ Why use lskv vs az

`az` is great for broad Azure management, but `lskv` is optimized for secret discovery/retrieval workflows.

Key benefits:
- faster secret name discovery with local cache (`find`)
- consistent `vault:secret` output contract for pipelines
- profile-based context for teams/environments
- direct batch retrieval (`find ... | lskv get -`)

### Cache benefit for vault-less profiles

If a profile has **no explicit `--vaults`**, `lskv cache refresh` discovers accessible vaults and stores their secret names locally.

That gives you:
- quick search across many vaults without repeatedly calling Azure APIs
- lower latency for incident/debug scenarios
- easier discovery when you know only part of a secret name, not the vault

Typical flow:

```bash
lskv profile init DEV --subscription-id <subscription-id>
lskv cache refresh
lskv find payment
lskv get <vault:secret>
```

---

## 6) 🧭 Commands and behavior

### Profiles

```bash
lskv profile list
lskv profile show
lskv profile switch DEV
lskv profile delete DEV
```

### Cache

```bash
lskv cache refresh
lskv cache status
lskv cache clear
```

### List

```bash
lskv list vaults
lskv list secrets all
lskv list secrets <vault>
```

Behavior:
- `list secrets all` reads from cache only
- `list secrets <vault>` reads from cache first, then falls back to Azure API
- `list vaults` verifies **data-plane access** to each candidate vault
- `get <vault:secret>` works directly against Azure Key Vault (no profile required)

### Find / Get format contract

- `find` prints `vault:secret`
- `get -` reads `vault:secret` from stdin
- one-line values are printed in aligned columns
- multiline values (e.g., certs/keys) are printed as blocks under `vault:secret`

---

## 7) 🛠️ Example-driven workflows

### A) I know the vault and secret name already

```bash
lskv get kv-dev-app:DbPassword
```

### B) I only know part of the secret name

```bash
lskv profile switch DEV
lskv cache refresh
lskv find db-password
```

Then pick one result and fetch value:

```bash
lskv get kv-dev-app:db-password
```

### C) Get many matching secrets at once

```bash
lskv find connection-string | lskv get -
```

Filter to subset before retrieval:

```bash
lskv find traefik | grep -E "prod-" | lskv get -
```

### D) List secrets from one vault (cache-first)

```bash
lskv list secrets kv-prod1
```

### E) List all cached secrets in current profile

```bash
lskv list secrets all
```

### F) Incident response lookup

```bash
lskv profile switch PROD
lskv cache refresh
lskv find payment
lskv get kv-prod1:payment-api-key
```

---

## 8) 🛠️ Common use cases

### Incident response lookup
1. `lskv cache refresh`
2. `lskv find <keyword>`
3. `lskv get <vault:secret>`

### Compare same secret pattern across vaults

```bash
lskv find connection-string | lskv get -
```

### Environment isolation

Use `--vaults` in profile init to restrict scope per environment/team.

---

## 9) 📝 Notes

- set `NO_COLOR=1` to disable terminal colors
- if auth errors mention tenant mismatch, re-login with correct tenant
- default installer target is `/usr/local/bin`
