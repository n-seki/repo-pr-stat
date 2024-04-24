# repo-pr-stat

## Usage

```
$ go run . --help
NAME:
   Repo-PR-Stat - show pull request stat

USAGE:
   Repo-PR-Stat [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --owner value, -o value                        repository owner
   --repository value, -r value                   repository name
   --start value, -s value                        start time
   --end value, -e value                          end time
   --include-base value [ --include-base value ]  include base branch name
   --exclude-base value [ --exclude-base value ]  exclude base branch name
   --token value, -t value                        GitHub Access Token
   --help, -h                                     show help                 show help
```

You can also provide GitHub Access Token via ENV(key=GITHUB_ACCESS_TOKEN).

## Example

```
$ go run . -o DroidKaigi -r conference-app-2023 -s 2023-08-01T00:00:00+09:00 -e 2023-08-07T23:59:59+09:00
{
  "count": 43,
  "averageTimeBetweenCreateMerge": "31h45m30.558139534s",
  "averageTimeBetweenOpenMerge": "26h47m37s",
  "prCountPerUser": {
    ...
  }
}
```
