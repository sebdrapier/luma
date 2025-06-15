"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useFieldArray, useForm } from "react-hook-form";
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
import { PresetSelect } from "@/modules/presets/components/preset-select";
import { Loader, Plus, SaveIcon, Trash2 } from "lucide-react";
import type { FC } from "react";
import { useCreateShow, useUpdateShow } from "../show-api";
import type { Show, ShowStep } from "../show-types";

const stepSchema = z.object({
  preset_id: z.string().min(1, { message: "Preset is required" }),
  beats: z.coerce.number().int().min(1, { message: "Beats must be >= 1" }),
  fade: z.coerce.number().int(),
});

const schema = z.object({
  name: z.string().min(1, { message: "Name is required" }),
  beat_duration_ms: z.coerce
    .number()
    .min(1, { message: "Beat duration must be > 0" }),
  steps: z.array(stepSchema).min(1, { message: "At least one step" }),
});

type FormValues = z.infer<typeof schema>;

interface SequencerFormProps {
  show?: Show;
  onSubmit?: () => void;
}

function gcd(a: number, b: number): number {
  while (b !== 0) {
    const t = b;
    b = a % b;
    a = t;
  }
  return a;
}

export const SequencerForm: FC<SequencerFormProps> = ({ show, onSubmit }) => {
  const beatDefault =
    show && show.steps.length > 0
      ? show.steps.reduce(
          (acc, s) => gcd(acc, s.delay_ms),
          show.steps[0].delay_ms
        )
      : 1000;

  const form = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      name: show?.name ?? "",
      beat_duration_ms: beatDefault,
      steps: show?.steps.map((s) => ({
        preset_id: s.preset_id,
        beats: Math.round(s.delay_ms / beatDefault),
        fade: s.fade_ms,
      })) ?? [{ preset_id: "", beats: 1 }],
    },
  });

  const { fields, append, remove } = useFieldArray({
    control: form.control,
    name: "steps",
  });

  const createShow = useCreateShow();
  const updateShow = useUpdateShow();

  const handleSubmit = async (values: FormValues) => {
    const steps: ShowStep[] = values.steps.map((s) => ({
      preset_id: s.preset_id,
      delay_ms: s.beats * values.beat_duration_ms,
      fade_ms: s.fade,
    }));

    try {
      if (show) {
        await updateShow.mutateAsync({ id: show.id, name: values.name, steps });
        toast.success("Show updated");
      } else {
        await createShow.mutateAsync({ name: values.name, steps });
        toast.success("Show created");
      }
      onSubmit?.();
    } catch {
      toast.error("Failed to save show");
    }
  };

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-6">
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Show Name</FormLabel>
              <FormControl>
                <Input placeholder="My Show" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="beat_duration_ms"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Beat duration (ms)</FormLabel>
              <FormControl>
                <Input type="number" min={1} {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <div className="space-y-4">
          <FormLabel>Steps</FormLabel>
          {fields.map((field, index) => (
            <div key={field.id} className="grid grid-cols-9 gap-2 items-end">
              <FormField
                control={form.control}
                name={`steps.${index}.preset_id`}
                render={({ field }) => (
                  <FormItem className="col-span-4">
                    <FormLabel>Preset</FormLabel>
                    <FormControl>
                      <PresetSelect
                        value={field.value}
                        onSelect={(p) => field.onChange(p.id)}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name={`steps.${index}.beats`}
                render={({ field }) => (
                  <FormItem className="col-span-2">
                    <FormLabel>Beats</FormLabel>
                    <FormControl>
                      <Input type="number" min={1} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name={`steps.${index}.fade`}
                render={({ field }) => (
                  <FormItem className="col-span-2">
                    <FormLabel>Fade</FormLabel>
                    <FormControl>
                      <Input type="number" min={0} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <Button
                type="button"
                variant="destructive"
                size="icon"
                onClick={() => remove(index)}
              >
                <Trash2 className="w-4 h-4" />
              </Button>
            </div>
          ))}

          <Button
            type="button"
            className="w-full"
            variant="outline"
            onClick={() => append({ preset_id: "", beats: 1, fade: 0 })}
          >
            <Plus className="mr-2 h-4 w-4" />
            Step
          </Button>
        </div>

        <Button
          type="submit"
          className="w-full"
          disabled={form.formState.isSubmitting}
        >
          {form.formState.isSubmitting ? (
            <Loader className="animate-spin" />
          ) : (
            <SaveIcon />
          )}
          {show ? "Update Show" : "Create Show"}
        </Button>
      </form>
    </Form>
  );
};
