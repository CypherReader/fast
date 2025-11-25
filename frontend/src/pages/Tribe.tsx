import { Card, CardContent } from "@/components/ui/card";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";

const Tribe = () => {
    // Mock Data
    const tribe = {
        name: "The Spartans",
        collectiveScore: 85,
        members: [
            { id: 1, name: "Leonidas", status: "Fasting", hours: 18, discipline: 98, avatar: "L" },
            { id: 2, name: "Gorgo", status: "Fasting", hours: 16, discipline: 95, avatar: "G" },
            { id: 3, name: "Dilios", status: "Eating", hours: 0, discipline: 82, avatar: "D" },
            { id: 4, name: "Stelios", status: "Fasting", hours: 42, discipline: 99, avatar: "S" },
        ]
    };

    return (
        <div className="min-h-screen bg-background p-6 space-y-8">
            {/* Header */}
            <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
                <div>
                    <h1 className="text-4xl font-bold text-primary tracking-tight">{tribe.name}</h1>
                    <p className="text-muted-foreground">Tribe ID: #SPARTANS-300</p>
                </div>
                <div className="bg-card border border-primary/20 p-4 rounded-xl shadow-lg">
                    <div className="text-sm text-muted-foreground">Collective Discipline</div>
                    <div className="text-3xl font-bold text-primary">{tribe.collectiveScore}%</div>
                </div>
            </div>

            {/* Leaderboard */}
            <div className="grid gap-4">
                {tribe.members.map((member) => (
                    <Card key={member.id} className="border-border/50 bg-card/50 backdrop-blur">
                        <CardContent className="flex items-center p-4 gap-4">
                            <Avatar className="h-12 w-12 border-2 border-primary">
                                <AvatarFallback>{member.avatar}</AvatarFallback>
                            </Avatar>

                            <div className="flex-1 min-w-0">
                                <div className="flex items-center justify-between mb-1">
                                    <h3 className="font-bold truncate">{member.name}</h3>
                                    <Badge variant={member.status === "Fasting" ? "default" : "secondary"}>
                                        {member.status} {member.status === "Fasting" && `(${member.hours}h)`}
                                    </Badge>
                                </div>

                                <div className="space-y-1">
                                    <div className="flex justify-between text-xs text-muted-foreground">
                                        <span>Discipline Index</span>
                                        <span>{member.discipline}</span>
                                    </div>
                                    <Progress value={member.discipline} className="h-2" />
                                </div>
                            </div>
                        </CardContent>
                    </Card>
                ))}
            </div>

            {/* Invite Block */}
            <Card className="border-dashed border-2 border-muted-foreground/20 bg-transparent">
                <CardContent className="flex flex-col items-center justify-center p-8 text-center space-y-4">
                    <div className="p-3 bg-muted rounded-full">
                        <span className="text-2xl">👋</span>
                    </div>
                    <div>
                        <h3 className="font-bold text-lg">Grow your Tribe</h3>
                        <p className="text-muted-foreground text-sm max-w-xs mx-auto">
                            Invite friends to join {tribe.name}. The more you fast together, the lower your price drops.
                        </p>
                    </div>
                </CardContent>
            </Card>
        </div>
    );
};

export default Tribe;
