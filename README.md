# Handles Server

A very simple server that reponds to Bluesky handle verification requests for
the domain(s) that the server is exposed to, as an alternative to managing
claims via DNS. Configure with a provider, and let it run.

[![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https%3A%2F%2Fgithub.com%2Fshrink-inc%2Fhandles-server&env=HANDLES_PROVIDER)

## Providers

A provider is configured using the `HANDLES_PROVIDER` Environment Variable. A
`key` identifies the provider and the `value` configures how the provider will
behave (always wrap the value in quotes (i.e: `HANDLES_PROVIDER="value"`)).

```shell
HANDLES_PROVIDER="provider:configuration"
```

### `HandleMap`

`HandleMap` is an example provider, it serves the comma-separated handle->did
values it has been configured with.

```shell
$ HANDLES_PROVIDER="map:alice.at.example.com->did:plc:example1,bob.at.example.com->did:plc:example2" \
npm run dev

[00:00:00.000] INFO (0000): Resolved configuration to provider 'map'
[00:00:00.000] DEBUG (0000): Successfully parsed a list of handles.
    handles: [
      [
        "alice.at.example.com",
        "did:plc:example1"
      ],
      [
        "bob.at.example.com",
        "did:plc:example2"
      ]
    ]
[00:00:00.000] INFO (0000): Instantiated 'map'
[00:00:00.000] INFO (0000): Listening on 3000
```

Make a request passing in the target handle as the `Host`.

```shell
$ curl http://localhost:3000/.well-known/atproto-did --header "Host: alice.at.example.com"
did:plc:example1
```

### Postgres

...

### Google Sheets

...
