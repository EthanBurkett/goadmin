import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import api from "@/lib/api";

export interface Migration {
  version: string;
  name: string;
  description: string;
  applied_at?: string;
  rolled_back?: boolean;
}

export interface MigrationHistory {
  id: number;
  migration_version: string;
  operation: string;
  success: boolean;
  error_message?: string;
  duration_ms: number;
  executed_at: string;
}

export interface MigrationStatus {
  current_version: string;
  pending_migrations: Migration[];
  applied_migrations: Migration[];
  total_applied: number;
  total_pending: number;
}

export function useMigrations() {
  return useQuery<Migration[]>({
    queryKey: ["migrations"],
    queryFn: async () => {
      const response = await api.get<Migration[]>("/migrations");
      return response;
    },
  });
}

export function useMigrationStatus() {
  return useQuery<MigrationStatus>({
    queryKey: ["migrations", "status"],
    queryFn: async () => {
      const response = await api.get<MigrationStatus>("/migrations/status");
      return response;
    },
    refetchInterval: 5000, // Refresh every 5 seconds
  });
}

export function useCurrentVersion() {
  return useQuery<string>({
    queryKey: ["migrations", "current"],
    queryFn: async () => {
      const response = await api.get<{
        version: string;
      }>("/migrations/current");
      return response.version;
    },
  });
}

export function useApplyMigrations() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      const response = await api.post("/migrations/apply");
      return response;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["migrations"] });
      queryClient.invalidateQueries({ queryKey: ["migrations", "status"] });
      queryClient.invalidateQueries({ queryKey: ["migrations", "current"] });
    },
  });
}

export function useRollbackMigration() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      const response = await api.post("/migrations/rollback");
      return response;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["migrations"] });
      queryClient.invalidateQueries({ queryKey: ["migrations", "status"] });
      queryClient.invalidateQueries({ queryKey: ["migrations", "current"] });
    },
  });
}
