import kyClient from "@/integration/ky";
import type { ChannelValue, Preset } from "@/modules/presets/preset-types";
import type { Show } from "@/modules/shows/show-types";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

async function getShows() {
  return kyClient.get("shows").json<Show[]>();
}

async function getShow(id: string) {
  return kyClient.get(`shows/${id}`).json<Show>();
}

async function createShow(show: Omit<Show, "id">) {
  return kyClient
    .post("shows", {
      json: show,
    })
    .json<Show>();
}

async function updateShow(id: string, show: Omit<Show, "id">) {
  return kyClient
    .put(`shows/${id}`, {
      json: show,
    })
    .json<Show>();
}

async function deleteShow(id: string) {
  return kyClient.delete(`shows/${id}`).json<void>();
}

export const showsQueryOptions = {
  all: () => ({
    queryKey: ["shows"] as const,
    queryFn: getShows,
  }),

  detail: (id: string) => ({
    queryKey: ["shows", id] as const,
    queryFn: () => getShow(id),
    enabled: !!id,
  }),
};

export function useShows() {
  return useQuery(showsQueryOptions.all());
}

export function useShow(id: string) {
  return useQuery(showsQueryOptions.detail(id));
}

export function useCreateShow() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: createShow,
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ["shows"] });
      queryClient.setQueryData(["shows", data.id], data);
    },
  });
}

export function useUpdateShow() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, ...data }: Show) => updateShow(id, data),
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({ queryKey: ["shows"] });
      queryClient.setQueryData(["shows", variables.id], data);
    },
  });
}

export function useDeleteShow() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: deleteShow,
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: ["shows"] });
      queryClient.removeQueries({ queryKey: ["shows", id] });
    },
  });
}

export function formatShowForPlayback(
  show: Show,
  presets: Record<string, Preset>
) {
  const steps = show.steps.map((step) => {
    const preset = presets[step.preset_id];
    if (!preset) {
      throw new Error(`Preset ${step.preset_id} not found`);
    }

    const channels: Record<string, number> = {};
    preset.channels.forEach((ch: ChannelValue) => {
      channels[ch.dmx_address.toString()] = ch.value;
    });

    return {
      preset: channels,
      duration: step.duration,
    };
  });

  return { steps };
}
