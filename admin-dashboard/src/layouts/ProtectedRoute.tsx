import React, { useEffect } from 'react';
import { Navigate, Outlet } from 'react-router-dom';
import { STAFF_ROLES, useAuthStore } from '../store/useAuthStore';
import type { Role } from '../types';

export const homeFor = (role?: Role) => (role === 'DRIVER' ? '/shipping' : '/dashboard');

interface ProtectedRouteProps {
  allowedRoles?: Role[];
}

export const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ allowedRoles }) => {
  const { isAuthenticated, user, logout } = useAuthStore();
  const isStaff = !!user && STAFF_ROLES.includes(user.role);

  // A customer session must never stay in the admin app.
  useEffect(() => {
    if (isAuthenticated && !isStaff) logout();
  }, [isAuthenticated, isStaff, logout]);

  if (!isAuthenticated || !user || !isStaff) {
    return <Navigate to="/login" replace />;
  }
  if (allowedRoles && !allowedRoles.includes(user.role)) {
    return <Navigate to={homeFor(user.role)} replace />;
  }

  return <Outlet />;
};
