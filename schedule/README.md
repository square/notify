# schedule (WIP)

scheduling library for oncall rotations etc

Currently supports

1. Multiple entries for a schedule specified in different timezones
1. Overrides
1. Command line interface

### Install

```
go get -u github.com/syamp/schedule/...
```

### Usage

Looking up schedule for "si" for level "primary" for current time
```
schedule --level primary --name si --path examples
us1
```

Looking up schedule for "si" for level "primary" for specified time at location
```
schedule -date="2015-03-12 10:05AM" -location="Asia/Calcutta" -level primary -name si                                                   [/Users/s/schedule]
ind3
```
See examples for input format for si.yaml

### Todo

1. API support
1. Database support
