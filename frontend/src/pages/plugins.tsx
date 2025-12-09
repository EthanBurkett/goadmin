import { useState } from "react";
import {
  usePlugins,
  useStartPlugin,
  useStopPlugin,
  useReloadPlugin,
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
  Play,
  Square,
  RefreshCw,
  AlertCircle,
  CheckCircle,
  Loader2,
} from "lucide-react";
import { Alert, AlertDescription } from "@/components/ui/alert";

export default function PluginsPage() {
  const { data: plugins = [], isLoading, error } = usePlugins();
  const startPlugin = useStartPlugin();
  const stopPlugin = useStopPlugin();
  const reloadPlugin = useReloadPlugin();
  const [selectedPluginId, setSelectedPluginId] = useState<string | null>(null);

  const handleStart = (pluginId: string) => {
    setSelectedPluginId(pluginId);
    startPlugin.mutate(pluginId);
  };

  const handleStop = (pluginId: string) => {
    setSelectedPluginId(pluginId);
    stopPlugin.mutate(pluginId);
  };

  const handleReload = (pluginId: string) => {
    setSelectedPluginId(pluginId);
    reloadPlugin.mutate(pluginId);
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

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="w-8 h-8 animate-spin" />
      </div>
    );
  }

  if (error) {
    return (
      <Alert variant="destructive">
        <AlertCircle className="h-4 w-4" />
        <AlertDescription>
          Error loading plugins: {(error as Error).message}
        </AlertDescription>
      </Alert>
    );
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Plugin Management</h1>
        <p className="text-muted-foreground">
          Manage and monitor installed plugins
        </p>
      </div>

      {plugins.length === 0 ? (
        <Card>
          <CardHeader>
            <CardTitle>No Plugins Found</CardTitle>
            <CardDescription>
              No plugins are currently installed. Place plugin .so files in the
              plugins directory and restart the server.
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
                  <TableHead>Name</TableHead>
                  <TableHead>Version</TableHead>
                  <TableHead>Author</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Loaded At</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {plugins.map((plugin) => (
                  <TableRow key={plugin.id}>
                    <TableCell>
                      <div className="flex flex-col">
                        <span className="font-medium">{plugin.name}</span>
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
                      <span className="text-sm">Plugin</span>
                    </TableCell>
                    <TableCell>{getStateBadge(plugin.state)}</TableCell>
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
                            <span className="ml-1">Start</span>
                          </Button>
                        ) : null}

                        {plugin.state === "started" ? (
                          <>
                            <Button
                              size="sm"
                              variant="outline"
                              onClick={() => handleReload(plugin.id)}
                              disabled={
                                reloadPlugin.isPending &&
                                selectedPluginId === plugin.id
                              }
                            >
                              {reloadPlugin.isPending &&
                              selectedPluginId === plugin.id ? (
                                <Loader2 className="w-4 h-4 animate-spin" />
                              ) : (
                                <RefreshCw className="w-4 h-4" />
                              )}
                              <span className="ml-1">Reload</span>
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
                              <span className="ml-1">Stop</span>
                            </Button>
                          </>
                        ) : null}
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
