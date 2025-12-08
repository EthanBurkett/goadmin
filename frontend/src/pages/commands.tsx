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
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  useCommands,
  useCreateCommand,
  useUpdateCommand,
  useDeleteCommand,
  type CustomCommand,
} from "@/hooks/useCommands";
import {
  Terminal,
  Plus,
  Trash2,
  Edit,
  Power,
  Check,
  X,
  Search,
} from "lucide-react";

function Commands() {
  const { data: commands, isLoading } = useCommands();
  const createCommand = useCreateCommand();
  const updateCommand = useUpdateCommand();
  const deleteCommand = useDeleteCommand();

  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [selectedCommand, setSelectedCommand] = useState<CustomCommand | null>(
    null
  );

  // Filters
  const [searchQuery, setSearchQuery] = useState("");
  const [statusFilter, setStatusFilter] = useState<
    "all" | "active" | "disabled"
  >("all");
  const [powerFilter, setPowerFilter] = useState<
    "all" | "low" | "medium" | "high" | "owner"
  >("all");

  const [formData, setFormData] = useState({
    name: "",
    usage: "",
    description: "",
    rconCommand: "",
    minArgs: 0,
    maxArgs: -1,
    minPower: 0,
    permissions: "",
    requirementType: "both" as "permission" | "power" | "both",
  });

  const resetForm = () => {
    setFormData({
      name: "",
      usage: "",
      description: "",
      rconCommand: "",
      minArgs: 0,
      maxArgs: -1,
      minPower: 0,
      permissions: "",
      requirementType: "both",
    });
  };

  const handleCreateCommand = async () => {
    try {
      const permissions = formData.permissions
        ? formData.permissions.split(",").map((p) => p.trim())
        : [];

      await createCommand.mutateAsync({
        name: formData.name,
        usage: formData.usage,
        description: formData.description,
        rconCommand: formData.rconCommand,
        minArgs: formData.minArgs,
        maxArgs: formData.maxArgs,
        minPower: formData.minPower,
        permissions,
        requirementType: formData.requirementType,
      });
      setIsCreateDialogOpen(false);
      resetForm();
    } catch (error) {
      console.error("Failed to create command:", error);
    }
  };

  const handleEditCommand = async () => {
    if (!selectedCommand) return;

    try {
      const permissions = formData.permissions
        ? formData.permissions.split(",").map((p) => p.trim())
        : [];

      await updateCommand.mutateAsync({
        id: selectedCommand.id,
        data: {
          name: formData.name || undefined,
          usage: formData.usage || undefined,
          description: formData.description || undefined,
          rconCommand: formData.rconCommand || undefined,
          minArgs: formData.minArgs,
          maxArgs: formData.maxArgs,
          minPower: formData.minPower,
          permissions,
          requirementType: formData.requirementType,
        },
      });
      setIsEditDialogOpen(false);
      setSelectedCommand(null);
      resetForm();
    } catch (error) {
      console.error("Failed to update command:", error);
    }
  };

  const handleDeleteCommand = async (id: number) => {
    if (confirm("Are you sure you want to delete this command?")) {
      try {
        await deleteCommand.mutateAsync(id);
      } catch (error) {
        console.error("Failed to delete command:", error);
      }
    }
  };

  const handleToggleEnabled = async (cmd: CustomCommand) => {
    try {
      await updateCommand.mutateAsync({
        id: cmd.id,
        data: { enabled: !cmd.enabled },
      });
    } catch (error) {
      console.error("Failed to toggle command:", error);
    }
  };

  const openEditDialog = (cmd: CustomCommand) => {
    setSelectedCommand(cmd);
    const permissions = cmd.permissions
      ? JSON.parse(cmd.permissions).join(", ")
      : "";
    setFormData({
      name: cmd.name,
      usage: cmd.usage,
      description: cmd.description,
      rconCommand: cmd.rconCommand,
      minArgs: cmd.minArgs,
      maxArgs: cmd.maxArgs,
      minPower: cmd.minPower,
      permissions,
      requirementType: (cmd.requirementType || "both") as
        | "permission"
        | "power"
        | "both",
    });
    setIsEditDialogOpen(true);
  };

  const getPowerColor = (power: number) => {
    if (power >= 80) return "text-red-500";
    if (power >= 50) return "text-amber-500";
    if (power >= 20) return "text-blue-500";
    return "text-muted-foreground";
  };

  // Filter commands
  const filteredCommands = useMemo(() => {
    if (!commands) return [];

    return commands.filter((cmd) => {
      // Search filter
      if (searchQuery) {
        const query = searchQuery.toLowerCase();
        const matchesSearch =
          cmd.name.toLowerCase().includes(query) ||
          cmd.description.toLowerCase().includes(query) ||
          cmd.usage.toLowerCase().includes(query);
        if (!matchesSearch) return false;
      }

      // Status filter
      if (statusFilter === "active" && !cmd.enabled) return false;
      if (statusFilter === "disabled" && cmd.enabled) return false;

      // Power filter
      if (powerFilter !== "all") {
        if (powerFilter === "low" && (cmd.minPower < 1 || cmd.minPower >= 20))
          return false;
        if (
          powerFilter === "medium" &&
          (cmd.minPower < 20 || cmd.minPower >= 50)
        )
          return false;
        if (powerFilter === "high" && (cmd.minPower < 50 || cmd.minPower >= 80))
          return false;
        if (powerFilter === "owner" && cmd.minPower < 80) return false;
      }

      return true;
    });
  }, [commands, searchQuery, statusFilter, powerFilter]);

  return (
    <ProtectedRoute>
      <DashboardLayout>
        <div className="p-8 space-y-6 bg-background min-h-screen">
          <div className="flex justify-between items-center">
            <div>
              <h1 className="text-4xl font-bold text-foreground mb-2">
                Custom Commands
              </h1>
              <p className="text-muted-foreground">
                Create and manage in-game chat commands with permissions
              </p>
            </div>
            <Dialog
              open={isCreateDialogOpen}
              onOpenChange={setIsCreateDialogOpen}
            >
              <DialogTrigger asChild>
                <Button>
                  <Plus className="h-4 w-4 mr-2" />
                  Create Command
                </Button>
              </DialogTrigger>
              <DialogContent className="max-w-2xl">
                <DialogHeader>
                  <DialogTitle>Create New Command</DialogTitle>
                  <DialogDescription>
                    Create a custom in-game command that players can use in chat
                  </DialogDescription>
                </DialogHeader>
                <div className="space-y-4">
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <Label htmlFor="name">Command Name (without !)</Label>
                      <Input
                        id="name"
                        value={formData.name}
                        onChange={(e) =>
                          setFormData({ ...formData, name: e.target.value })
                        }
                        placeholder="xp"
                      />
                    </div>
                    <div>
                      <Label htmlFor="minPower">Min Power Level (0-100)</Label>
                      <Input
                        id="minPower"
                        type="number"
                        min="0"
                        max="100"
                        value={formData.minPower}
                        onChange={(e) =>
                          setFormData({
                            ...formData,
                            minPower: parseInt(e.target.value) || 0,
                          })
                        }
                      />
                    </div>
                  </div>
                  <div>
                    <Label htmlFor="usage">Usage Example</Label>
                    <Input
                      id="usage"
                      value={formData.usage}
                      onChange={(e) =>
                        setFormData({ ...formData, usage: e.target.value })
                      }
                      placeholder="!xp <player>"
                    />
                  </div>
                  <div>
                    <Label htmlFor="rconCommand">RCON Command Template</Label>
                    <Input
                      id="rconCommand"
                      value={formData.rconCommand}
                      onChange={(e) =>
                        setFormData({
                          ...formData,
                          rconCommand: e.target.value,
                        })
                      }
                      placeholder="clientkick {playerId:arg0} {argsFrom:1}"
                      className="font-mono"
                    />
                    <div className="mt-2 text-xs text-muted-foreground space-y-1">
                      <div className="font-semibold">
                        Available placeholders:
                      </div>
                      <div className="grid grid-cols-2 gap-x-4 gap-y-1 pl-2">
                        <div>
                          <code className="bg-muted px-1 rounded">
                            &#123;arg0&#125;
                          </code>{" "}
                          - First argument
                        </div>
                        <div>
                          <code className="bg-muted px-1 rounded">
                            &#123;arg1&#125;
                          </code>{" "}
                          - Second argument
                        </div>
                        <div>
                          <code className="bg-muted px-1 rounded">
                            &#123;player&#125;
                          </code>{" "}
                          - Command user's name
                        </div>
                        <div>
                          <code className="bg-muted px-1 rounded">
                            &#123;guid&#125;
                          </code>{" "}
                          - Command user's GUID
                        </div>
                        <div>
                          <code className="bg-muted px-1 rounded">
                            &#123;playerId:arg0&#125;
                          </code>{" "}
                          - Resolve player name to ID
                        </div>
                        <div>
                          <code className="bg-muted px-1 rounded">
                            &#123;argsFrom:1&#125;
                          </code>{" "}
                          - Join args from index 1
                        </div>
                      </div>
                    </div>
                  </div>
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <Label htmlFor="minArgs">Min Arguments</Label>
                      <Input
                        id="minArgs"
                        type="number"
                        min="0"
                        value={formData.minArgs}
                        onChange={(e) =>
                          setFormData({
                            ...formData,
                            minArgs: parseInt(e.target.value) || 0,
                          })
                        }
                      />
                    </div>
                    <div>
                      <Label htmlFor="maxArgs">
                        Max Arguments (-1 = unlimited)
                      </Label>
                      <Input
                        id="maxArgs"
                        type="number"
                        min="-1"
                        value={formData.maxArgs}
                        onChange={(e) =>
                          setFormData({
                            ...formData,
                            maxArgs: parseInt(e.target.value) || -1,
                          })
                        }
                      />
                    </div>
                  </div>
                  <div>
                    <Label htmlFor="permissions">
                      Required Permissions (comma-separated)
                    </Label>
                    <Input
                      id="permissions"
                      value={formData.permissions}
                      onChange={(e) =>
                        setFormData({
                          ...formData,
                          permissions: e.target.value,
                        })
                      }
                      placeholder="kick, ban"
                    />
                  </div>
                  <div>
                    <Label>Requirement Type</Label>
                    <div className="space-y-2 mt-2">
                      <label className="flex items-center space-x-2 cursor-pointer">
                        <input
                          type="radio"
                          name="create-requirementType"
                          value="permission"
                          checked={formData.requirementType === "permission"}
                          onChange={(e) =>
                            setFormData({
                              ...formData,
                              requirementType: e.target.value as "permission",
                            })
                          }
                          className="h-4 w-4"
                        />
                        <span className="text-sm">Require permission only</span>
                      </label>
                      <label className="flex items-center space-x-2 cursor-pointer">
                        <input
                          type="radio"
                          name="create-requirementType"
                          value="power"
                          checked={formData.requirementType === "power"}
                          onChange={(e) =>
                            setFormData({
                              ...formData,
                              requirementType: e.target.value as "power",
                            })
                          }
                          className="h-4 w-4"
                        />
                        <span className="text-sm">
                          Require power level only
                        </span>
                      </label>
                      <label className="flex items-center space-x-2 cursor-pointer">
                        <input
                          type="radio"
                          name="create-requirementType"
                          value="both"
                          checked={formData.requirementType === "both"}
                          onChange={(e) =>
                            setFormData({
                              ...formData,
                              requirementType: e.target.value as "both",
                            })
                          }
                          className="h-4 w-4"
                        />
                        <span className="text-sm">
                          Require both power and permission
                        </span>
                      </label>
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
                      placeholder="Gives XP to a player"
                    />
                  </div>
                  <Button onClick={handleCreateCommand} className="w-full">
                    Create Command
                  </Button>
                </div>
              </DialogContent>
            </Dialog>
          </div>

          {/* Filters */}
          <Card>
            <CardContent className="pt-6">
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div>
                  <Label htmlFor="search" className="text-sm">
                    Search Commands
                  </Label>
                  <div className="relative mt-1">
                    <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                    <Input
                      id="search"
                      placeholder="Search by name, description..."
                      value={searchQuery}
                      onChange={(e) => setSearchQuery(e.target.value)}
                      className="pl-9"
                    />
                  </div>
                </div>
                <div>
                  <Label htmlFor="status-filter" className="text-sm">
                    Status
                  </Label>
                  <Select
                    value={statusFilter}
                    onValueChange={(value) =>
                      setStatusFilter(value as typeof statusFilter)
                    }
                  >
                    <SelectTrigger id="status-filter" className="mt-1">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="all">All Commands</SelectItem>
                      <SelectItem value="active">Active Only</SelectItem>
                      <SelectItem value="disabled">Disabled Only</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div>
                  <Label htmlFor="power-filter" className="text-sm">
                    Power Level
                  </Label>
                  <Select
                    value={powerFilter}
                    onValueChange={(value) =>
                      setPowerFilter(value as typeof powerFilter)
                    }
                  >
                    <SelectTrigger id="power-filter" className="mt-1">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="all">All Levels</SelectItem>
                      <SelectItem value="low">Low (0-19)</SelectItem>
                      <SelectItem value="medium">Medium (20-49)</SelectItem>
                      <SelectItem value="high">High (50-79)</SelectItem>
                      <SelectItem value="owner">Owner (80-100)</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>
              <div className="mt-4 text-sm text-muted-foreground">
                Showing {filteredCommands?.length || 0} of{" "}
                {commands?.length || 0} commands
              </div>
            </CardContent>
          </Card>

          {/* Commands Grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {isLoading ? (
              <div className="col-span-full text-center text-muted-foreground">
                Loading commands...
              </div>
            ) : filteredCommands && filteredCommands.length > 0 ? (
              filteredCommands.map((cmd) => (
                <Card key={cmd.id} className={!cmd.enabled ? "opacity-50" : ""}>
                  <CardHeader>
                    <div className="flex justify-between items-start">
                      <div className="flex items-center space-x-2">
                        <Terminal className="h-5 w-5 text-primary" />
                        <CardTitle className="font-mono">!{cmd.name}</CardTitle>
                      </div>
                      <div className="flex items-center space-x-2">
                        {cmd.isBuiltIn && (
                          <Badge variant="secondary" className="text-blue-500">
                            Built-in
                          </Badge>
                        )}
                        {cmd.enabled ? (
                          <Badge variant="outline" className="text-green-500">
                            <Check className="h-3 w-3 mr-1" />
                            Active
                          </Badge>
                        ) : (
                          <Badge
                            variant="outline"
                            className="text-muted-foreground"
                          >
                            <X className="h-3 w-3 mr-1" />
                            Disabled
                          </Badge>
                        )}
                      </div>
                    </div>
                    <CardDescription className="font-mono text-xs">
                      {cmd.usage}
                    </CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div>
                      <div className="text-sm text-muted-foreground mb-1">
                        Description
                      </div>
                      <div className="text-sm">
                        {cmd.description || "No description"}
                      </div>
                    </div>
                    <div>
                      <div className="text-sm text-muted-foreground mb-1">
                        RCON Template
                      </div>
                      <div className="text-sm font-mono bg-muted/30 p-2 rounded border border-border">
                        {cmd.rconCommand || (cmd.isBuiltIn ? "(Go callback function)" : "(empty)")}
                      </div>
                    </div>
                    <div className="grid grid-cols-2 gap-2 text-sm">
                      <div>
                        <div className="text-muted-foreground">Min Args</div>
                        <div>{cmd.minArgs}</div>
                      </div>
                      <div>
                        <div className="text-muted-foreground">Max Args</div>
                        <div>
                          {cmd.maxArgs === -1 ? "Unlimited" : cmd.maxArgs}
                        </div>
                      </div>
                    </div>
                    <div className="flex items-center space-x-2">
                      <Power
                        className={`h-4 w-4 ${getPowerColor(cmd.minPower)}`}
                      />
                      <span
                        className={`text-sm ${getPowerColor(cmd.minPower)}`}
                      >
                        Power: {cmd.minPower}
                      </span>
                    </div>
                    {cmd.permissions &&
                      JSON.parse(cmd.permissions).length > 0 && (
                        <div>
                          <div className="text-sm text-muted-foreground mb-1">
                            Permissions
                          </div>
                          <div className="flex flex-wrap gap-1">
                            {JSON.parse(cmd.permissions).map((perm: string) => (
                              <Badge
                                key={perm}
                                variant="outline"
                                className="text-xs"
                              >
                                {perm}
                              </Badge>
                            ))}
                          </div>
                        </div>
                      )}
                    <div className="flex space-x-2">
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => handleToggleEnabled(cmd)}
                        className="flex-1"
                        disabled={cmd.isBuiltIn}
                      >
                        {cmd.enabled ? "Disable" : "Enable"}
                      </Button>
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => openEditDialog(cmd)}
                        disabled={cmd.isBuiltIn}
                      >
                        <Edit className="h-3 w-3" />
                      </Button>
                      <Button
                        variant="destructive"
                        size="sm"
                        onClick={() => handleDeleteCommand(cmd.id)}
                        disabled={cmd.isBuiltIn}
                      >
                        <Trash2 className="h-3 w-3" />
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              ))
            ) : (
              <div className="col-span-full text-center text-muted-foreground">
                No commands found. Create your first command to get started.
              </div>
            )}
          </div>

          {/* Edit Dialog */}
          <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
            <DialogContent className="max-w-2xl">
              <DialogHeader>
                <DialogTitle>Edit Command</DialogTitle>
                <DialogDescription>
                  Update command settings and permissions
                </DialogDescription>
              </DialogHeader>
              <div className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <Label htmlFor="edit-name">Command Name (without !)</Label>
                    <Input
                      id="edit-name"
                      value={formData.name}
                      onChange={(e) =>
                        setFormData({ ...formData, name: e.target.value })
                      }
                    />
                  </div>
                  <div>
                    <Label htmlFor="edit-minPower">
                      Min Power Level (0-100)
                    </Label>
                    <Input
                      id="edit-minPower"
                      type="number"
                      min="0"
                      max="100"
                      value={formData.minPower}
                      onChange={(e) =>
                        setFormData({
                          ...formData,
                          minPower: parseInt(e.target.value) || 0,
                        })
                      }
                    />
                  </div>
                </div>
                <div>
                  <Label htmlFor="edit-usage">Usage Example</Label>
                  <Input
                    id="edit-usage"
                    value={formData.usage}
                    onChange={(e) =>
                      setFormData({ ...formData, usage: e.target.value })
                    }
                  />
                </div>
                <div>
                  <Label htmlFor="edit-rconCommand">
                    RCON Command Template
                  </Label>
                  <Input
                    id="edit-rconCommand"
                    value={formData.rconCommand}
                    onChange={(e) =>
                      setFormData({ ...formData, rconCommand: e.target.value })
                    }
                    className="font-mono"
                  />
                  <div className="mt-2 text-xs text-muted-foreground space-y-1">
                    <div className="font-semibold">Available placeholders:</div>
                    <div className="grid grid-cols-2 gap-x-4 gap-y-1 pl-2">
                      <div>
                        <code className="bg-muted px-1 rounded">
                          &#123;arg0&#125;
                        </code>{" "}
                        - First argument
                      </div>
                      <div>
                        <code className="bg-muted px-1 rounded">
                          &#123;arg1&#125;
                        </code>{" "}
                        - Second argument
                      </div>
                      <div>
                        <code className="bg-muted px-1 rounded">
                          &#123;player&#125;
                        </code>{" "}
                        - Command user's name
                      </div>
                      <div>
                        <code className="bg-muted px-1 rounded">
                          &#123;guid&#125;
                        </code>{" "}
                        - Command user's GUID
                      </div>
                      <div>
                        <code className="bg-muted px-1 rounded">
                          &#123;playerId:arg0&#125;
                        </code>{" "}
                        - Resolve player name to ID
                      </div>
                      <div>
                        <code className="bg-muted px-1 rounded">
                          &#123;argsFrom:1&#125;
                        </code>{" "}
                        - Join args from index 1
                      </div>
                    </div>
                  </div>
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <Label htmlFor="edit-minArgs">Min Arguments</Label>
                    <Input
                      id="edit-minArgs"
                      type="number"
                      min="0"
                      value={formData.minArgs}
                      onChange={(e) =>
                        setFormData({
                          ...formData,
                          minArgs: parseInt(e.target.value) || 0,
                        })
                      }
                    />
                  </div>
                  <div>
                    <Label htmlFor="edit-maxArgs">
                      Max Arguments (-1 = unlimited)
                    </Label>
                    <Input
                      id="edit-maxArgs"
                      type="number"
                      min="-1"
                      value={formData.maxArgs}
                      onChange={(e) =>
                        setFormData({
                          ...formData,
                          maxArgs: parseInt(e.target.value) || -1,
                        })
                      }
                    />
                  </div>
                </div>
                <div>
                  <Label htmlFor="edit-permissions">
                    Required Permissions (comma-separated)
                  </Label>
                  <Input
                    id="edit-permissions"
                    value={formData.permissions}
                    onChange={(e) =>
                      setFormData({ ...formData, permissions: e.target.value })
                    }
                  />
                </div>
                <div>
                  <Label>Requirement Type</Label>
                  <div className="space-y-2 mt-2">
                    <label className="flex items-center space-x-2 cursor-pointer">
                      <input
                        type="radio"
                        name="edit-requirementType"
                        value="permission"
                        checked={formData.requirementType === "permission"}
                        onChange={(e) =>
                          setFormData({
                            ...formData,
                            requirementType: e.target.value as "permission",
                          })
                        }
                        className="h-4 w-4"
                      />
                      <span className="text-sm">Require permission only</span>
                    </label>
                    <label className="flex items-center space-x-2 cursor-pointer">
                      <input
                        type="radio"
                        name="edit-requirementType"
                        value="power"
                        checked={formData.requirementType === "power"}
                        onChange={(e) =>
                          setFormData({
                            ...formData,
                            requirementType: e.target.value as "power",
                          })
                        }
                        className="h-4 w-4"
                      />
                      <span className="text-sm">Require power level only</span>
                    </label>
                    <label className="flex items-center space-x-2 cursor-pointer">
                      <input
                        type="radio"
                        name="edit-requirementType"
                        value="both"
                        checked={formData.requirementType === "both"}
                        onChange={(e) =>
                          setFormData({
                            ...formData,
                            requirementType: e.target.value as "both",
                          })
                        }
                        className="h-4 w-4"
                      />
                      <span className="text-sm">
                        Require both power and permission
                      </span>
                    </label>
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
                <Button onClick={handleEditCommand} className="w-full">
                  Update Command
                </Button>
              </div>
            </DialogContent>
          </Dialog>
        </div>
      </DashboardLayout>
    </ProtectedRoute>
  );
}

export default Commands;
