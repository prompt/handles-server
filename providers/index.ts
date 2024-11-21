import { FindsDecentralizedIDOfHandle } from "./../handles";
import { logger } from "./../logging";
import { z } from "zod";

import { HandleMap } from "./map";
import { PostgresHandles } from "./pg";

// A list of keys of available providers
const providerKeys = z.enum(["map", "pg"]);

type ProviderKey = z.infer<typeof providerKeys>;

// Providers and their instantiators
export const providers: Record<
  ProviderKey,
  (config: string) => FindsDecentralizedIDOfHandle
> = {
  map: (config: string) => HandleMap.fromEnvironmentVariable(config),
  pg: (config: string) =>
    PostgresHandles.fromEnvironmentVariable(process.env[config] || config),
};

export function instantiateProviderFromStringConfiguration(
  configFromEnvironment?: string,
): FindsDecentralizedIDOfHandle {
  const [providerKey, configuration] = z
    .string({
      message:
        "A configuration is required in the Environment Variable HANDLES_PROVIDER.",
    })
    .transform((value) => {
      const [provider, ...config] = value.split(":");
      return [provider, config.join(":")];
    })
    .pipe(
      z.tuple([providerKeys, z.string()], {
        message: "HANDLES_PROVIDER must be a provider:configuration tuple.",
      }),
    )
    .parse(configFromEnvironment);

  logger.info(`Resolved configuration to provider '${providerKey}'`);

  const provider = providers[providerKey](configuration);

  logger.info(`Instantiated '${providerKey}' provider`);

  return provider;
}
