import kyClient from "@/integration/ky";
import { queryClient } from "@/integration/tanstack-query";
import type { Fixture } from "@/modules/fixtures/fixture-types";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

async function getFixtures() {
  return kyClient.get("fixtures").json<Fixture[]>();
}

async function getFixture(id: string) {
  return kyClient.get(`fixtures/${id}`).json<Fixture>();
}

async function createFixture(fixture: Omit<Fixture, "id">) {
  return kyClient
    .post("fixtures", {
      json: fixture,
    })
    .json<Fixture>();
}

async function updateFixture(id: string, fixture: Omit<Fixture, "id">) {
  return kyClient
    .put(`fixtures/${id}`, {
      json: fixture,
    })
    .json<Fixture>();
}

async function deleteFixture(id: string) {
  return kyClient.delete(`fixtures/${id}`).json<void>();
}

export const fixturesQueryOptions = {
  all: () => ({
    queryKey: ["fixtures"] as const,
    queryFn: getFixtures,
  }),

  detail: (id: string) => ({
    queryKey: ["fixtures", id] as const,
    queryFn: () => getFixture(id),
    enabled: !!id,
  }),
};

export const useCreateFixtureMutation = () => ({
  mutationFn: createFixture,
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ["fixtures"] });
  },
});

export const useUpdateFixtureMutation = () => ({
  mutationFn: ({ id, ...data }: Fixture) => updateFixture(id, data),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ["fixtures"] });
  },
});

export const useDeleteFixtureMutation = () => ({
  mutationFn: deleteFixture,
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ["fixtures"] });
  },
});

export function useFixtures() {
  return useQuery(fixturesQueryOptions.all());
}

export function useFixture(id: string) {
  return useQuery(fixturesQueryOptions.detail(id));
}

export function useCreateFixture() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: createFixture,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["fixtures"] });
    },
  });
}

export function useUpdateFixture() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, ...data }: Fixture) => updateFixture(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["fixtures"] });
      queryClient.invalidateQueries({ queryKey: ["fixtures", variables.id] });
    },
  });
}

export function useDeleteFixture() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: deleteFixture,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["fixtures"] });
    },
  });
}
