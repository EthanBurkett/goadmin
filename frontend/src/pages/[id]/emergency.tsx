import { useDisabledCommands, useReenableCommand } from "@/hooks/useEmergency";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { AlertTriangle, CheckCircle, Clock, ShieldAlert } from "lucide-react";
import { formatDistanceToNow } from "date-fns";

export default function EmergencyPage() {
  const { data: disabledCommands, isLoading } = useDisabledCommands();
  const reenableMutation = useReenableCommand();

  const handleReenable = (command: string) => {
    if (
      confirm(
        `Are you sure you want to manually re-enable the "${command}" command?`
      )
    ) {
      reenableMutation.mutate(command);
    }
  };

  const formatTimeRemaining = (reenableAt: string) => {
    const reenableTime = new Date(reenableAt);
    const now = new Date();

    if (reenableTime <= now) {
      return "Re-enabling soon...";
    }

    return `Re-enables ${formatDistanceToNow(reenableTime, {
      addSuffix: true,
    })}`;
  };

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div>
          <h1 className="text-4xl font-bold">Emergency Command Status</h1>
          <p className="text-muted-foreground mt-2">
            Monitor and manage automatically disabled commands
          </p>
        </div>
        <Skeleton className="h-48" />
      </div>
    );
  }

  const disabledCommandsList = disabledCommands
    ? Object.entries(disabledCommands)
    : [];

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-4xl font-bold">Emergency Command Status</h1>
        <p className="text-muted-foreground mt-2">
          Monitor and manage automatically disabled commands
        </p>
      </div>

      {/* Status Overview */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <ShieldAlert className="h-5 w-5 text-orange-500" />
              System Status
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">
                  Disabled Commands
                </span>
                <Badge
                  variant={
                    disabledCommandsList.length > 0 ? "destructive" : "outline"
                  }
                >
                  {disabledCommandsList.length}
                </Badge>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">
                  Protection Status
                </span>
                {disabledCommandsList.length > 0 ? (
                  <Badge variant="default" className="bg-orange-500">
                    <AlertTriangle className="h-3 w-3 mr-1" />
                    Active
                  </Badge>
                ) : (
                  <Badge variant="outline" className="text-green-500">
                    <CheckCircle className="h-3 w-3 mr-1" />
                    Normal
                  </Badge>
                )}
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <AlertTriangle className="h-5 w-5 text-yellow-500" />
              About Emergency Shutdown
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">
              Commands are automatically disabled when abuse patterns are
              detected (e.g., excessive ban loops). This protects your server
              from malicious activity and prevents command spam.
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Disabled Commands */}
      {disabledCommandsList.length === 0 ? (
        <Card>
          <CardContent className="pt-6">
            <div className="text-center py-8">
              <CheckCircle className="h-12 w-12 text-green-500 mx-auto mb-4" />
              <h3 className="text-lg font-semibold mb-2">All Systems Normal</h3>
              <p className="text-muted-foreground">
                No commands are currently disabled. Your server is operating
                normally.
              </p>
            </div>
          </CardContent>
        </Card>
      ) : (
        <div className="space-y-4">
          <h2 className="text-2xl font-bold">Disabled Commands</h2>
          {disabledCommandsList.map(([command, info]) => (
            <Card key={command} className="border-destructive">
              <CardHeader>
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <CardTitle className="flex items-center gap-2">
                      <AlertTriangle className="h-5 w-5 text-destructive" />
                      Command: {command}
                    </CardTitle>
                    <CardDescription className="mt-2">
                      {info.reason}
                    </CardDescription>
                  </div>
                  <Badge variant="destructive">Disabled</Badge>
                </div>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {/* Details Grid */}
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
                    <div>
                      <span className="text-muted-foreground">
                        Disabled At:
                      </span>
                      <p className="font-medium">
                        {new Date(info.disabledAt).toLocaleString()}
                      </p>
                    </div>
                    <div>
                      <span className="text-muted-foreground">
                        Disabled By:
                      </span>
                      <p className="font-medium">{info.disabledBy}</p>
                    </div>
                    {info.autoRenable && (
                      <div className="col-span-1 md:col-span-2">
                        <span className="text-muted-foreground flex items-center gap-2">
                          <Clock className="h-4 w-4" />
                          Auto Re-enable:
                        </span>
                        <p className="font-medium">
                          {formatTimeRemaining(info.reenableAt)}
                        </p>
                      </div>
                    )}
                  </div>

                  {/* Actions */}
                  <div className="flex gap-2 pt-2 border-t">
                    <Button
                      variant="default"
                      size="sm"
                      onClick={() => handleReenable(command)}
                      disabled={reenableMutation.isPending}
                    >
                      {reenableMutation.isPending
                        ? "Re-enabling..."
                        : "Re-enable Now"}
                    </Button>
                    {info.autoRenable && (
                      <Button variant="outline" size="sm" disabled>
                        <Clock className="h-4 w-4 mr-2" />
                        Wait for Auto Re-enable
                      </Button>
                    )}
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
