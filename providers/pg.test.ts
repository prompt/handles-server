import { PostgresHandles } from "./pg";

const examplePostgresHandler = new PostgresHandles(
  {
    query: (q: string, v: string[]) =>
      Promise.resolve({
        rows:
          v.join() == "alice.at.example.com"
            ? [{ handle: "alice.at.example.com", did: "did:plc:example1" }]
            : [],
      }),
  },
  "handles",
);

describe("PostgresHandles", () => {
  test("Finds matching handle", async () => {
    expect(
      await examplePostgresHandler.findDecentralizedIDofHandle(
        "alice.at.example.com",
      ),
    ).toBe("did:plc:example1");
  });

  test("Finds matching handle with different casing", async () => {
    expect(
      await examplePostgresHandler.findDecentralizedIDofHandle(
        "ALICE.aT.example.com",
      ),
    ).toBe("did:plc:example1");
  });

  test("Does not provide did for handle not in database", async () => {
    expect(
      await examplePostgresHandler.findDecentralizedIDofHandle(
        "bob.at.example.com",
      ),
    ).toBe(null);
  });
});
