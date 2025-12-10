import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import api from "../lib/api";

export interface ShutdownInfo {
  command: string;
  reason: string;
  disabledAt: string;
  disabledBy: string;
  userId?: number;
  reenableAt: string;
  autoRenable: boolean;
}

export interface DisabledCommandsResponse {
  [command: string]: ShutdownInfo;
}

export const useDisabledCommands = () => {
  return useQuery({
    queryKey: ["emergency", "disabled"],
    queryFn: () => api.get<DisabledCommandsResponse>("/emergency/disabled"),
    refetchInterval: 10000, // Refresh every 10 seconds
  });
};

export const useReenableCommand = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (command: string) =>
      api.post(`/emergency/reenable/${command}`, {}),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["emergency", "disabled"] });
    },
  });
};
