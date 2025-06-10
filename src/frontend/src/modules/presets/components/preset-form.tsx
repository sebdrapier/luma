"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { type FC } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";

import { Loader, SaveIcon } from "lucide-react";
import { useUpdatePreset } from "../preset-api";
import type { Preset } from "../preset-types";

const schema = z.object({
  name: z.string().min(1, { message: "Name is required" }),
  description: z.string().optional(),
});

type PresetFormValues = z.infer<typeof schema>;

interface PresetFormProps {
  preset?: Preset;
  onSubmit?: () => void;
}

export const PresetForm: FC<PresetFormProps> = ({ preset, onSubmit }) => {
  const isEdit = !!preset;

  const form = useForm<PresetFormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      name: preset?.name ?? "",
      description: preset?.description ?? "",
    },
  });

  const updatePreset = useUpdatePreset();

  const handleSubmit = async (values: PresetFormValues) => {
    try {
      if (isEdit && preset) {
        await updatePreset.mutateAsync({ ...preset, ...values });
        toast.success("Preset updated!");
      }
      onSubmit?.();
    } catch {
      toast.error("An error occurred while saving the preset.");
    }
  };

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Name</FormLabel>
              <FormControl>
                <Input placeholder="My preset" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="description"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Description</FormLabel>
              <FormControl>
                <Input placeholder="Optional…" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <Button type="submit" className="w-full" disabled={form.formState.isSubmitting}>
          {form.formState.isSubmitting ? (
            <Loader className="animate-spin" />
          ) : (
            <SaveIcon />
          )}
          Update Preset
        </Button>
      </form>
    </Form>
  );
};
