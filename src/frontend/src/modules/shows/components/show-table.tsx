import { DeleteModal } from "@/components/delete-modal";
import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useDmxWebSocketContext } from "@/providers/ws-provider";
import { Play, Trash2 } from "lucide-react";
import type { FC } from "react";
import { useDeleteShow } from "../show-api";
import type { Show } from "../show-types";
import { EditShowSheet } from "./edit-show-sheet";

interface ShowTableProps {
  shows: Show[];
}

export const ShowTable: FC<ShowTableProps> = ({ shows }) => {
  const deleteMutation = useDeleteShow();
  const { runShow, stopShow, showStarted } = useDmxWebSocketContext();

  const handleRun = (show: Show) => {
    runShow(show.id, false);
  };

  const handleStop = () => {
    stopShow();
  };

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className="w-[200px]">Nom</TableHead>
          <TableHead>Étapes</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {shows.map((show) => {
          const isRunning = showStarted?.show_id === show.id;
          return (
            <TableRow key={show.id}>
              <TableCell className="font-medium">{show.name}</TableCell>
              <TableCell>{show.steps.length}</TableCell>
              <TableCell className="text-right space-x-2">
                <Button
                  size="icon"
                  variant="outline"
                  onClick={() => (isRunning ? handleStop() : handleRun(show))}
                  aria-label={`${isRunning ? "Stop" : "Run"} show ${show.name}`}
                >
                  <Play className="w-4 h-4" />
                </Button>

                <EditShowSheet show={show} />

                <DeleteModal
                  onConfirm={() => deleteMutation.mutate(show.id)}
                  title={`Supprimer le show "${show.name}"`}
                  description="Cela supprimera définitivement le show."
                >
                  <Button
                    variant="destructive"
                    size="icon"
                    onClick={(e) => e.stopPropagation()}
                    aria-label={`Delete show ${show.name}`}
                  >
                    <Trash2 className="w-4 h-4" />
                  </Button>
                </DeleteModal>
              </TableCell>
            </TableRow>
          );
        })}
      </TableBody>
    </Table>
  );
};
