import type { Preset } from "@/modules/presets/preset-types";
import { useDmxWebSocketContext } from "@/providers/ws-provider";
import { type FC } from "react";
import { PerformanceButton } from "../components/performance-button";

interface PresetGridProps {
  presets: Preset[];
}

export const PresetGrid: FC<PresetGridProps> = ({ presets }) => {
  const { dmxState, presetApplied, applyPreset } = useDmxWebSocketContext();

  const activePresetId = dmxState?.active_preset_id ?? presetApplied?.preset_id;

  return (
    <section>
      <h2 className="text-xl font-semibold space-y-2">Presets</h2>
      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4">
        {presets.map((preset) => (
          <PerformanceButton
            key={preset.id}
            id={preset.id}
            onClick={() => applyPreset(preset.id)}
            isActive={activePresetId === preset.id}
            name={preset.name}
            description={preset.description ?? ""}
          />
        ))}
      </div>
    </section>
  );
};
