import { ProtectedRoute } from "@/components/ProtectedRoute";
import { useStatus } from "@/hooks/useStatus";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Server, Users, Map, Gamepad2 } from "lucide-react";

function Status() {
  const { data: status, isLoading, isError } = useStatus();

  return (
    <ProtectedRoute requiredPermission="status.view">
      <div className="space-y-6 bg-background min-h-screen">
        <div>
          <h1 className="text-4xl font-bold text-foreground mb-2">
            Server Status
          </h1>
          <p className="text-muted-foreground">Real-time server information</p>
        </div>

        {isLoading ? (
          <div className="flex items-center justify-center h-64">
            <p className="text-foreground text-lg">Loading server status...</p>
          </div>
        ) : isError ? (
          <div className="flex items-center justify-center h-64">
            <p className="text-destructive text-lg">
              Error loading server status.
            </p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <Card className="bg-card border-border">
              <CardHeader>
                <div className="flex items-center space-x-2">
                  <Server className="h-5 w-5 text-primary" />
                  <CardTitle className="text-foreground">
                    Server Information
                  </CardTitle>
                </div>
                <CardDescription className="text-muted-foreground">
                  Basic server details
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <div className="text-sm text-muted-foreground mb-1">
                    Hostname
                  </div>
                  <div className="text-foreground font-mono bg-muted/30 p-2 rounded border border-border">
                    {status?.hostname}
                  </div>
                </div>
                <div>
                  <div className="text-sm text-muted-foreground mb-1">
                    Address
                  </div>
                  <div className="text-foreground font-mono bg-muted/30 p-2 rounded border border-border">
                    {status?.address}
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card className="bg-card border-border">
              <CardHeader>
                <div className="flex items-center space-x-2">
                  <Gamepad2 className="h-5 w-5 text-primary" />
                  <CardTitle className="text-foreground">
                    Game Details
                  </CardTitle>
                </div>
                <CardDescription className="text-muted-foreground">
                  Current map and game mode
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-2">
                    <Map className="h-4 w-4 text-muted-foreground" />
                    <span className="text-muted-foreground">Current Map:</span>
                  </div>
                  <Badge variant="outline">{status?.map}</Badge>
                </div>
                {status?.gametype && (
                  <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-2">
                      <Gamepad2 className="h-4 w-4 text-muted-foreground" />
                      <span className="text-muted-foreground">Game Type:</span>
                    </div>
                    <Badge variant="outline">{status?.gametype}</Badge>
                  </div>
                )}
              </CardContent>
            </Card>

            <Card className="bg-card border-border md:col-span-2">
              <CardHeader>
                <div className="flex items-center space-x-2">
                  <Users className="h-5 w-5 text-primary" />
                  <CardTitle className="text-foreground">
                    Player Statistics
                  </CardTitle>
                </div>
                <CardDescription className="text-muted-foreground">
                  Overview of connected players
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                  <div className="text-center p-6 bg-muted/30 rounded-lg border border-border">
                    <div className="text-4xl font-bold text-foreground mb-2">
                      {status?.players?.length || 0}
                    </div>
                    <div className="text-sm text-muted-foreground">
                      Players Online
                    </div>
                  </div>
                  <div className="text-center p-6 bg-muted/30 rounded-lg border border-border">
                    <div className="text-4xl font-bold text-foreground mb-2">
                      {status?.players?.length
                        ? Math.round(
                            status.players.reduce((avg, p) => avg + p.ping, 0) /
                              status.players.length
                          )
                        : 0}
                      ms
                    </div>
                    <div className="text-sm text-muted-foreground">
                      Average Ping
                    </div>
                  </div>
                  <div className="text-center p-6 bg-muted/30 rounded-lg border border-border">
                    <div className="text-4xl font-bold text-foreground mb-2">
                      {status?.players?.length
                        ? Math.max(...status.players.map((p) => p.score))
                        : 0}
                    </div>
                    <div className="text-sm text-muted-foreground">
                      Top Score
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        )}
      </div>
    </ProtectedRoute>
  );
}

export default Status;
