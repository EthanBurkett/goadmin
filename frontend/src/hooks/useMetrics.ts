import { useQuery } from "@tanstack/react-query";
import api from "@/lib/api";

export interface SystemMetrics {
  db_connections: number;
  db_idle_conns: number;
  db_open_conns: number;
  db_wait_count: number;
  db_wait_duration_ms: number;
  db_max_idle_closed: number;
  db_max_life_closed: number;
  total_audit_logs: number;
  archived_audit_logs: number;
  audit_success_rate: number;
  total_users: number;
  active_users: number;
  pending_users: number;
  total_reports: number;
  pending_reports: number;
  total_bans: number;
  active_bans: number;
  total_commands: number;
  custom_commands: number;
  plugin_commands: number;
  cache_size: number;
  uptime_seconds: number;
}

export const useSystemMetrics = () => {
  return useQuery<SystemMetrics>({
    queryKey: ["metrics", "system"],
    queryFn: async () => {
      const response = await api.get<SystemMetrics>("/metrics/json");
      return response;
    },
    refetchInterval: 30 * 1000, // Refresh every 30 seconds
    staleTime: 15 * 1000,
  });
};
