export interface FindsDecentralizedIDOfHandle {
  findDecentralizedIDofHandle(fqdn: string): Promise<string | null>;
}

export class HandleMap implements FindsDecentralizedIDOfHandle {
  constructor(private readonly identities: Record<string, string>) {}

  async findDecentralizedIDofHandle(handle: string) {
    return this.identities[handle] || null;
  }
}
