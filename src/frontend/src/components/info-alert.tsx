import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import type { FC } from "react";

interface InfoAlertProps {
  title: string;
  message: string;
}

export const InfoAlert: FC<InfoAlertProps> = ({ title, message }) => {
  return (
    <Alert>
      <AlertTitle>{title}</AlertTitle>
      <AlertDescription>{message}</AlertDescription>
    </Alert>
  );
};
