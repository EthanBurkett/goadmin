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
import { Badge } from "@/components/ui/badge";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
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
  useReports,
  useActionReport,
  useDeleteReport,
  useTempBans,
  useRevokeTempBan,
  type Report,
} from "@/hooks/useReports";
import { useHasPermission } from "@/hooks/useRbac";
import {
  AlertTriangle,
  Ban,
  CheckCircle,
  Clock,
  XCircle,
  Trash2,
  Shield,
  Calendar,
} from "lucide-react";
import { formatDistanceToNow } from "date-fns";

function Reports() {
  const { data: reports, isLoading: reportsLoading } = useReports();
  const { data: tempBans, isLoading: tempBansLoading } = useTempBans();
  const actionReport = useActionReport();
  const deleteReport = useDeleteReport();
  const revokeTempBan = useRevokeTempBan();
  const { hasPermission } = useHasPermission();

  const [selectedReport, setSelectedReport] = useState<Report | null>(null);
  const [isActionDialogOpen, setIsActionDialogOpen] = useState(false);
  const [actionType, setActionType] = useState<"dismiss" | "ban" | "tempban">(
    "dismiss"
  );
  const [actionReason, setActionReason] = useState("");
  const [tempBanDuration, setTempBanDuration] = useState("24");

  const canViewReports = hasPermission("reports.view");
  const canActionReports = hasPermission("reports.action");

  const handleOpenActionDialog = (report: Report) => {
    setSelectedReport(report);
    setActionReason("");
    setActionType("dismiss");
    setTempBanDuration("24");
    setIsActionDialogOpen(true);
  };

  const handleAction = async () => {
    if (!selectedReport) return;

    try {
      await actionReport.mutateAsync({
        id: selectedReport.id,
        action: actionType,
        duration:
          actionType === "tempban" ? parseInt(tempBanDuration) : undefined,
        reason: actionReason,
      });
      setIsActionDialogOpen(false);
      setSelectedReport(null);
    } catch (error) {
      console.error("Failed to action report:", error);
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm("Are you sure you want to delete this report?")) return;
    try {
      await deleteReport.mutateAsync(id);
    } catch (error) {
      console.error("Failed to delete report:", error);
    }
  };

  const handleRevokeTempBan = async (id: number) => {
    if (!confirm("Are you sure you want to revoke this temporary ban?")) return;
    try {
      await revokeTempBan.mutateAsync(id);
    } catch (error) {
      console.error("Failed to revoke temp ban:", error);
    }
  };

  const getStatusBadge = (status: Report["status"]) => {
    switch (status) {
      case "pending":
        return (
          <Badge variant="outline" className="text-yellow-500">
            <Clock className="h-3 w-3 mr-1" />
            Pending
          </Badge>
        );
      case "actioned":
        return (
          <Badge variant="outline" className="text-green-500">
            <CheckCircle className="h-3 w-3 mr-1" />
            Actioned
          </Badge>
        );
      case "dismissed":
        return (
          <Badge variant="outline" className="text-gray-500">
            <XCircle className="h-3 w-3 mr-1" />
            Dismissed
          </Badge>
        );
      default:
        return null;
    }
  };

  const pendingReports = reports?.filter((r) => r.status === "pending") || [];
  const reviewedReports = reports?.filter((r) => r.status !== "pending") || [];
  const activeTempBans = tempBans?.filter((b) => b.active) || [];

  return (
    <ProtectedRoute requiredPermission="reports.view">
      <div className="space-y-6 bg-background min-h-screen">
        <div>
          <h1 className="text-4xl font-bold text-foreground mb-2">
            Reports & Moderation
          </h1>
          <p className="text-muted-foreground">
            Review player reports and manage temporary bans
          </p>
        </div>

        {!canViewReports ? (
          <Card>
            <CardContent className="pt-6">
              <div className="text-center text-muted-foreground">
                <Shield className="h-12 w-12 mx-auto mb-2 opacity-50" />
                <p>You don't have permission to view reports</p>
              </div>
            </CardContent>
          </Card>
        ) : (
          <>
            {/* Pending Reports */}
            <Card>
              <CardHeader>
                <div className="flex items-center justify-between">
                  <div>
                    <CardTitle className="flex items-center gap-2">
                      <AlertTriangle className="h-5 w-5 text-yellow-500" />
                      Pending Reports
                    </CardTitle>
                    <CardDescription>
                      Reports awaiting admin review ({pendingReports.length})
                    </CardDescription>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                {reportsLoading ? (
                  <div className="text-center py-8 text-muted-foreground">
                    Loading reports...
                  </div>
                ) : pendingReports.length === 0 ? (
                  <div className="text-center py-8 text-muted-foreground">
                    <CheckCircle className="h-12 w-12 mx-auto mb-2 opacity-50" />
                    <p>No pending reports</p>
                  </div>
                ) : (
                  <div className="space-y-4">
                    {pendingReports.map((report) => (
                      <Card
                        key={report.id}
                        className="border-l-4 border-l-yellow-500"
                      >
                        <CardContent className="pt-6">
                          <div className="space-y-4">
                            <div className="flex items-start justify-between">
                              <div className="space-y-2 flex-1">
                                <div className="flex items-center gap-2">
                                  <Badge variant="destructive">
                                    Reported: {report.reportedName}
                                  </Badge>
                                  {getStatusBadge(report.status)}
                                  <span className="text-xs text-muted-foreground">
                                    #{report.id}
                                  </span>
                                </div>
                                <div className="text-sm">
                                  <span className="text-muted-foreground">
                                    Reporter:
                                  </span>{" "}
                                  {report.reporterName}
                                </div>
                                <div className="text-sm">
                                  <span className="text-muted-foreground">
                                    Reason:
                                  </span>{" "}
                                  {report.reason}
                                </div>
                                <div className="text-xs text-muted-foreground">
                                  {formatDistanceToNow(
                                    new Date(report.createdAt),
                                    { addSuffix: true }
                                  )}
                                </div>
                              </div>
                              {canActionReports && (
                                <div className="flex gap-2">
                                  <Button
                                    size="sm"
                                    onClick={() =>
                                      handleOpenActionDialog(report)
                                    }
                                  >
                                    Take Action
                                  </Button>
                                  <Button
                                    size="sm"
                                    variant="destructive"
                                    onClick={() => handleDelete(report.id)}
                                  >
                                    <Trash2 className="h-3 w-3" />
                                  </Button>
                                </div>
                              )}
                            </div>
                            <div className="flex gap-2 text-xs text-muted-foreground">
                              <code className="bg-muted px-2 py-1 rounded">
                                Reported GUID: {report.reportedGuid}
                              </code>
                              <code className="bg-muted px-2 py-1 rounded">
                                Reporter GUID: {report.reporterGuid}
                              </code>
                            </div>
                          </div>
                        </CardContent>
                      </Card>
                    ))}
                  </div>
                )}
              </CardContent>
            </Card>

            {/* Active Temp Bans */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Ban className="h-5 w-5 text-red-500" />
                  Active Temporary Bans
                </CardTitle>
                <CardDescription>
                  Currently active temporary bans ({activeTempBans.length})
                </CardDescription>
              </CardHeader>
              <CardContent>
                {tempBansLoading ? (
                  <div className="text-center py-8 text-muted-foreground">
                    Loading temp bans...
                  </div>
                ) : activeTempBans.length === 0 ? (
                  <div className="text-center py-8 text-muted-foreground">
                    No active temporary bans
                  </div>
                ) : (
                  <div className="space-y-4">
                    {activeTempBans.map((ban) => (
                      <Card key={ban.id}>
                        <CardContent className="pt-6">
                          <div className="flex items-start justify-between">
                            <div className="space-y-2 flex-1">
                              <div className="flex items-center gap-2">
                                <Badge variant="destructive">
                                  {ban.playerName}
                                </Badge>
                                <Badge variant="outline">
                                  <Calendar className="h-3 w-3 mr-1" />
                                  Expires{" "}
                                  {formatDistanceToNow(
                                    new Date(ban.expiresAt),
                                    {
                                      addSuffix: true,
                                    }
                                  )}
                                </Badge>
                              </div>
                              <div className="text-sm">
                                <span className="text-muted-foreground">
                                  Reason:
                                </span>{" "}
                                {ban.reason}
                              </div>
                              {ban.bannedBy && (
                                <div className="text-sm">
                                  <span className="text-muted-foreground">
                                    Banned by:
                                  </span>{" "}
                                  {ban.bannedBy.username}
                                </div>
                              )}
                              <div className="text-xs text-muted-foreground">
                                Created{" "}
                                {formatDistanceToNow(new Date(ban.createdAt), {
                                  addSuffix: true,
                                })}
                              </div>
                              <code className="text-xs bg-muted px-2 py-1 rounded block w-fit">
                                GUID: {ban.playerGuid}
                              </code>
                            </div>
                            {canActionReports && (
                              <Button
                                size="sm"
                                variant="outline"
                                onClick={() => handleRevokeTempBan(ban.id)}
                              >
                                Revoke Ban
                              </Button>
                            )}
                          </div>
                        </CardContent>
                      </Card>
                    ))}
                  </div>
                )}
              </CardContent>
            </Card>

            {/* Reviewed Reports */}
            <Card>
              <CardHeader>
                <CardTitle>Reviewed Reports</CardTitle>
                <CardDescription>
                  Previously actioned or dismissed reports (
                  {reviewedReports.length})
                </CardDescription>
              </CardHeader>
              <CardContent>
                {reviewedReports.length === 0 ? (
                  <div className="text-center py-8 text-muted-foreground">
                    No reviewed reports
                  </div>
                ) : (
                  <div className="space-y-4">
                    {reviewedReports.map((report) => (
                      <Card key={report.id} className="opacity-75">
                        <CardContent className="pt-6">
                          <div className="space-y-2">
                            <div className="flex items-center gap-2">
                              <Badge variant="outline">
                                {report.reportedName}
                              </Badge>
                              {getStatusBadge(report.status)}
                              <span className="text-xs text-muted-foreground">
                                #{report.id}
                              </span>
                            </div>
                            <div className="text-sm">
                              <span className="text-muted-foreground">
                                Reporter:
                              </span>{" "}
                              {report.reporterName}
                            </div>
                            <div className="text-sm">
                              <span className="text-muted-foreground">
                                Reason:
                              </span>{" "}
                              {report.reason}
                            </div>
                            {report.actionTaken && (
                              <div className="text-sm">
                                <span className="text-muted-foreground">
                                  Action Taken:
                                </span>{" "}
                                {report.actionTaken}
                              </div>
                            )}
                            {report.reviewedBy && (
                              <div className="text-sm">
                                <span className="text-muted-foreground">
                                  Reviewed by:
                                </span>{" "}
                                {report.reviewedBy.username}
                              </div>
                            )}
                            <div className="text-xs text-muted-foreground">
                              {formatDistanceToNow(new Date(report.createdAt), {
                                addSuffix: true,
                              })}
                            </div>
                          </div>
                        </CardContent>
                      </Card>
                    ))}
                  </div>
                )}
              </CardContent>
            </Card>
          </>
        )}
      </div>

      {/* Action Dialog */}
      <Dialog open={isActionDialogOpen} onOpenChange={setIsActionDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Take Action on Report</DialogTitle>
            <DialogDescription>
              Review and take action on report #{selectedReport?.id}
            </DialogDescription>
          </DialogHeader>
          {selectedReport && (
            <div className="space-y-4">
              <div className="p-4 bg-muted rounded-lg space-y-2">
                <div className="text-sm">
                  <span className="font-semibold">Reported Player:</span>{" "}
                  {selectedReport.reportedName}
                </div>
                <div className="text-sm">
                  <span className="font-semibold">Reporter:</span>{" "}
                  {selectedReport.reporterName}
                </div>
                <div className="text-sm">
                  <span className="font-semibold">Reason:</span>{" "}
                  {selectedReport.reason}
                </div>
              </div>

              <div className="space-y-2">
                <Label htmlFor="actionType">Action Type</Label>
                <Select
                  value={actionType}
                  onValueChange={(value) =>
                    setActionType(value as "dismiss" | "ban" | "tempban")
                  }
                >
                  <SelectTrigger id="actionType">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="dismiss">Dismiss Report</SelectItem>
                    <SelectItem value="tempban">Temporary Ban</SelectItem>
                    <SelectItem value="ban">Permanent Ban</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              {actionType === "tempban" && (
                <div className="space-y-2">
                  <Label htmlFor="duration">Duration (hours)</Label>
                  <Input
                    id="duration"
                    type="number"
                    min="1"
                    value={tempBanDuration}
                    onChange={(e) => setTempBanDuration(e.target.value)}
                    placeholder="24"
                  />
                </div>
              )}

              <div className="space-y-2">
                <Label htmlFor="actionReason">
                  {actionType === "dismiss"
                    ? "Reason for Dismissal"
                    : "Ban Reason"}
                </Label>
                <Textarea
                  id="actionReason"
                  value={actionReason}
                  onChange={(e) => setActionReason(e.target.value)}
                  placeholder="Enter reason..."
                  rows={3}
                />
              </div>
            </div>
          )}
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setIsActionDialogOpen(false)}
            >
              Cancel
            </Button>
            <Button
              onClick={handleAction}
              disabled={!actionReason.trim() || actionReport.isPending}
            >
              {actionReport.isPending ? "Processing..." : "Confirm Action"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </ProtectedRoute>
  );
}

export default Reports;
