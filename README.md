# GB-WEB
A web portal for GB mobile adapter connections via serial (and wifi). (not yet)

Dummy Servers based of [REON](https://github.com/REONTeam/dummy-servers) work. (pretty much 1 for 1 atm)

HTTP Folder structur, is copied from [REON](https://github.com/REONTeam/dummy-servers) as well.


Both DNS server and POP run at the same time. Should work on Windows, Linux and MacOS.

// Current Progress (These should be working)
* Dummy DNS (WIP)
* Dummy Pop (WIP)
* HTTP server (very basic features)

// Todo
* Maybe add more things


// How to use

Will run on Port 53 and 110 by default
```shell
go run main.go

##To specify custom DNS port
go run main.go -dport 55
```