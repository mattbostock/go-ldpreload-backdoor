# LD_PRELOAD libc hooking using Go

This is an experiment to use Go in a shared library to wrap a libc function and
start a TCP server (a 'backdoor') allowing arbitrary commands to be
run from a client such as telnet or netcat.

This is a toy intended for educational purposes to demonstrate some of
Go's capabilities.

Works on Linux only and requires Go version 1.5 or above in order to build the
shared library.

## Rationale

In writing this, I have four aims:

- to try out [Go's new build modes][], which allow Go to be compiled to a
  shared library that can be called from C

- to experiment with `LD_PRELOAD` exploits

- to experiment with calling C from Go

- to learn some C ;)

[Go's new build modes]: https://docs.google.com/document/d/1nr-TQHw_er6GOQRsF6T43GGhFDelrAP0NqSS_00RgZQ

## Usage

As this is an experiment, the backdoor will only listen on localhost.

    GO15VENDOREXPERIMENT=1 go build -buildmode=c-shared -o backdoor.so main.go
    LD_PRELOAD=./backdoor.so top

In a separate console, while `top` is running:

    nc localhost 4444
    [...type your commands here...]

## Limitations

- Only works on Linux
- Only works with binaries that call libc's `strrchr` function. I'd ideally
  like to hook `__libc_start_main` instead. The binaries I tested with are `ps`
  and `top` as provided by Ubuntu Trusty LTS.
