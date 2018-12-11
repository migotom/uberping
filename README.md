# Uberping

## Options

```
Uberping.

Usage:
  uping [options] [<hosts>...]
  uping -h | --help
  uping --version

Options:
  -d <tests-interval>      Interval between tests, if provided uping will perform tests indefinitely, e.g. every -I 1m, -I 1m30s, -I 1h30m10s
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

## Installation

blah blah

## Features
 
### Implemented:

- ping hosts using unprivileged udp or privileged icmp
- probe hosts using netcat like establishing tcp connection for specified service port
- take hosts to test from command line, file, database (currently only postgresql) and external REST API
- save test results to file, database and external REST API
- ability to combine input sources and outputs, eg. load hosts from file and database (list of hosts are refreshed before each tests iteration)
- run tests in parallel (configurable amount of test workers)
- print output live or groupped (may be needed to more human readable result from parallel tests)
- load settings from config TOML file (searching sequence below)
- ablity to run in continous mode with user defined intervals between tests
- DB/API connection retries

### Not yet implemented:

- ARP protocol
- fallback to other protocol in case of failure
- Windows support
- better customization
- more advanced ping options
- mode probes

### Config loading sequence (the first least important):

- Application defaults
- System (/etc/uping/config.toml, /Library/Application Support/Uping/config.toml)
- Home (~/.uping.toml, ~/Library/Application Support/Uping/config.toml)
- Command line -C option

### Note on Linux Support:

For use unprivileged ping via UDP on linux for regular (non super-user):

```
sudo sysctl -w net.ipv4.ping_group_range="0 2147483647"
```

If wish to use ICMP raw sockets mode as regualr user:

```
setcap cap_net_raw=+ep /bin/uping
```

### Credits

Uberping is highly inspired by [go-ping](https://github.com/sparrc/go-ping/).

Application was developed by Tomasz Kolaj and is licensed under Apache License Version 2.0.
Please reports bugs at https://github.com/migotom/uberping/issues.
