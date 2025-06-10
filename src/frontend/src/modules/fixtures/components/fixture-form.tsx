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

import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible";
import { cn } from "@/lib/utils";
import { ChevronRight, Loader, Plus, SaveIcon, X } from "lucide-react";
import type { FC } from "react";
import {
  useCreateFixture,
  useFixtures,
  useUpdateFixture,
} from "../fixture-api";
import type { Fixture } from "../fixture-types";

const channelSchema = z.object({
  name: z.string().min(1, "Name is required"),
  description: z.string().optional(),
  min: z.coerce.number().min(0).max(255),
  max: z.coerce.number().min(0).max(255),
  channel_address: z.coerce.number().min(1, "Channel must be >= 1"),
});

const schema = z.object({
  name: z.string().min(1, { message: "Name is required" }),
  description: z.string().optional(),
  type: z.string().min(1, { message: "Type is required" }),
  channels: z.array(channelSchema).min(1, "At least one channel is required"),
});

type FixtureFormValues = z.infer<typeof schema>;

interface FixtureFormProps {
  fixture?: Fixture;
  onSubmit?: () => void;
}

export const FixtureForm: FC<FixtureFormProps> = ({ fixture, onSubmit }) => {
  const { data: allFixtures } = useFixtures();

  const getNextChannelAddress = () => {
    const currentChannels = form.getValues("channels");

    if (currentChannels.length > 0) {
      const last = currentChannels[currentChannels.length - 1];
      return last.channel_address + 1;
    }

    const allChannels =
      allFixtures?.flatMap((f) => f.channels)?.map((c) => c.channel_address) ||
      [];

    const maxAddress = allChannels.length > 0 ? Math.max(...allChannels) : 0;

    return maxAddress + 1;
  };

  const isEdit = !!fixture;

  const form = useForm<FixtureFormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      name: fixture?.name || "",
      description: fixture?.description || "",
      type: fixture?.type || "",
      channels: fixture?.channels || [],
    },
  });

  const { fields, append, remove } = useFieldArray({
    control: form.control,
    name: "channels",
  });

  const createFixture = useCreateFixture();
  const updateFixture = useUpdateFixture();

  const handleSubmit = async (values: FixtureFormValues) => {
    try {
      if (isEdit && fixture) {
        await updateFixture.mutateAsync({ ...fixture, ...(values as Fixture) });
        toast.success("Fixture updated");
      } else {
        await createFixture.mutateAsync(values as Fixture);
        toast.success("Fixture created");
      }
      onSubmit?.();
    } catch {
      toast.error("An error occurred while saving the fixture");
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
              <FormLabel>Name</FormLabel>
              <FormControl>
                <Input placeholder="e.g. LED PAR 64" {...field} />
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
                <Input placeholder="Optional details..." {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="type"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Type</FormLabel>
              <FormControl>
                <Input placeholder="e.g. RGBW" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <div className="space-y-4">
          <FormLabel>Channels</FormLabel>

          {fields.map((field, index) => {
            const channelName =
              form.watch(`channels.${index}.name`) || `Channel ${index + 1}`;
            return (
              <Collapsible key={field.id} className="border rounded-md">
                <CollapsibleTrigger asChild>
                  <button
                    type="button"
                    className={cn(
                      "w-full flex items-center justify-between px-4 py-2 text-left hover:bg-muted transition"
                    )}
                  >
                    <span className="font-medium">{channelName}</span>
                    <ChevronRight className="transition-transform group-data-[state=open]:rotate-90" />
                  </button>
                </CollapsibleTrigger>
                <CollapsibleContent className="space-y-4 p-4">
                  <div className="flex items-end gap-2">
                    <FormField
                      control={form.control}
                      name={`channels.${index}.name`}
                      render={({ field }) => (
                        <FormItem className="grow">
                          <FormLabel>Name</FormLabel>
                          <FormControl>
                            <Input {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <Button
                      variant="destructive"
                      size="icon"
                      type="button"
                      onClick={() => remove(index)}
                    >
                      <X />
                    </Button>
                  </div>

                  <div className="grid grid-cols-3 gap-2">
                    <FormField
                      control={form.control}
                      name={`channels.${index}.channel_address`}
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Address</FormLabel>
                          <FormControl>
                            <Input type="number" min={1} {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={form.control}
                      name={`channels.${index}.min`}
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Min</FormLabel>
                          <FormControl>
                            <Input type="number" min={0} max={255} {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={form.control}
                      name={`channels.${index}.max`}
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Max</FormLabel>
                          <FormControl>
                            <Input type="number" min={0} max={255} {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>
                </CollapsibleContent>
              </Collapsible>
            );
          })}

          <Button
            type="button"
            className="w-full"
            variant="outline"
            onClick={() =>
              append({
                name: "",
                description: "",
                min: 0,
                max: 255,
                channel_address: getNextChannelAddress(),
              })
            }
          >
            <Plus className="mr-2 h-4 w-4" />
            Channel
          </Button>
        </div>

        <Button type="submit" className="w-full" disabled={form.formState.isSubmitting}>
          {form.formState.isSubmitting ? (
            <Loader className="animate-spin" />
          ) : (
            <SaveIcon />
          )}
          {isEdit ? "Update Fixture" : "Create Fixture"}
        </Button>
      </form>
    </Form>
  );
};
