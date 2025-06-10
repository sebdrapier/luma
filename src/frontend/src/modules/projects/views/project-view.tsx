import { ErrorAlert } from "@/components/error-alert";
import { InfoAlert } from "@/components/info-alert";
import { ProjectForm } from "../components/project-form";
import { ProjectFormSkeleton } from "../components/project-form-skeleton";
import { useProject } from "../project-api";

export const ProjectView = () => {
  const { data: project, isLoading, isError } = useProject();

  if (isLoading) return <ProjectFormSkeleton />;

  if (isError) {
    return (
      <div className="container mx-auto space-y-6">
        <ErrorAlert
          title="Error loading project"
          message="An error occurred while fetching the project. Please try again
            later."
        />
      </div>
    );
  }

  if (!project) {
    return (
      <div className="container mx-auto space-y-6">
        <InfoAlert
          title="No project available"
          message="No project found. Create one to get started."
        />
      </div>
    );
  }

  return (
    <div className="container mx-auto space-y-6">
      <h1 className="text-2xl font-extrabold tracking-tight w-max">Project</h1>
      <ProjectForm project={project} />
    </div>
  );
};
