"use client";

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
import { type FC } from "react";
import { useDeletePreset } from "../preset-api";
import type { Preset } from "../preset-types";
import { EditPresetSheet } from "./edit-preset-sheet";

interface PresetTableProps {
  presets: Preset[];
}

export const PresetTable: FC<PresetTableProps> = ({ presets }) => {
  const deleteMutation = useDeletePreset();

  const { applyPreset } = useDmxWebSocketContext();

  const handlePreview = async (preset: Preset) => {
    applyPreset(preset.id);
  };

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className="w-[200px]">Nom</TableHead>
          <TableHead>Description</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {presets.map((preset) => (
          <TableRow key={preset.id}>
            <TableCell className="font-medium">{preset.name}</TableCell>
            <TableCell>{preset.description}</TableCell>
            <TableCell className="text-right space-x-2">
              <Button
                size="icon"
                variant="outline"
                onClick={() => handlePreview(preset)}
                aria-label={`Preview preset ${preset.name}`}
              >
                <Play className="w-4 h-4" />
              </Button>

              <EditPresetSheet preset={preset} />

              <DeleteModal
                onConfirm={() => deleteMutation.mutate(preset.id)}
                title={`Delete preset "${preset.name}"`}
                description="This will permanently remove the preset."
              >
                <Button
                  variant="destructive"
                  size="icon"
                  onClick={(e) => e.stopPropagation()}
                  aria-label={`Delete preset ${preset.name}`}
                >
                  <Trash2 className="w-4 h-4" />
                </Button>
              </DeleteModal>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
};
