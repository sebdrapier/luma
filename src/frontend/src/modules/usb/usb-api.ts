import kyClient from "@/integration/ky";
import { useQuery } from "@tanstack/react-query";

export async function getUSBInterfaces(): Promise<string[]> {
  return kyClient.get("usb/interfaces").json<string[]>();
}

export const usbQueryOptions = {
  all: () => ({
    queryKey: ["usbInterfaces"] as const,
    queryFn: getUSBInterfaces,
  }),
};

export function useUSBInterfaces() {
  return useQuery(usbQueryOptions.all());
}
