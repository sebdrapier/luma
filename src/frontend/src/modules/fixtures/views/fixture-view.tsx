import { ErrorAlert } from "@/components/error-alert";
import { InfoAlert } from "@/components/info-alert";
import { Skeleton } from "@/components/ui/skeleton";
import { NewFixtureSheet } from "../components/new-fixture-sheet";
import { useFixtures } from "../fixture-api";
import { FixtureTable } from "../components/fixture-table";

export const FixtureView = () => {
  const { data: fixtures, isLoading, isError } = useFixtures();

  if (isLoading) {
    return <Skeleton className="h-9 w-full rounded-md" />;
  }

  if (isError) {
    return (
      <ErrorAlert
        title="Error loading fixtures"
        message="Failed to load DMX fixtures. Please check your backend connection or logs for details."
      />
    );
  }

  if (!fixtures || fixtures.length === 0) {
    return (
      <InfoAlert
        title="No fixtures found"
        message="No DMX fixtures are currently available. Add or configure fixtures to see them here."
      />
    );
  }

  return (
    <div className="container mx-auto space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-extrabold tracking-tight w-max">
          Fixtures
        </h1>
        <NewFixtureSheet />
      </div>

      <FixtureTable fixtures={fixtures}/>
    </div>
  );
};
