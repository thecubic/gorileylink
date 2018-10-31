# gorileylink
RileyLink BLE driver in Go

## example

```
$ go get github.com/thecubic/gorileylink/cmd/grl-battery && sudo ./grl-battery --rileylink cc:cc:cc:cc:cc:cc
2018/10/30 23:33:28 connecting to cc:cc:cc:cc:cc:cc
2018/10/30 23:33:29 found a RileyLink: cc:cc:cc:cc:cc:cc
2018/10/30 23:33:29 connected to cc:cc:cc:cc:cc:cc
88:6b:0f:8e:ee:b0: RileyLink ble_rfspy 2.0 81%
2018/10/30 23:33:29 disconnected from cc:cc:cc:cc:cc:cc
```
