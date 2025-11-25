import { useState } from "react";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Clock, Trophy, Heart } from "lucide-react";

const Community = () => {
  const [activeTab, setActiveTab] = useState("feed");

  const feedItems = [
    {
      id: 1,
      user: "Sarah M.",
      initials: "SM",
      action: "completed an 18-hour fast",
      time: "2 hours ago",
      likes: 24,
    },
    {
      id: 2,
      user: "Mike R.",
      initials: "MR",
      action: "shared a healthy breakfast",
      time: "4 hours ago",
      likes: 18,
      image: true,
    },
    {
      id: 3,
      user: "Alex K.",
      initials: "AK",
      action: "reached a 7-day streak",
      time: "6 hours ago",
      likes: 42,
    },
    {
      id: 4,
      user: "Jordan P.",
      initials: "JP",
      action: "completed a 20-hour fast",
      time: "8 hours ago",
      likes: 31,
    },
  ];

  const leaderboard = [
    { rank: 1, name: "Emma W.", hours: 126, badge: "🏆" },
    { rank: 2, name: "David L.", hours: 118, badge: "🥈" },
    { rank: 3, name: "Sophie T.", hours: 112, badge: "🥉" },
    { rank: 4, name: "Chris M.", hours: 98, badge: "" },
    { rank: 5, name: "Taylor B.", hours: 94, badge: "" },
    { rank: 6, name: "Morgan F.", hours: 87, badge: "" },
    { rank: 7, name: "Alex K.", hours: 84, badge: "" },
    { rank: 8, name: "You", hours: 76, badge: "", highlight: true },
  ];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="animate-fade-in">
        <h1 className="text-2xl font-bold bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
          Community
        </h1>
        <p className="text-sm text-muted-foreground">Connect with fellow fasters</p>
      </div>

      {/* Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-2">
          <TabsTrigger value="feed">Feed</TabsTrigger>
          <TabsTrigger value="leaderboard">Leaderboard</TabsTrigger>
        </TabsList>

        <TabsContent value="feed" className="space-y-4 mt-6">
          {feedItems.map((item, index) => (
            <Card 
              key={item.id} 
              className="border-primary/20 animate-fade-in-up hover:border-primary/40 transition-all duration-300 hover:shadow-lg hover:shadow-primary/10 hover:scale-[1.02]"
              style={{ animationDelay: `${index * 0.1}s` }}
            >
              <CardHeader className="pb-3">
                <div className="flex items-start gap-3">
                  <Avatar>
                    <AvatarFallback className="bg-primary/20 text-primary">
                      {item.initials}
                    </AvatarFallback>
                  </Avatar>
                  <div className="flex-1">
                    <div className="flex items-center gap-2">
                      <span className="font-semibold text-sm">{item.user}</span>
                      <Badge variant="secondary" className="text-xs">
                        <Clock className="h-3 w-3 mr-1" />
                        {item.time}
                      </Badge>
                    </div>
                    <p className="text-sm text-muted-foreground mt-1">{item.action}</p>
                  </div>
                </div>
              </CardHeader>
              {item.image && (
                <div className="px-4 pb-3">
                  <div className="h-48 bg-muted rounded-lg flex items-center justify-center">
                    <span className="text-sm text-muted-foreground">Meal Photo</span>
                  </div>
                </div>
              )}
              <CardContent className="pt-0">
                <button className="flex items-center gap-2 text-sm text-muted-foreground hover:text-primary transition-colors">
                  <Heart className="h-4 w-4" />
                  {item.likes} likes
                </button>
              </CardContent>
            </Card>
          ))}
        </TabsContent>

        <TabsContent value="leaderboard" className="space-y-3 mt-6">
          <div className="flex justify-between items-center mb-4">
            <h3 className="text-sm font-semibold">This Week</h3>
            <Badge variant="outline" className="text-xs">
              <Trophy className="h-3 w-3 mr-1" />
              Global
            </Badge>
          </div>

          {leaderboard.map((user, index) => (
            <Card
              key={user.rank}
              className={`border-primary/20 animate-fade-in-up hover:border-primary/40 transition-all duration-300 hover:shadow-lg hover:scale-[1.02] ${
                user.highlight ? "bg-primary/5 border-primary/40 glow-primary" : ""
              }`}
              style={{ animationDelay: `${index * 0.05}s` }}
            >
              <CardContent className="p-4">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-4">
                    <div
                      className={`text-lg font-bold w-8 text-center ${
                        user.rank <= 3 ? "text-primary" : "text-muted-foreground"
                      }`}
                    >
                      {user.badge || `#${user.rank}`}
                    </div>
                    <div>
                      <div className="font-semibold text-sm">{user.name}</div>
                      <div className="text-xs text-muted-foreground">
                        {user.hours} hours fasted
                      </div>
                    </div>
                  </div>
                  {user.highlight && (
                    <Badge variant="secondary" className="text-xs">
                      You
                    </Badge>
                  )}
                </div>
              </CardContent>
            </Card>
          ))}
        </TabsContent>
      </Tabs>
    </div>
  );
};

export default Community;
