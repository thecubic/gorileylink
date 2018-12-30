# gorileylink
RileyLink BLE driver and utilities in Go.

Currently only works on Linux, but should work with Mac OS X after minor platform abstraction

## Caveats

**NO WARRANTY EXPRESSED OR IMPLIED**

There are many, this list is not exhaustive.  

### `root` access required

One of the most annoying is that it requires extended capabilities to function and thus everything must be called as `root`, or with explicit granted capability privileges:

```
$ sudo setcap cap_net_admin,cap_net_raw+ep ./grl-demo
$ ./grl-demo SWEETBREAD-TWO
```

This appears to be an inherited caveat of the BT library (`github.com/currantlabs/ble`); likely related to the requirement that all things start as a "scan", vs. a direct (address) connect which by design does not need privs

### Existing BT connections drop during use

Existing connections will be terminated temporarily (e.g. your wireless mouse will stop working and your heavy metal music will blare out of your speakers instead of your BT headset)

### Temporary name confusion

For some reason, a `rename(A, B)` followed immediately by a `rename(B, A)` will not work, and it will not appear as `B` until you interact with it as `A` for something.  This may be due to value caching on some layer.

## Connecting

A RileyLink can be identified either by its current custom name (e.g. `SWEETBREAD-TWO`), or its BT address (e.g. `88:6B:0F:FF:FF:FF`)

## Utilities

`gorileylink` ships with some useful utilities.  To build these, it's recommended that you be in a `~/go/bin` directory (or wherever you like to store custom-built utilities, you do you), and then e.g. for `grl-demo`:

```
$ go get -v -u github.com/thecubic/gorileylink/cmd/grl-demo
```

This results in a `grl-demo` binary.  `go get` can just be `go build` after the module is fetched (and it has to be when developing on a local head)

This will either need to be called as `root` (via `sudo`), or can be called as a regular user provided granted capabilities:

```
$ sudo setcap cap_net_admin,cap_net_raw+ep ~/go/bin/grl-demo
```

The utilities accept identifying RileyLinks by name or by address.  Should there be two with the same name, the "first" matching one wins.

### `grl-rename`: Naming Utility

**WARNING: there is a weird here; the RileyLink may have both names temporarily**

You can rename a RileyLink to any accepted name from address or name input, or just report the current name by not providing a new name

```
$ go build github.com/thecubic/gorileylink/cmd/grl-rename
$ sudo ~/go/bin/grl-rename SWEETBREAD-TWO DaveyLink
INFO[0001] Renamed                                       customName=DaveyLink customNameBefore=SWEETBREAD-TWO rileylink=SWEETBREAD-TWO
$ # see WARNING above
$ sudo ~/go/bin/grl-rename DaveyLink
INFO[0003] Report Name                                   customName=DaveyLink rileylink=DaveyLink
```

### `grl-info`: Quick supervisor-chip (BLE) information

This reads the immediate information from the device; the BT address, the custom name, the BLE firmware version, and the battery level which is really not reliable at this time (here it says 80% but it's 100% and plugged-in)

```
$ go build github.com/thecubic/gorileylink/cmd/grl-info
$ sudo ~/go/bin/grl-info  SWEETBREAD-TWO
SWEETBREAD-TWO @ 88:6b:0f:ff:ff:ff: SWEETBREAD-TWO ble_rfspy 2.0 80%
```

### `grl-leds`: Set diagnostic LED mode

This sets the diagnostic LED mode, so that both blues will light up when it is talking to the the pump (or not).  Like `grl-rename`, providing no new value results in just fetching the existing value.

```
$ go build github.com/thecubic/gorileylink/cmd/grl-leds 
$ sudo ~/go/bin/grl-leds SWEETBREAD-TWO
INFO[0001] LED Mode                                      leds=off rileylink=SWEETBREAD-TWO
$ sudo ~/go/bin/grl-leds SWEETBREAD-TWO on
INFO[0002] LED Mode                                      leds=on rileylink=SWEETBREAD-TWO
```

### `grl-demo`: Demo application

Effectively a tour of supported features.

```
$ go build github.com/thecubic/gorileylink/cmd/grl-demo
$ sudo ~/go/bin/grl-demo -debug SWEETBREAD-TWO
DEBU[0005] connection succeeded                          nameoraddress=SWEETBREAD-TWO
DEBU[0005] bind as RileyLink succeeded                   nameoraddress=SWEETBREAD-TWO
DEBU[0005] BLE Subscription Successful                  
INFO[0005] Battery Level                                 batteryLevel=81
INFO[0005] Custom Name                                   customName=SWEETBREAD-TWO
INFO[0005] BLE Version                                   bleversion="ble_rfspy 2.0"
INFO[0005] State: OK                                    
INFO[0005] Radio Version                                 radioversion="subg_rfspy 2.2"
INFO[0005] Statistics                                    collected="2018-12-30 11:23:02.215360983 -0800 PST m=+5.675019184" crcfails=0 packetsrecv=0 packetsxmit=0 recvfifooverflows=0 recvoverflows=0 spisyncfails=0 uptime=43m3.886s
INFO[0005] starting LED dance                           
DEBU[0005] step                                          green=on
DEBU[0005] step + wait                                   blue=on
...
DEBU[0010] step                                          blue=off
DEBU[0010] step + wait                                   green=off

```
