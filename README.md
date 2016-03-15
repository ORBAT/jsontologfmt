# jsontologfmt

`jsontologfmt` reads [log15](https://gopkg.in/inconshreveable/log15.v2) compatible JSON logs from stdin and pretty prints them to stdout. The names of different keys (like where to find the timestamp and log level) are configurable, but the logging level must be a number from 0 to 4, with 0 being CRIT and 4 being DEBUG.

# Command line flags

```
  -h, --help          Show context-sensitive help (also try --help-long and --help-man).
  -t, --timekey="t"   Which JSON key contains the timestamp
  -a, --timelayout="2006-01-02T15:04:05.999999999Z07:00"
                      What layout the time is in. Defaults to ISO 8601, i.e. Go's time.RFC3339Nano
  -v, --lvlkey="lvl"  Which JSON key contains the level. The level value must be an integer from 0 to 4, with 0 being 'crit' and 4 being 'debug'
  -m, --msgkey="msg"  Which JSON key contains the message
  -l, --level=info    Maximum log level to output
 ```