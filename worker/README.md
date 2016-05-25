# worker/heapy.go

All goes the same as generic `item/heapy.go`, except for `update()`.


## Run

```sh

-r int
      number of requests (default 100)
-w int
      number of workers (default 5)


$ go run heapy.go -w=10 -r=1000
```
