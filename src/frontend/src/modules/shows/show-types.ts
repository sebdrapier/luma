export interface Show {
  id: string;
  name: string;
  steps: ShowStep[];
}

export interface ShowStep {
  preset_id: string;
  duration: number;
  fade_ms: number;
}
