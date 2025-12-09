import { DashboardLayout } from "@/components/DashboardLayout";
import { ProtectedRoute } from "@/components/ProtectedRoute";
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
  useApplyMigrations,
  useMigrationStatus,
  useRollbackMigration,
} from "@/hooks/useMigrations";
import {
  AlertCircle,
  CheckCircle2,
  Database,
  PlayCircle,
  RotateCcw,
  AlertTriangle,
} from "lucide-react";
import { toast } from "sonner";
import { formatDistanceToNow } from "date-fns";

export default function Migrations() {
  const { data: status, isLoading } = useMigrationStatus();
  const applyMigrations = useApplyMigrations();
  const rollbackMigration = useRollbackMigration();

  const handleApply = async () => {
    try {
      await applyMigrations.mutateAsync();
      toast.success("Migrations applied successfully");
    } catch (error: unknown) {
      toast.error(
        error instanceof Error ? error.message : "Failed to apply migrations"
      );
    }
  };

  const handleRollback = async () => {
    if (
      !confirm(
        "Are you sure you want to rollback the last migration? This action cannot be undone."
      )
    ) {
      return;
    }

    try {
      await rollbackMigration.mutateAsync();
      toast.success("Migration rolled back successfully");
    } catch (error: unknown) {
      toast.error(
        error instanceof Error ? error.message : "Failed to rollback migration"
      );
    }
  };

  if (isLoading) {
    return (
      <ProtectedRoute requiredPermission="migrations.manage">
        <DashboardLayout>
          <div className="p-8 space-y-6 bg-background min-h-screen">
            <div className="flex items-center justify-center h-64">
              <div className="text-center">
                <Database className="w-12 h-12 mx-auto mb-4 text-muted-foreground animate-pulse" />
                <p className="text-muted-foreground">
                  Loading migration status...
                </p>
              </div>
            </div>
          </div>
        </DashboardLayout>
      </ProtectedRoute>
    );
  }

  return (
    <ProtectedRoute requiredPermission="migrations.manage">
      <DashboardLayout>
        <div className="p-8 space-y-6 bg-background min-h-screen">
          <div>
            <h1 className="text-4xl font-bold text-foreground mb-2">
              Database Migrations
            </h1>
            <p className="text-muted-foreground">
              Manage database schema versions and migrations
            </p>
          </div>

          {/* Current Status */}
          <div className="grid gap-4 md:grid-cols-3">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">
                  Current Version
                </CardTitle>
                <Database className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">
                  {status?.current_version || "None"}
                </div>
                <p className="text-xs text-muted-foreground">
                  Active migration version
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">
                  Applied Migrations
                </CardTitle>
                <CheckCircle2 className="h-4 w-4 text-green-500" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">
                  {status?.total_applied || 0}
                </div>
                <p className="text-xs text-muted-foreground">
                  Successfully applied
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">
                  Pending Migrations
                </CardTitle>
                <AlertCircle className="h-4 w-4 text-yellow-500" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">
                  {status?.total_pending || 0}
                </div>
                <p className="text-xs text-muted-foreground">
                  Waiting to apply
                </p>
              </CardContent>
            </Card>
          </div>

          {/* Pending Migrations Warning */}
          {status && status.total_pending > 0 && (
            <Card className="border-yellow-500/50 bg-yellow-500/10">
              <CardHeader>
                <CardTitle className="flex items-center gap-2 text-yellow-600">
                  <AlertTriangle className="h-5 w-5" />
                  Pending Migrations Detected
                </CardTitle>
                <CardDescription>
                  There {status.total_pending === 1 ? "is" : "are"}{" "}
                  {status.total_pending} pending{" "}
                  {status.total_pending === 1 ? "migration" : "migrations"} that
                  need to be applied to update the database schema.
                </CardDescription>
              </CardHeader>
            </Card>
          )}

          {/* Actions */}
          <Card>
            <CardHeader>
              <CardTitle>Migration Actions</CardTitle>
              <CardDescription>
                Apply pending migrations or rollback the most recent migration
              </CardDescription>
            </CardHeader>
            <CardContent className="flex gap-3">
              <Button
                onClick={handleApply}
                disabled={
                  applyMigrations.isPending ||
                  (status?.total_pending || 0) === 0
                }
                className="gap-2"
              >
                <PlayCircle className="h-4 w-4" />
                Apply All Pending
              </Button>
              <Button
                onClick={handleRollback}
                disabled={
                  rollbackMigration.isPending ||
                  (status?.total_applied || 0) === 0
                }
                variant="destructive"
                className="gap-2"
              >
                <RotateCcw className="h-4 w-4" />
                Rollback Last
              </Button>
            </CardContent>
          </Card>

          {/* Applied Migrations */}
          {status && status.applied_migrations.length > 0 && (
            <Card>
              <CardHeader>
                <CardTitle>Applied Migrations</CardTitle>
                <CardDescription>
                  Migrations that have been successfully applied to the database
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {status.applied_migrations.map((migration) => (
                    <div
                      key={migration.version}
                      className="flex items-start justify-between border-b pb-3 last:border-0"
                    >
                      <div className="space-y-1">
                        <div className="flex items-center gap-2">
                          <Badge variant="default">{migration.version}</Badge>
                          <span className="font-medium">{migration.name}</span>
                          {migration.rolled_back && (
                            <Badge variant="destructive">Rolled Back</Badge>
                          )}
                        </div>
                        <p className="text-sm text-muted-foreground">
                          {migration.description}
                        </p>
                        {migration.applied_at && (
                          <p className="text-xs text-muted-foreground">
                            Applied{" "}
                            {formatDistanceToNow(
                              new Date(migration.applied_at),
                              {
                                addSuffix: true,
                              }
                            )}
                          </p>
                        )}
                      </div>
                      <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0" />
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          )}

          {/* Pending Migrations */}
          {status && status.pending_migrations.length > 0 && (
            <Card>
              <CardHeader>
                <CardTitle>Pending Migrations</CardTitle>
                <CardDescription>
                  Migrations waiting to be applied
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {status.pending_migrations.map((migration) => (
                    <div
                      key={migration.version}
                      className="flex items-start justify-between border-b pb-3 last:border-0"
                    >
                      <div className="space-y-1">
                        <div className="flex items-center gap-2">
                          <Badge variant="outline">{migration.version}</Badge>
                          <span className="font-medium">{migration.name}</span>
                        </div>
                        <p className="text-sm text-muted-foreground">
                          {migration.description}
                        </p>
                      </div>
                      <AlertCircle className="h-5 w-5 text-yellow-500 shrink-0" />
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          )}

          {/* Safety Notice */}
          <Card className="border-red-500/50 bg-red-500/10">
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-red-600">
                <AlertCircle className="h-5 w-5" />
                Important Safety Notice
              </CardTitle>
              <CardDescription>
                Always backup your database before applying or rolling back
                migrations. Rollback operations are destructive and may result
                in data loss. Test migrations in a development environment
                first.
              </CardDescription>
            </CardHeader>
          </Card>
        </div>
      </DashboardLayout>
    </ProtectedRoute>
  );
}
