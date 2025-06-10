import type { Fixture } from "@/modules/fixtures/fixture-types";
import type { Preset } from "@/modules/presets/preset-types";
import type { Show } from "@/modules/shows/show-types";

export interface Project {
  id: string;
  name: string;
  usb_interface: string;
  fixtures: Fixture[];
  presets: Preset[];
  shows: Show[];
}
