# Handles Server

A very simple server that verifies Bluesky (atproto) handles using the
[HTTPS well-known Method][atproto/resolution/well-known]; an alternative to
managing many DNS records.

## Quickstart

```console
curl -LO https://github.com/prompt/handles-server/releases/download/v1/handles-server-linux
chmod +x handles-server-linux
DID_PROVIDER="memory" MEMORY_DOMAINS="example.com" MEMORY_DIDS="alice.example.com@did:plc:001" \
./handles-server-linux
```

## Implementation

A `handle` is a hostname (e.g: `alice.example.com`) which the server may or may
not be able to provide a Decentralized ID for. A handle is made up of a `domain`
(e.g: `example.com`) and a `username` (e.g: `alice`). A provider
(`ProvidesDecentralizedIDs`) is responsible for getting a Decentralized ID from
a handle.

## Providers

- [x] Postgres
- [x] Memory
- [ ] Google Sheets
- [ ] Filesystem

## Configuration

| Environment Variable       | Description                                                | Example                                |
| -------------------------- | ---------------------------------------------------------- | -------------------------------------- |
| **`DID_PROVIDER`**         | **Required** Name of a supported provider                  | `postgres` `memory`                    |
| `REDIRECT_DID_TEMPLATE`    | URL template for redirects when a DID is found             | `https://bsky.app/profile/{did}`       |
| `REDIRECT_HANDLE_TEMPLATE` | URL template for redirects when a DID is not found         | `https://example.com/?handle={handle}` |
| `CHECK_DOMAIN_PARAMETER`   | Query parameter used by check domain endpoint (`/domainz`) | `handle` `hostname` `domain`           |

### `memory` provider

| Environment Variable | Description                                           | Example                         |
| -------------------- | ----------------------------------------------------- | ------------------------------- |
| **`MEMORY_DIDS`**    | **Required** Comma separated list of handle@did pairs | `alice.example.com@did:plc:001` |
| **`MEMORY_DOMAINS`** | **Required** Comma separate list of supported domains | `example.com,example.net`       |

### `postgres` provider

| Environment Variable     | Description                            | Example                                      |
| ------------------------ | -------------------------------------- | -------------------------------------------- |
| **`DATABASE_URL`**       | **Required** Postgres database URL     | `postgres://postgres@localhost:5432/handles` |
| `DATABASE_TABLE_DIDS`    | Table containing `handle` + `did` rows | `dids` `active_handles`                      |
| `DATABASE_TABLE_DOMAINS` | Table containing `domain` rows         | `domains` `active_domains`                   |

### URL templates

A string containing zero or more tokens which are replaced when rendering.

| Token               | Value                                           | Example(s)                 |
| ------------------- | ----------------------------------------------- | -------------------------- |
| `{handle}`          | Formatted handle from the request               | `alice.example.com`        |
| `{did}`             | Decentralized ID found for the request's handle | `did:plc:example001` ` `   |
| `{handle.domain}`   | Top level domain from the handle                | `example.com`              |
| `{handle.username}` | Username part of the handle                     | `alice` `bob`              |
| `{request.scheme}`  | Request's scheme                                | `https` `http`             |
| `{request.host}`    | Request's host                                  | `alice.example.com`        |
| `{request.path}`    | Path included in the request                    | `/hello-world` ` `         |
| `{request.query}`   | Query included in the request                   | `greeting=Hello+World` ` ` |

[atproto/resolution/well-known]: https://atproto.com/specs/handle#handle-resolution
