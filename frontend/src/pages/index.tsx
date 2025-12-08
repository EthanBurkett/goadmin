import { ProtectedRoute } from "@/components/ProtectedRoute";
import { DashboardLayout } from "@/components/DashboardLayout";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useStatus } from "@/hooks/useStatus";

function App() {
  const { data: status } = useStatus();

  return (
    <ProtectedRoute>
      <DashboardLayout>
        <div className="p-8 space-y-8 bg-background min-h-screen">
          <div>
            <h1 className="text-4xl font-bold text-foreground mb-2">
              Dashboard
            </h1>
            <p className="text-muted-foreground">
              Welcome to your CoD4 server admin panel
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <Card className="bg-card border-border">
              <CardHeader>
                <CardTitle className="text-foreground">Server Status</CardTitle>
                <CardDescription className="text-muted-foreground">
                  Current server information
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-2">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Hostname:</span>
                  <span className="text-foreground font-mono text-sm">
                    {status?.hostname || "Loading..."}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Map:</span>
                  <span className="text-foreground">
                    {status?.map || "Loading..."}
                  </span>
                </div>
                {status?.gametype && (
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Gametype:</span>
                    <span className="text-foreground">
                      {status?.gametype || "Loading..."}
                    </span>
                  </div>
                )}
              </CardContent>
            </Card>

            <Card className="bg-card border-border">
              <CardHeader>
                <CardTitle className="text-foreground">
                  Players Online
                </CardTitle>
                <CardDescription className="text-muted-foreground">
                  Currently connected players
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="text-4xl font-bold text-foreground">
                  {status?.players?.length ?? "..."}
                </div>
                <p className="text-sm text-muted-foreground mt-2">
                  Active players
                </p>
              </CardContent>
            </Card>

            <Card className="bg-card border-border">
              <CardHeader>
                <CardTitle className="text-foreground">Quick Actions</CardTitle>
                <CardDescription className="text-muted-foreground">
                  Manage your server
                </CardDescription>
              </CardHeader>
              <CardContent>
                <p className="text-muted-foreground text-sm">
                  Use the sidebar to navigate to different management areas
                </p>
              </CardContent>
            </Card>
          </div>
        </div>
      </DashboardLayout>
    </ProtectedRoute>
  );
}

export default App;
