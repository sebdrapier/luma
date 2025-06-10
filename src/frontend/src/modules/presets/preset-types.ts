export interface Preset {
  id: string;
  name: string;
  description: string;
  channels: ChannelValue[];
}

export interface ChannelValue {
  dmx_address: number;
  value: number;
}