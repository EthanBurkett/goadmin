import { useState, useMemo } from "react";
import { ProtectedRoute } from "@/components/ProtectedRoute";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  useServerStats,
  useSystemStats,
  usePlayerStats,
} from "@/hooks/useStats";
import {
  LineChart,
  Line,
  AreaChart,
  Area,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import { Activity, Server, Users, TrendingUp } from "lucide-react";
import { format } from "date-fns";

type TimeRange = "1h" | "6h" | "24h" | "7d" | "30d";

function Analytics() {
  const [timeRange, setTimeRange] = useState<TimeRange>("24h");

  // Calculate time range - memoized to prevent recalculation on every render
  const { start, end } = useMemo(() => {
    const end = new Date();
    const start = new Date();

    switch (timeRange) {
      case "1h":
        start.setHours(start.getHours() - 1);
        break;
      case "6h":
        start.setHours(start.getHours() - 6);
        break;
      case "24h":
        start.setHours(start.getHours() - 24);
        break;
      case "7d":
        start.setDate(start.getDate() - 7);
        break;
      case "30d":
        start.setDate(start.getDate() - 30);
        break;
    }

    return {
      start: start.toISOString(),
      end: end.toISOString(),
    };
  }, [timeRange]);

  const { data: serverStats, isLoading: serverLoading } = useServerStats(
    start,
    end
  );
  const { data: systemStats, isLoading: systemLoading } = useSystemStats(
    start,
    end
  );
  const { data: playerStats, isLoading: playerLoading } = usePlayerStats(
    start,
    end
  );

  const formatTime = (timestamp: string) => {
    return format(new Date(timestamp), "HH:mm");
  };

  const formatDate = (timestamp: string) => {
    return format(new Date(timestamp), "MMM dd HH:mm");
  };

  const formatMemoryDisplay = (bytes: number) => {
    const kb = bytes / 1024;
    if (kb < 1024) {
      return `${Math.round(kb)} KB`;
    }
    const mb = bytes / (1024 * 1024);
    if (mb < 1024) {
      return `${Math.round(mb)} MB`;
    }
    return `${(mb / 1024).toFixed(2)} GB`;
  };

  return (
    <ProtectedRoute requiredPermission="status.view">
      <div className="space-y-6 bg-background min-h-screen">
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-4xl font-bold text-foreground mb-2">
              Server Analytics
            </h1>
            <p className="text-muted-foreground">
              Real-time server metrics and performance data
            </p>
          </div>
          <Select
            value={timeRange}
            onValueChange={(v) => setTimeRange(v as TimeRange)}
          >
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Select time range" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="1h">Last Hour</SelectItem>
              <SelectItem value="6h">Last 6 Hours</SelectItem>
              <SelectItem value="24h">Last 24 Hours</SelectItem>
              <SelectItem value="7d">Last 7 Days</SelectItem>
              <SelectItem value="30d">Last 30 Days</SelectItem>
            </SelectContent>
          </Select>
        </div>

        {/* Stats Overview Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                Current Players
              </CardTitle>
              <Users className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {serverStats && serverStats.length > 0
                  ? serverStats[serverStats.length - 1].playerCount
                  : 0}
                {serverStats &&
                  serverStats.length > 0 &&
                  ` / ${serverStats[serverStats.length - 1].maxPlayers}`}
              </div>
              <p className="text-xs text-muted-foreground">
                Active players on server
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Server FPS</CardTitle>
              <Activity className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {serverStats && serverStats.length > 0
                  ? serverStats[serverStats.length - 1].fps
                  : 0}
              </div>
              <p className="text-xs text-muted-foreground">Server frame rate</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                Average Ping
              </CardTitle>
              <TrendingUp className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {playerStats && playerStats.length > 0
                  ? Math.round(playerStats[playerStats.length - 1].avgPing)
                  : 0}
                ms
              </div>
              <p className="text-xs text-muted-foreground">Network latency</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Current Map</CardTitle>
              <Server className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold truncate">
                {serverStats && serverStats.length > 0
                  ? serverStats[serverStats.length - 1].mapName
                  : "Unknown"}
              </div>
              <p className="text-xs text-muted-foreground">
                {serverStats && serverStats.length > 0
                  ? serverStats[serverStats.length - 1].gametype
                  : ""}
              </p>
            </CardContent>
          </Card>
        </div>

        {/* Player Count Chart */}
        <Card>
          <CardHeader>
            <CardTitle>Player Activity</CardTitle>
            <CardDescription>Player count over time</CardDescription>
          </CardHeader>
          <CardContent>
            {serverLoading ? (
              <div className="h-[300px] flex items-center justify-center text-muted-foreground">
                Loading...
              </div>
            ) : !serverStats || serverStats.length === 0 ? (
              <div className="h-[300px] flex items-center justify-center text-muted-foreground">
                No data available for selected time range
              </div>
            ) : (
              <ResponsiveContainer width="100%" height={300}>
                <AreaChart data={serverStats}>
                  <defs>
                    <linearGradient
                      id="colorPlayers"
                      x1="0"
                      y1="0"
                      x2="0"
                      y2="1"
                    >
                      <stop offset="5%" stopColor="#3b82f6" stopOpacity={0.8} />
                      <stop
                        offset="95%"
                        stopColor="#3b82f6"
                        stopOpacity={0.1}
                      />
                    </linearGradient>
                  </defs>
                  <CartesianGrid
                    strokeDasharray="3 3"
                    className="stroke-muted"
                  />
                  <XAxis
                    dataKey="timestamp"
                    tickFormatter={
                      timeRange === "7d" || timeRange === "30d"
                        ? formatDate
                        : formatTime
                    }
                    className="text-xs"
                  />
                  <YAxis className="text-xs" />
                  <Tooltip
                    contentStyle={{
                      backgroundColor: "var(--card)",
                      border: "1px solid var(--border)",
                      borderRadius: "6px",
                    }}
                    labelFormatter={(label) => format(new Date(label), "PPpp")}
                  />
                  <Legend />
                  <Area
                    type="monotone"
                    dataKey="playerCount"
                    stroke="#3b82f6"
                    strokeWidth={2}
                    fillOpacity={1}
                    fill="url(#colorPlayers)"
                    name="Players"
                  />
                </AreaChart>
              </ResponsiveContainer>
            )}
          </CardContent>
        </Card>

        {/* Server Performance */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <Card>
            <CardHeader>
              <CardTitle>Server FPS</CardTitle>
              <CardDescription>Server performance over time</CardDescription>
            </CardHeader>
            <CardContent>
              {serverLoading ? (
                <div className="h-[250px] flex items-center justify-center text-muted-foreground">
                  Loading...
                </div>
              ) : !serverStats || serverStats.length === 0 ? (
                <div className="h-[250px] flex items-center justify-center text-muted-foreground">
                  No data available for selected time range
                </div>
              ) : (
                <ResponsiveContainer width="100%" height={250}>
                  <LineChart data={serverStats}>
                    <CartesianGrid
                      strokeDasharray="3 3"
                      className="stroke-muted"
                    />
                    <XAxis
                      dataKey="timestamp"
                      tickFormatter={
                        timeRange === "7d" || timeRange === "30d"
                          ? formatDate
                          : formatTime
                      }
                      className="text-xs"
                    />
                    <YAxis className="text-xs" />
                    <Tooltip
                      contentStyle={{
                        backgroundColor: "var(--card)",
                        border: "1px solid var(--border)",
                        borderRadius: "6px",
                      }}
                    />
                    <Legend />
                    <Line
                      type="monotone"
                      dataKey="fps"
                      stroke="#10b981"
                      strokeWidth={2}
                      name="FPS"
                    />
                  </LineChart>
                </ResponsiveContainer>
              )}
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Player Statistics</CardTitle>
              <CardDescription>Average ping and score</CardDescription>
            </CardHeader>
            <CardContent>
              {playerLoading ? (
                <div className="h-[250px] flex items-center justify-center text-muted-foreground">
                  Loading...
                </div>
              ) : !playerStats || playerStats.length === 0 ? (
                <div className="h-[250px] flex items-center justify-center text-muted-foreground">
                  No data available for selected time range
                </div>
              ) : (
                <ResponsiveContainer width="100%" height={250}>
                  <BarChart data={playerStats}>
                    <CartesianGrid
                      strokeDasharray="3 3"
                      className="stroke-muted"
                    />
                    <XAxis
                      dataKey="timestamp"
                      tickFormatter={
                        timeRange === "7d" || timeRange === "30d"
                          ? formatDate
                          : formatTime
                      }
                      className="text-xs"
                    />
                    <YAxis className="text-xs" />
                    <Tooltip
                      contentStyle={{
                        backgroundColor: "var(--card)",
                        border: "1px solid var(--border)",
                        borderRadius: "6px",
                        color: "var(--foreground)",
                      }}
                      labelStyle={{
                        color: "var(--foreground)",
                      }}
                      itemStyle={{
                        color: "var(--foreground)",
                      }}
                      labelFormatter={(label) =>
                        format(new Date(label), "PPpp")
                      }
                      cursor={{ fill: "var(--muted)" }}
                    />
                    <Legend />
                    <Bar
                      dataKey="avgPing"
                      fill="#f59e0b"
                      name="Avg Ping (ms)"
                    />
                    <Bar dataKey="avgScore" fill="#8b5cf6" name="Avg Score" />
                  </BarChart>
                </ResponsiveContainer>
              )}
            </CardContent>
          </Card>
        </div>

        {/* System Stats (if available) */}
        {systemStats && systemStats.length > 0 && (
          <Card>
            <CardHeader>
              <CardTitle>System Resources</CardTitle>
              <CardDescription>Active game server memory usage</CardDescription>
            </CardHeader>
            <CardContent>
              {systemLoading ? (
                <div className="h-[250px] flex items-center justify-center text-muted-foreground">
                  Loading...
                </div>
              ) : !systemStats || systemStats.length === 0 ? (
                <div className="h-[250px] flex items-center justify-center text-muted-foreground">
                  No data available for selected time range
                </div>
              ) : (
                <ResponsiveContainer width="100%" height={250}>
                  <LineChart data={systemStats}>
                    <CartesianGrid
                      strokeDasharray="3 3"
                      className="stroke-muted"
                    />
                    <XAxis
                      dataKey="timestamp"
                      tickFormatter={
                        timeRange === "7d" || timeRange === "30d"
                          ? formatDate
                          : formatTime
                      }
                      className="text-xs"
                    />
                    <YAxis
                      className="text-xs"
                      tickFormatter={(value) =>
                        `${Math.round(value / (1024 * 1024))} MB`
                      }
                    />
                    <Tooltip
                      contentStyle={{
                        backgroundColor: "var(--card)",
                        border: "1px solid var(--border)",
                        borderRadius: "6px",
                      }}
                      formatter={(value: number, name: string) => [
                        formatMemoryDisplay(value),
                        name,
                      ]}
                      labelFormatter={(label) =>
                        format(new Date(label), "PPpp")
                      }
                    />
                    <Legend formatter={(value) => value} />
                    <Line
                      type="monotone"
                      dataKey="memoryUsed"
                      stroke="#ef4444"
                      strokeWidth={2}
                      name="Active Memory Usage"
                    />
                  </LineChart>
                </ResponsiveContainer>
              )}
            </CardContent>
          </Card>
        )}
      </div>
    </ProtectedRoute>
  );
}

export default Analytics;
