import type { ServerContextType } from "@/providers/ServerProvider";
import { createContext, useContext } from "react";

export const ServerContext = createContext<ServerContextType | undefined>(
  undefined
);

export function useServerContext() {
  const context = useContext(ServerContext);
  if (context === undefined) {
    throw new Error("useServerContext must be used within a ServerProvider");
  }
  return context;
}
