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
import { Edit2 } from "lucide-react";
import { useState, type FC } from "react";
import type { Show } from "../show-types";
import { ShowForm } from "./show-form";

interface EditShowSheetProps {
  show: Show;
}

export const EditShowSheet: FC<EditShowSheetProps> = ({ show }) => {
  const [open, setOpen] = useState(false);

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>
        <Button
          size="icon"
          variant="outline"
          onClick={(e) => e.stopPropagation()}
        >
          <Edit2 className="w-4 h-4" />
        </Button>
      </SheetTrigger>
      <SheetContent className="flex flex-col h-dvh min-h-0">
        <SheetHeader>
          <SheetTitle>Edit show</SheetTitle>
          <SheetDescription>Update name and steps</SheetDescription>
        </SheetHeader>
        <div className="flex-1 overflow-y-auto min-h-0 p-4">
          <ShowForm show={show} onSubmit={() => setOpen(false)} />
        </div>
      </SheetContent>
    </Sheet>
  );
};
