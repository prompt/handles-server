import pino, { Logger } from "pino";

const environment = process.env["NODE_ENV"] || "production";

export const logger: Logger = ["production", "test"].includes(environment)
  ? pino({ level: "warn" })
  : pino({
      transport: {
        target: "pino-pretty",
        options: { colorize: true },
      },
      level: "debug",
    });

export { type Logger };
