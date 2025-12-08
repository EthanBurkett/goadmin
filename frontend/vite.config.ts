import { defineConfig, type HmrContext, type ViteDevServer } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import * as routes from "./startup/routes";
import path from "path";

type ServicesType = Record<
  string,
  {
    start: () => void;
    update?: (server: ViteDevServer) => void;
  }
>;

const ServiceStartup = (services: ServicesType) => {
  return {
    name: "service-startup",
    configureServer: function () {
      Object.entries(services).forEach(([_, service]) => {
        service.start();
      });
    },
    handleHotUpdate: function ({ server }: HmrContext) {
      Object.entries(services).forEach(([_, service]) => {
        if (service.update) service.update(server);
      });
    },
  };
};
// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    react(),
    tailwindcss(),
    ServiceStartup({
      routes,
    }),
  ],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
      "@app": path.resolve(__dirname, ".."),
    },
  },
});
