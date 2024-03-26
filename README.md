# skv

`skv` stands for "simple key value" (for lack of better naming)

It is a rudimentary in-memory key value store that only supports string -> string mappings.

## Why?

Just a quick project to demonstrate the basics of networking:
- Implementing a TCP/IP server
- Creating a TCP protocol
- Encoding & decoding protocol messages
- Graceful shutdowns

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

In cases where the value doesn't exist. For example in the case of getting a key that doesn't exist, a `null` value is returned.

If there was an error, an `err` is sent back with an error code and message - `err:<error-code><error-string>`

## Pipelining

Currently, pipelining is not supported.
