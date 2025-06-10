import { Button } from "@/components/ui/button";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";
import { useState } from "react";
import { FixtureForm } from "./fixture-form";
import { Plus } from "lucide-react";

export const NewFixtureSheet = () => {
  const [open, setOpen] = useState(false);

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>
        <Button variant="default">
          <Plus />
          Fixture
        </Button>
      </SheetTrigger>
      <SheetContent  className="flex flex-col h-dvh min-h-0" >
        <SheetHeader>
          <SheetTitle>Create a new fixture</SheetTitle>
          <SheetDescription>
            Fill out the details below to add a new DMX fixture to your setup.
          </SheetDescription>
        </SheetHeader>
        <div className="flex-1 overflow-y-auto min-h-0 p-4">
          <FixtureForm onSubmit={() => setOpen(false)} />
        </div>
      </SheetContent>
    </Sheet>
  );
};
