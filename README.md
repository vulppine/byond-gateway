byond-rest
==========

A REST server for BYOND's Dream Daemon.

Allows for transmission of an internal server state
of any kind, and for this state to be served over
HTTP, or alternatively, sent as a remote procedure call.

Building
--------

###Windows

You will require [TDM-GCC](https://jmeubank.github.io/tdm-gcc)
to compile the Dream Daemon library.

To compile the server, run `go build -o byond-rest.exe`.

To compile the library:
 - Install TDM-GCC from above.
 - Open up a MinGW command prompt.
 - Go to the directory that you've downloaded the library code to.
 - Run `go env -w GOARCH=386 CGO_ENABLED=1`.
 - Run `go build -o byond-socks.dll -tags netgo -buildmode=c-shared` in `library/`.
 
#### This isn't building on Windows.

Run `go mod init byond-socks` in your build directory. Make an issue otherwise
if you've already done this.

###Linux

To compile the server, run `go build -o byond-rest`.

To compile the library, run `CGO_ENABLED=1 GOARCH=386 go build -o byond-socks.so -tags netgo -buildmode=c-shared ./` in `library/`

Where these will be located depended is the implementer's choice.

Usage
-----

You can use this in a couple of ways.

**Dynamically called from Dream Daemon**

    /proc/start_rest()
        shell("[byond-rest executable]", "[port]", "[rpc_port]", "[rpc_call]")

**Started separate from Dream Daemon**

    BYOND_REST_PORT=[ port ]
    BYOND_REST_RPC_PORT=[ rpc_port ]
    BYOND_REST_RPC_CALL=[ rpc_call ]

    byond-rest

After starting, you can send the server messages like so:

    /proc/send_rest(S)
        call("[byond-socks library]", "SendAndClose")("127.0.0.1", "[port]", json_encode(S))

where S is a BYOND associative array.

You can access the server's current state via `http://localhost:3621/api/status`.

If you wish to perform the remote procedure call as indicated by
rpc_call, you will need to include a variable `status` that holds
an integer. Whenever this is changed, the server will automatically
call rpc_call to a locally hosted server at rpc_port.

JSON RPC is used to make RPCs.

Example code has been included in the `examples` folder. In there, you can find:
 - a Discord bot
 - an example implementation of byond-rest in DM

Using `examples/dm/byond-rest.dm` on a Windows system requires [byond-extools](https://github.com/MCHSL/extools).

Why is this split into parts?
-----------------------------

I'm not very confident that this could be *only* a library.
I already had some issues trying to use some Go stdlib functions while
making the library included (net.Listen causes a crash, while
net.ListenTCP does not). Perhaps someday (alternatively, you could
try to do this yourself if you really want)

License
-------

Flipp Syder, MIT License, 2021
