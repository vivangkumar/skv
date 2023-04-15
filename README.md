# skv

`skv` stands for "simple key value" (for lack of better naming)

It is a rudimentary in-memory key value store that only supports string -> string mappings.

## Why?

For the lolz really.

I wanted to implement something basic and gradually build it up to implement features of a modern kv store such as replication and failovers.

Both of which require some form of network protocol and network level code implementations.

For the most part, the intention is to make use of Go's pretty
extensive stdlib when implementing this.

## Protocol

`skv` runs over a very basic  string protocol over TCP.

The operations it supports are `set`, `get` and `delete`.

All commands are terminated with `\r\n` (carriage return, new line)

### set

```
set:<key>:<value>\r\n
```

### get

```
get:<key>\r\n
```

### delete

```
del:<key>\r\n
```

### ok & err

`skv` will send back an `ok` to indicate that an operation
has been carried out successfully.

Where the operation returns a value (`get`), the return value comes
after the `ok` - `ok:<value>`.

If there was an error, an `err` is sent back with an error string - `err:<error-string>`

## Replication

I envisage this as synchronous replication starting with a leader follower pattern.

Operations performed on the leader will reflect in followers.

## Failovers

This will employ a simple leader-follower failover strategy where a follower is promoted as the leader in case the leader is unreachable.

## Pipelining

Allows making multiple requests to the server reducing the RTT.
