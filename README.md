# GB-WEB
A web portal for GB mobile adapter connections via serial (and wifi). (not yet)

Dummy Servers based of [REON](https://github.com/REONTeam/dummy-servers) work. (pretty much 1 for 1 atm)


Both DNS server and POP run at the same time

// Current Progress
* Dummy DNS (WIP)
* Dummy Pop (WIP)


// Todo
* HTTP server (very basic features)


// How to use

Will run on Port 53 and 110 by default
```go
go run .\main.go

// To specify custom DNS port
go run .\main.go -dport 55
```