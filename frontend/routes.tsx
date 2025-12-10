import type {RouteObject} from "react-router-dom";
const routes: RouteObject[] = [];
import Index from "./src/pages/index";
routes.push({ path: "/", element: <Index /> });
import Login from "./src/pages/login";
routes.push({ path: "/login", element: <Login /> });
import Metrics from "./src/pages/metrics";
routes.push({ path: "/metrics", element: <Metrics /> });
import Plugins from "./src/pages/plugins";
routes.push({ path: "/plugins", element: <Plugins /> });
import Servers from "./src/pages/servers";
routes.push({ path: "/servers", element: <Servers /> });
import IdLayout from "./src/pages/[id]/layout";
import IdAnalytics from "./src/pages/[id]/analytics";
import IdAudit from "./src/pages/[id]/audit";
import IdCommands from "./src/pages/[id]/commands";
import IdConsole from "./src/pages/[id]/console";
import IdEmergency from "./src/pages/[id]/emergency";
import IdGroups from "./src/pages/[id]/groups";
import IdIndex from "./src/pages/[id]/index";
import IdMigrations from "./src/pages/[id]/migrations";
import IdPlayers from "./src/pages/[id]/players";
import IdRbac from "./src/pages/[id]/rbac";
import IdReports from "./src/pages/[id]/reports";
import IdStatus from "./src/pages/[id]/status";
import IdWebhooks from "./src/pages/[id]/webhooks";
routes.push({
  path: "/:id",
  element: <IdLayout />,
  children: [
    { path: "analytics", element: <IdAnalytics /> },
    { path: "audit", element: <IdAudit /> },
    { path: "commands", element: <IdCommands /> },
    { path: "console", element: <IdConsole /> },
    { path: "emergency", element: <IdEmergency /> },
    { path: "groups", element: <IdGroups /> },
    { index: true, element: <IdIndex /> },
    { path: "migrations", element: <IdMigrations /> },
    { path: "players", element: <IdPlayers /> },
    { path: "rbac", element: <IdRbac /> },
    { path: "reports", element: <IdReports /> },
    { path: "status", element: <IdStatus /> },
    { path: "webhooks", element: <IdWebhooks /> },
  ],
});
export default routes;
