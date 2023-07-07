export const dashboardRoutes = {
  DASHBOARD: '/dashboard',
  ALARMS: '/dashboard/alarms',
  ADMIN: '/dashboard/admin',
}

export const dashboardLinks = [
  {
    name: 'Dashboard',
    link: dashboardRoutes.DASHBOARD,
  },
  {
    name: 'Alarms',
    link: dashboardRoutes.ALARMS,
  },
  {
    name: 'Admin',
    link: dashboardRoutes.ADMIN,
    adminRequired: true,
  },
]

export const dashboardRoutesParams = {
  ALARMS: 'alarms',
  ADMIN: 'admin',
}
