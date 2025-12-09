import React, { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import api from "../lib/api";
import { ServerContext } from "@/hooks/useServerContext";

export interface Server {
  id: number;
  name: string;
  host: string;
  port: number;
  rconPort: number;
  gamesMpPath: string;
  isActive: boolean;
  isDefault: boolean;
  description?: string;
  region?: string;
  maxPlayers: number;
  createdAt: string;
  updatedAt: string;
}

export interface ServerContextType {
  currentServer: Server | null;
  servers: Server[];
  loading: boolean;
  error: string | null;
  switchServer: (serverId: number) => void;
  refreshServers: () => Promise<Server[]>;
}

interface ServerProviderProps {
  children: React.ReactNode;
  disableRedirect?: boolean;
}

export function ServerProvider({
  children,
  disableRedirect = false,
}: ServerProviderProps) {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [currentServer, setCurrentServer] = useState<Server | null>(null);
  const [servers, setServers] = useState<Server[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const refreshServers = async () => {
    try {
      const activeServers = await api.get<Server[]>("/servers/active");
      setServers(activeServers);
      return activeServers;
    } catch (err) {
      console.error("Failed to fetch servers:", err);
      setError("Failed to load servers");
      return [];
    }
  };

  useEffect(() => {
    const loadServers = async () => {
      setLoading(true);
      setError(null);

      try {
        const activeServers = await refreshServers();

        if (id) {
          // Load specific server
          const serverId = parseInt(id, 10);
          const server = activeServers.find((s: Server) => s.id === serverId);

          if (server) {
            setCurrentServer(server);
          } else {
            // Server not found, redirect to default
            const defaultServer = activeServers.find(
              (s: Server) => s.isDefault
            );
            if (defaultServer) {
              navigate(`/${defaultServer.id}`, { replace: true });
            } else if (activeServers.length > 0) {
              navigate(`/${activeServers[0].id}`, { replace: true });
            }
          }
        } else {
          // No server ID in URL, redirect to default server (unless disabled)
          if (!disableRedirect) {
            const defaultServer = activeServers.find(
              (s: Server) => s.isDefault
            );
            if (defaultServer) {
              navigate(`/${defaultServer.id}`, { replace: true });
            } else if (activeServers.length > 0) {
              navigate(`/${activeServers[0].id}`, { replace: true });
            }
          }
        }
      } catch (err) {
        console.error("Error loading servers:", err);
        setError("Failed to load servers");
      } finally {
        setLoading(false);
      }
    };

    loadServers();
  }, [id, navigate, disableRedirect]);

  const switchServer = (serverId: number) => {
    const server = servers.find((s) => s.id === serverId);
    if (server) {
      setCurrentServer(server);
      // Navigate to the same page but with the new server ID
      const currentPath = window.location.pathname;

      // For global pages (/plugins, /servers), navigate to the selected server's dashboard
      const globalPages = ["/plugins", "/servers"];
      const isGlobalPage = globalPages.some(
        (page) => currentPath === page || currentPath.startsWith(page + "/")
      );

      if (isGlobalPage) {
        navigate(`/${serverId}`, { replace: true });
        return;
      }

      // For server-specific pages, navigate to the same page with new server ID
      const pathParts = currentPath.split("/").slice(2).join("/");
      navigate(`/${serverId}${pathParts ? "/" + pathParts : ""}`);
    }
  };

  return (
    <ServerContext.Provider
      value={{
        currentServer,
        servers,
        loading,
        error,
        switchServer,
        refreshServers,
      }}
    >
      {children}
    </ServerContext.Provider>
  );
}
