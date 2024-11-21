import pino, { Logger } from "pino";

export const logger: Logger = ["production", "test"].includes(
  process.env["NODE_ENV"] || "unknown",
)
  ? pino({ level: "warn" })
  : pino({
      transport: {
        target: "pino-pretty",
        options: { colorize: true },
      },
      level: "debug",
    });

export { type Logger };
