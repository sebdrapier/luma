import { Button } from "@/components/ui/button";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";
import type { Fixture } from "@/modules/fixtures/fixture-types";
import { Edit } from "lucide-react";
import { useState } from "react";
import { FixtureForm } from "./fixture-form";

interface EditFixtureSheetProps {
  fixture: Fixture;
}

export const EditFixtureSheet = ({ fixture }: EditFixtureSheetProps) => {
  const [open, setOpen] = useState(false);

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>
        <Button
          variant="outline"
          size="icon"
          onClick={(e) => e.stopPropagation()}
        >
          <Edit />
        </Button>
      </SheetTrigger>
      <SheetContent className="flex flex-col h-dvh min-h-0">
        <SheetHeader>
          <SheetTitle>Edit fixture</SheetTitle>
          <SheetDescription>
            Modify the details of your DMX fixture below.
          </SheetDescription>
        </SheetHeader>
        <div className="flex-1 overflow-y-auto min-h-0 p-4">
          <FixtureForm fixture={fixture} onSubmit={() => setOpen(false)} />
        </div>
      </SheetContent>
    </Sheet>
  );
};
