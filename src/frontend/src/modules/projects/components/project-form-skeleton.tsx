import { Skeleton } from "@/components/ui/skeleton";

export const ProjectFormSkeleton = () => {
  return (
    <div className="container mx-auto space-y-6">
      <Skeleton className="h-8 w-23" />
      <div className="grid gap-2">
        <Skeleton className="h-3 w-23" />
        <Skeleton className="h-9 w-full" />
      </div>
      <div className="grid gap-2">
        <Skeleton className="h-3 w-23" />
        <Skeleton className="h-9 w-full" />
      </div>

      <Skeleton className="h-9 w-full" />
    </div>
  );
};
