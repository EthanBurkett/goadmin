import { useState, useRef, useEffect } from "react";
import { ProtectedRoute } from "@/components/ProtectedRoute";
import { DashboardLayout } from "@/components/DashboardLayout";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  useSendCommand,
  useKickPlayer,
  useBanPlayer,
  useSayMessage,
  useCommandHistory,
} from "@/hooks/useRcon";
import { usePlayers } from "@/hooks/usePlayers";
import { Terminal, Send, UserX, Ban, MessageSquare } from "lucide-react";
import { getMatchingCommands } from "@/lib/cod4-commands";
import { ErrorBox } from "@/components/ErrorBox";

function Console() {
  const [command, setCommand] = useState("");
  const [showSuggestions, setShowSuggestions] = useState(false);
  const [selectedIndex, setSelectedIndex] = useState(0);
  const [kickPlayerId, setKickPlayerId] = useState("");
  const [kickReason, setKickReason] = useState("");
  const [banPlayerId, setBanPlayerId] = useState("");
  const [banReason, setBanReason] = useState("");
  const [sayMessage, setSayMessage] = useState("");
  const inputRef = useRef<HTMLInputElement>(null);
  const suggestionsRef = useRef<HTMLDivElement>(null);
  const [error, setError] = useState<string | null>(null);

  const { data: players } = usePlayers();
  const { data: history } = useCommandHistory();
  const sendCommandMutation = useSendCommand();
  const kickMutation = useKickPlayer();
  const banMutation = useBanPlayer();
  const sayMutation = useSayMessage();

  const suggestions = getMatchingCommands(command);

  const handleCommandChange = (value: string) => {
    setCommand(value);
    setShowSuggestions(value.length > 0);
    setSelectedIndex(0);
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (!showSuggestions || suggestions.length === 0) return;

    switch (e.key) {
      case "ArrowDown":
        e.preventDefault();
        setSelectedIndex((prev) =>
          prev < suggestions.length - 1 ? prev + 1 : prev
        );
        break;
      case "ArrowUp":
        e.preventDefault();
        setSelectedIndex((prev) => (prev > 0 ? prev - 1 : 0));
        break;
      case "Tab":
      case "Enter":
        if (e.key === "Tab" || (e.key === "Enter" && showSuggestions)) {
          e.preventDefault();
          selectSuggestion(suggestions[selectedIndex]);
        }
        break;
      case "Escape":
        setShowSuggestions(false);
        break;
    }
  };

  const selectSuggestion = (suggestion: (typeof suggestions)[0]) => {
    setCommand(suggestion.command);
    setShowSuggestions(false);
    inputRef.current?.focus();
  };

  const handleSubmitCommand = (e: React.FormEvent) => {
    e.preventDefault();
    setShowSuggestions(false);
    if (command.trim()) {
      const cmd = command.trim();
      sendCommandMutation.mutate(cmd, {
        onSuccess: () => {
          setCommand("");
        },
        onError: (e: Error) => {
          setError(e instanceof Error ? e.message : "Failed to send command");
        },
      });
    }
  };

  // Close suggestions when clicking outside
  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (
        inputRef.current &&
        !inputRef.current.contains(e.target as Node) &&
        suggestionsRef.current &&
        !suggestionsRef.current.contains(e.target as Node)
      ) {
        setShowSuggestions(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  return (
    <ProtectedRoute requiredPermission="rcon.command">
      <DashboardLayout>
        <div className="p-8 space-y-6 bg-background min-h-screen">
          <div>
            <h1 className="text-4xl font-bold text-foreground mb-2">
              Server Console
            </h1>
            <p className="text-muted-foreground">
              Execute RCON commands and manage server operations
            </p>
          </div>

          {/* Quick Actions */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <Dialog>
              <DialogTrigger asChild>
                <Button
                  variant="outline"
                  className="h-auto p-4 flex flex-col items-start gap-2"
                >
                  <UserX className="h-5 w-5 text-yellow-500" />
                  <div className="text-left">
                    <div className="font-semibold">Kick Player</div>
                    <div className="text-xs text-muted-foreground">
                      Remove player from server
                    </div>
                  </div>
                </Button>
              </DialogTrigger>
              <DialogContent className="bg-card border-border">
                <DialogHeader>
                  <DialogTitle className="text-foreground">
                    Kick Player
                  </DialogTitle>
                  <DialogDescription className="text-muted-foreground">
                    Remove a player from the server temporarily
                  </DialogDescription>
                </DialogHeader>
                <div className="space-y-4 py-4">
                  <div className="space-y-2">
                    <Label htmlFor="kick-id">Player</Label>
                    <Select
                      value={kickPlayerId}
                      onValueChange={setKickPlayerId}
                    >
                      <SelectTrigger id="kick-id">
                        <SelectValue placeholder="Select a player" />
                      </SelectTrigger>
                      <SelectContent>
                        {players && players.length > 0 ? (
                          players.map((player) => (
                            <SelectItem
                              key={player.id}
                              value={player.id.toString()}
                            >
                              {player.strippedName} (ID: {player.id})
                            </SelectItem>
                          ))
                        ) : (
                          <SelectItem value="no-players" disabled>
                            No players online
                          </SelectItem>
                        )}
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="kick-reason">Reason (optional)</Label>
                    <Input
                      id="kick-reason"
                      placeholder="Enter reason"
                      value={kickReason}
                      onChange={(e) => setKickReason(e.target.value)}
                    />
                  </div>
                </div>
                <DialogFooter>
                  <Button
                    onClick={() =>
                      kickMutation.mutate(
                        { playerId: kickPlayerId, reason: kickReason },
                        {
                          onSuccess: () => {
                            setKickPlayerId("");
                            setKickReason("");
                          },
                        }
                      )
                    }
                    disabled={!kickPlayerId || kickMutation.isPending}
                  >
                    {kickMutation.isPending ? "Kicking..." : "Kick Player"}
                  </Button>
                </DialogFooter>
              </DialogContent>
            </Dialog>

            <Dialog>
              <DialogTrigger asChild>
                <Button
                  variant="outline"
                  className="h-auto p-4 flex flex-col items-start gap-2"
                >
                  <Ban className="h-5 w-5 text-red-500" />
                  <div className="text-left">
                    <div className="font-semibold">Ban Player</div>
                    <div className="text-xs text-muted-foreground">
                      Permanently ban player
                    </div>
                  </div>
                </Button>
              </DialogTrigger>
              <DialogContent className="bg-card border-border">
                <DialogHeader>
                  <DialogTitle className="text-foreground">
                    Ban Player
                  </DialogTitle>
                  <DialogDescription className="text-muted-foreground">
                    Permanently ban a player from the server
                  </DialogDescription>
                </DialogHeader>
                <div className="space-y-4 py-4">
                  <div className="space-y-2">
                    <Label htmlFor="ban-id">Player</Label>
                    <Select value={banPlayerId} onValueChange={setBanPlayerId}>
                      <SelectTrigger id="ban-id">
                        <SelectValue placeholder="Select a player" />
                      </SelectTrigger>
                      <SelectContent>
                        {players && players.length > 0 ? (
                          players.map((player) => (
                            <SelectItem
                              key={player.id}
                              value={player.id.toString()}
                            >
                              {player.strippedName} (ID: {player.id})
                            </SelectItem>
                          ))
                        ) : (
                          <SelectItem value="no-players" disabled>
                            No players online
                          </SelectItem>
                        )}
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="ban-reason">Reason (optional)</Label>
                    <Input
                      id="ban-reason"
                      placeholder="Enter reason"
                      value={banReason}
                      onChange={(e) => setBanReason(e.target.value)}
                    />
                  </div>
                </div>
                <DialogFooter>
                  <Button
                    onClick={() =>
                      banMutation.mutate(
                        { playerId: banPlayerId, reason: banReason },
                        {
                          onSuccess: () => {
                            setBanPlayerId("");
                            setBanReason("");
                          },
                        }
                      )
                    }
                    disabled={!banPlayerId || banMutation.isPending}
                    variant="destructive"
                  >
                    {banMutation.isPending ? "Banning..." : "Ban Player"}
                  </Button>
                </DialogFooter>
              </DialogContent>
            </Dialog>

            <Dialog>
              <DialogTrigger asChild>
                <Button
                  variant="outline"
                  className="h-auto p-4 flex flex-col items-start gap-2"
                >
                  <MessageSquare className="h-5 w-5 text-blue-500" />
                  <div className="text-left">
                    <div className="font-semibold">Send Message</div>
                    <div className="text-xs text-muted-foreground">
                      Broadcast to all players
                    </div>
                  </div>
                </Button>
              </DialogTrigger>
              <DialogContent className="bg-card border-border">
                <DialogHeader>
                  <DialogTitle className="text-foreground">
                    Send Message
                  </DialogTitle>
                  <DialogDescription className="text-muted-foreground">
                    Broadcast a message to all players on the server
                  </DialogDescription>
                </DialogHeader>
                <div className="space-y-4 py-4">
                  <div className="space-y-2">
                    <Label htmlFor="say-message">Message</Label>
                    <Input
                      id="say-message"
                      placeholder="Enter message"
                      value={sayMessage}
                      onChange={(e) => setSayMessage(e.target.value)}
                    />
                  </div>
                </div>
                <DialogFooter>
                  <Button
                    onClick={() =>
                      sayMutation.mutate(sayMessage, {
                        onSuccess: () => {
                          setSayMessage("");
                        },
                      })
                    }
                    disabled={!sayMessage.trim() || sayMutation.isPending}
                  >
                    {sayMutation.isPending ? "Sending..." : "Send Message"}
                  </Button>
                </DialogFooter>
              </DialogContent>
            </Dialog>
          </div>

          {/* Command Input */}
          <Card className="bg-card border-border">
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-foreground">
                <Terminal className="h-5 w-5" />
                Command Console
              </CardTitle>
              <CardDescription className="text-muted-foreground">
                Execute custom RCON commands directly
              </CardDescription>
            </CardHeader>
            <CardContent>
              <form onSubmit={handleSubmitCommand} className="flex gap-2">
                <div className="flex-1 relative">
                  <Input
                    ref={inputRef}
                    placeholder="Enter RCON command... (type to see suggestions)"
                    value={command}
                    onChange={(e) => handleCommandChange(e.target.value)}
                    onKeyDown={handleKeyDown}
                    className="font-mono"
                    disabled={sendCommandMutation.isPending}
                    autoComplete="off"
                  />
                  {showSuggestions && suggestions.length > 0 && (
                    <div
                      ref={suggestionsRef}
                      className="absolute z-50 w-full mt-1 bg-card border border-border rounded-md shadow-lg max-h-[300px] overflow-y-auto"
                    >
                      {suggestions.slice(0, 10).map((suggestion, index) => (
                        <div
                          key={suggestion.command}
                          className={`px-3 py-2 cursor-pointer hover:bg-muted/50 ${
                            index === selectedIndex ? "bg-muted" : ""
                          }`}
                          onClick={() => selectSuggestion(suggestion)}
                          onMouseEnter={() => setSelectedIndex(index)}
                        >
                          <div className="font-mono text-sm text-foreground">
                            {suggestion.command}
                          </div>
                          <div className="text-xs text-muted-foreground mt-0.5">
                            {suggestion.description}
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
                <Button
                  type="submit"
                  disabled={!command.trim() || sendCommandMutation.isPending}
                >
                  <Send className="h-4 w-4 mr-2" />
                  {sendCommandMutation.isPending ? "Sending..." : "Send"}
                </Button>
              </form>
            </CardContent>
          </Card>
          {error && <ErrorBox error={error} onClose={() => setError("")} />}

          {/* Command History */}
          <Card className="bg-card border-border">
            <CardHeader>
              <CardTitle className="text-foreground">Command History</CardTitle>
              <CardDescription className="text-muted-foreground">
                Recent commands and their responses
              </CardDescription>
            </CardHeader>
            <CardContent>
              {!history || history.length === 0 ? (
                <div className="text-center py-12 text-muted-foreground">
                  No commands executed yet
                </div>
              ) : (
                <div className="max-h-[600px] overflow-y-auto space-y-3 pr-2">
                  {history.map((item) => (
                    <div
                      key={item.id}
                      className="border border-border rounded-lg p-4 space-y-2 bg-muted/20"
                    >
                      <div className="flex items-start justify-between gap-4">
                        <div className="flex-1 space-y-1">
                          <div className="flex items-center gap-2">
                            <Badge
                              variant={item.success ? "default" : "destructive"}
                            >
                              {item.success ? "Success" : "Error"}
                            </Badge>
                            <span className="text-xs text-muted-foreground">
                              {new Date(item.createdAt).toLocaleString()}
                            </span>
                          </div>
                          <div className="font-mono text-sm bg-muted p-2 rounded">
                            $ {item.command}
                          </div>
                        </div>
                      </div>
                      <div className="bg-muted/50 p-3 rounded font-mono text-sm whitespace-pre-wrap text-foreground">
                        {item.response}
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </DashboardLayout>
    </ProtectedRoute>
  );
}

export default Console;
