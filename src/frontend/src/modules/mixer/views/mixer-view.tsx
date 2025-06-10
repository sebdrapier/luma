import { WsStatus } from "@/components/status/ws-status";
import { MixersChannels } from "../components/mixer-channels";
import { MixerPreset } from "../components/mixer-preset";

export const MixerView = () => {
  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-extrabold tracking-tight w-max">Mixer</h1>
        <WsStatus />
      </div>

      <MixerPreset />

      <MixersChannels />
    </div>
  );
};
