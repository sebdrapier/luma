import { cn } from "@/lib/utils";
import { useDmxWebSocketContext } from "@/providers/ws-provider";
import { Badge } from "../ui/badge";

export const WsStatus = () => {
  const { isConnected } = useDmxWebSocketContext();

  return (
    <Badge
      className={cn(
        "rounded-full text-sm font-medium bg-red-100 text-red-800",
        {
          "bg-green-100 text-green-800": isConnected,
        }
      )}
    >
      {isConnected ? "Connected" : "Disconnected"}
    </Badge>
  );
};
