import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import api from "@/lib/api";

export interface CustomCommand {
  id: number;
  name: string;
  usage: string;
  description: string;
  rconCommand: string;
  minArgs: number;
  maxArgs: number;
  minPower: number;
  permissions: string; // JSON string
  requirementType: string; // "permission", "power", or "both"
  isBuiltIn: boolean;
  enabled: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface CreateCommandRequest {
  name: string;
  usage: string;
  description: string;
  rconCommand: string;
  minArgs: number;
  maxArgs: number;
  minPower: number;
  permissions: string[];
  requirementType: string;
}

export interface UpdateCommandRequest {
  name?: string;
  usage?: string;
  description?: string;
  rconCommand?: string;
  minArgs?: number;
  maxArgs?: number;
  minPower?: number;
  permissions?: string[];
  requirementType?: string;
  enabled?: boolean;
}

export const useCommands = () => {
  return useQuery<CustomCommand[]>({
    queryKey: ["commands"],
    queryFn: async () => {
      const response = await api.get<CustomCommand[]>("/commands");
      return response;
    },
  });
};

export const useCreateCommand = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (data: CreateCommandRequest) => {
      return await api.post("/commands", data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["commands"] });
    },
  });
};

export const useUpdateCommand = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async ({
      id,
      data,
    }: {
      id: number;
      data: UpdateCommandRequest;
    }) => {
      return await api.put(`/commands/${id}`, data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["commands"] });
    },
  });
};

export const useDeleteCommand = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: number) => {
      return await api.delete(`/commands/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["commands"] });
    },
  });
};
