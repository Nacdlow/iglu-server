# iglü Server

## Description

iglü server is the web server and control system for the smart home. Its
purpose is to control home appliances and Internet-connected devices. It 
is self-contained, handling access-control among other things, and should work
without Internet connection.

## Requirements

The following packages are required.

- git
- Go (1.12+)
- GNU Make
- go-bindata

### Installing go-bindata

```sh
$ go get -u github.com/go-bindata/go-bindata/...
$ go install github.com/go-bindata/go-bindata/...
```

## Building

You can build by running:
```sh
$ make
```

And to run the web server:
```
$ make
$ ./nacdlow-server run [--port 443] [--dev]
```

## Configuration

A configuration file is created `./config.toml`. New fields will be
automatically added. You might need to set the Dark Sky API key, which you can
find on our Internal Wiki.

## Plugins

Plugins are added in the `plugins/` folder. Plugins may also be downloaded
directly from the marketplace.
