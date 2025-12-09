import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import api from "@/lib/api";
import { toast } from "sonner";

export interface Server {
  id: number;
  name: string;
  host: string;
  port: number;
  rconPort: number;
  rconPassword?: string;
  gamesMpPath: string;
  isActive: boolean;
  isDefault: boolean;
  description?: string;
  region?: string;
  maxPlayers: number;
  createdAt: string;
  updatedAt: string;
}

export interface CreateServerData {
  name: string;
  host: string;
  port: number;
  rconPort: number;
  rconPassword: string;
  gamesMpPath: string;
  description?: string;
  region?: string;
  maxPlayers: number;
}

export interface UpdateServerData {
  name?: string;
  host?: string;
  port?: number;
  rconPort?: number;
  rconPassword?: string;
  gamesMpPath?: string;
  description?: string;
  region?: string;
  maxPlayers?: number;
}

export function useServers() {
  return useQuery({
    queryKey: ["servers"],
    queryFn: async () => {
      return await api.get<Server[]>("/servers");
    },
  });
}

export function useActiveServers() {
  return useQuery({
    queryKey: ["servers", "active"],
    queryFn: async () => {
      return await api.get<Server[]>("/servers/active");
    },
  });
}

export function useDefaultServer() {
  return useQuery({
    queryKey: ["servers", "default"],
    queryFn: async () => {
      return await api.get<Server>("/servers/default");
    },
  });
}

export function useServer(id: number) {
  return useQuery({
    queryKey: ["servers", id],
    queryFn: async () => {
      return await api.get<Server>(`/servers/${id}`);
    },
    enabled: !!id,
  });
}

export function useCreateServer() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: CreateServerData) => {
      return await api.post<Server>("/servers", data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["servers"] });
      toast.success("Server created successfully");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to create server");
    },
  });
}

export function useUpdateServer(id: number) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: UpdateServerData) => {
      return await api.put<Server>(`/servers/${id}`, data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["servers"] });
      queryClient.invalidateQueries({ queryKey: ["servers", id] });
      toast.success("Server updated successfully");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to update server");
    },
  });
}

export function useDeleteServer(id: number) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      return await api.delete<void>(`/servers/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["servers"] });
      toast.success("Server deleted successfully");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to delete server");
    },
  });
}

export function useSetDefaultServer(id: number) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      return await api.post<{ message: string }>(`/servers/${id}/default`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["servers"] });
      toast.success("Default server updated");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to set default server");
    },
  });
}

export function useActivateServer(id: number) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      return await api.post<{ message: string }>(`/servers/${id}/activate`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["servers"] });
      toast.success("Server activated");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to activate server");
    },
  });
}

export function useDeactivateServer(id: number) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      return await api.post<{ message: string }>(`/servers/${id}/deactivate`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["servers"] });
      toast.success("Server deactivated");
    },
    onError: (error: Error) => {
      toast.error(error.message || "Failed to deactivate server");
    },
  });
}
