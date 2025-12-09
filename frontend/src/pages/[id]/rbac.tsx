import { useState } from "react";
import type { ColumnDef } from "@tanstack/react-table";
import { ProtectedRoute } from "@/components/ProtectedRoute";
import { DataTable } from "@/components/DataTable";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  useRoles,
  usePermissions,
  useRbacUsers,
  useCreateRole,
  useCreatePermission,
  useDeleteRole,
  useDeletePermission,
  useAssignPermissionsToRole,
  useRemovePermissionFromRole,
  usePendingUsers,
  useApproveUser,
  useDenyUser,
  useDeleteUser,
  useAssignRoleToUser,
  useRemoveRoleFromUser,
  type Role,
  type Permission,
  type RbacUser as User,
  type PendingUser,
} from "@/hooks/useRbac";
import {
  Shield,
  Key,
  Users,
  Plus,
  MoreVertical,
  Trash2,
  UserPlus,
  Edit,
  UserCheck,
  UserX,
  Clock,
} from "lucide-react";
import { useAuth } from "@/hooks/useAuth";

function RBAC() {
  const { user: currentUser } = useAuth();
  const [activeTab, setActiveTab] = useState<
    "roles" | "permissions" | "users" | "pending"
  >("roles");
  const [newRoleName, setNewRoleName] = useState("");
  const [newRoleDesc, setNewRoleDesc] = useState("");
  const [newPermName, setNewPermName] = useState("");
  const [newPermDesc, setNewPermDesc] = useState("");
  const [isRoleDialogOpen, setIsRoleDialogOpen] = useState(false);
  const [isPermDialogOpen, setIsPermDialogOpen] = useState(false);
  const [editingRole, setEditingRole] = useState<Role | null>(null);
  const [isEditPermissionsOpen, setIsEditPermissionsOpen] = useState(false);
  const [approvingUser, setApprovingUser] = useState<PendingUser | null>(null);
  const [isApproveDialogOpen, setIsApproveDialogOpen] = useState(false);
  const [selectedRoleId, setSelectedRoleId] = useState<number | null>(null);
  const [managingUser, setManagingUser] = useState<User | null>(null);
  const [isManageRolesOpen, setIsManageRolesOpen] = useState(false);
  const [deletingUser, setDeletingUser] = useState<User | null>(null);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);

  // Queries
  const rolesQuery = useRoles();
  const permissionsQuery = usePermissions();
  const usersQuery = useRbacUsers();
  const pendingUsersQuery = usePendingUsers();

  // Mutations
  const createRoleMutation = useCreateRole();
  const createPermissionMutation = useCreatePermission();
  const deleteRoleMutation = useDeleteRole();
  const deletePermissionMutation = useDeletePermission();
  const assignPermissionMutation = useAssignPermissionsToRole();
  const removePermissionMutation = useRemovePermissionFromRole();
  const approveUserMutation = useApproveUser();
  const denyUserMutation = useDenyUser();
  const deleteUserMutation = useDeleteUser();
  const assignRoleToUserMutation = useAssignRoleToUser();
  const removeRoleFromUserMutation = useRemoveRoleFromUser();

  // Helper function to check if current user has a permission
  const hasPermission = (permissionName: string): boolean => {
    if (!currentUser?.roles) return false;
    return currentUser.roles.some((role) =>
      role.permissions?.some((perm) => perm.name === permissionName)
    );
  };

  // Table columns
  const roleColumns: ColumnDef<Role>[] = [
    {
      accessorKey: "id",
      header: "ID",
      cell: ({ row }) => <div className="font-mono">{row.getValue("id")}</div>,
    },
    {
      accessorKey: "name",
      header: "Role Name",
      cell: ({ row }) => (
        <div className="font-semibold">{row.getValue("name")}</div>
      ),
    },
    {
      accessorKey: "description",
      header: "Description",
      cell: ({ row }) => (
        <div className="text-muted-foreground">
          {row.getValue("description") || "No description"}
        </div>
      ),
    },
    {
      accessorKey: "permissions",
      header: "Permissions",
      cell: ({ row }) => {
        const permissions = row.getValue("permissions") as Permission[];
        return (
          <div className="flex flex-wrap gap-1">
            {permissions && permissions.length > 0 ? (
              permissions.map((perm) => (
                <Badge key={perm.id} variant="secondary" className="text-xs">
                  {perm.name}
                </Badge>
              ))
            ) : (
              <span className="text-muted-foreground text-sm">
                No permissions
              </span>
            )}
          </div>
        );
      },
    },
    {
      id: "actions",
      cell: ({ row }) => {
        const role = row.original;
        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="sm">
                <MoreVertical className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuLabel>Actions</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                onClick={() => {
                  setEditingRole(role);
                  setIsEditPermissionsOpen(true);
                }}
              >
                <Edit className="h-4 w-4 mr-2" />
                Edit Permissions
              </DropdownMenuItem>
              <DropdownMenuItem
                className="text-destructive"
                onClick={() => deleteRoleMutation.mutate(role.id)}
              >
                <Trash2 className="h-4 w-4 mr-2" />
                Delete Role
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        );
      },
    },
  ];

  const permissionColumns: ColumnDef<Permission>[] = [
    {
      accessorKey: "id",
      header: "ID",
      cell: ({ row }) => <div className="font-mono">{row.getValue("id")}</div>,
    },
    {
      accessorKey: "name",
      header: "Permission Name",
      cell: ({ row }) => (
        <Badge variant="outline">{row.getValue("name")}</Badge>
      ),
    },
    {
      accessorKey: "description",
      header: "Description",
      cell: ({ row }) => (
        <div className="text-muted-foreground">
          {row.getValue("description") || "No description"}
        </div>
      ),
    },
    {
      id: "actions",
      cell: ({ row }) => {
        const permission = row.original;
        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="sm">
                <MoreVertical className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuLabel>Actions</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                className="text-destructive"
                onClick={() => deletePermissionMutation.mutate(permission.id)}
              >
                <Trash2 className="h-4 w-4 mr-2" />
                Delete Permission
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        );
      },
    },
  ];

  const userColumns: ColumnDef<User>[] = [
    {
      accessorKey: "id",
      header: "ID",
      cell: ({ row }) => <div className="font-mono">{row.getValue("id")}</div>,
    },
    {
      accessorKey: "username",
      header: "Username",
      cell: ({ row }) => (
        <div className="font-semibold">{row.getValue("username")}</div>
      ),
    },
    {
      accessorKey: "roles",
      header: "Roles",
      cell: ({ row }) => {
        const roles = row.getValue("roles") as Role[];
        return (
          <div className="flex flex-wrap gap-1">
            {roles && roles.length > 0 ? (
              roles.map((role) => (
                <Badge key={role.id} className="text-xs">
                  {role.name}
                </Badge>
              ))
            ) : (
              <span className="text-muted-foreground text-sm">No roles</span>
            )}
          </div>
        );
      },
    },
    {
      id: "actions",
      cell: ({ row }) => {
        const user = row.original;
        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="sm">
                <MoreVertical className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuLabel>Actions</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                onClick={() => {
                  setManagingUser(user);
                  setIsManageRolesOpen(true);
                }}
              >
                <UserPlus className="h-4 w-4 mr-2" />
                Manage Roles
              </DropdownMenuItem>
              {hasPermission("users.delete") && (
                <DropdownMenuItem
                  className="text-destructive"
                  onClick={() => {
                    setDeletingUser(user);
                    setIsDeleteDialogOpen(true);
                  }}
                >
                  <Trash2 className="h-4 w-4 mr-2" />
                  Delete User
                </DropdownMenuItem>
              )}
            </DropdownMenuContent>
          </DropdownMenu>
        );
      },
    },
  ];

  const pendingUserColumns: ColumnDef<PendingUser>[] = [
    {
      accessorKey: "id",
      header: "ID",
      cell: ({ row }) => <div className="font-mono">{row.getValue("id")}</div>,
    },
    {
      accessorKey: "username",
      header: "Username",
      cell: ({ row }) => (
        <div className="font-semibold">{row.getValue("username")}</div>
      ),
    },
    {
      id: "actions",
      header: "Actions",
      cell: ({ row }) => {
        const user = row.original;
        return (
          <div className="flex gap-2">
            <Button
              size="sm"
              variant="default"
              onClick={() => {
                setApprovingUser(user);
                setIsApproveDialogOpen(true);
              }}
            >
              <UserCheck className="h-4 w-4 mr-2" />
              Approve
            </Button>
            <Button
              size="sm"
              variant="destructive"
              onClick={() => {
                if (
                  confirm(
                    `Are you sure you want to deny ${user.username}? This will delete their account.`
                  )
                ) {
                  denyUserMutation.mutate(user.id);
                }
              }}
            >
              <UserX className="h-4 w-4 mr-2" />
              Deny
            </Button>
          </div>
        );
      },
    },
  ];

  return (
    <ProtectedRoute requiredPermission="rbac.manage">
      <div className="space-y-6 bg-background min-h-screen">
        <div>
          <h1 className="text-4xl font-bold text-foreground mb-2">
            Role-Based Access Control
          </h1>
          <p className="text-muted-foreground">
            Manage roles, permissions, and user access
          </p>
        </div>

        {/* Tab Navigation */}
        <div className="flex gap-2 border-b border-border">
          <button
            onClick={() => setActiveTab("roles")}
            className={`px-4 py-2 font-medium transition-colors ${
              activeTab === "roles"
                ? "border-b-2 border-primary text-primary"
                : "text-muted-foreground hover:text-foreground"
            }`}
          >
            <Shield className="h-4 w-4 inline mr-2" />
            Roles
          </button>
          <button
            onClick={() => setActiveTab("permissions")}
            className={`px-4 py-2 font-medium transition-colors ${
              activeTab === "permissions"
                ? "border-b-2 border-primary text-primary"
                : "text-muted-foreground hover:text-foreground"
            }`}
          >
            <Key className="h-4 w-4 inline mr-2" />
            Permissions
          </button>
          <button
            onClick={() => setActiveTab("users")}
            className={`px-4 py-2 font-medium transition-colors ${
              activeTab === "users"
                ? "border-b-2 border-primary text-primary"
                : "text-muted-foreground hover:text-foreground"
            }`}
          >
            <Users className="h-4 w-4 inline mr-2" />
            Users
          </button>
          <button
            onClick={() => setActiveTab("pending")}
            className={`px-4 py-2 font-medium transition-colors relative flex flex-row items-center ${
              activeTab === "pending"
                ? "border-b-2 border-primary text-primary"
                : "text-muted-foreground hover:text-foreground"
            }`}
          >
            <Clock className="h-4 w-4 inline mr-2" />
            Pending Approvals
            {pendingUsersQuery.data && pendingUsersQuery.data.length > 0 && (
              <Badge className="ml-2 h-5 w-5 rounded-full p-0 flex items-center justify-center text-xs">
                {pendingUsersQuery.data.length}
              </Badge>
            )}
          </button>
        </div>

        {/* Roles Tab */}
        {activeTab === "roles" && (
          <Card className="bg-card border-border">
            <CardHeader>
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle className="text-foreground">Roles</CardTitle>
                  <CardDescription className="text-muted-foreground">
                    Manage user roles and their permissions
                  </CardDescription>
                </div>
                <Dialog
                  open={isRoleDialogOpen}
                  onOpenChange={setIsRoleDialogOpen}
                >
                  <DialogTrigger asChild>
                    <Button>
                      <Plus className="h-4 w-4 mr-2" />
                      Create Role
                    </Button>
                  </DialogTrigger>
                  <DialogContent className="bg-card border-border">
                    <DialogHeader>
                      <DialogTitle className="text-foreground">
                        Create New Role
                      </DialogTitle>
                      <DialogDescription className="text-muted-foreground">
                        Add a new role to the system
                      </DialogDescription>
                    </DialogHeader>
                    <div className="space-y-4 py-4">
                      <div className="space-y-2">
                        <Label htmlFor="role-name">Role Name</Label>
                        <Input
                          id="role-name"
                          placeholder="e.g., moderator"
                          value={newRoleName}
                          onChange={(e) => setNewRoleName(e.target.value)}
                        />
                      </div>
                      <div className="space-y-2">
                        <Label htmlFor="role-desc">Description</Label>
                        <Input
                          id="role-desc"
                          placeholder="Describe this role"
                          value={newRoleDesc}
                          onChange={(e) => setNewRoleDesc(e.target.value)}
                        />
                      </div>
                    </div>
                    <DialogFooter>
                      <Button
                        onClick={() =>
                          createRoleMutation.mutate(
                            {
                              name: newRoleName,
                              description: newRoleDesc,
                            },
                            {
                              onSuccess: () => {
                                setNewRoleName("");
                                setNewRoleDesc("");
                                setIsRoleDialogOpen(false);
                              },
                            }
                          )
                        }
                        disabled={!newRoleName || createRoleMutation.isPending}
                      >
                        {createRoleMutation.isPending
                          ? "Creating..."
                          : "Create Role"}
                      </Button>
                    </DialogFooter>
                  </DialogContent>
                </Dialog>
              </div>
            </CardHeader>
            <CardContent>
              {rolesQuery.isLoading ? (
                <div className="text-center py-12 text-muted-foreground">
                  Loading roles...
                </div>
              ) : rolesQuery.isError ? (
                <div className="text-center py-12 text-destructive">
                  Failed to load roles
                </div>
              ) : (
                <DataTable
                  columns={roleColumns}
                  data={rolesQuery.data || []}
                  searchKey="name"
                  searchPlaceholder="Search roles..."
                />
              )}
            </CardContent>
          </Card>
        )}

        {/* Permissions Tab */}
        {activeTab === "permissions" && (
          <Card className="bg-card border-border">
            <CardHeader>
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle className="text-foreground">Permissions</CardTitle>
                  <CardDescription className="text-muted-foreground">
                    Manage system permissions
                  </CardDescription>
                </div>
                <Dialog
                  open={isPermDialogOpen}
                  onOpenChange={setIsPermDialogOpen}
                >
                  <DialogTrigger asChild>
                    <Button>
                      <Plus className="h-4 w-4 mr-2" />
                      Create Permission
                    </Button>
                  </DialogTrigger>
                  <DialogContent className="bg-card border-border">
                    <DialogHeader>
                      <DialogTitle className="text-foreground">
                        Create New Permission
                      </DialogTitle>
                      <DialogDescription className="text-muted-foreground">
                        Add a new permission to the system
                      </DialogDescription>
                    </DialogHeader>
                    <div className="space-y-4 py-4">
                      <div className="space-y-2">
                        <Label htmlFor="perm-name">Permission Name</Label>
                        <Input
                          id="perm-name"
                          placeholder="e.g., players.manage"
                          value={newPermName}
                          onChange={(e) => setNewPermName(e.target.value)}
                        />
                      </div>
                      <div className="space-y-2">
                        <Label htmlFor="perm-desc">Description</Label>
                        <Input
                          id="perm-desc"
                          placeholder="Describe this permission"
                          value={newPermDesc}
                          onChange={(e) => setNewPermDesc(e.target.value)}
                        />
                      </div>
                    </div>
                    <DialogFooter>
                      <Button
                        onClick={() =>
                          createPermissionMutation.mutate(
                            {
                              name: newPermName,
                              description: newPermDesc,
                            },
                            {
                              onSuccess: () => {
                                setNewPermName("");
                                setNewPermDesc("");
                                setIsPermDialogOpen(false);
                              },
                            }
                          )
                        }
                        disabled={
                          !newPermName || createPermissionMutation.isPending
                        }
                      >
                        {createPermissionMutation.isPending
                          ? "Creating..."
                          : "Create Permission"}
                      </Button>
                    </DialogFooter>
                  </DialogContent>
                </Dialog>
              </div>
            </CardHeader>
            <CardContent>
              {permissionsQuery.isLoading ? (
                <div className="text-center py-12 text-muted-foreground">
                  Loading permissions...
                </div>
              ) : permissionsQuery.isError ? (
                <div className="text-center py-12 text-destructive">
                  Failed to load permissions
                </div>
              ) : (
                <DataTable
                  columns={permissionColumns}
                  data={permissionsQuery.data || []}
                  searchKey="name"
                  searchPlaceholder="Search permissions..."
                />
              )}
            </CardContent>
          </Card>
        )}

        {/* Users Tab */}
        {activeTab === "users" && (
          <Card className="bg-card border-border">
            <CardHeader>
              <CardTitle className="text-foreground">Users</CardTitle>
              <CardDescription className="text-muted-foreground">
                Manage user role assignments
              </CardDescription>
            </CardHeader>
            <CardContent>
              {usersQuery.isLoading ? (
                <div className="text-center py-12 text-muted-foreground">
                  Loading users...
                </div>
              ) : usersQuery.isError ? (
                <div className="text-center py-12 text-destructive">
                  Failed to load users
                </div>
              ) : (
                <DataTable
                  columns={userColumns}
                  data={usersQuery.data || []}
                  searchKey="username"
                  searchPlaceholder="Search users..."
                />
              )}
            </CardContent>
          </Card>
        )}

        {/* Pending Users Tab */}
        {activeTab === "pending" && (
          <Card className="bg-card border-border">
            <CardHeader>
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle className="text-foreground">
                    Pending User Approvals
                  </CardTitle>
                  <CardDescription className="text-muted-foreground">
                    Review and approve new user registrations
                  </CardDescription>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              {pendingUsersQuery.isLoading ? (
                <div className="text-center py-12 text-muted-foreground">
                  Loading pending users...
                </div>
              ) : pendingUsersQuery.isError ? (
                <div className="text-center py-12 text-destructive">
                  Failed to load pending users
                </div>
              ) : pendingUsersQuery.data &&
                pendingUsersQuery.data.length === 0 ? (
                <div className="text-center py-12 text-muted-foreground">
                  No pending user approvals
                </div>
              ) : (
                <DataTable
                  columns={pendingUserColumns}
                  data={pendingUsersQuery.data || []}
                  searchKey="username"
                  searchPlaceholder="Search pending users..."
                />
              )}
            </CardContent>
          </Card>
        )}

        {/* Approve User Dialog */}
        <Dialog
          open={isApproveDialogOpen}
          onOpenChange={setIsApproveDialogOpen}
        >
          <DialogContent className="bg-card border-border">
            <DialogHeader>
              <DialogTitle className="text-foreground">
                Approve User: {approvingUser?.username}
              </DialogTitle>
              <DialogDescription className="text-muted-foreground">
                Select a role to assign to this user
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4 py-4">
              <div className="space-y-2">
                <Label htmlFor="role-select">Assign Role</Label>
                <Select
                  value={selectedRoleId?.toString()}
                  onValueChange={(value) => setSelectedRoleId(Number(value))}
                >
                  <SelectTrigger className="bg-muted/30 border-border">
                    <SelectValue placeholder="Select a role..." />
                  </SelectTrigger>
                  <SelectContent>
                    {rolesQuery.data?.map((role) => (
                      <SelectItem key={role.id} value={role.id.toString()}>
                        {role.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>
            <DialogFooter>
              <Button
                variant="outline"
                onClick={() => {
                  setIsApproveDialogOpen(false);
                  setApprovingUser(null);
                  setSelectedRoleId(null);
                }}
              >
                Cancel
              </Button>
              <Button
                onClick={() => {
                  if (approvingUser && selectedRoleId) {
                    approveUserMutation.mutate(
                      {
                        userId: approvingUser.id,
                        roleId: selectedRoleId,
                      },
                      {
                        onSuccess: () => {
                          setIsApproveDialogOpen(false);
                          setApprovingUser(null);
                          setSelectedRoleId(null);
                        },
                      }
                    );
                  }
                }}
                disabled={!selectedRoleId}
              >
                <UserCheck className="h-4 w-4 mr-2" />
                Approve User
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        {/* Edit Role Permissions Dialog */}
        <Dialog
          open={isEditPermissionsOpen}
          onOpenChange={setIsEditPermissionsOpen}
        >
          <DialogContent className="bg-card border-border max-w-2xl">
            <DialogHeader>
              <DialogTitle className="text-foreground">
                Edit Permissions for {editingRole?.name}
              </DialogTitle>
              <DialogDescription className="text-muted-foreground">
                Add or remove permissions from this role
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4 py-4">
              <div>
                <Label className="text-sm font-medium">
                  Current Permissions
                </Label>
                <div className="mt-2 flex flex-wrap gap-2">
                  {editingRole?.permissions &&
                  editingRole.permissions.length > 0 ? (
                    editingRole.permissions.map((perm) => (
                      <Badge
                        key={perm.id}
                        variant="secondary"
                        className="gap-1"
                      >
                        {perm.name}
                        <button
                          onClick={() =>
                            removePermissionMutation.mutate({
                              roleId: editingRole.id,
                              permissionId: perm.id,
                            })
                          }
                          className="ml-1 hover:text-destructive"
                        >
                          ×
                        </button>
                      </Badge>
                    ))
                  ) : (
                    <span className="text-sm text-muted-foreground">
                      No permissions assigned
                    </span>
                  )}
                </div>
              </div>
              <div>
                <Label className="text-sm font-medium">
                  Available Permissions
                </Label>
                <div className="mt-2 flex flex-wrap gap-2">
                  {permissionsQuery.data
                    ?.filter(
                      (perm) =>
                        !editingRole?.permissions?.some((p) => p.id === perm.id)
                    )
                    .map((perm) => (
                      <Badge
                        key={perm.id}
                        variant="outline"
                        className="cursor-pointer hover:bg-primary/10"
                        onClick={() =>
                          assignPermissionMutation.mutate({
                            roleId: editingRole!.id,
                            permissionId: perm.id,
                          })
                        }
                      >
                        + {perm.name}
                      </Badge>
                    ))}
                </div>
              </div>
            </div>
            <DialogFooter>
              <Button
                variant="outline"
                onClick={() => {
                  setIsEditPermissionsOpen(false);
                  setEditingRole(null);
                }}
              >
                Close
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        {/* Manage User Roles Dialog */}
        <Dialog open={isManageRolesOpen} onOpenChange={setIsManageRolesOpen}>
          <DialogContent className="bg-card border-border max-w-2xl">
            <DialogHeader>
              <DialogTitle className="text-foreground">
                Manage Roles for {managingUser?.username}
              </DialogTitle>
              <DialogDescription className="text-muted-foreground">
                Add or remove roles from this user
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4 py-4">
              <div>
                <Label className="text-sm font-medium">Current Roles</Label>
                <div className="mt-2 flex flex-wrap gap-2">
                  {managingUser?.roles && managingUser.roles.length > 0 ? (
                    managingUser.roles.map((role) => (
                      <Badge
                        key={role.id}
                        variant="secondary"
                        className="gap-1"
                      >
                        {role.name}
                        <button
                          onClick={() =>
                            removeRoleFromUserMutation.mutate({
                              userId: managingUser.id,
                              roleId: role.id,
                            })
                          }
                          className="ml-1 hover:text-destructive"
                        >
                          ×
                        </button>
                      </Badge>
                    ))
                  ) : (
                    <span className="text-sm text-muted-foreground">
                      No roles assigned
                    </span>
                  )}
                </div>
              </div>
              <div>
                <Label className="text-sm font-medium">Available Roles</Label>
                <div className="mt-2 flex flex-wrap gap-2">
                  {rolesQuery.data
                    ?.filter(
                      (role) =>
                        !managingUser?.roles?.some((r) => r.id === role.id)
                    )
                    .map((role) => (
                      <Badge
                        key={role.id}
                        variant="outline"
                        className="cursor-pointer hover:bg-primary/10"
                        onClick={() =>
                          assignRoleToUserMutation.mutate({
                            userId: managingUser!.id,
                            roleId: role.id,
                          })
                        }
                      >
                        + {role.name}
                      </Badge>
                    ))}
                </div>
              </div>
            </div>
            <DialogFooter>
              <Button
                variant="outline"
                onClick={() => {
                  setIsManageRolesOpen(false);
                  setManagingUser(null);
                }}
              >
                Close
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        {/* Delete User Confirmation Dialog */}
        <Dialog open={isDeleteDialogOpen} onOpenChange={setIsDeleteDialogOpen}>
          <DialogContent className="bg-card border-border">
            <DialogHeader>
              <DialogTitle className="text-foreground">
                Delete User Account
              </DialogTitle>
              <DialogDescription className="text-muted-foreground">
                Are you sure you want to delete {deletingUser?.username}? This
                action cannot be undone.
              </DialogDescription>
            </DialogHeader>
            <DialogFooter>
              <Button
                variant="outline"
                onClick={() => {
                  setIsDeleteDialogOpen(false);
                  setDeletingUser(null);
                }}
              >
                Cancel
              </Button>
              <Button
                variant="destructive"
                onClick={() => {
                  if (deletingUser) {
                    deleteUserMutation.mutate(deletingUser.id, {
                      onSuccess: () => {
                        setIsDeleteDialogOpen(false);
                        setDeletingUser(null);
                      },
                    });
                  }
                }}
              >
                <Trash2 className="h-4 w-4 mr-2" />
                Delete User
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </ProtectedRoute>
  );
}

export default RBAC;
