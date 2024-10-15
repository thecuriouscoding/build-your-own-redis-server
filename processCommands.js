const parseCommand = (data) => {
    const lines = data
        .toString()
        .split("\r\n")
        .filter((line) => line != "");
    const command = lines[2].toUpperCase();
    const args = lines.slice(4).filter((_, index) => index % 2 === 0);
    return { command, args };
};

module.exports = {
    parseCommand
}