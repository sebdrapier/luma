import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import type { FC } from "react";
import type { Fixture } from "../fixture-types";

interface ChannelTableProps {
  fixture: Fixture;
}

export const ChannelTable: FC<ChannelTableProps> = ({ fixture }) => {
  return (
    <TableRow key={`${fixture.id}-subrow`} className="bg-muted">
      <TableCell colSpan={4}>
        <div className="p-2">
          <p className="font-medium mb-2">Channels</p>
          {fixture.channels.length === 0 ? (
            <p className="text-sm text-muted-foreground">No channels</p>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Description</TableHead>
                  <TableHead>Min</TableHead>
                  <TableHead>Max</TableHead>
                  <TableHead>Address</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {fixture.channels.map((channel, i) => (
                  <TableRow key={i}>
                    <TableCell>{channel.name}</TableCell>
                    <TableCell className="max-w-[500px] truncate">
                      {channel.description}
                    </TableCell>
                    <TableCell>{channel.min}</TableCell>
                    <TableCell>{channel.max}</TableCell>
                    <TableCell>{channel.channel_address}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </div>
      </TableCell>
    </TableRow>
  );
};
