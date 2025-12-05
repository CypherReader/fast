export const Footer = () => {
  return (
    <footer className="py-12 bg-background border-t border-border">
      <div className="container px-4">
        <div className="flex flex-col md:flex-row items-center justify-between gap-6 mb-8">
          {/* Logo */}
          <div className="flex items-center gap-2">
            <img src="/fasthero.png" alt="FastingHero" className="w-8 h-8 rounded-lg" />
            <span className="font-bold text-lg">FastingHero</span>
          </div>

          {/* Links */}
          <div className="flex items-center gap-6 text-sm text-muted-foreground">
            <a href="#" className="hover:text-foreground transition-colors">
              Privacy
            </a>
            <a href="#" className="hover:text-foreground transition-colors">
              Terms
            </a>
            <a href="#" className="hover:text-foreground transition-colors">
              Support
            </a>
          </div>

          {/* Copyright */}
          <p className="text-sm text-muted-foreground">
            © 2024 FastingHero. All rights reserved.
          </p>
        </div>

        {/* Legal disclaimer */}
        <div className="border-t border-border pt-6">
          <p className="text-xs text-muted-foreground/70 max-w-3xl mx-auto text-center leading-relaxed">
            *How the Vault Works: You pay $20/month. Complete 10 fasts = full $20 refunded.
            Complete 5 fasts = $10 refunded. Complete 0 fasts = $0 refunded.
            Refunds process automatically on the 1st of each month.
            Net cost ranges from $0-20/month based on your fasting consistency.
            Average user net cost: $1.80/month.
          </p>
        </div>
      </div>
    </footer>
  );
};
