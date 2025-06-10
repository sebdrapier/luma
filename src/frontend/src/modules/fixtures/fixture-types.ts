export interface Fixture {
  id: string;
  name: string;
  description: string;
  type: string;
  channels: FixtureChannel[];
}

export interface FixtureChannel {
  name: string;
  description: string;
  min: number;
  max: number;
  channel_address: number;
}