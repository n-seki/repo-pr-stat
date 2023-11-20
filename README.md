# repo-pr-stat

## Usage

```
$ go run . -help
NAME:
   Repo-PR-Stat - show pull request stat

USAGE:
   Repo-PR-Stat [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --owner value, -o value       repository owner
   --repository value, -r value  repository name
   --start value, -s value       start time
   --end value, -e value         end time
   --token value, -t value       GitHub Access Token
   --help, -h                    show help
```

You can also provide GitHub Access Token via ENV(key=GITHUB_ACCESS_TOKEN).

## Example

```
$ go run . -o DroidKaigi -r conference-app-2023 -s 2023-08-01T00:00:00+09:00 -e 2023-08-07T23:59:59+09:00
{
  "count": 45,
  "averageTimeBetweenCreateClose": "30h28m33.511111111s",
  "averageTimeBwtweenOpenClose": "25h43m54.333333333s",
  "prCountPerUser": {
    ...
  }
}
```
