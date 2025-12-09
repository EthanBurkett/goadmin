import { useState, useMemo } from "react";
import { ProtectedRoute } from "@/components/ProtectedRoute";
import { DashboardLayout } from "@/components/DashboardLayout";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  useGroups,
  useCreateGroup,
  useUpdateGroup,
  useDeleteGroup,
  useInGamePlayers,
  useAssignPlayerToGroup,
  type Group,
} from "@/hooks/useGroups";
import { useCommands } from "@/hooks/useCommands";
import { Shield, Users, Plus, Trash2, Edit } from "lucide-react";
import { cn } from "@/lib/utils";

function Groups() {
  const { data: groups, isLoading: groupsLoading } = useGroups();
  const { data: players, isLoading: playersLoading } = useInGamePlayers();
  const { data: commands } = useCommands();
  const createGroup = useCreateGroup();
  const updateGroup = useUpdateGroup();
  const deleteGroup = useDeleteGroup();
  const assignPlayer = useAssignPlayerToGroup();

  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [selectedGroup, setSelectedGroup] = useState<Group | null>(null);

  const [formData, setFormData] = useState({
    name: "",
    power: 0,
    permissions: [] as string[],
    description: "",
  });

  // Extract unique permissions from all commands
  const availablePermissions = useMemo(() => {
    if (!commands) return [];

    const permissionSet = new Set<string>();
    commands.forEach((cmd) => {
      if (cmd.permissions) {
        try {
          const perms = JSON.parse(cmd.permissions);
          perms.forEach((p: string) => permissionSet.add(p));
        } catch {
          // Ignore parsing errors
        }
      }
    });

    // Add "all" as a special permission
    return ["all", ...Array.from(permissionSet).sort()];
  }, [commands]);

  const resetForm = () => {
    setFormData({
      name: "",
      power: 0,
      permissions: [],
      description: "",
    });
  };

  const togglePermission = (permission: string) => {
    setFormData((prev) => ({
      ...prev,
      permissions: prev.permissions.includes(permission)
        ? prev.permissions.filter((p) => p !== permission)
        : [...prev.permissions, permission],
    }));
  };

  const handleCreateGroup = async () => {
    try {
      await createGroup.mutateAsync({
        name: formData.name,
        power: formData.power,
        permissions: formData.permissions,
        description: formData.description,
      });
      setIsCreateDialogOpen(false);
      resetForm();
    } catch (error) {
      console.error("Failed to create group:", error);
    }
  };

  const handleEditGroup = async () => {
    if (!selectedGroup) return;

    try {
      await updateGroup.mutateAsync({
        id: selectedGroup.id,
        data: {
          name: formData.name || undefined,
          power: formData.power,
          permissions: formData.permissions,
          description: formData.description || undefined,
        },
      });
      setIsEditDialogOpen(false);
      setSelectedGroup(null);
      resetForm();
    } catch (error) {
      console.error("Failed to update group:", error);
    }
  };

  const handleDeleteGroup = async (id: number) => {
    if (confirm("Are you sure you want to delete this group?")) {
      try {
        await deleteGroup.mutateAsync(id);
      } catch (error) {
        console.error("Failed to delete group:", error);
      }
    }
  };

  const handleAssignPlayer = async (
    playerId: number,
    groupId: number | null
  ) => {
    try {
      await assignPlayer.mutateAsync({ playerId, groupId });
    } catch (error) {
      console.error("Failed to assign player:", error);
    }
  };

  const openEditDialog = (group: Group) => {
    setSelectedGroup(group);
    const permissions = group.permissions ? JSON.parse(group.permissions) : [];
    setFormData({
      name: group.name,
      power: group.power,
      permissions,
      description: group.description,
    });
    setIsEditDialogOpen(true);
  };

  const getPowerColor = (power: number) => {
    if (power >= 80) return "text-red-500";
    if (power >= 50) return "text-amber-500";
    if (power >= 20) return "text-blue-500";
    return "text-muted-foreground";
  };

  return (
    <ProtectedRoute requiredPermission="groups.manage">
      <DashboardLayout>
        <div className="p-8 space-y-6 bg-background min-h-screen">
          <div className="flex justify-between items-center">
            <div>
              <h1 className="text-4xl font-bold text-foreground mb-2">
                In-Game Groups
              </h1>
              <p className="text-muted-foreground">
                Manage in-game admin groups and permissions (B3-style)
              </p>
            </div>
            <Dialog
              open={isCreateDialogOpen}
              onOpenChange={setIsCreateDialogOpen}
            >
              <DialogTrigger asChild>
                <Button>
                  <Plus className="h-4 w-4 mr-2" />
                  Create Group
                </Button>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>Create New Group</DialogTitle>
                  <DialogDescription>
                    Create a new in-game admin group with custom permissions
                  </DialogDescription>
                </DialogHeader>
                <div className="space-y-4">
                  <div>
                    <Label htmlFor="name">Group Name</Label>
                    <Input
                      id="name"
                      value={formData.name}
                      onChange={(e) =>
                        setFormData({ ...formData, name: e.target.value })
                      }
                      placeholder="SuperAdmin"
                    />
                  </div>
                  <div>
                    <Label htmlFor="power">Power Level (0-100)</Label>
                    <Input
                      id="power"
                      type="number"
                      min="0"
                      max="100"
                      value={formData.power}
                      onChange={(e) =>
                        setFormData({
                          ...formData,
                          power: parseInt(e.target.value) || 0,
                        })
                      }
                    />
                  </div>
                  <div>
                    <Label>Permissions</Label>
                    <div className="border rounded-md p-4 max-h-60 overflow-y-auto">
                      {availablePermissions.length > 0 ? (
                        <div className="space-y-2">
                          {availablePermissions.map((perm) => (
                            <label
                              key={perm}
                              className="flex items-center space-x-2 cursor-pointer hover:bg-muted/50 p-2 rounded"
                            >
                              <input
                                type="checkbox"
                                checked={formData.permissions.includes(perm)}
                                onChange={() => togglePermission(perm)}
                                className="h-4 w-4 rounded border-gray-300"
                              />
                              <span className="text-sm">
                                {perm === "all" ? (
                                  <span className="font-semibold text-primary">
                                    {perm} (grants all permissions)
                                  </span>
                                ) : (
                                  perm
                                )}
                              </span>
                            </label>
                          ))}
                        </div>
                      ) : (
                        <div className="text-sm text-muted-foreground">
                          No commands found. Create commands to define
                          permissions.
                        </div>
                      )}
                    </div>
                    <div className="text-xs text-muted-foreground mt-2">
                      Selected:{" "}
                      {formData.permissions.length === 0
                        ? "None"
                        : formData.permissions.join(", ")}
                    </div>
                  </div>
                  <div>
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
                      placeholder="Full server control"
                    />
                  </div>
                  <Button onClick={handleCreateGroup} className="w-full">
                    Create Group
                  </Button>
                </div>
              </DialogContent>
            </Dialog>
          </div>

          {/* Groups Grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {groupsLoading ? (
              <div className="col-span-full text-center text-muted-foreground">
                Loading groups...
              </div>
            ) : groups && groups.length > 0 ? (
              groups.map((group) => (
                <Card key={group.id}>
                  <CardHeader>
                    <div className="flex justify-between items-start">
                      <div className="flex items-center space-x-2">
                        <Shield className="h-5 w-5 text-primary" />
                        <CardTitle>{group.name}</CardTitle>
                      </div>
                      <Badge
                        className={cn(getPowerColor(group.power), "bg-sidebar")}
                      >
                        {group.power}
                      </Badge>
                    </div>
                    <CardDescription>{group.description}</CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div>
                      <div className="text-sm text-muted-foreground mb-1">
                        Permissions
                      </div>
                      <div className="flex flex-wrap gap-1">
                        {group.permissions &&
                        JSON.parse(group.permissions).length > 0 ? (
                          JSON.parse(group.permissions).map((perm: string) => (
                            <Badge
                              key={perm}
                              variant="outline"
                              className="text-xs"
                            >
                              {perm}
                            </Badge>
                          ))
                        ) : (
                          <span className="text-xs text-muted-foreground">
                            No permissions
                          </span>
                        )}
                      </div>
                    </div>
                    <div className="flex space-x-2">
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => openEditDialog(group)}
                        className="flex-1"
                      >
                        <Edit className="h-3 w-3 mr-1" />
                        Edit
                      </Button>
                      <Button
                        variant="destructive"
                        size="sm"
                        onClick={() => handleDeleteGroup(group.id)}
                      >
                        <Trash2 className="h-3 w-3" />
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              ))
            ) : (
              <div className="col-span-full text-center text-muted-foreground">
                No groups found. Create your first group to get started.
              </div>
            )}
          </div>

          {/* Players Section */}
          <Card>
            <CardHeader>
              <div className="flex items-center space-x-2">
                <Users className="h-5 w-5 text-primary" />
                <CardTitle>In-Game Players</CardTitle>
              </div>
              <CardDescription>
                Assign players to groups by their PB GUID
              </CardDescription>
            </CardHeader>
            <CardContent>
              {playersLoading ? (
                <div className="text-center text-muted-foreground">
                  Loading players...
                </div>
              ) : players && players.length > 0 ? (
                <div className="space-y-2">
                  {players.map((player) => (
                    <div
                      key={player.id}
                      className="flex items-center justify-between p-3 border border-border rounded-lg"
                    >
                      <div>
                        <div className="font-medium">{player.name}</div>
                        <div className="text-sm text-muted-foreground font-mono">
                          {player.guid}
                        </div>
                      </div>
                      <div className="flex items-center space-x-2">
                        <Select
                          value={player.groupId?.toString() || "none"}
                          onValueChange={(value) =>
                            handleAssignPlayer(
                              player.id,
                              value === "none" ? null : parseInt(value)
                            )
                          }
                        >
                          <SelectTrigger className="w-[180px]">
                            <SelectValue placeholder="Select group" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value="none">No Group</SelectItem>
                            {groups?.map((group) => (
                              <SelectItem
                                key={group.id}
                                value={group.id.toString()}
                              >
                                {group.name} ({group.power})
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                        {player.group && (
                          <Badge
                            className={cn(
                              getPowerColor(player.group.power),
                              "bg-background/50"
                            )}
                          >
                            {player.group.name}
                          </Badge>
                        )}
                      </div>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="text-center text-muted-foreground">
                  No in-game players found. Players will appear here when they
                  join the server.
                </div>
              )}
            </CardContent>
          </Card>

          {/* Edit Dialog */}
          <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Edit Group</DialogTitle>
                <DialogDescription>
                  Update group settings and permissions
                </DialogDescription>
              </DialogHeader>
              <div className="space-y-4">
                <div>
                  <Label htmlFor="edit-name">Group Name</Label>
                  <Input
                    id="edit-name"
                    value={formData.name}
                    onChange={(e) =>
                      setFormData({ ...formData, name: e.target.value })
                    }
                  />
                </div>
                <div>
                  <Label htmlFor="edit-power">Power Level (0-100)</Label>
                  <Input
                    id="edit-power"
                    type="number"
                    min="0"
                    max="100"
                    value={formData.power}
                    onChange={(e) =>
                      setFormData({
                        ...formData,
                        power: parseInt(e.target.value) || 0,
                      })
                    }
                  />
                </div>
                <div>
                  <Label>Permissions</Label>
                  <div className="border rounded-md p-4 max-h-60 overflow-y-auto">
                    {availablePermissions.length > 0 ? (
                      <div className="space-y-2">
                        {availablePermissions.map((perm) => (
                          <label
                            key={perm}
                            className="flex items-center space-x-2 cursor-pointer hover:bg-muted/50 p-2 rounded"
                          >
                            <input
                              type="checkbox"
                              checked={formData.permissions.includes(perm)}
                              onChange={() => togglePermission(perm)}
                              className="h-4 w-4 rounded border-gray-300"
                            />
                            <span className="text-sm">
                              {perm === "all" ? (
                                <span className="font-semibold text-primary">
                                  {perm} (grants all permissions)
                                </span>
                              ) : (
                                perm
                              )}
                            </span>
                          </label>
                        ))}
                      </div>
                    ) : (
                      <div className="text-sm text-muted-foreground">
                        No commands found. Create commands to define
                        permissions.
                      </div>
                    )}
                  </div>
                  <div className="text-xs text-muted-foreground mt-2">
                    Selected:{" "}
                    {formData.permissions.length === 0
                      ? "None"
                      : formData.permissions.join(", ")}
                  </div>
                </div>
                <div>
                  <Label htmlFor="edit-description">Description</Label>
                  <Textarea
                    id="edit-description"
                    value={formData.description}
                    onChange={(e) =>
                      setFormData({ ...formData, description: e.target.value })
                    }
                  />
                </div>
                <Button onClick={handleEditGroup} className="w-full">
                  Update Group
                </Button>
              </div>
            </DialogContent>
          </Dialog>
        </div>
      </DashboardLayout>
    </ProtectedRoute>
  );
}

export default Groups;
