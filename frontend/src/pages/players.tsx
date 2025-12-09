import { ProtectedRoute } from "@/components/ProtectedRoute";
import { DashboardLayout } from "@/components/DashboardLayout";
import { DataTable } from "@/components/DataTable";
import { usePlayers } from "@/hooks/usePlayers";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import type { ColumnDef } from "@tanstack/react-table";
import { ArrowUpDown } from "lucide-react";
import { Button } from "@/components/ui/button";

interface Player {
  id: number;
  score: number;
  ping: number;
  uuid: string;
  steamId: string;
  name: string;
  strippedName: string;
  address: string;
  rate: number;
  qPort: number;
}

const columns: ColumnDef<Player>[] = [
  {
    accessorKey: "id",
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          className="hover:bg-muted/50"
        >
          ID
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      );
    },
    cell: ({ row }) => <div className="font-mono">{row.getValue("id")}</div>,
  },
  {
    accessorKey: "strippedName",
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          className="hover:bg-muted/50"
        >
          Name
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      );
    },
    cell: ({ row }) => {
      const player = row.original;
      return (
        <div className="space-y-1">
          <div className="font-medium">{player.strippedName}</div>
          <div className="text-xs text-muted-foreground font-mono">
            {player.uuid}
          </div>
        </div>
      );
    },
  },
  {
    accessorKey: "score",
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          className="hover:bg-muted/50"
        >
          Score
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      );
    },
    cell: ({ row }) => <div className="font-mono">{row.getValue("score")}</div>,
  },
  {
    accessorKey: "ping",
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          className="hover:bg-muted/50"
        >
          Ping
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      );
    },
    cell: ({ row }) => {
      const ping = row.getValue("ping") as number;
      return (
        <Badge
          variant="outline"
          className={
            ping < 50
              ? "bg-green-500/20 text-green-400 border-green-500/50"
              : ping < 100
              ? "bg-yellow-500/20 text-yellow-400 border-yellow-500/50"
              : "bg-red-500/20 text-red-400 border-red-500/50"
          }
        >
          {ping}ms
        </Badge>
      );
    },
  },
  {
    accessorKey: "steamId",
    header: "Steam ID",
    cell: ({ row }) => (
      <div className="font-mono text-sm">{row.getValue("steamId")}</div>
    ),
  },
  {
    accessorKey: "address",
    header: "Address",
    cell: ({ row }) => (
      <div className="font-mono text-sm">{row.getValue("address")}</div>
    ),
  },
  {
    accessorKey: "rate",
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          className="hover:bg-muted/50"
        >
          Rate
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      );
    },
    cell: ({ row }) => {
      const rate = row.getValue("rate") as number;
      return <div className="font-mono">{rate.toLocaleString()}</div>;
    },
  },
];

function Players() {
  const players = usePlayers(1000);

  return (
    <ProtectedRoute requiredPermission="players.view">
      <DashboardLayout>
        <div className="p-8 space-y-6 bg-background min-h-screen">
          <div>
            <h1 className="text-4xl font-bold text-foreground mb-2">Players</h1>
            <p className="text-muted-foreground">
              View and manage online players
            </p>
          </div>

          <Card className="bg-card border-border">
            <CardHeader>
              <CardTitle className="text-foreground">Online Players</CardTitle>
              <CardDescription className="text-muted-foreground">
                {players.data?.length || 0} players currently connected
              </CardDescription>
            </CardHeader>
            <CardContent>
              {players.isLoading ? (
                <div className="flex items-center justify-center h-64">
                  <p className="text-foreground text-lg">Loading players...</p>
                </div>
              ) : players.isError ? (
                <div className="flex items-center justify-center h-64">
                  <p className="text-destructive text-lg">
                    Error loading players.
                  </p>
                </div>
              ) : (
                <DataTable
                  columns={columns}
                  data={players.data || []}
                  searchKey="strippedName"
                  searchPlaceholder="Search by name..."
                />
              )}
            </CardContent>
          </Card>
        </div>
      </DashboardLayout>
    </ProtectedRoute>
  );
}

export default Players;
