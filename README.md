# PART 4
This branch pertains to the addition of "Snapshot" persistence mode which will enable the application to create a snapshot at every configured seconds. Persistence mode by default is in-memory and flag `persistence` with value `snapshot` or `inmemory` can be utilised to enable the applicable functionality. Flag `snapshot-interval` can be utilised to configure snapshot interval.

# How to run
You can run ```go run . -"persistence"=snapshot -"snapshot-interval"=10s``` and then connect to server via redis-cli and try performing SET, GET, DEL and EXPIRE commands.

# Command Usage
1. SET - SET enables to push a key-value pair of string type onto the server.<br>

    ```SET [key] [value]```

2. GET - GET enables to get the value against a key from the server.<br>

    ```GET [key]```

3. LPUSH - LPUSH enables to insert elements on the left end of the list.<br>

    ```<𝘓𝘗𝘜𝘚𝘏 [𝘬𝘦𝘺] [𝘦𝘭𝘦𝘮𝘦𝘯𝘵𝘴...]>```

4. RPUSH - RPUSH enables to insert elements on the right end of the list.

    ```<𝘙𝘗𝘜𝘚𝘏 [𝘬𝘦𝘺] [𝘦𝘭𝘦𝘮𝘦𝘯𝘵𝘴...]>```

5. LRANGE - LRANGE enables to view the inserted elements of list.

    ```<𝘓𝘙A𝘕𝘎𝘌 [𝘬𝘦𝘺] [𝘴𝘵𝘢𝘳𝘵] [𝘦𝘯𝘥]>```

6. LPOP - LPOP will pop passed n number of elements from the left end of list

    ```<𝘓𝘗𝘖𝘗 [𝘬𝘦𝘺] [𝘯𝘶𝘮𝘣𝘦𝘳]>```

7. RPOP - RPOP will pop passed n number of elements from the right end of list.

    ```<𝘙𝘗𝘖𝘗 [𝘬𝘦𝘺] [𝘯𝘶𝘮𝘣𝘦𝘳]>```

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