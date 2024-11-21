export interface FindsDecentralizedIDOfHandle {
  findDecentralizedIDofHandle(fqdn: string): Promise<string | null>;
}
