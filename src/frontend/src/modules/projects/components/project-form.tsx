"use client";

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
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";

import { USBSelect } from "@/modules/usb/components/usb-select";
import { Loader, SaveIcon } from "lucide-react";
import type { FC } from "react";
import { useUpdateProjectSettings } from "../project-api";
import type { Project } from "../project-types";

const projectFormSchema = z.object({
  name: z.string().min(1, { message: "Project name is required" }),
  usb_interface: z
    .string()
    .min(1, { message: "USB interface selection is required" }),
});

type ProjectFormValues = z.infer<typeof projectFormSchema>;

interface ProjectFormProps {
  project: Project;
}

export const ProjectForm: FC<ProjectFormProps> = ({ project }) => {
  const updateProjet = useUpdateProjectSettings();

  const form = useForm<ProjectFormValues>({
    resolver: zodResolver(projectFormSchema),
    defaultValues: {
      name: project.name ?? "",
      usb_interface: project.usb_interface ?? "",
    },
  });

  async function onSubmit(values: ProjectFormValues) {
    try {
      await updateProjet.mutateAsync(values);
      toast.success("Project saved successfully");
    } catch {
      toast.error("Failed to save project");
    }
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Project Name</FormLabel>
              <FormControl>
                <Input placeholder="Project Name" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="usb_interface"
          render={({ field }) => (
            <FormItem>
              <FormLabel>USB Interface</FormLabel>
              <FormControl>
                <USBSelect value={field.value} onSelect={field.onChange} />
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
          Save Project
        </Button>
      </form>
    </Form>
  );
};
