# lookupblaster

Simple multi-threaded utility to lookup the ip addresses in a host range

## How to build

`go build` 

## Usage example

The following example shows how to look up all hostnames in the IP range 194.71.11.173/20, using 512 concurrent threads:
`./lookupblaster 194.71.11.173/20 512` 

