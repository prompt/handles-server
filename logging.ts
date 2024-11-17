import pino, { Logger } from "pino";

export const logger: Logger =
  process.env["NODE_ENV"] === "production"
    ? pino({ level: "warn" })
    : pino({
        transport: {
          target: "pino-pretty",
          options: { colorize: true },
        },
        level: "debug",
      });

export { type Logger };
