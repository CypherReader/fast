import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { PlayCircle } from "lucide-react";

const EXPERT_CONTENT = [
    {
        id: 1,
        name: "Dr. Jason Fung",
        role: "The Architect",
        topic: "The Insulin Switch",
        desc: "Why 'Calories In, Calories Out' fails. Lowering the insulin set-point is the only way to heal.",
        link: "https://www.youtube.com/watch?v=etNLeBTlGQk",
        color: "text-blue-400",
        bg: "bg-blue-500/10 border-blue-500/20"
    },
    {
        id: 2,
        name: "Dr. Mindy Pelz",
        role: "Hormone Expert",
        topic: "Fasting Like a Girl",
        desc: "Timing your fasts around your cycle to maximize progesterone. Essential for female physiology.",
        link: "https://www.youtube.com/watch?v=hT5DjSUt9_Q",
        color: "text-pink-400",
        bg: "bg-pink-500/10 border-pink-500/20"
    },
    {
        id: 3,
        name: "Dr. Pradip Jamnadas",
        role: "Cardiologist",
        topic: "Fasting for Survival",
        desc: "A surgeon's view on how fasting cleanses the endothelium and prevents cardiovascular decay.",
        link: "https://www.youtube.com/watch?v=qw8wMELZ5dI",
        color: "text-red-400",
        bg: "bg-red-500/10 border-red-500/20"
    },
    {
        id: 4,
        name: "Dr. Eric Berg",
        role: "Keto Specialist",
        topic: "Healthy Ketosis",
        desc: "Tactical advice on electrolytes, the 'Dawn Phenomenon', and what to eat to break your fast.",
        link: "https://www.youtube.com/watch?v=VlLKxrwnGTI",
        color: "text-green-400",
        bg: "bg-green-500/10 border-green-500/20"
    },
    {
        id: 5,
        name: "Dr. Annette Bosworth",
        role: "Internal Medicine",
        topic: "The Dr. Boz Ratio",
        desc: "Mathematics of metabolism. Using the Glucose/Ketone ratio to predict autophagy depth.",
        link: "https://www.youtube.com/watch?v=9AwrAVcU7A4",
        color: "text-orange-400",
        bg: "bg-orange-500/10 border-orange-500/20"
    }
];

export const KnowledgeHub = () => {
    return (
        <div className="space-y-4">
            <div className="grid gap-4">
                {EXPERT_CONTENT.map((expert) => (
                    <Card
                        key={expert.id}
                        className={`border transition-all duration-300 hover:scale-[1.02] hover:shadow-lg ${expert.bg} border-opacity-30`}
                    >
                        <CardHeader className="pb-2">
                            <div className="flex justify-between items-start">
                                <div>
                                    <Badge variant="outline" className={`mb-2 ${expert.color} border-current opacity-80`}>
                                        {expert.role}
                                    </Badge>
                                    <CardTitle className="text-lg font-bold text-foreground">
                                        {expert.name}
                                    </CardTitle>
                                </div>
                                <a
                                    href={expert.link}
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    className={`p-2 rounded-full hover:bg-white/10 transition-colors ${expert.color}`}
                                >
                                    <PlayCircle className="w-6 h-6" />
                                </a>
                            </div>
                        </CardHeader>
                        <CardContent>
                            <h4 className={`font-semibold mb-1 ${expert.color}`}>{expert.topic}</h4>
                            <p className="text-sm text-muted-foreground leading-relaxed">
                                {expert.desc}
                            </p>
                        </CardContent>
                    </Card>
                ))}
            </div>
        </div>
    );
};
