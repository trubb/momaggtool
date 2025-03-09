# momaggtool

A tool for dealing with trend following strategies.

Like aztest, but in CLI (and eventually TUI?) form.

```bash
go run ./cmd/ --help
NAME:
   momaggtool - A tool for dealing with trend following strategies

USAGE:
   momaggtool [global options] command [command options]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --funds FILE, --ff FILE  Load funds from toml-formatted FILE
   --help, -h               show help
```

```bash
go run ./cmd/ --ff funds.toml
```
