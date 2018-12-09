# Uberping

## Options

```
Uberping.

Usage:
  uping [options] [<hosts>...]
  uping -h | --help
  uping --version

Options:
  -I <tests-interval>      Interval between tests, if provided uping will perform tests indefinitely, e.g. every -I 1m, -I 1m30s, -I 1h30m10s
  -C <config-file>         Use configuration file, e.g. API endpoints, secrets, etc...
  -s                       Be silent and don't print output to stdout
  -g                       Print grouped results
  -p udp|icmp              Set type of ping packet, unprivileged udp or privileged icmp [default: icmp]
  -f                       Use fallback mode, uping will try to use next ping mode if selected by -p failed
  -c <count>               Number of pings to perform [default: 4]
  -i <ping-interval>       Interval between pings, e.g. -i 1s, -i 100ms [default: 1s]
  -t <host-timeout>        Timeout before probing one host terminates, regardless of how many pings performed, e.g. -t 1s, -t 100ms [default: <count> * 1s]
  -w <workers>             Number of parallel workers to run [default: 4]

Sources (may be combined):
  --source-db              Load hosts using database configured by -C <config-file>
  --source-api             Load hosts using external API configured by -C <config-file>
  --source-file <file-in>  Load hosts from file <file-in>

Outputs (may be combined):
  --out-db                 Save tests results database configured by -C <config-file>
  --out-api                Save tests results using external API configured by -C <config-file>
  --out-file <file-out>    Save tests results to file <file-out>
```
## Usage examples

blah blah

## Installation

blah blah

## Features
 
blah blah

## TODO

0.1:
- [x] add drivers for file, argv/stdout, api loader/saver
- [x] add argv and toml config parsers
- [x] add support for icmp/udp pinging
- [x] run tests tasks parallel with gorutines pool

0.2:
- [x] add driver for db loader/saver
- [x] add api tests
- [x] add searching for default config, linux: etc, home, -C, macosx: ... windows: ....
- [x] modify example config
- [x] update readme of uping command help

0.3:
- [x] add daemon mode with intervals (or/and nonstop option)
- [ ] add retry to db driver
- [ ] add gorutines for loaders/savers
- [ ] add db and new schema tests
- [ ] add windows config loading as well
- [ ] update example config, add comments describing API/DB fields

0.4:
- [ ] add arp protocol
- [ ] add fallback of protocol selection

0.5:
- [ ] polishing code, fix grammar mistakes, typos, etc,
- [ ] organize depedencies as third party modules
- [ ] add more/better comments
- [ ] add makefile
- [ ] improve readme (better description, features, config loading sequence, etc)

0.6 .. 1.0:
- [ ] better customization of api/db config schema, eg. custom json requests, template system for endpoints
- [ ] add more advanced ping options
- [ ] add more probes

### Credits

Application was developed by Tomasz Kolaj and is licensed under Apache License Version 2.0.
Please reports bugs at https://github.com/migotom/uberping/issues.
