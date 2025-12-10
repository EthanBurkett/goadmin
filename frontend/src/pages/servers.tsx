import { useState } from "react";
import { ProtectedRoute } from "@/components/ProtectedRoute";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarProvider,
} from "@/components/ui/sidebar";
import {
  useServers,
  useCreateServer,
  useUpdateServer,
  useDeleteServer,
  useSetDefaultServer,
  useActivateServer,
  useDeactivateServer,
  type CreateServerData,
  type UpdateServerData,
  type Server,
} from "@/hooks/useServers";
import {
  Server as ServerIcon,
  Star,
  Plus,
  Pencil,
  Trash2,
  Power,
  PowerOff,
  Shield,
  LogOut,
  Database,
  Plug,
  Activity,
} from "lucide-react";
import { ServerSelector } from "@/components/ServerSelector";
import { useAuthContext } from "@/hooks/useAuthContext";
import { ServerProvider } from "@/providers/ServerProvider";
import { useNavigate } from "react-router-dom";

function Servers() {
  const { data: servers, isLoading } = useServers();
  const createServerMutation = useCreateServer();
  const { user, logout } = useAuthContext();
  const navigate = useNavigate();

  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [editingServer, setEditingServer] = useState<Server | null>(null);
  const [selectedServerId, setSelectedServerId] = useState<number>(0);

  const [formData, setFormData] = useState<CreateServerData>({
    name: "",
    host: "",
    port: 28960,
    rconPort: 28960,
    rconPassword: "",
    gamesMpPath: "",
    description: "",
    region: "",
    maxPlayers: 64,
  });

  const updateServerMutation = useUpdateServer(editingServer?.id || 0);
  const deleteServerMutation = useDeleteServer(selectedServerId);
  const setDefaultServerMutation = useSetDefaultServer(selectedServerId);
  const activateServerMutation = useActivateServer(selectedServerId);
  const deactivateServerMutation = useDeactivateServer(selectedServerId);

  const handleCreateSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await createServerMutation.mutateAsync(formData);
    setIsCreateDialogOpen(false);
    resetForm();
  };

  const handleEditSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editingServer) return;

    const updateData: UpdateServerData = {
      name: formData.name,
      host: formData.host,
      port: formData.port,
      rconPort: formData.rconPort,
      rconPassword: formData.rconPassword || undefined,
      gamesMpPath: formData.gamesMpPath,
      description: formData.description,
      region: formData.region,
      maxPlayers: formData.maxPlayers,
    };

    await updateServerMutation.mutateAsync(updateData);
    setIsEditDialogOpen(false);
    setEditingServer(null);
    resetForm();
  };

  const resetForm = () => {
    setFormData({
      name: "",
      host: "",
      port: 28960,
      rconPort: 28960,
      rconPassword: "",
      gamesMpPath: "",
      description: "",
      region: "",
      maxPlayers: 64,
    });
  };

  const handleEdit = (server: Server) => {
    setEditingServer(server);
    setFormData({
      name: server.name,
      host: server.host,
      port: server.port,
      rconPort: server.rconPort,
      rconPassword: "",
      gamesMpPath: server.gamesMpPath,
      description: server.description || "",
      region: server.region || "",
      maxPlayers: server.maxPlayers,
    });
    setIsEditDialogOpen(true);
  };

  const handleDelete = (server: Server) => {
    if (confirm("Are you sure you want to delete this server?")) {
      setSelectedServerId(server.id);
      deleteServerMutation.mutate();
    }
  };

  const handleSetDefault = (server: Server) => {
    setSelectedServerId(server.id);
    setDefaultServerMutation.mutate();
  };

  const handleToggleActive = (server: Server) => {
    setSelectedServerId(server.id);
    if (server.isActive) {
      deactivateServerMutation.mutate();
    } else {
      activateServerMutation.mutate();
    }
  };

  const handleLogout = async () => {
    await logout();
  };

  return (
    <ProtectedRoute requiredPermission="servers.manage">
      <ServerProvider disableRedirect>
        <SidebarProvider>
          <div className="flex min-h-screen w-full bg-background">
            <Sidebar className="border-border">
              <SidebarHeader className="border-b border-border p-4 space-y-3">
                <div className="flex items-center space-x-2">
                  <Shield className="h-6 w-6 text-primary" />
                  <span className="text-lg font-bold text-foreground">
                    GoAdmin
                  </span>
                </div>
                <ServerSelector />
              </SidebarHeader>
              <SidebarContent>
                <div className="p-4 space-y-1">
                  <div className="text-xs font-semibold text-muted-foreground mb-2 px-3">
                    GLOBAL
                  </div>
                  <Button
                    variant="ghost"
                    className="w-full justify-start bg-accent text-accent-foreground"
                    onClick={() => navigate("/servers")}
                  >
                    <Database className="mr-2 h-4 w-4" />
                    Servers
                  </Button>
                  <Button
                    variant="ghost"
                    className="w-full justify-start"
                    onClick={() => navigate("/plugins")}
                  >
                    <Plug className="mr-2 h-4 w-4" />
                    Plugins
                  </Button>
                  <Button
                    variant="ghost"
                    className="w-full justify-start"
                    onClick={() => navigate("/metrics")}
                  >
                    <Activity className="mr-2 h-4 w-4" />
                    Metrics
                  </Button>
                </div>
              </SidebarContent>
              <SidebarFooter className="border-t border-border p-4">
                <div className="space-y-2">
                  <div className="text-sm text-muted-foreground">
                    Logged in as{" "}
                    <span className="text-foreground font-semibold">
                      {user?.username}
                    </span>
                  </div>
                  <Button
                    variant="outline"
                    className="w-full justify-start"
                    onClick={handleLogout}
                  >
                    <LogOut className="mr-2 h-4 w-4" />
                    Sign Out
                  </Button>
                </div>
              </SidebarFooter>
            </Sidebar>
            <main className="flex-1 overflow-auto">
              <div className="p-8">
                <div className="space-y-6">
                  <div className="flex justify-between items-center">
                    <div>
                      <h1 className="text-4xl font-bold text-foreground mb-2">
                        Server Management
                      </h1>
                      <p className="text-muted-foreground">
                        Manage your game server instances
                      </p>
                    </div>
                    <Dialog
                      open={isCreateDialogOpen}
                      onOpenChange={setIsCreateDialogOpen}
                    >
                      <DialogTrigger asChild>
                        <Button>
                          <Plus className="h-4 w-4 mr-2" />
                          Add Server
                        </Button>
                      </DialogTrigger>
                      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
                        <DialogHeader>
                          <DialogTitle>Add New Server</DialogTitle>
                          <DialogDescription>
                            Configure a new game server instance
                          </DialogDescription>
                        </DialogHeader>
                        <form onSubmit={handleCreateSubmit}>
                          <div className="grid gap-4 py-4">
                            <div className="grid grid-cols-2 gap-4">
                              <div className="space-y-2">
                                <Label htmlFor="name">Server Name *</Label>
                                <Input
                                  id="name"
                                  required
                                  value={formData.name}
                                  onChange={(e) =>
                                    setFormData({
                                      ...formData,
                                      name: e.target.value,
                                    })
                                  }
                                  placeholder="My CoD4 Server"
                                />
                              </div>
                              <div className="space-y-2">
                                <Label htmlFor="host">Host *</Label>
                                <Input
                                  id="host"
                                  required
                                  value={formData.host}
                                  onChange={(e) =>
                                    setFormData({
                                      ...formData,
                                      host: e.target.value,
                                    })
                                  }
                                  placeholder="127.0.0.1"
                                />
                              </div>
                            </div>

                            <div className="grid grid-cols-3 gap-4">
                              <div className="space-y-2">
                                <Label htmlFor="port">Game Port *</Label>
                                <Input
                                  id="port"
                                  type="number"
                                  required
                                  value={formData.port}
                                  onChange={(e) =>
                                    setFormData({
                                      ...formData,
                                      port: parseInt(e.target.value),
                                    })
                                  }
                                />
                              </div>
                              <div className="space-y-2">
                                <Label htmlFor="rconPort">RCON Port *</Label>
                                <Input
                                  id="rconPort"
                                  type="number"
                                  required
                                  value={formData.rconPort}
                                  onChange={(e) =>
                                    setFormData({
                                      ...formData,
                                      rconPort: parseInt(e.target.value),
                                    })
                                  }
                                />
                              </div>
                              <div className="space-y-2">
                                <Label htmlFor="maxPlayers">
                                  Max Players *
                                </Label>
                                <Input
                                  id="maxPlayers"
                                  type="number"
                                  required
                                  value={formData.maxPlayers}
                                  onChange={(e) =>
                                    setFormData({
                                      ...formData,
                                      maxPlayers: parseInt(e.target.value),
                                    })
                                  }
                                />
                              </div>
                            </div>

                            <div className="space-y-2">
                              <Label htmlFor="rconPassword">
                                RCON Password *
                              </Label>
                              <Input
                                id="rconPassword"
                                type="password"
                                required
                                value={formData.rconPassword}
                                onChange={(e) =>
                                  setFormData({
                                    ...formData,
                                    rconPassword: e.target.value,
                                  })
                                }
                                placeholder="Enter RCON password"
                              />
                            </div>

                            <div className="space-y-2">
                              <Label htmlFor="gamesMpPath">
                                games_mp.log Path *
                              </Label>
                              <Input
                                id="gamesMpPath"
                                required
                                value={formData.gamesMpPath}
                                onChange={(e) =>
                                  setFormData({
                                    ...formData,
                                    gamesMpPath: e.target.value,
                                  })
                                }
                                placeholder="/path/to/games_mp.log"
                              />
                            </div>

                            <div className="space-y-2">
                              <Label htmlFor="region">Region</Label>
                              <Input
                                id="region"
                                value={formData.region}
                                onChange={(e) =>
                                  setFormData({
                                    ...formData,
                                    region: e.target.value,
                                  })
                                }
                                placeholder="US East, EU West, etc."
                              />
                            </div>

                            <div className="space-y-2">
                              <Label htmlFor="description">Description</Label>
                              <Textarea
                                id="description"
                                value={formData.description}
                                onChange={(e) =>
                                  setFormData({
                                    ...formData,
                                    description: e.target.value,
                                  })
                                }
                                placeholder="Optional server description"
                              />
                            </div>
                          </div>
                          <DialogFooter>
                            <Button
                              type="submit"
                              disabled={createServerMutation.isPending}
                            >
                              {createServerMutation.isPending
                                ? "Creating..."
                                : "Create Server"}
                            </Button>
                          </DialogFooter>
                        </form>
                      </DialogContent>
                    </Dialog>
                  </div>

                  {/* Edit Server Dialog */}
                  <Dialog
                    open={isEditDialogOpen}
                    onOpenChange={setIsEditDialogOpen}
                  >
                    <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
                      <DialogHeader>
                        <DialogTitle>Edit Server</DialogTitle>
                        <DialogDescription>
                          Update server configuration
                        </DialogDescription>
                      </DialogHeader>
                      <form onSubmit={handleEditSubmit}>
                        <div className="grid gap-4 py-4">
                          <div className="grid grid-cols-2 gap-4">
                            <div className="space-y-2">
                              <Label htmlFor="edit-name">Server Name *</Label>
                              <Input
                                id="edit-name"
                                required
                                value={formData.name}
                                onChange={(e) =>
                                  setFormData({
                                    ...formData,
                                    name: e.target.value,
                                  })
                                }
                                placeholder="My CoD4 Server"
                              />
                            </div>
                            <div className="space-y-2">
                              <Label htmlFor="edit-host">Host *</Label>
                              <Input
                                id="edit-host"
                                required
                                value={formData.host}
                                onChange={(e) =>
                                  setFormData({
                                    ...formData,
                                    host: e.target.value,
                                  })
                                }
                                placeholder="127.0.0.1"
                              />
                            </div>
                          </div>

                          <div className="grid grid-cols-3 gap-4">
                            <div className="space-y-2">
                              <Label htmlFor="edit-port">Game Port *</Label>
                              <Input
                                id="edit-port"
                                type="number"
                                required
                                value={formData.port}
                                onChange={(e) =>
                                  setFormData({
                                    ...formData,
                                    port: parseInt(e.target.value),
                                  })
                                }
                              />
                            </div>
                            <div className="space-y-2">
                              <Label htmlFor="edit-rconPort">RCON Port *</Label>
                              <Input
                                id="edit-rconPort"
                                type="number"
                                required
                                value={formData.rconPort}
                                onChange={(e) =>
                                  setFormData({
                                    ...formData,
                                    rconPort: parseInt(e.target.value),
                                  })
                                }
                              />
                            </div>
                            <div className="space-y-2">
                              <Label htmlFor="edit-maxPlayers">
                                Max Players *
                              </Label>
                              <Input
                                id="edit-maxPlayers"
                                type="number"
                                required
                                value={formData.maxPlayers}
                                onChange={(e) =>
                                  setFormData({
                                    ...formData,
                                    maxPlayers: parseInt(e.target.value),
                                  })
                                }
                              />
                            </div>
                          </div>

                          <div className="space-y-2">
                            <Label htmlFor="edit-rconPassword">
                              RCON Password
                            </Label>
                            <Input
                              id="edit-rconPassword"
                              type="password"
                              value={formData.rconPassword}
                              onChange={(e) =>
                                setFormData({
                                  ...formData,
                                  rconPassword: e.target.value,
                                })
                              }
                              placeholder="Leave blank to keep current password"
                            />
                          </div>

                          <div className="space-y-2">
                            <Label htmlFor="edit-gamesMpPath">
                              games_mp.log Path *
                            </Label>
                            <Input
                              id="edit-gamesMpPath"
                              required
                              value={formData.gamesMpPath}
                              onChange={(e) =>
                                setFormData({
                                  ...formData,
                                  gamesMpPath: e.target.value,
                                })
                              }
                              placeholder="/path/to/games_mp.log"
                            />
                          </div>

                          <div className="space-y-2">
                            <Label htmlFor="edit-region">Region</Label>
                            <Input
                              id="edit-region"
                              value={formData.region}
                              onChange={(e) =>
                                setFormData({
                                  ...formData,
                                  region: e.target.value,
                                })
                              }
                              placeholder="US East, EU West, etc."
                            />
                          </div>

                          <div className="space-y-2">
                            <Label htmlFor="edit-description">
                              Description
                            </Label>
                            <Textarea
                              id="edit-description"
                              value={formData.description}
                              onChange={(e) =>
                                setFormData({
                                  ...formData,
                                  description: e.target.value,
                                })
                              }
                              placeholder="Optional server description"
                            />
                          </div>
                        </div>
                        <DialogFooter>
                          <Button
                            type="submit"
                            disabled={updateServerMutation.isPending}
                          >
                            {updateServerMutation.isPending
                              ? "Updating..."
                              : "Update Server"}
                          </Button>
                        </DialogFooter>
                      </form>
                    </DialogContent>
                  </Dialog>

                  <Card>
                    <CardHeader>
                      <CardTitle>Server Instances</CardTitle>
                      <CardDescription>
                        Manage and configure your game server instances
                      </CardDescription>
                    </CardHeader>
                    <CardContent>
                      {isLoading ? (
                        <div className="text-center py-8 text-muted-foreground">
                          Loading servers...
                        </div>
                      ) : !servers || servers.length === 0 ? (
                        <div className="text-center py-8 text-muted-foreground">
                          No servers configured. Add your first server to get
                          started.
                        </div>
                      ) : (
                        <Table>
                          <TableHeader>
                            <TableRow>
                              <TableHead>Name</TableHead>
                              <TableHead>Host:Port</TableHead>
                              <TableHead>Region</TableHead>
                              <TableHead>Max Players</TableHead>
                              <TableHead>Status</TableHead>
                              <TableHead>Actions</TableHead>
                            </TableRow>
                          </TableHeader>
                          <TableBody>
                            {servers.map((server) => (
                              <TableRow key={server.id}>
                                <TableCell className="font-medium">
                                  <div className="flex items-center gap-2">
                                    <ServerIcon className="h-4 w-4" />
                                    {server.name}
                                    {server.isDefault && (
                                      <Star className="h-3 w-3 fill-yellow-500 text-yellow-500" />
                                    )}
                                  </div>
                                </TableCell>
                                <TableCell className="font-mono text-sm">
                                  {server.host}:{server.port}
                                </TableCell>
                                <TableCell>{server.region || "—"}</TableCell>
                                <TableCell>{server.maxPlayers}</TableCell>
                                <TableCell>
                                  <Badge
                                    variant={
                                      server.isActive ? "default" : "secondary"
                                    }
                                  >
                                    {server.isActive ? "Active" : "Inactive"}
                                  </Badge>
                                </TableCell>
                                <TableCell>
                                  <div className="flex items-center gap-2">
                                    {!server.isDefault && (
                                      <Button
                                        size="sm"
                                        variant="ghost"
                                        onClick={() => handleSetDefault(server)}
                                        title="Set as default"
                                      >
                                        <Star className="h-4 w-4" />
                                      </Button>
                                    )}
                                    <Button
                                      size="sm"
                                      variant="ghost"
                                      onClick={() => handleToggleActive(server)}
                                      title={
                                        server.isActive
                                          ? "Deactivate"
                                          : "Activate"
                                      }
                                    >
                                      {server.isActive ? (
                                        <PowerOff className="h-4 w-4" />
                                      ) : (
                                        <Power className="h-4 w-4" />
                                      )}
                                    </Button>
                                    <Button
                                      size="sm"
                                      variant="ghost"
                                      onClick={() => handleEdit(server)}
                                    >
                                      <Pencil className="h-4 w-4" />
                                    </Button>
                                    {!server.isDefault && (
                                      <Button
                                        size="sm"
                                        variant="ghost"
                                        onClick={() => handleDelete(server)}
                                      >
                                        <Trash2 className="h-4 w-4" />
                                      </Button>
                                    )}
                                  </div>
                                </TableCell>
                              </TableRow>
                            ))}
                          </TableBody>
                        </Table>
                      )}
                    </CardContent>
                  </Card>
                </div>
              </div>
            </main>
          </div>
        </SidebarProvider>
      </ServerProvider>
    </ProtectedRoute>
  );
}

export default Servers;
