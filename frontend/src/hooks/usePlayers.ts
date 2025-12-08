import { useQuery } from "@tanstack/react-query";
import {
  type OfflinePlayer,
  type OnlinePlayer,
  type OnlinePlayerDump,
  type Player,
} from "@/types/players";
import api from "@/lib/api";

export const usePlayers = (refetchInterval?: number) => {
  return useQuery<OnlinePlayer[]>({
    queryKey: ["players"],
    queryFn: async () => {
      const response = await api.get<Player[]>("/players");
      return response;
    },
    refetchInterval,
  });
};

export const usePlayer = (playerID: string) => {
  return useQuery<OnlinePlayerDump | OfflinePlayer>({
    queryKey: ["player", playerID],
    queryFn: async () => {
      const response = await api.get<OnlinePlayerDump | OfflinePlayer>(
        `/players/${playerID}`
      );
      return response;
    },
    enabled: !!playerID,
  });
};
