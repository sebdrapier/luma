import { Button } from "@/components/ui/button";
import { PresetSelect } from "@/modules/presets/components/preset-select";
import type { Preset } from "@/modules/presets/preset-types";
import { useDmxWebSocketContext } from "@/providers/ws-provider";
import { RotateCcw } from "lucide-react";
import { useState } from "react";

import { NewPresetSheet } from "@/modules/presets/components/new-preset-sheet";

export const MixerPreset = () => {
  const [selectedPreset, setSelectedPreset] = useState<Preset | null>(null);
  const { applyPreset, blackout } = useDmxWebSocketContext();

  const handleSelectPreset = (preset: Preset) => {
    setSelectedPreset(preset);
    applyPreset(preset.id);
  };

  const handleReset = () => {
    blackout();
    setSelectedPreset(null);
  };

  return (
    <div className="w-full flex gap-2">
      <PresetSelect
        onSelect={handleSelectPreset}
        value={selectedPreset?.id ?? ""}
      />

      <Button onClick={handleReset} variant="outline" size="icon">
        <RotateCcw />
      </Button>

      <NewPresetSheet preset={selectedPreset} />
    </div>
  );
};
