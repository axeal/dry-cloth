# dry-cloth

A tool for cleaning up old Digital Ocean droplets.

```
$ dry-cloth -h
NAME:
   dry-cloth - clean up old droplets

USAGE:
   dry-cloth [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --access-token value  Digital Ocean access token [$DIGITALOCEAN_ACCESS_TOKEN]
   --preserve-tag value  Tag to prevent droplet deletion
   --max-age-days value  Maximum age of droplets to keep (default: 14)
   --dry-run             Dry run without deleting droplets (default: false)
   --help, -h            show help
```

## Usage

Provide your Digital Ocean access token via the argument `--access-token` or environment variable `DIGITALOCEAN_ACCESS_TOKEN`.

Set the maximum age of droplets to retain via the `--max-age-days` argument (default is 14).

To preserve droplets with a specific tag set the `--preserve-tag` argument.

To print a list of droplets that would be deleted, without performing any actual deletions, add the `--dry-run` flag.

e.g. to list the droplets older than 28 days, except those with the tag DoNotDelete, that would be deleted, with an access-token of `123abc`:

```
dry-cloth --access-token 123abc --preserve-tag DoNotDelete --max-age-days 28 --dry-run
```
