import { FindsDecentralizedIDOfHandle } from "@/handles";
import { z } from "zod";
import { logger } from "@/logging";

export class HandleMap implements FindsDecentralizedIDOfHandle {
  constructor(private readonly identities: Map<string, string>) {}

  async findDecentralizedIDofHandle(handle: string) {
    return this.identities.get(handle) || null;
  }

  static fromEnvironmentVariable(config: string): HandleMap {
    const handles = z
      .string()
      .transform((value) => value.split(","))
      .pipe(
        z.array(
          z
            .string()
            .transform((value) => value.split("->"))
            .pipe(
              z
                .array(z.string())
                .length(
                  2,
                  "Each list value must be in the format handle->did.",
                ),
            )
            .pipe(
              z.tuple([
                z
                  .string()
                  .regex(
                    /^([A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])(\.|$)){2,10}/,
                    {
                      message: "The handle-> must be a valid domain name.",
                    },
                  ),
                z
                  .string()
                  .regex(/^did:[a-z]+:[a-zA-Z0-9._:%-]*[a-zA-Z0-9._-]$/, {
                    message:
                      "The ->did must be valid according to the atproto spec.",
                  }),
              ]),
            ),
          {
            message:
              "A comma-separated list of handle->did values is required.",
          },
        ),
      )
      .parse(config);

    logger.child({ handles }).debug("Successfully parsed a list of handles.");

    return new HandleMap(new Map(handles));
  }
}
