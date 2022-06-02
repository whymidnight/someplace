# Operator CLI
This module seeks to optimize the management of the storefront by providing tools for instancing a storefront and maintaining its listings.



## Configuration
```bash
$ cat ./config.yml
OraclePath: ./oracle.key
```
### OraclePath
Path to `oracle.json` - Account to be used as parental authority when maintaining storefront lifecycle.

## Usage
### Instance a Storefront
Command:
```bash
go run main.go instance config.yml
```

Executes required initialization commands to spin up a storefront using params in `config.yml`.


### Report All Quests
Command:
```bash
go run main.go report_quests config.yml
```

Writes a CSV formatted table into `ListingsTable` path in `config.yml`.


### Synchronise Storefront Listing State
Command:
```bash
go run main.go sync_quests config.yml
```

Creates/Ammends listings of their respective records of `ListingsTable` in `config.yml`, writes a `ListingsTable` lockfile, and re-sources storefront listings for `ListingsTable`.


### Reset Quest Rewards
Command:
```bash
go run main.go reset_quest_rewards config.yml
```

Creates/Ammends the splits of tender to their respective addresses in `Splits` in `config.yml` for storefront.
