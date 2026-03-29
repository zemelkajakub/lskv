# lskv

`lskv` is a CLI for discovering and retrieving Azure Key Vault secrets.

It is built for:
- fast secret lookup from local cache
- profile-based work across subscriptions/environments
- pipeline-friendly batch retrieval (`find | get -`)

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
curl -fsSL https://raw.githubusercontent.com/zemelkajakub/lskv/main/scripts/install.sh | sudo bash
```

Install to user path instead:

```bash
curl -fsSL https://raw.githubusercontent.com/zemelkajakub/lskv/main/scripts/install.sh | bash -s -- --install-dir "$HOME/.local/bin"
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

## 4) ⚖️ Why use lskv vs az

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

## 5) 🧭 Commands and behavior

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

## 6) 🛠️ Common use cases

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

## 7) 📝 Notes

- set `NO_COLOR=1` to disable terminal colors
- if auth errors mention tenant mismatch, re-login with correct tenant
- default installer target is `/usr/local/bin`
