import { WsStatus } from "@/components/status/ws-status";
import { Button } from "@/components/ui/button";
import { useDmxWebSocketContext } from "@/providers/ws-provider";
import { PresetGrid } from "../components/preset-grid";
import { ShowGrid } from "../components/show-grid";

export const PerformanceView = () => {
  const { projectConfig, blackout } = useDmxWebSocketContext();

  if (!projectConfig) {
    return (
      <div className="flex items-center justify-center h-full">
        <span className="text-gray-500">Loading...</span>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-extrabold tracking-tight w-max">
          Performance
        </h1>
        <WsStatus />
      </div>

      <PresetGrid presets={projectConfig.presets} />

      <ShowGrid shows={projectConfig.shows} />

      <Button onClick={blackout} variant="destructive">
        Blackout
      </Button>
    </div>
  );
};
