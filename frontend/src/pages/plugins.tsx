import { useState } from "react";
import { ProtectedRoute } from "@/components/ProtectedRoute";
import {
  usePlugins,
  useStartPlugin,
  useStopPlugin,
  useHotReloadPlugin,
  usePluginMetrics,
  usePluginDependencies,
  useAllPluginMetrics,
} from "@/hooks/usePlugins";
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
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarProvider,
} from "@/components/ui/sidebar";
import { Progress } from "@/components/ui/progress";
import {
  Play,
  Square,
  AlertCircle,
  CheckCircle,
  Loader2,
  Shield,
  LogOut,
  Database,
  Plug,
  Activity,
  Zap,
  ChevronDown,
  ChevronUp,
  Network,
  Gauge,
} from "lucide-react";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { ServerSelector } from "@/components/ServerSelector";
import { useAuthContext } from "@/hooks/useAuthContext";
import { ServerProvider } from "@/providers/ServerProvider";
import { useNavigate } from "react-router-dom";

// Plugin details row component
function PluginDetailsRow({ pluginId }: { pluginId: string }) {
  const { data: metrics, isLoading: metricsLoading } =
    usePluginMetrics(pluginId);
  const { data: dependencies, isLoading: depsLoading } =
    usePluginDependencies(pluginId);

  if (metricsLoading || depsLoading) {
    return (
      <div className="p-4 flex items-center justify-center">
        <Loader2 className="w-4 h-4 animate-spin mr-2" />
        <span className="text-sm text-muted-foreground">
          Loading details...
        </span>
      </div>
    );
  }

  const hasDependencies =
    dependencies?.dependency_tree &&
    Object.keys(dependencies.dependency_tree).length > 0;

  return (
    <div className="p-4 space-y-4">
      {/* Resource Metrics */}
      {metrics && (
        <div className="space-y-3">
          <h4 className="text-sm font-semibold flex items-center gap-2">
            <Gauge className="w-4 h-4" />
            Resource Usage
          </h4>
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <div className="flex items-center justify-between text-sm">
                <span className="text-muted-foreground">Memory Usage</span>
                <span className="font-medium">
                  {metrics.MemoryUsageMB?.toFixed(2)} MB
                </span>
              </div>
              <Progress
                value={Math.min((metrics.MemoryUsageMB / 100) * 100, 100)}
              />
            </div>
            <div className="space-y-2">
              <div className="flex items-center justify-between text-sm">
                <span className="text-muted-foreground">Goroutines</span>
                <span className="font-medium">{metrics.GoroutineCount}</span>
              </div>
              <Progress
                value={Math.min((metrics.GoroutineCount / 50) * 100, 100)}
              />
            </div>
          </div>
          <div className="flex items-center gap-4 text-xs">
            <div className="flex items-center gap-2">
              <span className="text-muted-foreground">Last Checked:</span>
              <span>{new Date(metrics.LastChecked).toLocaleTimeString()}</span>
            </div>
            {metrics.ViolationCount > 0 && (
              <Badge variant="destructive" className="text-xs">
                {metrics.ViolationCount} violation
                {metrics.ViolationCount !== 1 ? "s" : ""}
              </Badge>
            )}
            {metrics.Throttled && (
              <Badge variant="secondary" className="text-xs">
                Throttled
              </Badge>
            )}
          </div>
        </div>
      )}

      {/* Dependencies */}
      {hasDependencies && (
        <div className="space-y-3">
          <h4 className="text-sm font-semibold flex items-center gap-2">
            <Network className="w-4 h-4" />
            Dependencies
          </h4>
          <div className="space-y-2">
            {Object.entries(dependencies.dependency_tree).map(
              ([plugin, deps]) => (
                <div key={plugin} className="text-sm">
                  <code className="bg-muted px-2 py-1 rounded text-xs">
                    {plugin}
                  </code>
                  {deps.length > 0 ? (
                    <span className="text-muted-foreground ml-2">
                      → depends on:{" "}
                      {deps.map((d) => (
                        <code
                          key={d}
                          className="bg-muted px-1 py-0.5 rounded text-xs ml-1"
                        >
                          {d}
                        </code>
                      ))}
                    </span>
                  ) : (
                    <span className="text-muted-foreground ml-2">
                      (no dependencies)
                    </span>
                  )}
                </div>
              )
            )}
          </div>
        </div>
      )}

      {!metrics && !hasDependencies && (
        <div className="text-sm text-muted-foreground text-center py-2">
          No additional details available
        </div>
      )}
    </div>
  );
}

function Plugins() {
  const { data: plugins = [], isLoading, error } = usePlugins();
  const { data: allMetrics = {} } = useAllPluginMetrics();
  const startPlugin = useStartPlugin();
  const stopPlugin = useStopPlugin();
  const hotReloadPlugin = useHotReloadPlugin();
  const [selectedPluginId, setSelectedPluginId] = useState<string | null>(null);
  const [expandedRows, setExpandedRows] = useState<Set<string>>(new Set());
  const { user, logout } = useAuthContext();
  const navigate = useNavigate();

  const handleStart = (pluginId: string) => {
    setSelectedPluginId(pluginId);
    startPlugin.mutate(pluginId);
  };

  const handleStop = (pluginId: string) => {
    setSelectedPluginId(pluginId);
    stopPlugin.mutate(pluginId);
  };

  const handleHotReload = (pluginId: string) => {
    setSelectedPluginId(pluginId);
    hotReloadPlugin.mutate(pluginId);
  };

  const toggleRowExpanded = (pluginId: string) => {
    const newExpanded = new Set(expandedRows);
    if (newExpanded.has(pluginId)) {
      newExpanded.delete(pluginId);
    } else {
      newExpanded.add(pluginId);
    }
    setExpandedRows(newExpanded);
  };

  const handleLogout = () => {
    logout();
  };

  const getStateBadge = (state: string) => {
    switch (state) {
      case "started":
        return (
          <Badge variant="default" className="bg-green-500">
            <CheckCircle className="w-3 h-3 mr-1" />
            Running
          </Badge>
        );
      case "stopped":
        return (
          <Badge variant="secondary">
            <Square className="w-3 h-3 mr-1" />
            Stopped
          </Badge>
        );
      case "loaded":
        return (
          <Badge variant="outline">
            <Loader2 className="w-3 h-3 mr-1" />
            Loaded
          </Badge>
        );
      case "error":
        return (
          <Badge variant="destructive">
            <AlertCircle className="w-3 h-3 mr-1" />
            Error
          </Badge>
        );
      default:
        return <Badge variant="outline">{state}</Badge>;
    }
  };

  const formatDate = (dateStr: string) => {
    if (!dateStr) return "N/A";
    const date = new Date(dateStr);
    return date.toLocaleString();
  };

  return (
    <SidebarProvider>
      <div className="flex h-screen w-full overflow-hidden">
        <Sidebar>
          <SidebarHeader className="border-b border-border p-4 space-y-3">
            <div className="flex items-center space-x-2">
              <Shield className="h-6 w-6 text-primary" />
              <span className="text-lg font-bold text-foreground">GoAdmin</span>
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
                className="w-full justify-start"
                onClick={() => navigate("/servers")}
              >
                <Database className="mr-2 h-4 w-4" />
                Servers
              </Button>
              <Button
                variant="ghost"
                className="w-full justify-start bg-accent text-accent-foreground"
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
          <SidebarFooter>
            <div className="space-y-2 p-2">
              <div className="px-3 py-2 rounded-lg bg-muted/50">
                <div className="flex items-center gap-2 mb-1">
                  <Shield className="h-3 w-3" />
                  <span className="text-xs font-medium">Logged in as</span>
                </div>
                <span className="text-sm font-semibold block">
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
              <div>
                <h1 className="text-4xl font-bold text-foreground mb-2">
                  Plugin Management
                </h1>
                <p className="text-muted-foreground">
                  Manage and monitor installed plugins
                </p>
              </div>

              {isLoading ? (
                <div className="flex items-center justify-center h-64">
                  <Loader2 className="w-8 h-8 animate-spin" />
                </div>
              ) : error ? (
                <Alert variant="destructive">
                  <AlertCircle className="h-4 w-4" />
                  <AlertDescription>
                    Error loading plugins: {(error as Error).message}
                  </AlertDescription>
                </Alert>
              ) : plugins.length === 0 ? (
                <Card>
                  <CardHeader>
                    <CardTitle>No Plugins Found</CardTitle>
                    <CardDescription>
                      No plugins are currently installed. Create plugins in the
                      plugins directory and rebuild the application.
                    </CardDescription>
                  </CardHeader>
                </Card>
              ) : (
                <Card>
                  <CardHeader>
                    <CardTitle>Installed Plugins ({plugins.length})</CardTitle>
                    <CardDescription>
                      View and control the lifecycle of installed plugins
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    <Table>
                      <TableHeader>
                        <TableRow>
                          <TableHead className="w-12"></TableHead>
                          <TableHead>Name</TableHead>
                          <TableHead>Version</TableHead>
                          <TableHead>Status</TableHead>
                          <TableHead>Memory</TableHead>
                          <TableHead>Goroutines</TableHead>
                          <TableHead>Loaded At</TableHead>
                          <TableHead className="text-right">Actions</TableHead>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {plugins.map((plugin) => {
                          const metrics = allMetrics[plugin.id];
                          const isExpanded = expandedRows.has(plugin.id);

                          return (
                            <>
                              <TableRow key={plugin.id}>
                                <TableCell>
                                  <Button
                                    variant="ghost"
                                    size="sm"
                                    onClick={() => toggleRowExpanded(plugin.id)}
                                    className="h-6 w-6 p-0"
                                  >
                                    {isExpanded ? (
                                      <ChevronUp className="h-4 w-4" />
                                    ) : (
                                      <ChevronDown className="h-4 w-4" />
                                    )}
                                  </Button>
                                </TableCell>
                                <TableCell>
                                  <div className="flex flex-col">
                                    <span className="font-medium">
                                      {plugin.name}
                                    </span>
                                    {plugin.error && (
                                      <span className="text-sm text-destructive">
                                        Error: {plugin.error}
                                      </span>
                                    )}
                                  </div>
                                </TableCell>
                                <TableCell>
                                  <code className="text-xs bg-muted px-2 py-1 rounded">
                                    v{plugin.version}
                                  </code>
                                </TableCell>
                                <TableCell>
                                  {getStateBadge(plugin.state)}
                                </TableCell>
                                <TableCell>
                                  {metrics &&
                                  metrics.MemoryUsageMB !== undefined ? (
                                    <div className="flex items-center gap-2">
                                      <span className="text-sm">
                                        {metrics.MemoryUsageMB.toFixed(1)} MB
                                      </span>
                                      {metrics.ViolationCount > 0 && (
                                        <Badge
                                          variant="destructive"
                                          className="text-xs"
                                        >
                                          {metrics.ViolationCount} violations
                                        </Badge>
                                      )}
                                    </div>
                                  ) : (
                                    <span className="text-muted-foreground text-sm">
                                      N/A
                                    </span>
                                  )}
                                </TableCell>
                                <TableCell>
                                  {metrics &&
                                  metrics.GoroutineCount !== undefined ? (
                                    <span className="text-sm">
                                      {metrics.GoroutineCount}
                                    </span>
                                  ) : (
                                    <span className="text-muted-foreground text-sm">
                                      N/A
                                    </span>
                                  )}
                                </TableCell>
                                <TableCell className="text-sm text-muted-foreground">
                                  {formatDate(plugin.loadedAt)}
                                </TableCell>
                                <TableCell className="text-right">
                                  <div className="flex items-center justify-end gap-2">
                                    {plugin.state === "stopped" ||
                                    plugin.state === "loaded" ? (
                                      <Button
                                        size="sm"
                                        variant="outline"
                                        onClick={() => handleStart(plugin.id)}
                                        disabled={
                                          startPlugin.isPending &&
                                          selectedPluginId === plugin.id
                                        }
                                      >
                                        {startPlugin.isPending &&
                                        selectedPluginId === plugin.id ? (
                                          <Loader2 className="w-4 h-4 animate-spin" />
                                        ) : (
                                          <Play className="w-4 h-4" />
                                        )}
                                      </Button>
                                    ) : null}

                                    {plugin.state === "started" ? (
                                      <>
                                        <Button
                                          size="sm"
                                          variant="outline"
                                          onClick={() =>
                                            handleHotReload(plugin.id)
                                          }
                                          disabled={
                                            hotReloadPlugin.isPending &&
                                            selectedPluginId === plugin.id
                                          }
                                          title="Hot Reload"
                                        >
                                          {hotReloadPlugin.isPending &&
                                          selectedPluginId === plugin.id ? (
                                            <Loader2 className="w-4 h-4 animate-spin" />
                                          ) : (
                                            <Zap className="w-4 h-4" />
                                          )}
                                        </Button>
                                        <Button
                                          size="sm"
                                          variant="outline"
                                          onClick={() => handleStop(plugin.id)}
                                          disabled={
                                            stopPlugin.isPending &&
                                            selectedPluginId === plugin.id
                                          }
                                        >
                                          {stopPlugin.isPending &&
                                          selectedPluginId === plugin.id ? (
                                            <Loader2 className="w-4 h-4 animate-spin" />
                                          ) : (
                                            <Square className="w-4 h-4" />
                                          )}
                                        </Button>
                                      </>
                                    ) : null}
                                  </div>
                                </TableCell>
                              </TableRow>
                              {isExpanded && (
                                <TableRow>
                                  <TableCell
                                    colSpan={8}
                                    className="bg-muted/50"
                                  >
                                    <PluginDetailsRow pluginId={plugin.id} />
                                  </TableCell>
                                </TableRow>
                              )}
                            </>
                          );
                        })}
                      </TableBody>
                    </Table>
                  </CardContent>
                </Card>
              )}
            </div>
          </div>
        </main>
      </div>
    </SidebarProvider>
  );
}

export default function PluginsPage() {
  return (
    <ProtectedRoute>
      <ServerProvider disableRedirect>
        <Plugins />
      </ServerProvider>
    </ProtectedRoute>
  );
}
