import { useState } from "react";
import { useAuditLogs, type AuditLogsFilters } from "@/hooks/useAudit";
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Skeleton } from "@/components/ui/skeleton";
import { formatDistanceToNow } from "@/lib/time";
import { ProtectedRoute } from "@/components/ProtectedRoute";

const ACTION_TYPES = [
  "all",
  "ban",
  "tempban",
  "unban",
  "kick",
  "rcon_command",
  "login",
  "logout",
  "login_failed",
  "user_approved",
  "user_rejected",
  "role_assigned",
  "role_removed",
  "report_created",
  "report_resolved",
  "report_dismissed",
];

const SOURCE_TYPES = ["all", "web_ui", "in_game", "api", "system"];

export default function AuditPage() {
  const [filters, setFilters] = useState<AuditLogsFilters>({
    limit: 50,
    offset: 0,
  });

  const [_searchUsername, setSearchUsername] = useState("");
  const [selectedAction, setSelectedAction] = useState("all");
  const [selectedSource, setSelectedSource] = useState("all");
  const [selectedSuccess, setSelectedSuccess] = useState("all");
  const [startDate, setStartDate] = useState("");
  const [endDate, setEndDate] = useState("");

  const { data, isLoading, error } = useAuditLogs(filters);

  const applyFilters = () => {
    const newFilters: AuditLogsFilters = {
      limit: 50,
      offset: 0,
    };

    if (selectedAction !== "all") newFilters.action = selectedAction;
    if (selectedSource !== "all") newFilters.source = selectedSource;
    if (selectedSuccess === "true") newFilters.success = true;
    if (selectedSuccess === "false") newFilters.success = false;
    if (startDate) newFilters.start_date = startDate;
    if (endDate) newFilters.end_date = endDate;

    setFilters(newFilters);
  };

  const resetFilters = () => {
    setSearchUsername("");
    setSelectedAction("all");
    setSelectedSource("all");
    setSelectedSuccess("all");
    setStartDate("");
    setEndDate("");
    setFilters({ limit: 50, offset: 0 });
  };

  const loadMore = () => {
    setFilters((prev) => ({
      ...prev,
      offset: (prev.offset || 0) + (prev.limit || 50),
    }));
  };

  const loadPrevious = () => {
    setFilters((prev) => ({
      ...prev,
      offset: Math.max(0, (prev.offset || 0) - (prev.limit || 50)),
    }));
  };

  const exportToCSV = () => {
    if (!data?.logs || data.logs.length === 0) return;

    const headers = [
      "Timestamp",
      "User",
      "IP Address",
      "Action",
      "Source",
      "Success",
      "Target",
      "Details",
    ];
    const rows = data.logs.map((log) => [
      new Date(log.createdAt).toISOString(),
      log.username || "System",
      log.ipAddress || "N/A",
      log.action,
      log.source,
      log.success ? "Yes" : "No",
      log.targetName || log.targetId || "N/A",
      log.errorMessage || log.metadata || log.result || "",
    ]);

    const csvContent = [
      headers.join(","),
      ...rows.map((row) => row.map((cell) => `"${cell}"`).join(",")),
    ].join("\n");

    const blob = new Blob([csvContent], { type: "text/csv" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `audit-logs-${new Date().toISOString()}.csv`;
    a.click();
    URL.revokeObjectURL(url);
  };

  const exportToJSON = () => {
    if (!data?.logs || data.logs.length === 0) return;

    const jsonContent = JSON.stringify(data.logs, null, 2);
    const blob = new Blob([jsonContent], { type: "application/json" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `audit-logs-${new Date().toISOString()}.json`;
    a.click();
    URL.revokeObjectURL(url);
  };

  return (
    <ProtectedRoute requiredPermission="audit.view">
      <div className="space-y-6">
        <div>
          <h1 className="text-4xl font-bold text-foreground mb-2">
            Audit Logs
          </h1>
          <p className="text-muted-foreground">
            View and filter all system actions and events
          </p>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>Filters</CardTitle>
            <CardDescription>
              Filter audit logs by action, source, date range, and more
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              <div className="space-y-2">
                <Label htmlFor="action">Action Type</Label>
                <Select
                  value={selectedAction}
                  onValueChange={setSelectedAction}
                >
                  <SelectTrigger id="action">
                    <SelectValue placeholder="Select action" />
                  </SelectTrigger>
                  <SelectContent>
                    {ACTION_TYPES.map((action) => (
                      <SelectItem key={action} value={action}>
                        {action === "all"
                          ? "All Actions"
                          : action.replace(/_/g, " ").toUpperCase()}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-2">
                <Label htmlFor="source">Source</Label>
                <Select
                  value={selectedSource}
                  onValueChange={setSelectedSource}
                >
                  <SelectTrigger id="source">
                    <SelectValue placeholder="Select source" />
                  </SelectTrigger>
                  <SelectContent>
                    {SOURCE_TYPES.map((source) => (
                      <SelectItem key={source} value={source}>
                        {source === "all"
                          ? "All Sources"
                          : source.replace(/_/g, " ").toUpperCase()}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-2">
                <Label htmlFor="success">Status</Label>
                <Select
                  value={selectedSuccess}
                  onValueChange={setSelectedSuccess}
                >
                  <SelectTrigger id="success">
                    <SelectValue placeholder="Select status" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All</SelectItem>
                    <SelectItem value="true">Success</SelectItem>
                    <SelectItem value="false">Failed</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-2">
                <Label htmlFor="start-date">Start Date</Label>
                <Input
                  id="start-date"
                  type="datetime-local"
                  value={startDate}
                  onChange={(e) => setStartDate(e.target.value)}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="end-date">End Date</Label>
                <Input
                  id="end-date"
                  type="datetime-local"
                  value={endDate}
                  onChange={(e) => setEndDate(e.target.value)}
                />
              </div>
            </div>

            <div className="flex gap-2 mt-4">
              <Button onClick={applyFilters}>Apply Filters</Button>
              <Button variant="outline" onClick={resetFilters}>
                Reset
              </Button>
              <div className="ml-auto flex gap-2">
                <Button
                  variant="outline"
                  onClick={exportToCSV}
                  disabled={!data?.logs?.length}
                >
                  Export CSV
                </Button>
                <Button
                  variant="outline"
                  onClick={exportToJSON}
                  disabled={!data?.logs?.length}
                >
                  Export JSON
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Audit Log Entries</CardTitle>
            <CardDescription>
              {data?.total !== undefined
                ? `Showing ${data.logs?.length || 0} of ${data.total} entries`
                : "Loading..."}
            </CardDescription>
          </CardHeader>
          <CardContent>
            {error && (
              <div className="text-red-500 p-4 border border-red-200 rounded">
                Error loading audit logs:{" "}
                {error instanceof Error ? error.message : "Unknown error"}
              </div>
            )}

            {isLoading ? (
              <div className="space-y-2">
                {[...Array(5)].map((_, i) => (
                  <Skeleton key={i} className="h-12 w-full" />
                ))}
              </div>
            ) : (
              <>
                <div className="rounded-md border">
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>Timestamp</TableHead>
                        <TableHead>User</TableHead>
                        <TableHead>IP Address</TableHead>
                        <TableHead>Action</TableHead>
                        <TableHead>Source</TableHead>
                        <TableHead>Status</TableHead>
                        <TableHead>Target</TableHead>
                        <TableHead>Details</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {!data?.logs || data.logs.length === 0 ? (
                        <TableRow>
                          <TableCell
                            colSpan={8}
                            className="text-center text-muted-foreground"
                          >
                            No audit logs found
                          </TableCell>
                        </TableRow>
                      ) : (
                        data.logs.map((log) => (
                          <TableRow key={log.id}>
                            <TableCell className="whitespace-nowrap">
                              <div className="text-sm">
                                {new Date(log.createdAt).toLocaleString()}
                              </div>
                              <div className="text-xs text-muted-foreground">
                                {formatDistanceToNow(new Date(log.createdAt))}
                              </div>
                            </TableCell>
                            <TableCell>
                              {log.username || (
                                <span className="text-muted-foreground">
                                  System
                                </span>
                              )}
                            </TableCell>
                            <TableCell className="font-mono text-xs">
                              {log.ipAddress || "N/A"}
                            </TableCell>
                            <TableCell>
                              <Badge variant="outline">
                                {log.action.replace(/_/g, " ").toUpperCase()}
                              </Badge>
                            </TableCell>
                            <TableCell>
                              <Badge variant="secondary">
                                {log.source.replace(/_/g, " ").toUpperCase()}
                              </Badge>
                            </TableCell>
                            <TableCell>
                              {log.success ? (
                                <Badge className="bg-green-500">Success</Badge>
                              ) : (
                                <Badge variant="destructive">Failed</Badge>
                              )}
                            </TableCell>
                            <TableCell>
                              {log.targetName || log.targetId || (
                                <span className="text-muted-foreground">
                                  N/A
                                </span>
                              )}
                            </TableCell>
                            <TableCell className="max-w-xs">
                              <div className="truncate text-sm">
                                {log.errorMessage ||
                                  log.metadata ||
                                  log.result || (
                                    <span className="text-muted-foreground">
                                      -
                                    </span>
                                  )}
                              </div>
                            </TableCell>
                          </TableRow>
                        ))
                      )}
                    </TableBody>
                  </Table>
                </div>

                {data && data.logs && data.logs.length > 0 && (
                  <div className="flex items-center justify-between mt-4">
                    <div className="text-sm text-muted-foreground">
                      Showing {(filters.offset || 0) + 1} to{" "}
                      {Math.min(
                        (filters.offset || 0) + (filters.limit || 50),
                        data.total
                      )}{" "}
                      of {data.total} entries
                    </div>
                    <div className="flex gap-2">
                      <Button
                        variant="outline"
                        onClick={loadPrevious}
                        disabled={(filters.offset || 0) === 0}
                      >
                        Previous
                      </Button>
                      <Button
                        variant="outline"
                        onClick={loadMore}
                        disabled={
                          (filters.offset || 0) + (filters.limit || 50) >=
                          data.total
                        }
                      >
                        Next
                      </Button>
                    </div>
                  </div>
                )}
              </>
            )}
          </CardContent>
        </Card>
      </div>
    </ProtectedRoute>
  );
}
