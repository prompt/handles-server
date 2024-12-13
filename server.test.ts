import { HandleMap } from "./providers/map";
import { handleRequest, HttpConfiguration, defaultFallbackUrl } from "./server";
import {
  createRequest,
  createResponse,
  RequestOptions,
  MockResponse,
} from "node-mocks-http";

const testHandleMap = new HandleMap(
  new Map([["registered.at.example.com", "did:plc:example1"]]),
);

// eslint-disable-next-line @typescript-eslint/no-explicit-any
type Response = MockResponse<any>;

async function performRequest(
  options: RequestOptions,
  configuration: HttpConfiguration = {
    fallbackUrl: defaultFallbackUrl,
  },
): Promise<Response> {
  const req = createRequest(options);
  const res = createResponse();

  await handleRequest(req, res, testHandleMap, configuration);

  return res;
}

describe("Handles Server", () => {
  test("Responds to healthcheck request", async () => {
    const res = await performRequest({ method: "GET", url: "/healthz" });

    expect(res._getData()).toBe("OK");
    expect(res._getStatusCode()).toBe(200);
  });

  test("Request to verify registered handle returns did.", async () => {
    const res = await performRequest({
      method: "GET",
      url: "/.well-known/atproto-did",
      headers: { Host: "registered.at.example.com" },
    });

    expect(res._getStatusCode()).toBe(200);
    expect(res._getData()).toBe("did:plc:example1");
  });

  test("Request to verify unregistered handle returns error", async () => {
    const res = await performRequest({
      method: "GET",
      url: "/.well-known/atproto-did",
      headers: { Host: "not-registered.at.example.com" },
    });

    expect(res._getStatusCode()).toBe(404);
  });

  test("Request to visit registered handle redirects to profile", async () => {
    const res = await performRequest({
      method: "GET",
      url: "/",
      headers: { Host: "registered.at.example.com" },
    });

    expect(res._getStatusCode()).toBe(307);
    expect(res._getHeaders().location).toBe(
      "https://bsky.app/profile/registered.at.example.com",
    );
  });

  test("Request to visit unregistered handle without fallback URL redirects to domain", async () => {
    const res = await performRequest({
      method: "GET",
      url: "/",
      headers: { Host: "not-registered.at.example.com" },
    });

    expect(res._getStatusCode()).toBe(307);
    expect(res._getHeaders().location).toBe(
      "https://at.example.com?utm_source=handles-server&utm_medium=http&utm_campaign=redirect&utm_term=not-registered.at.example.com",
    );
  });

  test("Request to visit unregistered handle redirects to fallback URL", async () => {
    const res = await performRequest(
      {
        method: "GET",
        url: "/",
        headers: { Host: "not-registered.at.example.com" },
      },
      {
        fallbackUrl: "https://example.com/domains/{domain}?handle={handle}",
      },
    );

    expect(res._getStatusCode()).toBe(307);
    expect(res._getHeaders().location).toBe(
      "https://example.com/domains/at.example.com?handle=not-registered.at.example.com",
    );
  });

  test("Replaces multiple instances of token in fallback URL", async () => {
    const res = await performRequest(
      {
        method: "GET",
        url: "/",
        headers: { Host: "not-registered.at.example.com" },
      },
      {
        fallbackUrl: "https://example.com/{handle}?handle={handle}",
      },
    );

    expect(res._getStatusCode()).toBe(307);
    expect(res._getHeaders().location).toBe(
      "https://example.com/not-registered.at.example.com?handle=not-registered.at.example.com",
    );
  });
});
