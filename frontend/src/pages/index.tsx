import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import api from "@/lib/api";

interface Server {
  id: number;
  isDefault: boolean;
}

export default function RootIndex() {
  const navigate = useNavigate();

  useEffect(() => {
    const redirectToDefaultServer = async () => {
      try {
        const servers = await api.get<Server[]>("/servers/active");
        const defaultServer = servers.find((s) => s.isDefault);

        if (defaultServer) {
          navigate(`/${defaultServer.id}`, { replace: true });
        } else if (servers.length > 0) {
          navigate(`/${servers[0].id}`, { replace: true });
        } else {
          // No servers available
          console.error("No active servers found");
        }
      } catch (error) {
        console.error("Failed to load servers:", error);
      }
    };

    redirectToDefaultServer();
  }, [navigate]);

  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="text-lg text-muted-foreground">Loading...</div>
    </div>
  );
}
