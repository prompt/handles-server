import {
  createServer,
  IncomingMessage,
  ServerResponse,
  OutgoingHttpHeaders,
} from "http";
import { logger } from "./logging";
import { FindsDecentralizedIDOfHandle } from "./handles";

export const defaultFallbackUrl =
  "https://{domain}?utm_source=handles-server&utm_medium=http&utm_campaign=redirect&utm_term={handle}";

function respond(
  res: ServerResponse,
  body: string,
  status: number = 200,
  headers: OutgoingHttpHeaders = {},
): void {
  res.writeHead(status, { "Content-Type": "text/plain", ...headers });
  res.write(body);
  res.end();
}

export type HttpConfiguration = {
  fallbackUrl: string;
};

export async function handleRequest(
  req: IncomingMessage,
  res: ServerResponse,
  handles: FindsDecentralizedIDOfHandle,
  options: HttpConfiguration,
) {
  if (req.url === "/healthz") {
    return respond(res, "OK", 200);
  }

  const handle = req.headers.host || "";
  const domain = handle.split(".").slice(1).join(".").toLowerCase();

  let log = logger.child({ handle, domain, url: req.url });

  if (!handle || !domain) {
    respond(res, "A request must include a Handle in the `host` header.", 400);
    return log.debug("Rejected request without a host header.", {
      status: 400,
    });
  }

  const did = await handles.findDecentralizedIDofHandle(handle);

  log = log.child({ did });

  log.debug("Obtained `did` from provider.");

  if (did && req.url === "/.well-known/atproto-did") {
    respond(res, did);
    return log.debug("Successfully verified Handle.");
  }

  if (!did && req.url === "/.well-known/atproto-did") {
    respond(res, "Handle cannot be verified, did not found.", 404);
    return log.debug("Denied verification request due to missing did.");
  }

  if (did) {
    respond(res, "", 307, { location: `https://bsky.app/profile/${handle}` });
    return log.debug("Redirected to Bluesky profile.");
  }

  let noProfileDestination = options.fallbackUrl;

  Object.entries({ domain, handle }).forEach(
    ([token, value]) =>
      (noProfileDestination = noProfileDestination.replaceAll(
        `{${token}}`,
        value,
      )),
  );

  respond(res, "", 307, { location: noProfileDestination });
  return log.debug(
    { noProfileDestination },
    "Redirected to no profile destination.",
  );
}

export function expose(handles: FindsDecentralizedIDOfHandle) {
  return createServer((req, res) =>
    handleRequest(req, res, handles, {
      fallbackUrl: process.env.HANDLES_FALLBACK_URL || defaultFallbackUrl,
    }),
  );
}
