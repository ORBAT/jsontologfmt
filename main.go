// package jsontologfmt provides a command line tool for reading log15-compatible (https://gopkg.in/inconshreveable/log15.v2) JSON log messages and outputting them using
// a more human-friendly format using the log15 terminal formatter.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"gopkg.in/inconshreveable/log15.v2"
)
import "gopkg.in/alecthomas/kingpin.v2"

var (
	app        = kingpin.New("jsontologfmt", "Displays log15 (https://gopkg.in/inconshreveable/log15.v2) compatible JSON logs using the log15 terminal formatter")
	logger     = log15.New("main", "jsontologfmt")
	keyNames   log15.RecordKeyNames
	timeLayout string
	logLvl     string
)

func init() {
	app.HelpFlag.Short('h')
	app.Flag("timekey", "Which JSON key contains the timestamp").
		Short('t').Default("t").StringVar(&keyNames.Time)
	app.Flag("timelayout", "What layout the time is in. Defaults to ISO 8601, i.e. Go's time.RFC3339Nano").
		Short('a').Default(time.RFC3339Nano).StringVar(&timeLayout)
	app.Flag("lvlkey", "Which JSON key contains the level. The level value must be an integer from 0 to 4, with 0 being 'crit' and 4 being 'debug'").
		Short('v').Default("lvl").StringVar(&keyNames.Lvl)
	app.Flag("msgkey", "Which JSON key contains the message").
		Short('m').Default("msg").StringVar(&keyNames.Msg)
	app.Flag("level", "Maximum log level to output").
		Short('l').Default("info").EnumVar(&logLvl, "debug", "info", "warn", "error", "crit")

	logger.SetHandler(log15.StderrHandler)
}

func mapToRecord(m map[string]interface{}) *log15.Record {
	ctx := make(context)
	r := new(log15.Record)
	r.KeyNames = keyNames

	for k, v := range m {
		switch k {
		case keyNames.Lvl:
			if jn, ok := v.(json.Number); ok {
				if val, err := jn.Int64(); err == nil && int(val) >= 0 && int(val) <= int(log15.LvlDebug) {
					r.Lvl = log15.Lvl(int(val))
				} else {
					logger.Warn("couldn't convert number under level key to int64?", "raw_val", v, "raw_type", fmt.Sprintf("%T", v), "error_key", keyNames.Lvl)
				}
			} else {
				logger.Warn("level key didn't contain a number", "raw_val", v, "raw_type", fmt.Sprintf("%T", v), "error_key", keyNames.Lvl)
				r.Lvl = log15.LvlInfo
			}

		case keyNames.Msg:
			if val, ok := v.(string); ok {
				r.Msg = val
			} else {
				logger.Warn("message key contained something that isn't a string?", "raw_val", v, "raw_type", fmt.Sprintf("%T", v), "error_key", keyNames.Msg)
			}

		case keyNames.Time:
			if val, ok := v.(string); ok {
				if t, err := time.Parse(timeLayout, val); err == nil {
					r.Time = t
				} else {
					logger.Warn("weird time in log message", "raw_val", v, "raw_type", fmt.Sprintf("%T", v), "err", err, "error_key", keyNames.Time)
				}
			} else {
				logger.Warn("expected to find a string under the time key", "raw_val", v, "raw_type", fmt.Sprintf("%T", v), "error_key", keyNames.Time)
			}

		default: // everything else is part of the actual log message, so shove it in the context
			ctx[k] = v
		}
	}

	r.Ctx = ctx.toArray()

	return r
}

func main() {
	_, err := app.Parse(os.Args[1:])
	app.FatalIfError(err, "")
	lvl, _ := log15.LvlFromString(logLvl)

	var m map[string]interface{}

	dec := json.NewDecoder(os.Stdin)
	dec.UseNumber()

	fmt := log15.LvlFilterHandler(lvl, log15.StdoutHandler)

	for err := dec.Decode(&m); err == nil; err = dec.Decode(&m) {
		fmt.Log(mapToRecord(m))
		clearMap(m)
	}

}

func clearMap(m map[string]interface{}) {
	for k, _ := range m {
		delete(m, k)
	}
}

// same as log15.Ctx, but copy-pasted since toArray() was unexported
type context map[string]interface{}

func (c context) toArray() []interface{} {
	arr := make([]interface{}, len(c)*2)

	i := 0
	for k, v := range c {
		arr[i] = k
		arr[i+1] = v
		i += 2
	}

	return arr
}
