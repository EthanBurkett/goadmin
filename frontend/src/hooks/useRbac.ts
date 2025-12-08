import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import api from "@/lib/api";

export interface Permission {
  id: number;
  name: string;
  description: string;
}

export interface Role {
  id: number;
  name: string;
  description: string;
  permissions?: Permission[];
}

export interface RbacUser {
  id: number;
  username: string;
  roles?: Role[];
}

export interface PendingUser {
  id: number;
  username: string;
  approved: boolean;
}

export function useRoles() {
  return useQuery({
    queryKey: ["roles"],
    queryFn: () => api.get<Role[]>("/rbac/roles"),
  });
}

export function usePermissions() {
  return useQuery({
    queryKey: ["permissions"],
    queryFn: () => api.get<Permission[]>("/rbac/permissions"),
  });
}

export function useRbacUsers() {
  return useQuery({
    queryKey: ["rbac-users"],
    queryFn: () => api.get<RbacUser[]>("/rbac/users"),
  });
}

export function useCreateRole() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: { name: string; description: string }) =>
      api.post("/rbac/roles", data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["roles"] });
    },
  });
}

export function useCreatePermission() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: { name: string; description: string }) =>
      api.post("/rbac/permissions", data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["permissions"] });
    },
  });
}

export function useDeleteRole() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => api.delete(`/rbac/roles/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["roles"] });
    },
  });
}

export function useDeletePermission() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => api.delete(`/rbac/permissions/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["permissions"] });
    },
  });
}

export function useAssignPermissionsToRole() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: { roleId: number; permissionId: number }) =>
      api.post(`/rbac/roles/${data.roleId}/permissions`, {
        permissionId: data.permissionId,
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["roles"] });
    },
  });
}

export function useRemovePermissionFromRole() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: { roleId: number; permissionId: number }) =>
      api.delete(`/rbac/roles/${data.roleId}/permissions/${data.permissionId}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["roles"] });
    },
  });
}

export function usePendingUsers() {
  return useQuery({
    queryKey: ["pending-users"],
    queryFn: () => api.get<PendingUser[]>("/rbac/users/pending"),
  });
}

export function useApproveUser() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: { userId: number; roleId: number }) =>
      api.post(`/rbac/users/${data.userId}/approve`, { roleId: data.roleId }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["pending-users"] });
      queryClient.invalidateQueries({ queryKey: ["rbac-users"] });
    },
  });
}

export function useDenyUser() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (userId: number) => api.post(`/rbac/users/${userId}/deny`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["pending-users"] });
    },
  });
}

export function useDeleteUser() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (userId: number) => api.delete(`/rbac/users/${userId}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["rbac-users"] });
    },
  });
}

export function useAssignRoleToUser() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: { userId: number; roleId: number }) =>
      api.post(`/rbac/users/${data.userId}/roles`, { roleId: data.roleId }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["rbac-users"] });
      queryClient.invalidateQueries({ queryKey: ["roles"] });
    },
  });
}

export function useRemoveRoleFromUser() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: { userId: number; roleId: number }) =>
      api.delete(`/rbac/users/${data.userId}/roles/${data.roleId}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["rbac-users"] });
      queryClient.invalidateQueries({ queryKey: ["roles"] });
    },
  });
}

// Hook to check permissions based on current user context
export function useHasPermission() {
  const { data: user } = useQuery({
    queryKey: ["auth", "me"],
    queryFn: async () => {
      try {
        const response = await api.get<{ user: RbacUser }>("/auth/me");
        return response.user;
      } catch {
        return null;
      }
    },
  });

  return {
    hasPermission: (permission: string): boolean => {
      if (!user?.roles) return false;

      // Check if any role has the permission
      return user.roles.some((role: Role) =>
        role.permissions?.some((p: Permission) => p.name === permission)
      );
    },
  };
}
