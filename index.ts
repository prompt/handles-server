import { expose } from "./server";
import { HandleMap } from "./handles";

const server = expose(
  new HandleMap({
    "alice.example.com": "did:plc:example",
  }),
);

if (process.env.VERCEL !== "1") {
  server.listen(process.env.PORT || 3000, () => console.log("Listening"));
}

export default server;
