import { Check, ChevronsUpDown, Server as ServerIcon } from "lucide-react";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useServerContext } from "@/hooks/useServerContext";

export function ServerSelector() {
  const { currentServer, servers, switchServer } = useServerContext();

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          variant="outline"
          role="combobox"
          className="w-full justify-between"
        >
          <div className="flex items-center gap-2">
            <ServerIcon className="h-4 w-4" />
            <span className="truncate">
              {currentServer?.name || "Select server..."}
            </span>
          </div>
          <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent
        className="w-[250px]"
        align="start"
        side="right"
        sideOffset={8}
        alignOffset={-8}
      >
        <DropdownMenuLabel>Select Server</DropdownMenuLabel>
        <DropdownMenuSeparator />
        {servers.map((server) => (
          <DropdownMenuItem
            key={server.id}
            onClick={() => switchServer(server.id)}
            className="cursor-pointer"
          >
            <Check
              className={cn(
                "mr-2 h-4 w-4",
                currentServer?.id === server.id ? "opacity-100" : "opacity-0"
              )}
            />
            <div className="flex flex-col flex-1">
              <span>{server.name}</span>
              {server.isDefault && (
                <span className="text-xs text-muted-foreground">Default</span>
              )}
            </div>
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
