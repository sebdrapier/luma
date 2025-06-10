import { cn } from "@/lib/utils";
import type { FC } from "react";

interface PerformanceButtonProps {
  id: string;
  onClick: () => void;
  isActive: boolean;
  name: string;
  description: string;
}

export const PerformanceButton: FC<PerformanceButtonProps> = ({
  isActive,
  onClick,
  id,
  name,
  description,
}) => {
  return (
    <button
      key={id}
      onClick={onClick}
      className={cn(
        "p-4 rounded-lg shadow-sm transition bg-card cursor-pointer",
        {
          "bg-primary text-primary-foreground": isActive,
        }
      )}
    >
      <h3 className="font-medium mb-1 truncate">{name}</h3>
      <p className="text-sm text-gray-500">{description}</p>
    </button>
  );
};
