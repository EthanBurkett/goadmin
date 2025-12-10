import { useSystemMetrics } from "@/hooks/useMetrics";
import { ProtectedRoute } from "@/components/ProtectedRoute";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Badge } from "@/components/ui/badge";
import {
  Database,
  Users,
  FileText,
  Ban,
  Terminal,
  Activity,
  Clock,
  Shield,
  LogOut,
  Plug,
} from "lucide-react";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarProvider,
} from "@/components/ui/sidebar";
import { Button } from "@/components/ui/button";
import { ServerSelector } from "@/components/ServerSelector";
import { useAuthContext } from "@/hooks/useAuthContext";
import { useNavigate } from "react-router-dom";
import { ServerProvider } from "@/providers/ServerProvider";

function MetricsContent() {
  const { data: metrics, isLoading } = useSystemMetrics();

  const formatUptime = (seconds: number) => {
    const days = Math.floor(seconds / 86400);
    const hours = Math.floor((seconds % 86400) / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    return `${days}d ${hours}h ${minutes}m`;
  };

  if (isLoading) {
    return (
      <div className="space-y-6">
        <h1 className="text-4xl font-bold">System Metrics</h1>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {[...Array(9)].map((_, i) => (
            <Skeleton key={i} className="h-32" />
          ))}
        </div>
      </div>
    );
  }

  if (!metrics) {
    return (
      <div className="space-y-6">
        <h1 className="text-4xl font-bold">System Metrics</h1>
        <Card>
          <CardContent className="pt-6">
            <p className="text-muted-foreground">
              Unable to load system metrics
            </p>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-4xl font-bold text-foreground mb-2">
          System Metrics
        </h1>
        <p className="text-muted-foreground">
          Real-time monitoring and performance statistics
        </p>
      </div>

      {/* Uptime */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Clock className="h-5 w-5" />
            System Uptime
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-3xl font-bold">
            {formatUptime(metrics.uptime_seconds)}
          </div>
          <div className="text-sm text-muted-foreground mt-1">
            {metrics.uptime_seconds.toLocaleString()} seconds
          </div>
        </CardContent>
      </Card>

      {/* Database Metrics */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Database className="h-5 w-5" />
            Database Connection Pool
          </CardTitle>
          <CardDescription>
            Connection statistics and performance
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div>
              <div className="text-sm text-muted-foreground">
                Open Connections
              </div>
              <div className="text-2xl font-bold">{metrics.db_open_conns}</div>
            </div>
            <div>
              <div className="text-sm text-muted-foreground">
                Idle Connections
              </div>
              <div className="text-2xl font-bold">{metrics.db_idle_conns}</div>
            </div>
            <div>
              <div className="text-sm text-muted-foreground">Wait Count</div>
              <div className="text-2xl font-bold">{metrics.db_wait_count}</div>
            </div>
            <div>
              <div className="text-sm text-muted-foreground">Wait Duration</div>
              <div className="text-2xl font-bold">
                {metrics.db_wait_duration_ms.toFixed(0)}ms
              </div>
            </div>
          </div>
          <div className="mt-4 pt-4 border-t grid grid-cols-2 gap-4">
            <div>
              <div className="text-sm text-muted-foreground">
                Max Idle Closed
              </div>
              <div className="text-lg font-semibold">
                {metrics.db_max_idle_closed}
              </div>
            </div>
            <div>
              <div className="text-sm text-muted-foreground">
                Max Lifetime Closed
              </div>
              <div className="text-lg font-semibold">
                {metrics.db_max_life_closed}
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {/* Users */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Users className="h-5 w-5" />
              Users
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">Total</span>
                <Badge variant="outline" className="text-lg">
                  {metrics.total_users}
                </Badge>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">Active</span>
                <Badge className="bg-green-500 text-lg">
                  {metrics.active_users}
                </Badge>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">Pending</span>
                <Badge variant="secondary" className="text-lg">
                  {metrics.pending_users}
                </Badge>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Audit Logs */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Activity className="h-5 w-5" />
              Audit Logs
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">Total</span>
                <Badge variant="outline" className="text-lg">
                  {metrics.total_audit_logs.toLocaleString()}
                </Badge>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">Archived</span>
                <Badge variant="secondary" className="text-lg">
                  {metrics.archived_audit_logs.toLocaleString()}
                </Badge>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">
                  Success Rate
                </span>
                <Badge className="bg-green-500 text-lg">
                  {metrics.audit_success_rate.toFixed(1)}%
                </Badge>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Reports */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <FileText className="h-5 w-5" />
              Reports
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">Total</span>
                <Badge variant="outline" className="text-lg">
                  {metrics.total_reports}
                </Badge>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">Pending</span>
                <Badge variant="secondary" className="text-lg">
                  {metrics.pending_reports}
                </Badge>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Bans */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Ban className="h-5 w-5" />
              Bans
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">Total</span>
                <Badge variant="outline" className="text-lg">
                  {metrics.total_bans}
                </Badge>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">Active</span>
                <Badge variant="destructive" className="text-lg">
                  {metrics.active_bans}
                </Badge>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Commands */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Terminal className="h-5 w-5" />
              Commands
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">Total</span>
                <Badge variant="outline" className="text-lg">
                  {metrics.total_commands}
                </Badge>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">Custom</span>
                <Badge variant="secondary" className="text-lg">
                  {metrics.custom_commands}
                </Badge>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">Plugin</span>
                <Badge className="bg-purple-500 text-lg">
                  {metrics.plugin_commands}
                </Badge>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

export default function MetricsPage() {
  return (
    <ProtectedRoute>
      <MetricsPageContent />
    </ProtectedRoute>
  );
}

function MetricsPageContent() {
  const { user, logout } = useAuthContext();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
  };

  return (
    <ServerProvider disableRedirect>
      <SidebarProvider>
        <div className="flex h-screen w-full overflow-hidden">
          <Sidebar>
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
                  className="w-full justify-start"
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
                  className="w-full justify-start bg-accent text-accent-foreground"
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
                  <span className="font-medium text-foreground">
                    {user?.username}
                  </span>
                </div>
                <Button
                  variant="outline"
                  className="w-full justify-start"
                  onClick={handleLogout}
                >
                  <LogOut className="mr-2 h-4 w-4" />
                  Logout
                </Button>
              </div>
            </SidebarFooter>
          </Sidebar>
          <main className="flex-1 overflow-y-auto p-8 bg-background">
            <MetricsContent />
          </main>
        </div>
      </SidebarProvider>
    </ServerProvider>
  );
}
