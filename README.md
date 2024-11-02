# About
Repository to build your own redis-server. It contains various handlers, persistence modes like `snapshot`, `aof`, `inmemory`.

# How to run
You can run ```go run .``` for inmemory persistence or ```go run . -"persistence"=snapshot -"snapshot-interval"=10s``` for snapshot persistence or ```go run . -"persistence"=aof``` for aof persistence mode. You can then connect to server via redis-cli and try performing various commands.

# Command Usage
1. SET - SET enables to push a key-value pair of string type onto the server.<br>

    ```SET [key] [value]```

2. GET - GET enables to get the value against a key from the server.<br>

    ```GET [key]```

3. EXPIRE - EXPIRE enables putting an expiry/TTL to the keys stored on server.<br>

    ```EXPIRE [key]```

4. INCR - INCR enables incrementing number value keys by 1.<br>

    ```INCR [key]```

5. DECR - INCR enables decrementing number value keys by 1.<br>

    ```DECR [key]```

6. TTL - TTL enables getting the left time to expire of a key.<br>

    ```TTL [key]```

7. LPUSH - LPUSH enables to insert elements on the left end of the list.<br>

    ```LPUSH [key] [elements...]```

8. RPUSH - RPUSH enables to insert elements on the right end of the list.

    ```RPUSH [key] [elements...]```

9. LRANGE - LRANGE enables to view the inserted elements of list.

    ```LRANGE [key] [start] [end]```

10. LPOP - LPOP will pop passed n number of elements from the left end of list

    ```LPOP [key] [number]```

11. RPOP - RPOP will pop passed n number of elements from the right end of list.

    ```RPOP [key] [number]```

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