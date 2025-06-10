import { ErrorAlert } from "@/components/error-alert";
import { InfoAlert } from "@/components/info-alert";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Skeleton } from "@/components/ui/skeleton";
import { type FC } from "react";
import { useUSBInterfaces } from "../usb-api";

export interface USBSelectProps {
  onSelect: (value: string) => void;
  value: string;
  id?: string;
  label?: string;
}

export const USBSelect: FC<USBSelectProps> = ({
  onSelect,
  value,
  id = "usb-select",
  label = "USB Interface",
}) => {
  const { data: usbs, isLoading, isError } = useUSBInterfaces();

  if (isLoading) return <Skeleton className="h-9 w-full rounded-md" />;

  if (isError)
    return (
      <ErrorAlert
        title="Error loading USB interfaces"
        message="Could not retrieve USB interfaces. Make sure the backend is running."
      />
    );

  if (!usbs)
    return (
      <InfoAlert
        title="No USB interfaces found"
        message="Please connect a DMX device and try again."
      />
    );

  return (
    <Select onValueChange={onSelect} defaultValue={value}>
      <SelectTrigger id={id} aria-label={label} className="w-full">
        <SelectValue
          placeholder={
            isLoading ? "Loading USB interfaces..." : "Select USB interface"
          }
        />
      </SelectTrigger>
      <SelectContent>
        {usbs?.map((usb) => (
          <SelectItem key={usb} value={usb}>
            {usb}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
};
