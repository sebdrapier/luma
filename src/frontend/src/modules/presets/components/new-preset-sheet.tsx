import { Button } from "@/components/ui/button";
import type { ChannelValue, Preset } from "@/modules/presets/preset-types";
import { useDmxWebSocketContext } from "@/providers/ws-provider";
import { Edit, PlusIcon } from "lucide-react";
import { useState, type FC } from "react";

import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";
import { NewPresetForm } from "@/modules/presets/components/new-preset-form";

interface NewPresetSheetProps {
  preset: Preset | null;
}

export const NewPresetSheet: FC<NewPresetSheetProps> = ({ preset }) => {
  const { dmxState } = useDmxWebSocketContext();

  const [open, setOpen] = useState(false);

  if (!dmxState || !dmxState.channels) return;

  const dmxChannels: ChannelValue[] =
    dmxState.channels.map((chan) => ({
      dmx_address: chan.address,
      value: chan.value,
    })) ?? [];

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>
        <Button>
          {preset ? (
            <>
              <Edit />
              Edit Preset
            </>
          ) : (
            <>
              <PlusIcon />
              New Preset
            </>
          )}
        </Button>
      </SheetTrigger>

      <SheetContent className="flex flex-col h-dvh min-h-0">
        <SheetHeader>
          <SheetTitle>Create preset</SheetTitle>
          <SheetDescription>Add name and description</SheetDescription>
        </SheetHeader>

        <div className="flex-1 overflow-y-auto min-h-0 p-4">
          <NewPresetForm
            preset={{
              id: "",
              name: "",
              description: "",
              ...preset,
              channels: dmxChannels,
            }}
            onSubmit={() => {
              setOpen(false);
            }}
          />
        </div>
      </SheetContent>
    </Sheet>
  );
};
