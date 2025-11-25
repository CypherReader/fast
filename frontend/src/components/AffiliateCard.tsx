import { Card, CardContent } from "./ui/card";
import { Button } from "./ui/button";
import { ExternalLink, Droplets } from "lucide-react";

const AffiliateCard = () => {
  return (
    <Card className="bg-gradient-to-r from-primary/10 to-secondary/10 border-primary/30">
      <CardContent className="p-4">
        <div className="flex items-start gap-3">
          <div className="p-2 rounded-lg bg-primary/20">
            <Droplets className="h-5 w-5 text-primary" />
          </div>
          <div className="flex-1">
            <h3 className="font-semibold text-sm mb-1">Maximize Your Fast</h3>
            <p className="text-xs text-muted-foreground mb-3">
              Stay hydrated with essential electrolytes
            </p>
            <Button
              size="sm"
              variant="outline"
              className="w-full border-primary/50 hover:bg-primary/10"
            >
              View Recommended Supplements
              <ExternalLink className="ml-2 h-3 w-3" />
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  );
};

export default AffiliateCard;
