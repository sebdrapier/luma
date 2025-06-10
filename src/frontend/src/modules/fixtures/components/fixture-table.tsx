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
import { Trash2 } from "lucide-react";
import { useState, type FC } from "react";
import { useDeleteFixture } from "../fixture-api";
import type { Fixture } from "../fixture-types";
import { ChannelTable } from "./channel-table";
import { EditFixtureSheet } from "./edit-fixture-sheet";

interface FixtureTableProps {
  fixtures: Fixture[];
}

export const FixtureTable: FC<FixtureTableProps> = ({ fixtures }) => {
  const deleteMutation = useDeleteFixture();
  const [openSubrows, setOpenSubrows] = useState<Record<string, boolean>>({});

  const toggleSubrow = (id: string) => {
    setOpenSubrows((prev) => ({ ...prev, [id]: !prev[id] }));
  };

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className="w-[200px]">Name</TableHead>
          <TableHead>Description</TableHead>
          <TableHead>Type</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {fixtures.map((fixture) => (
          <>
            <TableRow
              key={fixture.id}
              onClick={() => toggleSubrow(fixture.id)}
              className="cursor-pointer"
            >
              <TableCell className="font-medium">{fixture.name}</TableCell>
              <TableCell>{fixture.description}</TableCell>
              <TableCell>{fixture.type}</TableCell>
              <TableCell className="text-right space-x-2">
                <EditFixtureSheet fixture={fixture} />
                <DeleteModal
                  onConfirm={() => deleteMutation.mutate(fixture.id)}
                  title={`Delete "${fixture.name}"`}
                  description="This will permanently remove the fixture and all its channels."
                >
                  <Button
                    variant="destructive"
                    size="icon"
                    onClick={(e) => e.stopPropagation()}
                  >
                    <Trash2 className="w-4 h-4" />
                  </Button>
                </DeleteModal>
              </TableCell>
            </TableRow>
            {openSubrows[fixture.id] && <ChannelTable fixture={fixture} />}
          </>
        ))}
      </TableBody>
    </Table>
  );
};
