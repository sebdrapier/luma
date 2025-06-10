"use client";

import { Edit2 } from "lucide-react";
import { type FC, useState } from "react";

import { Button } from "@/components/ui/button";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";

import type { Preset } from "../preset-types";
import { PresetForm } from "./preset-form";

interface EditPresetSheetProps {
  preset: Preset;
}

export const EditPresetSheet: FC<EditPresetSheetProps> = ({ preset }) => {
  const [open, setOpen] = useState(false);

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>
        <Button size="icon" variant="outline">
          <Edit2 className="w-4 h-4" />
        </Button>
      </SheetTrigger>

      <SheetContent className="flex flex-col h-dvh min-h-0">
        <SheetHeader>
          <SheetTitle>Update preset</SheetTitle>
          <SheetDescription>Update name and description</SheetDescription>
        </SheetHeader>

        <div className="flex-1 overflow-y-auto min-h-0 p-4">
          <PresetForm
            preset={preset}
            onSubmit={() => {
              setOpen(false);
            }}
          />
        </div>
      </SheetContent>
    </Sheet>
  );
};
