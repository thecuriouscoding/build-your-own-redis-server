# RESP

Redis uses a text based protocol known as the RESP (Redis Serialization Protocol). Major points are:

Simple Strings: Replies that are not errors.

    Format: +OK\r\n

Errors: Error messages from the server.

    Format: -ERR unknown command\r\n

Integers: Numeric replies.

    Format: :1\r\n

Bulk Strings: Strings of arbitrary length (including binary data).

    Format: $6\r\nfoobar\r\n (The number represents the length of the string)

Arrays: Used for multiple bulk strings (like command arguments).

    Format: *3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n (An array of three bulk strings)

## Explanation

Example Command: SET

`SET` command is parsed as:
Client Sends:

```
*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n
```

- `*3`: Array with 3 elements.
- `$3`: Bulk string of length 3.
- `SET`: command.
- `$3`: Bulk string of length 3.
- `foo`: Key.
- `$3`: Bulk string of length 3.
- `bar`: Value.