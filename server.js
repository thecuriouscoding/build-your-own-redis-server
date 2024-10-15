const net = require("net");
const { parseCommand } = require("./processCommands");

const logger = require("./logger")("server");

const server = net.createServer();
const port = 6379;
const host = "127.0.0.1";

server.on("connection", (socket) => {
    socket.on("data", (data) => {
        let res;
        try {
            const { command, args } = parseCommand(data);
            logger.info(`Command is ${command}\nArguments are ${args}`)
            res = "+OK\r\n"
        } catch (err) {
            logger.error(err);
            res = "-ERR unknown command\r\n";
        }
        socket.write(res);
    });

    socket.on("end", () => {
        logger.info("Client disconnected");
    });
});

server.listen(port, host, () => {
    logger.info(`Server running at ${host}:${port}`);
});
