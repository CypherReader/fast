import { Outlet, useLocation } from "react-router-dom";
import { Home, TrendingUp, Users, User } from "lucide-react";
import { NavLink } from "./NavLink";
import { cn } from "@/lib/utils";

const Layout = () => {
  const location = useLocation();

  const navItems = [
    { path: "/", icon: Home, label: "Dashboard" },
    { path: "/progress", icon: TrendingUp, label: "Progress" },
    { path: "/tribe", icon: Users, label: "Tribe" },
    { path: "/profile", icon: User, label: "Profile" },
  ];

  return (
    <div className="min-h-screen bg-background pb-20 dark">
      <main className="container mx-auto max-w-md px-4 py-6">
        <Outlet />
      </main>

      {/* Bottom Navigation */}
      <nav className="fixed bottom-0 left-0 right-0 bg-card/95 backdrop-blur-lg border-t border-border shadow-2xl">
        <div className="container mx-auto max-w-md">
          <div className="flex items-center justify-around h-16">
            {navItems.map((item) => {
              const Icon = item.icon;
              const isActive = location.pathname === item.path;

              return (
                <NavLink
                  key={item.path}
                  to={item.path}
                  className={cn(
                    "flex flex-col items-center justify-center flex-1 h-full transition-all duration-300 relative group",
                    isActive
                      ? "text-primary"
                      : "text-muted-foreground hover:text-foreground hover:scale-110"
                  )}
                >
                  {isActive && (
                    <div className="absolute top-0 left-1/2 -translate-x-1/2 w-12 h-1 bg-gradient-to-r from-primary to-secondary rounded-full" />
                  )}
                  <Icon className={cn(
                    "h-6 w-6 mb-1 transition-all duration-300",
                    isActive && "animate-bounce"
                  )} />
                  <span className="text-xs font-medium">{item.label}</span>
                </NavLink>
              );
            })}
          </div>
        </div>
      </nav>
    </div>
  );
};

export default Layout;
