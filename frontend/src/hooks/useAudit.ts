import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import api from "@/lib/api";

export interface AuditLog {
  id: number;
  createdAt: string;
  userId: number | null;
  username: string;
  ipAddress: string;
  action: string;
  source: string;
  success: boolean;
  errorMessage?: string;
  targetType?: string;
  targetId?: string;
  targetName?: string;
  metadata?: string;
  result?: string;
  user?: {
    id: number;
    username: string;
  };
}

export interface AuditLogsResponse {
  logs: AuditLog[];
  total: number;
  limit: number;
  offset: number;
}

export interface AuditLogsFilters {
  user_id?: number;
  action?: string;
  source?: string;
  success?: boolean;
  target_type?: string;
  target_id?: string;
  start_date?: string;
  end_date?: string;
  limit?: number;
  offset?: number;
}

export const useAuditLogs = (filters?: AuditLogsFilters) => {
  return useQuery<AuditLogsResponse>({
    queryKey: ["audit", "logs", filters],
    queryFn: async () => {
      const params = new URLSearchParams();

      if (filters?.user_id)
        params.append("user_id", filters.user_id.toString());
      if (filters?.action) params.append("action", filters.action);
      if (filters?.source) params.append("source", filters.source);
      if (filters?.success !== undefined)
        params.append("success", filters.success.toString());
      if (filters?.target_type)
        params.append("target_type", filters.target_type);
      if (filters?.target_id) params.append("target_id", filters.target_id);
      if (filters?.start_date) params.append("start_date", filters.start_date);
      if (filters?.end_date) params.append("end_date", filters.end_date);
      if (filters?.limit) params.append("limit", filters.limit.toString());
      if (filters?.offset) params.append("offset", filters.offset.toString());

      const queryString = params.toString();
      const url = queryString ? `/audit/logs?${queryString}` : "/audit/logs";

      const response = await api.get<AuditLogsResponse>(url);
      return response;
    },
    staleTime: 30 * 1000, // 30 seconds
  });
};

export const useRecentAuditLogs = (limit = 100) => {
  return useQuery<{ logs: AuditLog[] }>({
    queryKey: ["audit", "logs", "recent", limit],
    queryFn: async () => {
      const response = await api.get<{ logs: AuditLog[] }>(
        `/audit/logs/recent?limit=${limit}`
      );
      return response;
    },
    staleTime: 30 * 1000,
  });
};

export const useAuditLogsByUser = (userId: number, limit = 100) => {
  return useQuery<{ logs: AuditLog[] }>({
    queryKey: ["audit", "logs", "user", userId, limit],
    queryFn: async () => {
      const response = await api.get<{ logs: AuditLog[] }>(
        `/audit/logs/user/${userId}?limit=${limit}`
      );
      return response;
    },
    enabled: !!userId,
    staleTime: 30 * 1000,
  });
};

export const useAuditLogsByAction = (action: string, limit = 100) => {
  return useQuery<{ logs: AuditLog[] }>({
    queryKey: ["audit", "logs", "action", action, limit],
    queryFn: async () => {
      const response = await api.get<{ logs: AuditLog[] }>(
        `/audit/logs/action/${action}?limit=${limit}`
      );
      return response;
    },
    enabled: !!action,
    staleTime: 30 * 1000,
  });
};

export interface AuditStats {
  total: number;
  archived: number;
  by_action: Array<{ Action: string; Count: number }>;
  by_source: Array<{ Source: string; Count: number }>;
  success_rate: number;
  oldest_log: string;
  newest_log: string;
}

export const useAuditStats = () => {
  return useQuery<AuditStats>({
    queryKey: ["audit", "stats"],
    queryFn: async () => {
      const response = await api.get<AuditStats>("/audit/stats");
      return response;
    },
    staleTime: 60 * 1000, // 1 minute
  });
};

export const useArchiveAuditLogs = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (retentionDays: number = 90) => {
      const response = await api.post<{
        archived: number;
        retention_days: number;
      }>(`/audit/archive?retention_days=${retentionDays}`, {});
      return response;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["audit"] });
    },
  });
};

export const usePurgeArchivedLogs = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async () => {
      const response = await api.post<{ purged: number }>("/audit/purge", {});
      return response;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["audit"] });
    },
  });
};
