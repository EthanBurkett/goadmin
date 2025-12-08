import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import api from "@/lib/api";

interface CommandResponse {
  response: string;
}

export interface CommandHistory {
  id: number;
  userId: number;
  command: string;
  response: string;
  success: boolean;
  createdAt: string;
}

export function useCommandHistory() {
  return useQuery({
    queryKey: ["command-history"],
    queryFn: () => api.get<CommandHistory[]>("/rcon/history"),
  });
}

export function useSendCommand() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (command: string) => {
      return await api.post<CommandResponse>("/rcon/command", { command });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["command-history"] });
    },
  });
}

export function useKickPlayer() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async ({
      playerId,
      reason,
    }: {
      playerId: string;
      reason?: string;
    }) => {
      return await api.post<CommandResponse>("/rcon/kick", {
        playerId,
        reason,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["command-history"] });
    },
  });
}

export function useBanPlayer() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async ({
      playerId,
      reason,
    }: {
      playerId: string;
      reason?: string;
    }) => {
      return await api.post<CommandResponse>("/rcon/ban", { playerId, reason });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["command-history"] });
    },
  });
}

export function useSayMessage() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (message: string) => {
      return await api.post<CommandResponse>("/rcon/say", { message });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["command-history"] });
    },
  });
}

export type { CommandResponse };
