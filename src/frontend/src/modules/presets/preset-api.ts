import kyClient from "@/integration/ky";
import type { Preset } from "@/modules/presets/preset-types";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

async function getPresets() {
  return kyClient.get("presets").json<Preset[]>();
}

async function getPreset(id: string) {
  return kyClient.get(`presets/${id}`).json<Preset>();
}

async function createPreset(preset: Omit<Preset, "id">) {
  return kyClient
    .post("presets", {
      json: preset,
    })
    .json<Preset>();
}

async function updatePreset(id: string, preset: Omit<Preset, "id">) {
  return kyClient
    .put(`presets/${id}`, {
      json: preset,
    })
    .json<Preset>();
}

async function deletePreset(id: string) {
  return kyClient.delete(`presets/${id}`).json<void>();
}

export const presetsQueryOptions = {
  all: () => ({
    queryKey: ["presets"] as const,
    queryFn: getPresets,
  }),

  detail: (id: string) => ({
    queryKey: ["presets", id] as const,
    queryFn: () => getPreset(id),
    enabled: !!id,
  }),
};

export function usePresets() {
  return useQuery(presetsQueryOptions.all());
}

export function usePreset(id: string) {
  return useQuery(presetsQueryOptions.detail(id));
}

export function useCreatePreset() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: createPreset,
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ["presets"] });
      queryClient.setQueryData(["presets", data.id], data);
    },
  });
}

export function useUpdatePreset() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, ...data }: Preset) => updatePreset(id, data),
    onSuccess: (data, variables) => {
      queryClient.invalidateQueries({ queryKey: ["presets"] });
      queryClient.setQueryData(["presets", variables.id], data);
    },
  });
}

export function useDeletePreset() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: deletePreset,
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: ["presets"] });
      queryClient.removeQueries({ queryKey: ["presets", id] });
    },
  });
}

export async function applyPreset(preset: Preset) {
  const channels: Record<string, number> = {};
  preset.channels.forEach((ch) => {
    channels[ch.dmx_address.toString()] = ch.value;
  });

  return channels;
}
