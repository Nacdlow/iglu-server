# Nacdlow Server

## Description

Nacdlow server is the web server and control system for the smart home. Its
purpose is to control home appliances and Internet-connected devices. It 
is self-contained, handling access-control among other things, and should work
without Internet connection.

## Cloning and running

```sh
$ git clone git@gitlab.com:group-nacdlow/nacdlow-server.git
$ cd nacdlow-server
```

And to run the web server
```
$ make
$ ./nacdlow-server run [--port 443]
```

## Configuration

A configuration file is created `./config.toml`. New fields will be
automatically added. You might need to set the Dark Sky API key, which you can
find on our Internal Wiki.
