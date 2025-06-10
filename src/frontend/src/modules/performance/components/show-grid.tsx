import type { Show } from "@/modules/shows/show-types";
import { useDmxWebSocketContext } from "@/providers/ws-provider";
import { type FC } from "react";
import { PerformanceButton } from "../components/performance-button";

interface ShowGridProps {
  shows: Show[];
}

export const ShowGrid: FC<ShowGridProps> = ({ shows }) => {
  const { dmxState, runShow, stopShow, showStarted } = useDmxWebSocketContext();

  const activeShowId = dmxState?.active_show_id ?? showStarted?.show_id;

  return (
    <section>
      <h2 className="text-xl font-semibold space-y-2">Shows</h2>
      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4">
        {shows.map((show) => (
          <PerformanceButton
            key={show.id}
            id={show.id}
            onClick={() =>
              activeShowId === show.id ? stopShow() : runShow(show.id, true)
            }
            isActive={activeShowId === show.id}
            name={show.name}
            description={`Étapes: ${show.steps.length}`}
          />
        ))}
      </div>
    </section>
  );
};
