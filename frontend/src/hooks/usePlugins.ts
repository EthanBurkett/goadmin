import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import api from "@/lib/api";
import { toast } from "sonner";

export interface PluginStatus {
  id: string;
  name: string;
  version: string;
  state: "loaded" | "started" | "stopped" | "error";
  enabled: boolean;
  loadedAt: string;
  error?: string;
}

export interface PluginMetrics {
  PluginID: string;
  MemoryUsageMB: number;
  GoroutineCount: number;
  LastChecked: string;
  ViolationCount: number;
  Throttled: boolean;
}

export interface PluginDependencyTree {
  plugin_id: string;
  dependency_tree: Record<string, string[]>;
}

// Get all plugins
export function usePlugins() {
  return useQuery<PluginStatus[]>({
    queryKey: ["plugins"],
    queryFn: async () => {
      const response = await api.get<{
        plugins: PluginStatus[];
      }>("/plugins");
      return response.plugins || [];
    },
  });
}

// Get single plugin status
export function usePlugin(pluginId: string) {
  return useQuery<PluginStatus>({
    queryKey: ["plugins", pluginId],
    queryFn: async () => {
      const response = await api.get<PluginStatus>(`/plugins/${pluginId}`);
      return response;
    },
    enabled: !!pluginId,
  });
}

// Start plugin
export function useStartPlugin() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (pluginId: string) => {
      const response = await api.post<{ message: string }>(
        `/plugins/${pluginId}/start`
      );
      return response;
    },
    onSuccess: (_, pluginId) => {
      queryClient.invalidateQueries({ queryKey: ["plugins"] });
      queryClient.invalidateQueries({ queryKey: ["plugins", pluginId] });
      toast.success("Plugin started successfully");
    },
    onError: (error: unknown) => {
      toast.error(
        error instanceof Error ? error.message : "Failed to start plugin"
      );
    },
  });
}

// Stop plugin
export function useStopPlugin() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (pluginId: string) => {
      const response = await api.post<{ message: string }>(
        `/plugins/${pluginId}/stop`
      );
      return response;
    },
    onSuccess: (_, pluginId) => {
      queryClient.invalidateQueries({ queryKey: ["plugins"] });
      queryClient.invalidateQueries({ queryKey: ["plugins", pluginId] });
      toast.success("Plugin stopped successfully");
    },
    onError: (error: unknown) => {
      toast.error(
        error instanceof Error ? error.message : "Failed to stop plugin"
      );
    },
  });
}

// Reload plugin
export function useReloadPlugin() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (pluginId: string) => {
      const response = await api.post<{ message: string }>(
        `/plugins/${pluginId}/reload`
      );
      return response;
    },
    onSuccess: (_, pluginId) => {
      queryClient.invalidateQueries({ queryKey: ["plugins"] });
      queryClient.invalidateQueries({ queryKey: ["plugins", pluginId] });
      toast.success("Plugin reloaded successfully");
    },
    onError: (error: unknown) => {
      toast.error(
        error instanceof Error ? error.message : "Failed to reload plugin"
      );
    },
  });
}

// Hot-reload plugin
export function useHotReloadPlugin() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (pluginId: string) => {
      const response = await api.post<{ message: string }>(
        `/plugins/${pluginId}/hot-reload`
      );
      return response;
    },
    onSuccess: (_, pluginId) => {
      queryClient.invalidateQueries({ queryKey: ["plugins"] });
      queryClient.invalidateQueries({ queryKey: ["plugins", pluginId] });
      queryClient.invalidateQueries({ queryKey: ["plugin-metrics", pluginId] });
      toast.success("Plugin hot-reloaded successfully");
    },
    onError: (error: unknown) => {
      toast.error(
        error instanceof Error ? error.message : "Failed to hot-reload plugin"
      );
    },
  });
}

// Get plugin metrics
export function usePluginMetrics(pluginId: string) {
  return useQuery<PluginMetrics>({
    queryKey: ["plugin-metrics", pluginId],
    queryFn: async () => {
      const response = await api.get<PluginMetrics>(
        `/plugins/${pluginId}/metrics`
      );
      return response;
    },
    enabled: !!pluginId,
    refetchInterval: 30000, // Refetch every 30 seconds
  });
}

// Get all plugin metrics
export function useAllPluginMetrics() {
  return useQuery<Record<string, PluginMetrics>>({
    queryKey: ["plugin-metrics", "all"],
    queryFn: async () => {
      const response = await api.get<{
        metrics: Record<string, PluginMetrics>;
      }>("/plugins/metrics/all");
      return response.metrics || {};
    },
    refetchInterval: 30000, // Refetch every 30 seconds
  });
}

// Get plugin dependencies
export function usePluginDependencies(pluginId: string) {
  return useQuery<PluginDependencyTree>({
    queryKey: ["plugin-dependencies", pluginId],
    queryFn: async () => {
      const response = await api.get<PluginDependencyTree>(
        `/plugins/${pluginId}/dependencies`
      );
      return response;
    },
    enabled: !!pluginId,
  });
}
