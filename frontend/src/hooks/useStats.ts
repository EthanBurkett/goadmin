import { useQuery } from "@tanstack/react-query";
import api from "@/lib/api";

export interface ServerStat {
  id: number;
  timestamp: string;
  playerCount: number;
  maxPlayers: number;
  mapName: string;
  gametype: string;
  hostname: string;
  fps: number;
  uptime: number;
  createdAt: string;
}

export interface SystemStat {
  id: number;
  timestamp: string;
  cpuUsage: number;
  memoryUsed: number;
  memoryTotal: number;
  createdAt: string;
}

export interface PlayerStat {
  id: number;
  timestamp: string;
  totalKills: number;
  totalDeaths: number;
  avgPing: number;
  avgScore: number;
  createdAt: string;
}

export const useServerStats = (start?: string, end?: string) => {
  return useQuery<ServerStat[]>({
    queryKey: ["server-stats", start, end],
    queryFn: async () => {
      const params = new URLSearchParams();
      if (start) params.append("start", start);
      if (end) params.append("end", end);
      const response = await api.get<ServerStat[]>(
        `/rcon/stats/server?${params.toString()}`
      );
      return response;
    },
    refetchInterval: 60000, // Refetch every 60 seconds
    staleTime: 30000, // Consider data stale after 30 seconds
  });
};

export const useSystemStats = (start?: string, end?: string) => {
  return useQuery<SystemStat[]>({
    queryKey: ["system-stats", start, end],
    queryFn: async () => {
      const params = new URLSearchParams();
      if (start) params.append("start", start);
      if (end) params.append("end", end);
      const response = await api.get<SystemStat[]>(
        `/rcon/stats/system?${params.toString()}`
      );
      return response;
    },
    refetchInterval: 60000, // Refetch every 60 seconds
    staleTime: 30000, // Consider data stale after 30 seconds
  });
};

export const usePlayerStats = (start?: string, end?: string) => {
  return useQuery<PlayerStat[]>({
    queryKey: ["player-stats", start, end],
    queryFn: async () => {
      const params = new URLSearchParams();
      if (start) params.append("start", start);
      if (end) params.append("end", end);
      const response = await api.get<PlayerStat[]>(
        `/rcon/stats/players?${params.toString()}`
      );
      return response;
    },
    refetchInterval: 60000, // Refetch every 60 seconds
    staleTime: 30000, // Consider data stale after 30 seconds
  });
};
