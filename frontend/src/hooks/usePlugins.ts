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
