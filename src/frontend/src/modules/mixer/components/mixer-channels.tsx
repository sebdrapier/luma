import { Label } from "@/components/ui/label";
import { Slider } from "@/components/ui/slider";
import { useDmxWebSocketContext } from "@/providers/ws-provider";

export const MixersChannels = () => {
  const { projectConfig, dmxState, updateChannel } = useDmxWebSocketContext();

  if (!projectConfig) {
    return <div className="p-4">Chargement des fixtures…</div>;
  }

  const getChannelValue = (address: number) =>
    dmxState?.channels?.find((c) => c.address === address)?.value ?? 0;
  return (
    <div className="p-4 grid gap-6">
      {projectConfig.fixtures.map((fixture) => (
        <div key={fixture.id} className="space-y-2">
          <h2 className="font-semibold">{fixture.name}</h2>
          <div className="flex gap-6 overflow-x-auto">
            {fixture.channels.map((channel) => {
              const currentValue = getChannelValue(channel.channel_address);

              return (
                <div
                  key={channel.channel_address}
                  className="flex flex-col items-center gap-2 w-16"
                >
                  <p className="font-mono">{currentValue}</p>
                  <Slider
                    value={[currentValue]}
                    min={channel.min}
                    max={channel.max}
                    orientation="vertical"
                    onValueChange={([value]) =>
                      updateChannel(channel.channel_address, value)
                    }
                    className="[&>:last-child>span]:h-6 [&>:last-child>span]:w-4 [&>:last-child>span]:rounded"
                    aria-label={channel.name}
                  />
                  <Label className="text-muted-foreground text-xs text-center max-w-[4rem] truncate">
                    {channel.name}
                  </Label>
                </div>
              );
            })}
          </div>
        </div>
      ))}
    </div>
  );
};
