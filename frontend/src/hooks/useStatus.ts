import { useQuery } from "@tanstack/react-query";
import api from "@/lib/api";

interface StatusResponse {
  hostname: string;
  address: string;
  map: string;
  gametype: string;
  players: Array<{
    id: number;
    score: number;
    ping: number;
    name: string;
    strippedName: string;
  }>;
}

export function useStatus() {
  return useQuery<StatusResponse>({
    queryKey: ["status"],
    queryFn: () => api.get<StatusResponse>("/status"),
    refetchInterval: 5000,
  });
}

export type { StatusResponse };
