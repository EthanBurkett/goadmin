import { useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import {
  useWebhooks,
  useCreateWebhook,
  useUpdateWebhook,
  useDeleteWebhook,
  useTestWebhook,
  useWebhookDeliveries,
  WEBHOOK_EVENTS,
  type WebhookRequest,
  type Webhook,
} from "@/hooks/useWebhooks";
import {
  Loader2,
  Plus,
  Trash2,
  Edit,
  Send,
  CheckCircle,
  XCircle,
  Clock,
  ExternalLink,
} from "lucide-react";
import { toast } from "sonner";
import { DataTable } from "@/components/DataTable";
import type { ColumnDef } from "@tanstack/react-table";
import { formatDistanceToNow } from "date-fns";
import { ProtectedRoute } from "@/components/ProtectedRoute";

export default function WebhooksPage() {
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editingWebhook, setEditingWebhook] = useState<Webhook | null>(null);
  const [selectedWebhookId, setSelectedWebhookId] = useState<number | null>(
    null
  );
  const [deliveriesDialogOpen, setDeliveriesDialogOpen] = useState(false);

  const { data: webhooks, isLoading } = useWebhooks();
  const createWebhook = useCreateWebhook();
  const updateWebhook = useUpdateWebhook();
  const deleteWebhook = useDeleteWebhook();
  const testWebhook = useTestWebhook();
  const { data: deliveries } = useWebhookDeliveries(selectedWebhookId || 0);

  const [formData, setFormData] = useState<WebhookRequest>({
    name: "",
    url: "",
    secret: "",
    events: [],
    active: true,
    max_retries: 3,
    retry_delay: 60,
    timeout: 10,
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      if (editingWebhook) {
        await updateWebhook.mutateAsync({
          id: editingWebhook.id,
          webhook: formData,
        });
        toast.success("Webhook updated successfully");
      } else {
        await createWebhook.mutateAsync(formData);
        toast.success("Webhook created successfully");
      }
      setDialogOpen(false);
      resetForm();
    } catch {
      toast.error(
        editingWebhook ? "Failed to update webhook" : "Failed to create webhook"
      );
    }
  };

  const handleEdit = (webhook: Webhook) => {
    setEditingWebhook(webhook);
    setFormData({
      name: webhook.name,
      url: webhook.url,
      secret: webhook.secret,
      events: webhook.events,
      active: webhook.active,
      max_retries: webhook.max_retries,
      retry_delay: webhook.retry_delay,
      timeout: webhook.timeout,
    });
    setDialogOpen(true);
  };

  const handleDelete = async (id: number) => {
    if (!confirm("Are you sure you want to delete this webhook?")) return;

    try {
      await deleteWebhook.mutateAsync(id);
      toast.success("Webhook deleted successfully");
    } catch {
      toast.error("Failed to delete webhook");
    }
  };

  const handleTest = async (id: number) => {
    try {
      await testWebhook.mutateAsync(id);
      toast.success("Test webhook queued for delivery");
    } catch {
      toast.error("Failed to send test webhook");
    }
  };

  const resetForm = () => {
    setFormData({
      name: "",
      url: "",
      secret: "",
      events: [],
      active: true,
      max_retries: 3,
      retry_delay: 60,
      timeout: 10,
    });
    setEditingWebhook(null);
  };

  const toggleEvent = (event: string) => {
    setFormData((prev) => ({
      ...prev,
      events: prev.events.includes(event)
        ? prev.events.filter((e) => e !== event)
        : [...prev.events, event],
    }));
  };

  const columns: ColumnDef<Webhook>[] = [
    {
      accessorKey: "name",
      header: "Name",
      cell: ({ row }) => <div className="font-medium">{row.original.name}</div>,
    },
    {
      accessorKey: "url",
      header: "URL",
      cell: ({ row }) => (
        <div className="flex items-center gap-2">
          <span className="text-sm text-muted-foreground max-w-[300px] truncate">
            {row.original.url}
          </span>
          <ExternalLink className="h-3 w-3 text-muted-foreground" />
        </div>
      ),
    },
    {
      accessorKey: "events",
      header: "Events",
      cell: ({ row }) => (
        <div className="flex flex-wrap gap-1">
          {row.original.events.slice(0, 3).map((event) => (
            <Badge key={event} variant="secondary" className="text-xs">
              {event}
            </Badge>
          ))}
          {row.original.events.length > 3 && (
            <Badge variant="secondary" className="text-xs">
              +{row.original.events.length - 3}
            </Badge>
          )}
        </div>
      ),
    },
    {
      accessorKey: "active",
      header: "Status",
      cell: ({ row }) => (
        <Badge variant={row.original.active ? "default" : "secondary"}>
          {row.original.active ? "Active" : "Inactive"}
        </Badge>
      ),
    },
    {
      accessorKey: "created_at",
      header: "Created",
      cell: ({ row }) => (
        <span className="text-sm text-muted-foreground">
          {formatDistanceToNow(new Date(row.original.created_at), {
            addSuffix: true,
          })}
        </span>
      ),
    },
    {
      id: "actions",
      header: "Actions",
      cell: ({ row }) => (
        <div className="flex items-center gap-2">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => {
              setSelectedWebhookId(row.original.id);
              setDeliveriesDialogOpen(true);
            }}
          >
            <Clock className="h-4 w-4" />
          </Button>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => handleTest(row.original.id)}
            disabled={testWebhook.isPending}
          >
            <Send className="h-4 w-4" />
          </Button>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => handleEdit(row.original)}
          >
            <Edit className="h-4 w-4" />
          </Button>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => handleDelete(row.original.id)}
            disabled={deleteWebhook.isPending}
          >
            <Trash2 className="h-4 w-4" />
          </Button>
        </div>
      ),
    },
  ];

  return (
    <ProtectedRoute requiredPermission="webhooks.manage">
      <div className="space-y-6 bg-background min-h-screen">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-4xl font-bold text-foreground mb-2">
              Webhooks
            </h1>
            <p className="text-muted-foreground">
              Configure webhooks to receive real-time events from GoAdmin
            </p>
          </div>
          <Dialog
            open={dialogOpen}
            onOpenChange={(open) => {
              setDialogOpen(open);
              if (!open) resetForm();
            }}
          >
            <DialogTrigger asChild>
              <Button>
                <Plus className="h-4 w-4 mr-2" />
                Create Webhook
              </Button>
            </DialogTrigger>
            <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
              <DialogHeader>
                <DialogTitle>
                  {editingWebhook ? "Edit Webhook" : "Create Webhook"}
                </DialogTitle>
                <DialogDescription>
                  Configure a webhook endpoint to receive event notifications
                </DialogDescription>
              </DialogHeader>
              <form onSubmit={handleSubmit} className="space-y-6">
                <div className="space-y-2">
                  <Label htmlFor="name">Name</Label>
                  <Input
                    id="name"
                    value={formData.name}
                    onChange={(e) =>
                      setFormData({ ...formData, name: e.target.value })
                    }
                    placeholder="My Discord Webhook"
                    required
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="url">URL</Label>
                  <Input
                    id="url"
                    type="url"
                    value={formData.url}
                    onChange={(e) =>
                      setFormData({ ...formData, url: e.target.value })
                    }
                    placeholder="https://discord.com/api/webhooks/..."
                    required
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="secret">Secret</Label>
                  <Input
                    id="secret"
                    value={formData.secret}
                    onChange={(e) =>
                      setFormData({ ...formData, secret: e.target.value })
                    }
                    placeholder="your-secret-key"
                    required
                  />
                  <p className="text-xs text-muted-foreground">
                    Used for HMAC SHA256 signature verification
                  </p>
                </div>

                <div className="space-y-2">
                  <Label>Events</Label>
                  <div className="grid grid-cols-2 gap-2">
                    {WEBHOOK_EVENTS.map((event) => (
                      <div
                        key={event.value}
                        className="flex items-center space-x-2"
                      >
                        <input
                          type="checkbox"
                          id={event.value}
                          checked={formData.events.includes(event.value)}
                          onChange={() => toggleEvent(event.value)}
                          className="rounded border-gray-300"
                        />
                        <Label
                          htmlFor={event.value}
                          className="text-sm font-normal cursor-pointer"
                        >
                          {event.label}
                        </Label>
                      </div>
                    ))}
                  </div>
                </div>

                <div className="grid grid-cols-3 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="max_retries">Max Retries</Label>
                    <Input
                      id="max_retries"
                      type="number"
                      min="0"
                      max="10"
                      value={formData.max_retries}
                      onChange={(e) =>
                        setFormData({
                          ...formData,
                          max_retries: parseInt(e.target.value),
                        })
                      }
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="retry_delay">Retry Delay (s)</Label>
                    <Input
                      id="retry_delay"
                      type="number"
                      min="1"
                      value={formData.retry_delay}
                      onChange={(e) =>
                        setFormData({
                          ...formData,
                          retry_delay: parseInt(e.target.value),
                        })
                      }
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="timeout">Timeout (s)</Label>
                    <Input
                      id="timeout"
                      type="number"
                      min="1"
                      value={formData.timeout}
                      onChange={(e) =>
                        setFormData({
                          ...formData,
                          timeout: parseInt(e.target.value),
                        })
                      }
                    />
                  </div>
                </div>

                <div className="flex items-center space-x-2">
                  <Switch
                    id="active"
                    checked={formData.active}
                    onCheckedChange={(checked) =>
                      setFormData({ ...formData, active: checked })
                    }
                  />
                  <Label htmlFor="active">Active</Label>
                </div>

                <div className="flex justify-end gap-2">
                  <Button
                    type="button"
                    variant="outline"
                    onClick={() => {
                      setDialogOpen(false);
                      resetForm();
                    }}
                  >
                    Cancel
                  </Button>
                  <Button
                    type="submit"
                    disabled={
                      createWebhook.isPending || updateWebhook.isPending
                    }
                  >
                    {(createWebhook.isPending || updateWebhook.isPending) && (
                      <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                    )}
                    {editingWebhook ? "Update" : "Create"}
                  </Button>
                </div>
              </form>
            </DialogContent>
          </Dialog>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>Webhook Endpoints</CardTitle>
            <CardDescription>
              Manage webhook endpoints and view their delivery status
            </CardDescription>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="flex items-center justify-center py-8">
                <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
              </div>
            ) : webhooks && webhooks.length > 0 ? (
              <DataTable columns={columns} data={webhooks} />
            ) : (
              <div className="text-center py-8 text-muted-foreground">
                No webhooks configured. Create your first webhook to get
                started.
              </div>
            )}
          </CardContent>
        </Card>

        <Dialog
          open={deliveriesDialogOpen}
          onOpenChange={setDeliveriesDialogOpen}
        >
          <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto">
            <DialogHeader>
              <DialogTitle>Delivery History</DialogTitle>
              <DialogDescription>
                View webhook delivery attempts and their status
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4">
              {deliveries && deliveries.length > 0 ? (
                deliveries.map((delivery) => (
                  <Card key={delivery.id}>
                    <CardContent className="pt-6">
                      <div className="flex items-start justify-between">
                        <div className="space-y-2 flex-1">
                          <div className="flex items-center gap-2">
                            <Badge variant="outline">
                              {delivery.event_type}
                            </Badge>
                            {delivery.status === "delivered" && (
                              <CheckCircle className="h-4 w-4 text-green-500" />
                            )}
                            {delivery.status === "failed" && (
                              <XCircle className="h-4 w-4 text-red-500" />
                            )}
                            {delivery.status === "pending" && (
                              <Clock className="h-4 w-4 text-yellow-500" />
                            )}
                            <Badge
                              variant={
                                delivery.status === "delivered"
                                  ? "default"
                                  : delivery.status === "failed"
                                  ? "destructive"
                                  : "secondary"
                              }
                            >
                              {delivery.status}
                            </Badge>
                          </div>
                          <div className="text-sm text-muted-foreground">
                            Attempts: {delivery.attempt_count}
                            {delivery.response_code &&
                              ` • Response: ${delivery.response_code}`}
                          </div>
                          {delivery.error_message && (
                            <div className="text-sm text-red-500">
                              Error: {delivery.error_message}
                            </div>
                          )}
                          <div className="text-xs text-muted-foreground">
                            {formatDistanceToNow(
                              new Date(delivery.created_at),
                              { addSuffix: true }
                            )}
                          </div>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))
              ) : (
                <div className="text-center py-8 text-muted-foreground">
                  No delivery history available
                </div>
              )}
            </div>
          </DialogContent>
        </Dialog>
      </div>
    </ProtectedRoute>
  );
}
