import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import api from "@/lib/api";

export interface Group {
  id: number;
  name: string;
  power: number;
  permissions: string; // JSON string
  description: string;
  createdAt: string;
  updatedAt: string;
}

export interface InGamePlayer {
  id: number;
  guid: string;
  name: string;
  groupId: number | null;
  group?: Group;
  enabled: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface CreateGroupRequest {
  name: string;
  power: number;
  permissions: string[];
  description: string;
}

export interface UpdateGroupRequest {
  name?: string;
  power?: number;
  permissions?: string[];
  description?: string;
}

export interface CreateInGamePlayerRequest {
  guid: string;
  name: string;
  groupId?: number | null;
}

export const useGroups = () => {
  return useQuery<Group[]>({
    queryKey: ["groups"],
    queryFn: async () => {
      const response = await api.get<Group[]>("/groups");
      return response;
    },
  });
};

export const useCreateGroup = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (data: CreateGroupRequest) => {
      return await api.post("/groups", data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["groups"] });
    },
  });
};

export const useUpdateGroup = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async ({
      id,
      data,
    }: {
      id: number;
      data: UpdateGroupRequest;
    }) => {
      return await api.put(`/groups/${id}`, data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["groups"] });
    },
  });
};

export const useDeleteGroup = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: number) => {
      return await api.delete(`/groups/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["groups"] });
      queryClient.invalidateQueries({ queryKey: ["in-game-players"] });
    },
  });
};

export const useInGamePlayers = () => {
  return useQuery<InGamePlayer[]>({
    queryKey: ["in-game-players"],
    queryFn: async () => {
      const response = await api.get<InGamePlayer[]>("/groups/players");
      return response;
    },
  });
};

export const useCreateInGamePlayer = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (data: CreateInGamePlayerRequest) => {
      return await api.post("/groups/players", data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["in-game-players"] });
    },
  });
};

export const useAssignPlayerToGroup = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async ({
      playerId,
      groupId,
    }: {
      playerId: number;
      groupId: number | null;
    }) => {
      return await api.put(`/groups/players/${playerId}/assign`, {
        playerId,
        groupId,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["in-game-players"] });
    },
  });
};

export const useRemovePlayerFromGroup = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (playerId: number) => {
      return await api.delete(`/groups/players/${playerId}/group`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["in-game-players"] });
    },
  });
};
