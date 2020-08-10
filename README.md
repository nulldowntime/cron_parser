# Crontab Line Parser

Author: Daniel Siechniewicz <daniel@nulldowntime.com>

## Intro

This is a command line application that parses a single line in crontab format and returns a table with all time definitions expanded to their numeric equivalents, plus a command "as is".

## Setting up

No external (non-standard) packages are used. As long as go is installed it should be possible to compile and/or run this utility. go.mod is provided to be able to use this outside of go path, although it uses a fake module path.

### Run in DEV

go run parse_cron.go "*/15 0 1,15 * 1-5 /usr/bin/find"

### Build

go build -o parse_cron

### Run binary

./parse_cron "*/15 0 1,15 * 1-5 /usr/bin/find"

## Caveats/Missing features

* Comments are not handled
* The command is not sanitized/escaped at the moment
* No effort is made to fix any data, it's either handled as is, or it's a failure
* Similarly, nothing is sorted, duplicate list items are not removed, etc
* There are certainly more cases to be handled, these are just the preliminary ones (i.e. should a too large interval be ignored, as it is now, or an error?)
* Some valid configurations are not supported, especially where there exists some ambiguity in the crontab format, like day of the week number 0
* Special strings (@yearly, etc) are not supported
