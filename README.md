# yaja (Yet Another JAbber implementation)

[![chat on our conference room](https://camo.githubusercontent.com/a839cc0a3d4dac7ec82237381b165dd144365b6d/68747470733a2f2f74696e7975726c2e636f6d2f6a6f696e7468656d7563)](https://conversations.im/j/yaja@conference.chat.sum7.eu)

## Features
- Messages XML Library (first version - PR are welcome)
	- Full RFC 6120 (XMPP - Core)
- Client Library (WIP)
	-	Stream: TLS Required
		- SASL-Auth (PLAIN, DIGEST-MD5)
	- Read & Decode (recv xml)
	- Send (send xml)
	- No OTR (never implemented -> nowadays OMEMO or PGP prefered)
	- No OMEMO support (not implemented by me: library only for bots and testing)
- Client (planned)
- Daemon
	- Tester (based on Client Library)
		- Bidirected Messaging
		- Check IPv4 & IPv6
		- TLS Version
	- Server (maybe broken, planned)
		- get certificate by lets encrypt
		- registration (for every possible TLS domain)

## Library

### Messages
all implementation of all comman (RFCs and XEPs) xml element

**Version**
- RFC 6120 (XMPP - Core)

### Client

**Planned**
- auto decoding of XML (with some auto answer  e.g. ping)
- SendRecv to get answer of a request
- Register

## Run

```
A small standalone command line round about jabber (e.g tester WIP: client & server)

Usage:
  yaja [command]

Available Commands:
  daemon      daemon of yaja
  help        Help about any command

Flags:
  -h, --help   help for yaja

Use "yaja [command] --help" for more information about a command.

```

### Daemons
```
daemon of yaja

Usage:
  yaja daemon [command]

Available Commands:
  server      runs xmpp server
  tester      runs xmpp tester server

Flags:
  -h, --help   help for daemon

Use "yaja daemon [command] --help" for more information about a command.
```
#### Tester
Website for displaying: [genofire/yaja-tester-viewer](https://github.com/genofire/yaja-tester-viewer/tree/master)

(dirty and based on [Freifunk Meshviewer](https://github.com/ffrgb/meshviewer/))

Demo: [tester.chat.sum7.eu](https://tester.chat.sum7.eu)

```
runs xmpp tester server

Usage:
  yaja daemon tester [flags]

Examples:
yaja daemon tester -c /etc/yaja.conf

Flags:
  -c, --config string   path to configuration file (default "yaja-tester.conf")
  -h, --help            help for tester

```
**Features**
- notification of disconnect by server (domain)
	- manage by bot command `admin (add|del) <JID>` and `admin list`)
- auto accept subscription by every user to every bot
- `ping` to `pong` by every user to every bot (for self check)

**Planned**
- improve chat bot implementation
- improve notification (add my self, not only other by admins)
- add new accounts/server
- other checks, maybe like [running](https://conversations.im/compliance/) - [source-code](https://github.com/iNPUTmice/ComplianceTester)
	- software and version of xmpp servers

**Inspiration by**

*Sorry i did not like Java on my server*
- ServerStatus of iNPUTmice: [running](https://status.conversations.im/) - [source-code](https://github.com/iNPUTmice/ServerStatus)

#### Server
```
runs xmpp server

Usage:
  yaja daemon server [flags]

Examples:
yaja daemon server -c /etc/yaja.conf

Flags:
  -c, --config string   path to configuration file (default "yaja-server.conf")
  -h, --help            help for server
```

## Inspiration by source-code structures (but rewritten)
- **server side:** [tam7t](https://github.com/tam7t/xmpp) a fork of [agl](https://github.com/agl)'s work
- **client side:** [mattn](https://github.com/mattn/go-xmpp) (original by russ cox)
