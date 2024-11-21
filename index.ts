import { instantiateProviderFromStringConfiguration } from "./providers";
import { expose } from "./server";
import { logger } from "./logging";

const provider = instantiateProviderFromStringConfiguration(
  process.env.HANDLES_PROVIDER,
);

const server = expose(provider);

if (process.env.VERCEL !== "1") {
  const port = process.env.PORT || 3000;
  server.listen(port, () => logger.info(`Listening on ${port}`));
}

export default server;
