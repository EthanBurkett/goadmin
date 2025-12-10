import { useEffect, useRef, useState, useCallback } from "react";
import api from "@/lib/api";
import type { AuditLog } from "./useAudit";

interface StreamStats {
  connected_clients: number;
  broadcast_buffer: {
    capacity: number;
    size: number;
  };
  timestamp: string;
}

export function useAuditStream(onLog?: (log: AuditLog) => void) {
  const [isConnected, setIsConnected] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout>(null);
  const reconnectAttemptsRef = useRef(0);
  const maxReconnectAttempts = 5;
  const baseReconnectDelay = 1000; // Start with 1 second

  const connect = useCallback(() => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      console.log("[AuditStream] Already connected, skipping");
      return; // Already connected
    }

    console.log("[AuditStream] Attempting to connect...");

    try {
      // Create WebSocket URL
      // Note: WebSocket doesn't support cookies in the initial handshake,
      // so we'll rely on the server accepting the connection without auth
      // and authenticate via the existing session
      const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
      const host = window.location.host;
      const wsUrl = `${protocol}//${host}/audit/stream`;

      console.log("[AuditStream] Connecting to:", wsUrl);

      const ws = new WebSocket(wsUrl);

      ws.onopen = () => {
        console.log("[AuditStream] WebSocket connected successfully");
        setIsConnected(true);
        setError(null);
        reconnectAttemptsRef.current = 0;

        // Send ping every 30 seconds to keep connection alive
        const pingInterval = setInterval(() => {
          if (ws.readyState === WebSocket.OPEN) {
            ws.send(JSON.stringify({ type: "ping" }));
          } else {
            clearInterval(pingInterval);
          }
        }, 30000);
      };

      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);

          // Handle pong messages
          if (data.type === "pong") {
            return;
          }

          // Handle audit log messages
          if (onLog && data.id) {
            onLog(data as AuditLog);
          }
        } catch (err) {
          console.error("[AuditStream] Failed to parse message:", err);
        }
      };

      ws.onerror = (event) => {
        console.error("[AuditStream] WebSocket error:", event);
        setError("WebSocket connection error");
      };

      ws.onclose = (event) => {
        console.log(
          "[AuditStream] WebSocket closed:",
          event.code,
          event.reason
        );
        setIsConnected(false);
        wsRef.current = null;

        // Attempt to reconnect with exponential backoff
        if (reconnectAttemptsRef.current < maxReconnectAttempts) {
          const delay =
            baseReconnectDelay * Math.pow(2, reconnectAttemptsRef.current);
          console.log(
            `[AuditStream] Reconnecting in ${delay}ms (attempt ${
              reconnectAttemptsRef.current + 1
            }/${maxReconnectAttempts})`
          );

          reconnectTimeoutRef.current = setTimeout(() => {
            reconnectAttemptsRef.current++;
            connect();
          }, delay);
        } else {
          setError("Maximum reconnection attempts reached");
        }
      };

      wsRef.current = ws;
    } catch (err) {
      console.error("[AuditStream] Failed to create WebSocket:", err);
      setError(err instanceof Error ? err.message : "Failed to connect");
    }
  }, [onLog]);

  const disconnect = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
    }
    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }
    setIsConnected(false);
  }, []);

  useEffect(() => {
    const enabled = onLog !== undefined; // Only connect if callback is provided

    if (!enabled) {
      console.log("[AuditStream] Skipping connection - no callback provided");
      return;
    }

    connect();

    return () => {
      disconnect();
    };
  }, [connect, disconnect, onLog]);

  return {
    isConnected,
    error,
    reconnect: connect,
    disconnect,
  };
}

// Hook to get stream statistics
export function useAuditStreamStats() {
  const [stats, setStats] = useState<StreamStats | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchStats = async () => {
      try {
        setIsLoading(true);
        console.log("[AuditStream] Fetching stream stats...");
        const response = await api.get<StreamStats>("/audit/stream/stats");
        console.log("[AuditStream] Stream stats:", response);
        setStats(response);
        setError(null);
      } catch (err) {
        console.error("[AuditStream] Failed to fetch stats:", err);
        setError(err instanceof Error ? err.message : "Failed to fetch stats");
        setStats(null); // Clear stats on error
      } finally {
        setIsLoading(false);
      }
    };

    fetchStats();
    const interval = setInterval(fetchStats, 10000); // Refresh every 10 seconds

    return () => clearInterval(interval);
  }, []);

  return { stats, isLoading, error };
}
