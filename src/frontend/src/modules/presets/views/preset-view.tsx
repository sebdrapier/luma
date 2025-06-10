import { ErrorAlert } from "@/components/error-alert";
import { InfoAlert } from "@/components/info-alert";
import { Skeleton } from "@/components/ui/skeleton";
import { usePresets } from "../preset-api";
import { PresetTable } from "../components/preset-table";

export const PresetView = () => {
  const { data: presets, isLoading, isError } = usePresets();

  if (isLoading) {
    return <Skeleton className="h-9 w-full rounded-md" />;
  }

  if (isError) {
    return (
      <ErrorAlert
        title="Error loading presets"
        message="Failed to load presets. Please check your backend connection or logs for details."
      />
    );
  }

  if (!presets || presets.length === 0) {
    return (
      <InfoAlert
        title="No presets found"
        message="No presets are currently available. Add or configure presets to see them here."
      />
    );
  }

  return (
    <div className="container mx-auto space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-extrabold tracking-tight w-max">
          Presets
        </h1>
      </div>

      <PresetTable presets={presets} />
    </div>
  );
};
