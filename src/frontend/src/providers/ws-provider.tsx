import type { Fixture } from "@/modules/fixtures/fixture-types";
import type { Preset } from "@/modules/presets/preset-types";
import type { Show } from "@/modules/shows/show-types";
import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
  type ReactNode,
} from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";

export const OutgoingMessageType = {
  APPLY_PRESET: "apply_preset",
  RUN_SHOW: "run_show",
  STOP_SHOW: "stop_show",
  UPDATE_CHANNEL: "update_channel",
  BLACKOUT: "blackout",
} as const;
export type OutgoingMessageType =
  (typeof OutgoingMessageType)[keyof typeof OutgoingMessageType];

export const IncomingMessageType = {
  PRESET_APPLIED: "preset_applied",
  SHOW_STARTED: "show_started",
  SHOW_STOPPED: "show_stopped",
  CHANNEL_UPDATE: "channel_update",
  BLACKOUT_EVENT: "blackout",
  DMX_STATE: "dmx_state",
  DMX_UPDATE: "dmx_update",
  PROJECT_CONFIG: "project_config",
  MONITORING_STARTED: "monitoring_started",
  MONITORING_STOPPED: "monitoring_stopped",
  ERROR: "error",
} as const;
export type IncomingMessageType =
  (typeof IncomingMessageType)[keyof typeof IncomingMessageType];

export type WSMessage<T = any> = { type: IncomingMessageType; payload: T };

export interface PresetAppliedPayload {
  preset_id: string;
  channels: Record<string, number>;
}
export interface ShowStartedPayload {
  show_id: string;
  steps: number;
  loop: boolean;
}
export type ShowStoppedPayload = object;
export interface ChannelUpdatePayload {
  dmx_address: number;
  value: number;
}
export type BlackoutPayload = object;
export interface DMXStatePayload {
  channels: { address: number; value: number }[];
  active_preset_id?: string;
  active_show_id?: string;
  show_step?: number;
  show_loop?: boolean;
  timestamp: number;
}
export interface ProjectConfigPayload {
  project_id: string;
  project_name: string;
  fixtures: Fixture[];
  presets: Preset[];
  shows: Show[];
}
export interface ErrorPayload {
  error: string;
  details?: string;
}

function useInternalDmxWebSocket(url: string, reconnectOnClose = true) {
  const { sendJsonMessage, lastJsonMessage, readyState } = useWebSocket(url, {
    onOpen: () => console.log("WebSocket connecté"),
    onError: (e) => console.error("WS Erreur", e),
    shouldReconnect: () => reconnectOnClose,
    reconnectAttempts: 10,
    retryOnError: true,
    share: true,
  });

  const [presetApplied, setPresetApplied] =
    useState<PresetAppliedPayload | null>(null);
  const [showStarted, setShowStarted] = useState<ShowStartedPayload | null>(
    null
  );
  const [dmxState, setDmxState] = useState<DMXStatePayload | null>(null);
  const [projectConfig, setProjectConfig] =
    useState<ProjectConfigPayload | null>(null);
  const [error, setError] = useState<ErrorPayload | null>(null);

  useEffect(() => {
    if (!lastJsonMessage) return;
    const msg = lastJsonMessage as WSMessage<any>;
    switch (msg.type) {
      case IncomingMessageType.PRESET_APPLIED:
        setPresetApplied(msg.payload as PresetAppliedPayload);
        break;
      case IncomingMessageType.SHOW_STARTED:
        setShowStarted(msg.payload as ShowStartedPayload);
        break;
      case IncomingMessageType.DMX_STATE:
        setDmxState(msg.payload as DMXStatePayload);
        break;
      case IncomingMessageType.DMX_UPDATE:
        setDmxState(msg.payload as DMXStatePayload);
        break;
      case IncomingMessageType.PROJECT_CONFIG:
        setProjectConfig(msg.payload as ProjectConfigPayload);
        break;
      case IncomingMessageType.ERROR:
        setError(msg.payload as ErrorPayload);
        break;
      default:
        break;
    }
  }, [lastJsonMessage]);

  const applyPreset = useCallback(
    (presetId: string) => {
      sendJsonMessage({
        type: OutgoingMessageType.APPLY_PRESET,
        payload: { preset_id: presetId },
      });
    },
    [sendJsonMessage]
  );

  const runShow = useCallback(
    (showId: string, loop = false) => {
      sendJsonMessage({
        type: OutgoingMessageType.RUN_SHOW,
        payload: { show_id: showId, loop },
      });
    },
    [sendJsonMessage]
  );

  const stopShow = useCallback(() => {
    sendJsonMessage({ type: OutgoingMessageType.STOP_SHOW, payload: {} });
  }, [sendJsonMessage]);

  const updateChannel = useCallback(
    (dmx_address: number, value: number) => {
      sendJsonMessage({
        type: OutgoingMessageType.UPDATE_CHANNEL,
        payload: { dmx_address, value },
      });
    },
    [sendJsonMessage]
  );

  const blackout = useCallback(() => {
    sendJsonMessage({ type: OutgoingMessageType.BLACKOUT, payload: {} });
  }, [sendJsonMessage]);

  return {
    applyPreset,
    runShow,
    stopShow,
    updateChannel,
    blackout,
    presetApplied,
    showStarted,
    dmxState,
    projectConfig,
    error,
    isConnected: readyState === ReadyState.OPEN,
  };
}

export interface DmxWebSocketContextValue {
  applyPreset: (presetId: string) => void;
  runShow: (showId: string, loop?: boolean) => void;
  stopShow: () => void;
  updateChannel: (dmx_address: number, value: number) => void;
  blackout: () => void;
  presetApplied: PresetAppliedPayload | null;
  showStarted: ShowStartedPayload | null;
  dmxState: DMXStatePayload | null;
  projectConfig: ProjectConfigPayload | null;
  error: ErrorPayload | null;
  isConnected: boolean;
}

const DmxWebSocketContext = createContext<DmxWebSocketContextValue | undefined>(
  undefined
);

export interface DmxWebSocketProviderProps {
  url: string;
  children: ReactNode;
}

export function DmxWebSocketProvider({
  url,
  children,
}: DmxWebSocketProviderProps) {
  const ws = useInternalDmxWebSocket(url);
  return (
    <DmxWebSocketContext.Provider value={ws}>
      {children}
    </DmxWebSocketContext.Provider>
  );
}

export function useDmxWebSocketContext(): DmxWebSocketContextValue {
  const context = useContext(DmxWebSocketContext);
  if (!context)
    throw new Error(
      "useDmxWebSocketContext must be used within DmxWebSocketProvider"
    );
  return context;
}
