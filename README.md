# arp-watch

The `arp-watch` tool detects changes to the local ARP (Address Resolution Protocol) cache. It currently supports Linux, UNIX, and MacOS. 

## Installation

Make sure you have Go installed and your `GOPATH` set.

```
go get github.com/st3v/arp-watch
```

## Usage

The `arp-watch` tool will report changes to MAC addresses inside the local ARP cache. A change is reported using an event key in the form of `net.arp.<key>.<action>`. Thereby, `<key>` represents a sanitized IP or an IP alias (see config options). `<action>` represents the type of change, i.e.:

* `set`: a previously unknown entry has been detected in the ARP cache
* `unset`: a previously observed entry has disappered from the ARP cache
* `changed`:  a previously observed entry has changed in the ARP cache (i.e. the MAC address for a given IP changed)

The tool reports change events on the command line and optionally sends them as a metric to a [Cloud Foundry Metron agent](http://docs.cloudfoundry.org/loggregator/architecture.html).

Example:

```
$ arp-watch -stateFilePath /tmp/arp-watch.state -configPath /tmp/arp-watch.config
net.arp.192_168_0_1.set: '' -> '00:11::22:33:44:55'
net.arp.192_168_0_2.changed: '66:77::88:99:aa:bb' -> '66:77::88:99:aa:bb'
net.arp.192_168_0_3.unset: 'aa:bb::cc:dd:ee:ff' -> ''
```

### Command Line Parameters

```
$ arp-watch --help
Usage of ./arp-watch:
  -configPath string
    	Path to config file. Optional.
  -stateFilePath string
    	Path to state file. Optional.
```

If `-configPath` is not being passed, the tool will check the ARP cache exactly once and return immediately afterwards. The tool will not apply filters and aliases (see below).

If `-stateFilePath` is not being passed, the tool will initially report the currently cached addresses as being *set*. (i.e. there is not initial state). If the `-stateFilePath` flag is set and it points to an existing state file, the tool will set its initial state accordingly. Upon termination the tool will write its state to the specified `-statefilePath`.

### Config File Options

The config file uses JSON as its format. The following options are available. All of them are optional.

**`frequency`**

If `frequency` is set the tool will keep running indefinitely and continuously check for changes in the ARP cache. The `frequency` defines how often the tool checks the ARP cache.  It's a string that is a sequnce of decimal numbers each with and optional fraction and a unit suffix, e.g. "300ms", "1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h". If `frequency` is not set, the tool will check the ARP cache once and terminate immediately.

**`filters`**

This is an array of strings, each representing an IP that should be considered when watching the ARP cache. If `filters` is not set or set to an empty array, the tool will consider all cached IPs.

**`aliases`**

This option can be used to map IP addresses to arbitrary keys, e.g. `'192.168.0.1': 'node-1'`. The tool uses this mapping when reporting changes in the ARP cache, e.g. `net.arp.node-1.set`.  If a given IP does not have an alias, the tool will report the IP with its dots (`.`) replaced by dashes (`-`), e.g. `net.arp.192-168-0-1.set`.

**`metron.endpoint`**

This option specifies the endpoint for an existing Metron agent. Everytime the tool detects a change in the ARP cache, it will send a corresponding state change metric to Metron. `net.arp.<key>.<action>` will be used as the format for the metric key (see above). The value for the metric is always `1` and its unit is `count`.

**`metron.origin`**

This option specifies the origin string used when sending metrics to the Metron agent.

Example Config File:

```
{
  "metron": {
    "endpoint": "localhost:3457",
    "origin": "node-1"
  },
  "frequency": "1s",
  "filters": [
    "192.168.0.1", 
    "192.168.0.2" 
  ],
  "aliases": {
    "192.168.0.1": "host-1", 
    "192.168.0.2": "host-2"
  }
}
```

## Licensing
Translator is licensed under the Apache License, Version 2.0. See
[LICENSE](https://github.com/st3v/arp-watch/blob/master/LICENSE) for the full
license text.
