import type { ReactNode } from "react";
import { Navigate, useNavigate } from "react-router-dom";
import { useAuthContext } from "@/hooks/useAuthContext";
import { useEffect, useCallback } from "react";
import { toast } from "sonner";

interface ProtectedRouteProps {
  children: ReactNode;
  requiredPermission?: string;
}

export function ProtectedRoute({
  children,
  requiredPermission,
}: ProtectedRouteProps) {
  const { isAuthenticated, isLoading, user } = useAuthContext();
  const navigate = useNavigate();

  // Check if user has the required permission
  const hasPermission = useCallback(
    (permission: string) => {
      if (!user?.roles) return false;
      return user.roles.some((role) =>
        role.permissions?.some((perm) => perm.name === permission)
      );
    },
    [user]
  );

  useEffect(() => {
    if (
      !isLoading &&
      isAuthenticated &&
      requiredPermission &&
      !hasPermission(requiredPermission)
    ) {
      toast.error("You don't have permission to access this page");
      navigate(-1);
    }
  }, [isLoading, isAuthenticated, requiredPermission, navigate, hasPermission]);

  if (isLoading) {
    return (
      <div className="min-h-screen bg-neutral-900 flex items-center justify-center">
        <p className="text-white text-lg">Loading...</p>
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  if (requiredPermission && !hasPermission(requiredPermission)) {
    return (
      <div className="min-h-screen bg-neutral-900 flex items-center justify-center">
        <p className="text-white text-lg">Checking permissions...</p>
      </div>
    );
  }

  return <>{children}</>;
}
