const enabledLogs = (process.env.DEBUG || "*")
  .split(",")
  .map((e) => e.trim());

const logger = (logFrom) => {
  const log = (mode, message) => {
    const logMessage = `${new Date().toISOString()} ${mode} [${logFrom}]: ${message}`;

    if (mode === "error") {
      console.error(logMessage);
      return;
    }

    if (enabledLogs.includes("*") || enabledLogs.includes(logFrom)) {
      console[mode](logMessage);
    }
  };

  return {
    info: (message) => log("log", message),
    error: (message) => log("error", message),
    warn: (message) => log("warn", message),
    debug: (message) => log("debug", message),
  };
};

module.exports = logger;
