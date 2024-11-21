import { HandleMap } from "./map";
import { ZodError } from "zod";

describe("HandleMap", () => {
  test("Provides did for handle", async () => {
    const handles = new HandleMap(
      new Map([
        ["alice.example.com", "did:plc:example1"],
        ["bob.example.com", "did:plc:example2"],
      ]),
    );

    expect(await handles.findDecentralizedIDofHandle("alice.example.com")).toBe(
      "did:plc:example1",
    );
    expect(await handles.findDecentralizedIDofHandle("bob.example.com")).toBe(
      "did:plc:example2",
    );
  });

  test("Provides null for no handle", async () => {
    const handles = new HandleMap(new Map());

    expect(await handles.findDecentralizedIDofHandle("alice.example.com")).toBe(
      null,
    );
  });

  test("Creates map from configuration value", async () => {
    const handles = HandleMap.fromEnvironmentVariable(
      "alice.example.com->did:plc:example1,bob.example.com->did:plc:example2",
    );

    expect(await handles.findDecentralizedIDofHandle("alice.example.com")).toBe(
      "did:plc:example1",
    );
    expect(await handles.findDecentralizedIDofHandle("bob.example.com")).toBe(
      "did:plc:example2",
    );
  });

  test("Throws error when given incorrect configuration value", async () => {
    expect(() =>
      HandleMap.fromEnvironmentVariable(
        "alice.example.com@did:plc:example1,bob.example.com->did:plc:example2",
      ),
    ).toThrow(ZodError);
  });
});
