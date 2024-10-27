# PART 3
This branch pertains to the addition of handling of List Related Commands i.e. LPUSH, RPUSH, LRANGE, LPOP and RPOP. In-memory store has been used.

# How to run
1. You can run ```go run .``` and then connect to server via redis-cli and try performing SET, GET, DEL and EXPIRE commands.

# Command Usage
1. LPUSH - LPUSH enables to insert elements on the left end of the list.<br>

    ```<𝘓𝘗𝘜𝘚𝘏 [𝘬𝘦𝘺] [𝘦𝘭𝘦𝘮𝘦𝘯𝘵𝘴...]>```

2. RPUSH - RPUSH enables to insert elements on the right end of the list.

    ```<𝘙𝘗𝘜𝘚𝘏 [𝘬𝘦𝘺] [𝘦𝘭𝘦𝘮𝘦𝘯𝘵𝘴...]>```

3. LRANGE - LRANGE enables to view the inserted elements of list.

    ```<𝘓𝘙A𝘕𝘎𝘌 [𝘬𝘦𝘺] [𝘴𝘵𝘢𝘳𝘵] [𝘦𝘯𝘥]>```

4. LPOP - LPOP will pop passed n number of elements from the left end of list

    ```<𝘓𝘗𝘖𝘗 [𝘬𝘦𝘺] [𝘯𝘶𝘮𝘣𝘦𝘳]>```

5. RPOP - RPOP will pop passed n number of elements from the right end of list.

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