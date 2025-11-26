import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Users, ArrowRight, Shield } from "lucide-react";

interface Tribe {
    id: string;
    name: string;
    members: number;
    description: string;
    tags: string[];
}

export const TribeDiscovery = () => {
    const recommendedTribes: Tribe[] = [
        {
            id: '1',
            name: 'Keto Warriors',
            members: 1240,
            description: 'High fat, low carb, high discipline. We crush cravings together.',
            tags: ['Keto', 'Weight Loss']
        },
        {
            id: '2',
            name: 'OMAD Club',
            members: 850,
            description: 'One Meal A Day. The ultimate efficiency hack for busy pros.',
            tags: ['OMAD', 'Productivity']
        },
        {
            id: '3',
            name: '72h Reset',
            members: 320,
            description: 'Monthly 3-day fasts for deep autophagy and immune reset.',
            tags: ['Advanced', 'Health']
        }
    ];

    return (
        <div className="space-y-4">
            <div className="bg-gradient-to-r from-purple-900/40 to-blue-900/40 p-6 rounded-xl border border-purple-500/20">
                <h3 className="text-xl font-bold text-white mb-2 flex items-center">
                    <Shield className="w-5 h-5 mr-2 text-purple-400" />
                    Find Your Tribe
                </h3>
                <p className="text-slate-300 text-sm mb-4">
                    Fasting is 3x easier with accountability. Join a squad that matches your goals.
                </p>
            </div>

            <div className="grid gap-4">
                {recommendedTribes.map((tribe) => (
                    <Card key={tribe.id} className="bg-slate-900 border-slate-800 hover:border-slate-700 transition-all">
                        <CardContent className="p-4">
                            <div className="flex justify-between items-start mb-2">
                                <div>
                                    <h4 className="font-bold text-white text-lg">{tribe.name}</h4>
                                    <div className="flex items-center text-xs text-slate-400 mt-1">
                                        <Users className="w-3 h-3 mr-1" />
                                        {tribe.members.toLocaleString()} members
                                    </div>
                                </div>
                                <Button size="sm" className="bg-purple-600 hover:bg-purple-700 text-white">
                                    Join <ArrowRight className="w-3 h-3 ml-1" />
                                </Button>
                            </div>

                            <p className="text-sm text-slate-400 mb-3">
                                {tribe.description}
                            </p>

                            <div className="flex gap-2">
                                {tribe.tags.map(tag => (
                                    <span key={tag} className="text-xs bg-slate-800 text-slate-300 px-2 py-1 rounded border border-slate-700">
                                        {tag}
                                    </span>
                                ))}
                            </div>
                        </CardContent>
                    </Card>
                ))}
            </div>
        </div>
    );
};
