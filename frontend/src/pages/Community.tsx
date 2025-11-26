import { useState } from "react";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Clock, Trophy, Heart } from "lucide-react";
import { api } from "@/services/api";
import TribesTab from "@/components/TribesTab";

const Community = () => {
  const [activeTab, setActiveTab] = useState("feed");
  const [feedItems, setFeedItems] = useState<any[]>([]);

  const fetchFeed = async () => {
    try {
      const res = await api.get('/social/feed');
      // Transform backend events to UI format
      const items = res.data.map((event: any) => ({
        id: event.id,
        user: event.user_name || "Unknown User",
        initials: (event.user_name || "U").substring(0, 2).toUpperCase(),
        action: formatAction(event),
        time: new Date(event.created_at).toLocaleString(), // Simple formatting
        likes: 0, // Mock for now
      }));
      setFeedItems(items);
    } catch (e) {
      console.error("Failed to fetch feed", e);
    }
  };

  const [leaderboardItems, setLeaderboardItems] = useState<any[]>([]);
  const fetchLeaderboard = async () => {
    try {
      const res = await api.get('/leaderboard/');
      // Transform
      const items = res.data.map((entry: any) => ({
        rank: entry.rank,
        name: entry.user_name || "Unknown",
        hours: Math.round(entry.total_fasting_hours),
        badge: entry.rank === 1 ? "🏆" : entry.rank === 2 ? "🥈" : entry.rank === 3 ? "🥉" : "",
        highlight: false // TODO: Check if current user
      }));
      setLeaderboardItems(items);
    } catch (e) {
      console.error("Failed to fetch leaderboard", e);
    }
  };

  const [gamificationProfile, setGamificationProfile] = useState<any>(null);
  const fetchGamificationProfile = async () => {
    try {
      const res = await api.get('/gamification/profile');
      setGamificationProfile(res.data);
    } catch (e) {
      console.error("Failed to fetch gamification profile", e);
    }
  };

  const formatAction = (event: any) => {
    switch (event.type) {
      case "fast_completed":
        const data = JSON.parse(JSON.stringify(event.data)); // data is already object from axios
        return `completed a ${data.duration_hours || '?'}h fast`;
      case "keto_logged":
        return "logged keto metrics";
      case "tribe_joined":
        return "joined a tribe";
      default:
        return "did something awesome";
    }
  };

  if (activeTab === "feed" && feedItems.length === 0) {
    fetchFeed();
  }

  if (activeTab === "leaderboard" && leaderboardItems.length === 0) {
    fetchLeaderboard();
  }

  if (activeTab === "progress" && !gamificationProfile) {
    fetchGamificationProfile();
  }






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
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="feed">Feed</TabsTrigger>
          <TabsTrigger value="tribes">Tribes</TabsTrigger>
          <TabsTrigger value="leaderboard">Leaderboard</TabsTrigger>
          <TabsTrigger value="progress">My Progress</TabsTrigger>
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

        <TabsContent value="tribes" className="mt-6">
          <TribesTab />
        </TabsContent>

        <TabsContent value="leaderboard" className="space-y-3 mt-6">
          <div className="flex justify-between items-center mb-4">
            <h3 className="text-sm font-semibold">This Week</h3>
            <Badge variant="outline" className="text-xs">
              <Trophy className="h-3 w-3 mr-1" />
              Global
            </Badge>
          </div>

          {leaderboardItems.map((user, index) => (
            <Card
              key={user.rank}
              className={`border-primary/20 animate-fade-in-up hover:border-primary/40 transition-all duration-300 hover:shadow-lg hover:scale-[1.02] ${user.highlight ? "bg-primary/5 border-primary/40 glow-primary" : ""
                }`}
              style={{ animationDelay: `${index * 0.05}s` }}
            >
              <CardContent className="p-4">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-4">
                    <div
                      className={`text-lg font-bold w-8 text-center ${user.rank <= 3 ? "text-primary" : "text-muted-foreground"
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

        <TabsContent value="progress" className="space-y-6 mt-6">
          {gamificationProfile ? (
            <>
              {/* Streak Section */}
              <Card className="border-primary/20 bg-gradient-to-br from-primary/5 to-transparent">
                <CardHeader>
                  <h3 className="text-lg font-semibold flex items-center gap-2">
                    <Trophy className="h-5 w-5 text-yellow-500" />
                    Current Streak
                  </h3>
                </CardHeader>
                <CardContent>
                  <div className="flex items-end gap-2">
                    <span className="text-4xl font-bold text-primary">
                      {gamificationProfile.streak?.current_streak || 0}
                    </span>
                    <span className="text-muted-foreground mb-1">days</span>
                  </div>
                  <p className="text-sm text-muted-foreground mt-2">
                    Longest streak: {gamificationProfile.streak?.longest_streak || 0} days
                  </p>
                </CardContent>
              </Card>

              {/* Badges Section */}
              <div>
                <h3 className="text-lg font-semibold mb-4">My Badges</h3>
                <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                  {gamificationProfile.badges?.map((badge: any) => (
                    <Card key={badge.badge_id} className="border-primary/20 hover:border-primary/40 transition-all">
                      <CardContent className="p-4 flex flex-col items-center text-center gap-2">
                        <div className="text-4xl mb-2">{badge.badge_info?.icon || "🏅"}</div>
                        <div className="font-semibold text-sm">{badge.badge_info?.name || badge.badge_id}</div>
                        <div className="text-xs text-muted-foreground">{badge.badge_info?.description}</div>
                        <div className="text-[10px] text-muted-foreground mt-1">
                          Earned: {new Date(badge.earned_at).toLocaleDateString()}
                        </div>
                      </CardContent>
                    </Card>
                  ))}
                  {(!gamificationProfile.badges || gamificationProfile.badges.length === 0) && (
                    <div className="col-span-full text-center text-muted-foreground py-8">
                      No badges earned yet. Keep fasting to unlock them!
                    </div>
                  )}
                </div>
              </div>
            </>
          ) : (
            <div className="text-center py-8">Loading profile...</div>
          )}
        </TabsContent>
      </Tabs>
    </div>
  );
};

export default Community;
