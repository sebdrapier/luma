"use client";

import { Button } from "@/components/ui/button";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";
import { Plus } from "lucide-react";
import { useState } from "react";
import { SequencerForm } from "./sequencer-form";

export const NewShowSheet = () => {
  const [open, setOpen] = useState(false);

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>
        <Button>
          <Plus />
          Show
        </Button>
      </SheetTrigger>
      <SheetContent className="flex flex-col h-dvh min-h-0">
        <SheetHeader>
          <SheetTitle>Create show sequence</SheetTitle>
          <SheetDescription>
            Choose beat duration and presets for each step
          </SheetDescription>
        </SheetHeader>
        <div className="flex-1 overflow-y-auto min-h-0 p-4">
          <SequencerForm onSubmit={() => setOpen(false)} />
        </div>
      </SheetContent>
    </Sheet>
  );
};
