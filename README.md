# Uberping

## Options

```
Uberping.

Usage:
  uping [options] [<hosts>...]
  uping -h | --help
  uping --version

Options:
  -C <config-file>         Use configuration file, eg. API endpoints, secrets, etc...
  -s                       Be silent and don't print output to stdout
  -g                       Print grouped results
  -p udp|icmp              Set type of ping packet, unprivileged udp or privileged icmp [default: icmp]
  -f				       Use fallback mode, uping will try to use next ping mode if selected by -p failed
  -c <count>               Number of pings to perform [default: 4]
  -i <ping-interval>       Interval between pings, eg. -i 1s, -i 100ms [default: 1s]
  -t <host-timeout>        Timeout before probing one host terminates, regardless of how many pings perfomed, eg. -t 1s, -t 100ms [default: <count> * 1s]
  -w <workers>             Number of paraller workers to run [default: 4]

Sources (may be combined):
  --source-api             Load hosts using external API configured by -C <config-file>
  --source-file <file-in>  Load hosts from file <file-in>

Outputs (may be combined):
  --out-api                Save tests results using external API configured by -C <config-file>
  --out-file <file-out>    Save tests results to file <file-out>
```

## Installation

bla bla

## TODO

0.2:
- add driver for db loader/saver
- add api (and future db) driver tests

0.3:
- add arp protocol
- add fallback of protocol selection

0.4:
- add daemon mode with intervals (or/and non stop option)

0.5:
- polishing code, fix grammar mistakes, typos, etc,
- add more/better comments
- add makefile
- improve readme

0.6 .. 1.0:
- better customization of api/db config schema, eg. custom json requests, template system for endpoints
- add more advanced ping options
- add more probes

### Credits

Application was developed by Tomasz Kolaj and is licensed under Apache License Version 2.0.
Please reports bugs at https://github.com/migotom/uberping/issues.
