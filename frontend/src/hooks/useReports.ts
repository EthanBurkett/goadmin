import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import api from "@/lib/api";

export interface Report {
  id: number;
  reporterName: string;
  reporterGuid: string;
  reportedName: string;
  reportedGuid: string;
  reason: string;
  status: "pending" | "reviewed" | "actioned" | "dismissed";
  actionTaken?: string;
  reviewedByUserId?: number;
  reviewedBy?: {
    id: number;
    username: string;
  };
  createdAt: string;
  updatedAt: string;
}

export interface TempBan {
  id: number;
  playerName: string;
  playerGuid: string;
  reason: string;
  bannedByUser?: number;
  bannedBy?: {
    id: number;
    username: string;
  };
  expiresAt: string;
  active: boolean;
  createdAt: string;
  updatedAt: string;
}

export function useReports() {
  return useQuery({
    queryKey: ["reports"],
    queryFn: async () => {
      const response = await api.get<Report[]>("/reports");
      return response;
    },
  });
}

export function usePendingReports() {
  return useQuery({
    queryKey: ["reports", "pending"],
    queryFn: async () => {
      const response = await api.get<Report[]>("/reports/pending");
      return response;
    },
  });
}

export function useActionReport() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      id,
      action,
      duration,
      reason,
    }: {
      id: number;
      action: "dismiss" | "ban" | "tempban";
      duration?: number;
      reason: string;
    }) => {
      const response = await api.post(`/reports/${id}/action`, {
        action,
        duration,
        reason,
      });
      return response;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["reports"] });
      queryClient.invalidateQueries({ queryKey: ["tempbans"] });
    },
  });
}

export function useDeleteReport() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: number) => {
      const response = await api.delete(`/reports/${id}`);
      return response;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["reports"] });
    },
  });
}

export function useTempBans() {
  return useQuery({
    queryKey: ["tempbans"],
    queryFn: async () => {
      const response = await api.get<TempBan[]>("/tempbans");
      return response;
    },
  });
}

export function useActiveTempBans() {
  return useQuery({
    queryKey: ["tempbans", "active"],
    queryFn: async () => {
      const response = await api.get<TempBan[]>("/tempbans/active");
      return response;
    },
  });
}

export function useRevokeTempBan() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: number) => {
      const response = await api.post(`/tempbans/${id}/revoke`);
      return response;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["tempbans"] });
    },
  });
}
