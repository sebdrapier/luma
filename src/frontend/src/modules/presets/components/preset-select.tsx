import { ErrorAlert } from "@/components/error-alert";
import { InfoAlert } from "@/components/info-alert";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Skeleton } from "@/components/ui/skeleton";
import type { FC } from "react";
import { usePresets } from "../preset-api";
import type { Preset } from "../preset-types";

export interface PresetSelectProps {
  onSelect: (preset: Preset) => void;
  value: string;
}

export const PresetSelect: FC<PresetSelectProps> = ({ onSelect, value }) => {
  const { data: presets, isLoading, isError } = usePresets();

  if (isLoading) return <Skeleton className="h-9 w-full rounded-md" />;

  if (isError)
    return (
      <ErrorAlert
        title="Error loading presets"
        message="Could not retrieve presets. Make sure the backend is running."
      />
    );

  if (!presets || presets.length === 0)
    return (
      <InfoAlert
        title="No presets found"
        message="Please connect a DMX device and try again."
      />
    );

  return (
    <Select
      value={value}
      onValueChange={(presetId: string) => {
        const preset = presets.find((p) => p.id === presetId);
        if (preset) onSelect(preset);
      }}
    >
      <SelectTrigger className="w-full">
        <SelectValue placeholder="Select preset" />
      </SelectTrigger>
      <SelectContent>
        {presets.map((preset) => (
          <SelectItem key={preset.id} value={preset.id}>
            {preset.name}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
};
