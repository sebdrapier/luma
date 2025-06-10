import { ErrorAlert } from "@/components/error-alert";
import { InfoAlert } from "@/components/info-alert";
import { Skeleton } from "@/components/ui/skeleton";
import { ShowTable } from "../components/show-table";
import { useShows } from "../show-api";

export const ShowView = () => {
  const { data: shows, isLoading, isError } = useShows();

  if (isLoading) {
    return <Skeleton className="h-9 w-full rounded-md" />;
  }

  if (isError) {
    return (
      <ErrorAlert
        title="Error loading shows"
        message="Failed to load shows. Please check your backend connection or logs for details."
      />
    );
  }

  if (!shows) {
    return (
      <InfoAlert
        title="No shows found"
        message="No shows are currently available. Add or configure shows to see them here."
      />
    );
  }

  return (
    <div className="container mx-auto space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-extrabold tracking-tight w-max">Shows</h1>
      </div>

      <ShowTable shows={shows} />
    </div>
  );
};
