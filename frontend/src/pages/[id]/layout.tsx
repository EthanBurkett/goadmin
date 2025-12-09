import { ServerProvider } from "@/providers/ServerProvider";
import { DashboardLayout } from "@/components/DashboardLayout";
import { Outlet } from "react-router-dom";

export default function ServerLayout() {
  return (
    <ServerProvider>
      <DashboardLayout>
        <div className="p-8">
          <Outlet />
        </div>
      </DashboardLayout>
    </ServerProvider>
  );
}
