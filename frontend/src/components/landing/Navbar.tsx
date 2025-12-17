import { Button } from "@/components/ui/button";
import { DollarSign } from "lucide-react";

export const Navbar = () => {
  return (
    <nav className="fixed top-0 left-0 right-0 z-50 bg-background/80 backdrop-blur-xl border-b border-border">
      <div className="container px-4">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <div className="flex items-center gap-2">
            <img src="/fasthero.png" alt="FastingHero" className="w-8 h-8 rounded-lg" />
            <span className="font-bold text-lg">FastingHero</span>
          </div>

          {/* CTA */}
          <div className="flex items-center gap-4">
            <Button variant="ghost" size="sm" onClick={() => window.location.href = '/login'}>
              Login
            </Button>
            {/* HIDDEN FOR V2: Start Your Vault */}
            <Button variant="hero" size="sm" onClick={() => window.location.href = '/onboarding'}>
              Get Started
            </Button>
          </div>
        </div>
      </div>
    </nav>
  );
};
