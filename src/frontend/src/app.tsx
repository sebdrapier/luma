import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { ProjectView } from "@/modules/projects/views/project-view";
import {
  Grid3x3,
  Lightbulb,
  PlayCircle,
  Save,
  Settings,
  Sliders,
} from "lucide-react";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "./components/ui/tooltip";
import { FixtureView } from "./modules/fixtures/views/fixture-view";
import { MixerView } from "./modules/mixer/views/mixer-view";
import { PerformanceView } from "./modules/performance/views/performance-view";
import { PresetView } from "./modules/presets/views/preset-view";
import { ShowView } from "./modules/shows/views/show-view";

const tabs = [
  {
    value: "project",
    label: "Project",
    icon: Settings,
    content: <ProjectView />,
  },
  {
    value: "fixtures",
    label: "Fixtures",
    icon: Lightbulb,
    content: <FixtureView />,
  },
  { value: "presets", label: "Presets", icon: Save, content: <PresetView /> },
  { value: "shows", label: "Shows", icon: PlayCircle, content: <ShowView /> },
  { value: "mixer", label: "Mixer", icon: Sliders, content: <MixerView /> },
  {
    value: "performance",
    label: "Performance",
    icon: Grid3x3,
    content: <PerformanceView />,
  },
];

function App() {
  return (
    <main className="min-h-dvh">
      <Tabs
        defaultValue={tabs[0].value}
        orientation="vertical"
        className="w-full flex-row h-full"
      >
        <TabsList className="flex-col h-dvh">
          {tabs.map(({ value, label, icon: Icon }) => (
            <TooltipProvider delayDuration={0} key={value}>
              <Tooltip>
                <TooltipTrigger asChild>
                  <span>
                    <TabsTrigger value={value} className="py-3 min-w-9">
                      <Icon size={16} aria-hidden="true" />
                    </TabsTrigger>
                  </span>
                </TooltipTrigger>
                <TooltipContent side="right" className="px-2 py-1 text-xs">
                  {label}
                </TooltipContent>
              </Tooltip>
            </TooltipProvider>
          ))}
        </TabsList>

        <div className="grow p-4 overflow-auto max-h-dvh">
          {tabs.map(({ value, content }) => (
            <TabsContent key={value} value={value}>
              {content}
            </TabsContent>
          ))}
        </div>
      </Tabs>
    </main>
  );
}

export default App;
