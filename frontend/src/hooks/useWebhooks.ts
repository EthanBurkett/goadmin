import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import api from "@/lib/api";

export interface Webhook {
  id: number;
  name: string;
  url: string;
  secret: string;
  events: string[];
  active: boolean;
  max_retries: number;
  retry_delay: number;
  timeout: number;
  created_at: string;
  updated_at: string;
}

export interface WebhookDelivery {
  id: number;
  webhook_id: number;
  event_type: string;
  payload: string;
  status: "pending" | "delivered" | "failed";
  attempt_count: number;
  response_code: number | null;
  error_message: string | null;
  delivered_at: string | null;
  next_retry_at: string | null;
  created_at: string;
}

export interface WebhookRequest {
  name: string;
  url: string;
  secret: string;
  events: string[];
  active: boolean;
  max_retries?: number;
  retry_delay?: number;
  timeout?: number;
}

export function useWebhooks() {
  return useQuery<Webhook[]>({
    queryKey: ["webhooks"],
    queryFn: async () => {
      const response = await api.get<Webhook[]>("/webhooks");
      return response;
    },
  });
}

export function useWebhook(id: number) {
  return useQuery<Webhook>({
    queryKey: ["webhooks", id],
    queryFn: async () => {
      const response = await api.get<Webhook>(`/webhooks/${id}`);
      return response;
    },
    enabled: !!id,
  });
}

export function useWebhookDeliveries(id: number) {
  return useQuery<WebhookDelivery[]>({
    queryKey: ["webhooks", id, "deliveries"],
    queryFn: async () => {
      const response = await api.get<WebhookDelivery[]>(
        `/webhooks/${id}/deliveries`
      );
      return response;
    },
    enabled: !!id,
  });
}

export function useCreateWebhook() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (webhook: WebhookRequest) => {
      const response = await api.post("/webhooks", webhook);
      return response;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["webhooks"] });
    },
  });
}

export function useUpdateWebhook() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      id,
      webhook,
    }: {
      id: number;
      webhook: WebhookRequest;
    }) => {
      const response = await api.put(`/webhooks/${id}`, webhook);
      return response;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["webhooks"] });
    },
  });
}

export function useDeleteWebhook() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: number) => {
      await api.delete(`/webhooks/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["webhooks"] });
    },
  });
}

export function useTestWebhook() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: number) => {
      const response = await api.post(`/webhooks/${id}/test`);
      return response;
    },
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({
        queryKey: ["webhooks", id, "deliveries"],
      });
    },
  });
}

export const WEBHOOK_EVENTS = [
  { value: "player.banned", label: "Player Banned" },
  { value: "player.unbanned", label: "Player Unbanned" },
  { value: "player.kicked", label: "Player Kicked" },
  { value: "report.created", label: "Report Created" },
  { value: "report.actioned", label: "Report Actioned" },
  { value: "user.approved", label: "User Approved" },
  { value: "user.rejected", label: "User Rejected" },
  { value: "server.online", label: "Server Online" },
  { value: "server.offline", label: "Server Offline" },
  { value: "security.alert", label: "Security Alert" },
] as const;
