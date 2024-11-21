import { FindsDecentralizedIDOfHandle } from "./../handles";
import { logger } from "./../logging";
import { Pool } from "pg";

export type QueriesDatabaseForHandle = {
  query(
    queryText: string,
    values: string[],
  ): Promise<{ rows: Array<{ handle: string; did: string }> }>;
};

export class PostgresHandles implements FindsDecentralizedIDOfHandle {
  constructor(private readonly database: QueriesDatabaseForHandle) {}

  async findDecentralizedIDofHandle(handle: string) {
    const { rows } = await this.database.query(
      "SELECT handle, did FROM handles WHERE handle = $1",
      [handle.toLowerCase()],
    );

    logger.debug({ handle, rows }, "Successfully queried database.");

    return rows.length === 1 ? rows[0].did : null;
  }

  static fromEnvironmentVariable(config: string): PostgresHandles {
    const client = new Pool({ connectionString: config });

    client.connect();

    return new PostgresHandles(client);
  }
}
