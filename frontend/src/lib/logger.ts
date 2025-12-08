import { createLogger, format, transports } from "winston";
import config from "@app/config.json";

const isDev = config.environment === "development";

const defaultMeta = { service: "frontend" };

// Custom colors for each log level
const customColors = {
  error: "red",
  warn: "yellow",
  info: "cyan",
  debug: "magenta",
};

import winston from "winston";
winston.addColors(customColors);

const devFormat = format.combine(
  format.timestamp({ format: "HH:mm:ss" }),
  format.errors({ stack: true }),
  format.printf(({ timestamp, level, message, service, stack, ...meta }) => {
    const magenta = "\x1b[35m";
    const reset = "\x1b[0m";
    const gray = "\x1b[90m";

    let coloredLevel = level;
    if (level.includes("error")) coloredLevel = `\x1b[31m${level}${reset}`;
    else if (level.includes("warn")) coloredLevel = `\x1b[33m${level}${reset}`;
    else if (level.includes("info")) coloredLevel = `\x1b[36m${level}${reset}`;
    else if (level.includes("debug")) coloredLevel = `\x1b[35m${level}${reset}`;

    const serviceName = `${magenta}[${service}]${reset}`;

    const metaKeys = Object.keys(meta).filter(
      (key) => !["service", "timestamp", "level", "message"].includes(key)
    );
    const metaStr = metaKeys.length
      ? ` ${gray}${JSON.stringify(
          metaKeys.reduce((obj, key) => ({ ...obj, [key]: meta[key] }), {})
        )}${reset}`
      : "";

    const logLine = `${gray}${timestamp}${reset} ${coloredLevel} ${serviceName} ${message}${metaStr}`;

    return stack ? `${logLine}\n${stack}` : logLine;
  })
);

const prodFormat = format.combine(
  format.timestamp(),
  format.errors({ stack: true }),
  format.json()
);

const logger = createLogger({
  level: isDev ? "debug" : "info",
  defaultMeta,
  format: isDev ? devFormat : prodFormat,
  transports: [new transports.Console()],
});

export default logger;

export const info = (msg: string, meta = {}) => logger.info(msg, meta);
export const warn = (msg: string, meta = {}) => logger.warn(msg, meta);
export const error = (msg: string, meta = {}) => logger.error(msg, meta);
export const debug = (msg: string, meta = {}) => logger.debug(msg, meta);
