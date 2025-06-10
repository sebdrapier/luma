import kyClient from "@/integration/ky";
import type { Project } from "@/modules/projects/project-types";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

async function getProject() {
  return kyClient.get("projects").json<Project>();
}

async function createProject(project: Omit<Project, "id">) {
  return kyClient
    .post("projects", {
      json: project,
    })
    .json<Project>();
}

async function updateProject(project: Project) {
  return kyClient
    .put("projects", {
      json: project,
    })
    .json<Project>();
}

async function deleteProject() {
  return kyClient.delete("projects").json<void>();
}

export const projectQueryOptions = {
  current: () => ({
    queryKey: ["project"] as const,
    queryFn: getProject,
    staleTime: 1000 * 60 * 5,
  }),
};

export function useProject() {
  return useQuery(projectQueryOptions.current());
}

export function useCreateProject() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: createProject,
    onSuccess: (data) => {
      queryClient.setQueryData(["project"], data);
      queryClient.invalidateQueries({ queryKey: ["fixtures"] });
      queryClient.invalidateQueries({ queryKey: ["presets"] });
      queryClient.invalidateQueries({ queryKey: ["shows"] });
    },
  });
}

export function useUpdateProject() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: updateProject,
    onSuccess: (data) => {
      queryClient.setQueryData(["project"], data);
    },
  });
}

export function useDeleteProject() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: deleteProject,
    onSuccess: () => {
      queryClient.removeQueries({ queryKey: ["project"] });
      queryClient.removeQueries({ queryKey: ["fixtures"] });
      queryClient.removeQueries({ queryKey: ["presets"] });
      queryClient.removeQueries({ queryKey: ["shows"] });
    },
  });
}

export function useProjectName() {
  const { data: project } = useProject();
  return project?.name || "Untitled Project";
}

export function useProjectUSBInterface() {
  const { data: project } = useProject();
  return project?.usb_interface || "/dev/ttyUSB0";
}

export function useUpdateProjectSettings() {
  const queryClient = useQueryClient();
  const { data: currentProject } = useProject();

  return useMutation({
    mutationFn: async (settings: { name?: string; usb_interface?: string }) => {
      if (!currentProject) throw new Error("No project loaded");

      const updatedProject = {
        ...currentProject,
        ...settings,
      };

      return updateProject(updatedProject);
    },
    onSuccess: (data) => {
      queryClient.setQueryData(["project"], data);
    },
  });
}
