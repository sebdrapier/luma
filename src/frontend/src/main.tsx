import App from "@/app.tsx";
import { ThemeProvider } from "@/components/theme-provider.tsx";
import { Toaster } from "@/components/ui/sonner";
import { queryClient } from "@/integration/tanstack-query.ts";
import "@/styles/tailwind.css";
import "@/styles/overrides.css";

import { QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { DmxWebSocketProvider } from "./providers/ws-provider";

function makeWebSocketUrl(path = "/ws/control") {
  // pick ws vs wss depending on page protocol
  const protocol = window.location.protocol === "https:" ? "wss" : "ws";
  const { hostname, port } = window.location;
  // only include “:port” if there actually is one
  let portPart = port ? `${port}` : "";
  if (port === "5173") portPart = "3000";
  return `${protocol}://${hostname}:${portPart}${path}`;
}

//"ws://localhost:3000/ws/control"

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
      <QueryClientProvider client={queryClient}>
        <DmxWebSocketProvider url={makeWebSocketUrl()}>
          <App />
        </DmxWebSocketProvider>
        <Toaster position="top-center" richColors />
        <ReactQueryDevtools initialIsOpen={false} />
      </QueryClientProvider>
    </ThemeProvider>
  </StrictMode>
);
