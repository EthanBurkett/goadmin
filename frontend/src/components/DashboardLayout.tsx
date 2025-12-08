import {
  Home,
  Users,
  Server,
  Terminal,
  Shield,
  LogOut,
  BarChart3,
  Command,
  Group,
  AlertTriangle,
  FileText,
} from "lucide-react";
import { Link, useLocation } from "react-router-dom";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
} from "@/components/ui/sidebar";
import { Button } from "@/components/ui/button";
import { useAuthContext } from "@/hooks/useAuthContext";

const navigation = [
  { name: "Dashboard", href: "/", icon: Home, permission: null },
  {
    name: "Players",
    href: "/players",
    icon: Users,
    permission: "players.view",
  },
  {
    name: "Server Status",
    href: "/status",
    icon: Server,
    permission: "status.view",
  },
  {
    name: "Analytics",
    href: "/analytics",
    icon: BarChart3,
    permission: "status.view",
  },
  {
    name: "Console",
    href: "/console",
    icon: Terminal,
    permission: "rcon.command",
  },
  {
    name: "Groups",
    href: "/groups",
    icon: Group,
    permission: "rbac.manage",
  },
  {
    name: "Commands",
    href: "/commands",
    icon: Command,
    permission: "rbac.manage",
  },
  {
    name: "Reports",
    href: "/reports",
    icon: AlertTriangle,
    permission: "reports.view",
  },
  {
    name: "Audit Logs",
    href: "/audit",
    icon: FileText,
    permission: "rbac.manage",
  },
  { name: "RBAC", href: "/rbac", icon: Shield, permission: "rbac.manage" },
];

interface DashboardLayoutProps {
  children: React.ReactNode;
}

export function DashboardLayout({ children }: DashboardLayoutProps) {
  const location = useLocation();
  const { user, logout } = useAuthContext();

  const handleLogout = async () => {
    await logout();
  };

  // Helper to check if user has permission
  const hasPermission = (permission: string | null) => {
    if (!permission) return true;
    if (!user?.roles) return false;
    return user.roles.some((role) =>
      role.permissions?.some((perm) => perm.name === permission)
    );
  };

  // Filter navigation items based on user permissions
  const visibleNavigation = navigation.filter((item) =>
    hasPermission(item.permission)
  );

  return (
    <SidebarProvider>
      <div className="flex min-h-screen w-full bg-background">
        <Sidebar className="border-border">
          <SidebarHeader className="border-b border-border p-4">
            <div className="flex items-center space-x-2">
              <Shield className="h-6 w-6 text-primary" />
              <span className="text-lg font-bold text-foreground">GoAdmin</span>
            </div>
          </SidebarHeader>
          <SidebarContent>
            <SidebarGroup>
              <SidebarGroupLabel className="text-muted-foreground">
                Navigation
              </SidebarGroupLabel>
              <SidebarGroupContent>
                <SidebarMenu>
                  {visibleNavigation.map((item) => (
                    <SidebarMenuItem key={item.name}>
                      <SidebarMenuButton
                        asChild
                        isActive={location.pathname === item.href}
                        className="text-foreground/90 hover:text-foreground hover:bg-muted"
                      >
                        <Link to={item.href}>
                          <item.icon className="h-4 w-4" />
                          <span>{item.name}</span>
                        </Link>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                  ))}
                </SidebarMenu>
              </SidebarGroupContent>
            </SidebarGroup>
          </SidebarContent>
          <SidebarFooter className="border-t border-border p-4">
            <div className="space-y-2">
              <div className="text-sm text-muted-foreground">
                Logged in as{" "}
                <span className="text-foreground font-semibold">
                  {user?.username}
                </span>
              </div>
              <Button
                variant="outline"
                className="w-full justify-start"
                onClick={handleLogout}
              >
                <LogOut className="mr-2 h-4 w-4" />
                Sign Out
              </Button>
            </div>
          </SidebarFooter>
        </Sidebar>
        <main className="flex-1 overflow-auto">{children}</main>
      </div>
    </SidebarProvider>
  );
}
