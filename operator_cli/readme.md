# Operator CLI
This module seeks to optimize the management of the storefront by providing tools for instancing a storefront and maintaining its listings.



## Configuration
```bash
$ cat ./config.yml
OraclePath: ./oracle.key
TreasuryMint: 2Wob5Y6FWYMaxgHuitMLpKxFCC5zJwQUG2C46dm8wtzL
ListingsTable: ./listings.csv
```
### OraclePath
Path to `oracle.json` - Account to be used as parental authority when maintaining storefront lifecycle.

### TreasuryMint
`PublicKey` of Storefront Treasury Mint. E.g. - `$BALLZ`.

### ListingsTable
Path to `file` that will contain storefront listing metadata.



## Usage
### Instance a Storefront
Command:
```bash
go run main.go config.yml instance
```

Executes required initialization commands to spin up a storefront using params in `config.yml`.


### Report All Storefront Listings
Command:
```bash
go run main.go config.yml report
```

Writes a CSV formatted table into `ListingsTable` path in `config.yml`.


### Synchronise Storefront Listing State
Command:
```bash
go run main.go config.yml sync_listings
```

Creates/Ammends listings of their respective records of `ListingsTable` in `config.yml`, writes a `ListingsTable` lockfile, and re-sources storefront listings for `ListingsTable`.
