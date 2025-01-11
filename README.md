# `shamir-msg` - Shamir secret sharing for everyone

Brings the Shamir secret sharing algorithm from Hashicorp Vault to a CLI utility for usage outside Vault. Shamir breaks a secret key into shards such that a certain threshold is required to re-assemble the keys. This allows splitting a key into cryptographically isolated keys where no one key gives information about the final secret.

## Installation

```
go install github.com/pk910/shamir-msg
```

## Usage

```
$ shamir-msg --help
NAME:
   shamir-msg - Split and recombine a secret using Shamir's Secret Sharing

USAGE:
   shamir-msg [global options] [command [command options]]

COMMANDS:
   split    Split a secret into shards
   combine  Combine shards into the original secret
   run      Run the tool in interactive mode (default)
   help, h  Shows a list of commands or help for one command
```
## Examples

### Splitting a secret key
```
NAME:
   shamir-msg split - Split a secret into shards

USAGE:
   shamir-msg split [command [command options]]

OPTIONS:
   --keys value, -k value        The number of total key shards to generate (default: 3)
   --threshold value, -t value   The min number of shards needed to re-assemble the secret (default: 2)
   --secret value, -s value      The secret to split into shards
   --group-size value, -g value  Split output into groups of this size and separate with a space (default: 6)
   --quiet, -q                   Suppress everything except the secret output (default: false)
   --help, -h                    show help
```

Example:
```
$ shamir-msg split -k 5 -t 3 -s "supersecret"
##### Shamir Secret Sharing Tool #####
Split and recombine a secret using Shamir's Secret Sharing

Total keys shards:  5
Required for reconstruction:  3

Shard 1: 6VV4EF ICUMAL 978ZAU T8ST2
Shard 2: LIJBTL BCKIJL SDDTZB EF44S
Shard 3: LPMUYB AUFEP4 IRC5V8 4YSF2
Shard 4: UL8G3I 2CRJ6P 25JDXV 2SDU2
Shard 5: S8BZGC A9PRNV W9DUY8 FMXFS
```

### Combining a secret key
```
NAME:
   shamir-msg combine - Combine shards into the original secret

USAGE:
   shamir-msg combine [command [command options]]

OPTIONS:
   --shard value, -s value [ --shard value, -s value ]  The shards to recombine into the secret
   --quiet, -q                                          Suppress everything except the secret output (default: false)
   --help, -h                                           show help
```

```
$ shamir-msg combine -s "6VV4EF ICUMAL 978ZAU T8ST2" -s "LPMUYB AUFEP4 IRC5V8 4YSF2" -s "UL8G3I2CRJ6P25JDXV2SDU2"
##### Shamir Secret Sharing Tool #####
Split and recombine a secret using Shamir's Secret Sharing

Secret: supersecret

```

### Interactive mode

Simply run `shamir-msg` for a interactive version of the tool.

## Disclaimer

This tool uses the shamir secret sharing implementation from [HashiCorp/Vault](https://github.com/hashicorp/vault/tree/main/shamir)