# profile

- lskv profile init DEV --subscription-id 0000-1111
- lskv profile init PROD --subscription-id 2222-3333 --description "Production Subscription" --vaults prod1-kv,prod2-kv
- lskv profile list
- lskv profile show
- lskv profile show DEV
- lskv profile switch DEV
- lskv profile delete DEV
- lskv profile delete PROD DEV

# cache

- lskv cache status
- lskv cache refresh
- lskv cache clear

# find

- lskv find traefik

# list

- lskv list secrets all
- lskv list secrets prod1-kv

# get

- lskv get vault:secret
- lskv find traefik | grep -E "prod*" | lskv get -