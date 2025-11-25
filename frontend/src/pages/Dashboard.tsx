import { useState, useEffect } from "react";
import { Flame, Droplet, Brain, Users, Zap, Lock, Share2, Activity, Target } from "lucide-react";
import FastingJourney from "@/components/bio/FastingJourney";
import MetricCard from "@/components/bio/MetricCard";
import PricingMechanism from "@/components/bio/PricingMechanism";
import { fastingApi } from "@/api/client";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

const INITIAL_PRICE = 50.00;
const MIN_PRICE = 1.00;

const Dashboard = () => {
  const [fastingHours, setFastingHours] = useState(16.5);
  const [ketoneLevel] = useState(1.8);
  const [disciplineScore] = useState(65);
  const [price, setPrice] = useState(INITIAL_PRICE);

  // Simulate dynamic pricing calculation
  useEffect(() => {
    const calculatedPrice = INITIAL_PRICE - ((INITIAL_PRICE - MIN_PRICE) * (disciplineScore / 100));
    setPrice(Math.max(MIN_PRICE, calculatedPrice));
  }, [disciplineScore]);

  // Fetch real fasting data
  useEffect(() => {
    const fetchData = async () => {
      try {
        const res = await fastingApi.getCurrent();
        if (res.data) {
          const start = new Date(res.data.start_time).getTime();
          const now = new Date().getTime();
          const hours = (now - start) / (1000 * 60 * 60);
          setFastingHours(hours);
        } else {
          setFastingHours(0);
        }
      } catch (e) {
        console.error("Failed to fetch fasting data", e);
      }
    };
    fetchData();
    const interval = setInterval(fetchData, 60000); // Update every minute
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="min-h-screen bg-slate-950 text-slate-200 font-sans selection:bg-cyan-500/30 p-4 space-y-6">

      {/* Header */}
      <header className="bg-slate-900 border-b border-slate-800 p-4 sticky top-0 z-50 rounded-xl mb-6 shadow-lg shadow-black/20">
        <div className="flex justify-between items-center">
          <div className="flex items-center space-x-2">
            <Brain className="text-cyan-400 w-6 h-6" />
            <span className="text-xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-cyan-400 to-blue-500">
              NEURO-FAST
            </span>
          </div>
          <div className="flex items-center space-x-4">
            <div className="flex flex-col items-end">
              <span className="text-[10px] text-slate-400 uppercase tracking-wider">Lazy Tax</span>
              <span className={`text-lg font-mono font-bold ${price <= 5 ? 'text-green-400' : 'text-red-400'}`}>
                ${price.toFixed(2)}
              </span>
            </div>
          </div>
        </div>
      </header>

      <Tabs defaultValue="focus" className="w-full space-y-6">
        <TabsList className="grid w-full grid-cols-3 bg-slate-900 p-1 rounded-xl border border-slate-800">
          <TabsTrigger value="focus" className="data-[state=active]:bg-cyan-950 data-[state=active]:text-cyan-400 rounded-lg transition-all">
            <Target className="w-4 h-4 mr-2" /> Focus
          </TabsTrigger>
          <TabsTrigger value="tribe" className="data-[state=active]:bg-purple-950 data-[state=active]:text-purple-400 rounded-lg transition-all">
            <Users className="w-4 h-4 mr-2" /> Tribe
          </TabsTrigger>
          <TabsTrigger value="you" className="data-[state=active]:bg-slate-800 data-[state=active]:text-white rounded-lg transition-all">
            <Activity className="w-4 h-4 mr-2" /> You
          </TabsTrigger>
        </TabsList>

        {/* TAB 1: FOCUS - The Main Fasting Interface */}
        <TabsContent value="focus" className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">

          {/* Hero Timer Card */}
          <div className="bg-slate-900/50 border border-slate-800 rounded-2xl p-8 flex flex-col items-center justify-center relative overflow-hidden">
            <div className="absolute inset-0 bg-gradient-to-b from-cyan-500/5 to-transparent pointer-events-none"></div>
            <div className="w-48 h-48 rounded-full border-4 border-slate-800 flex items-center justify-center relative mb-4">
              <div className={`absolute inset-0 rounded-full border-4 ${fastingHours > 0 ? 'border-cyan-500 animate-spin-slow' : 'border-slate-700'} border-t-transparent opacity-70`}></div>
              <div className="text-center z-10">
                <span className="text-4xl font-bold text-white font-mono">{fastingHours.toFixed(1)}</span>
                <span className="block text-xs text-slate-400 uppercase tracking-widest mt-1">Hours</span>
              </div>
            </div>

            <div className="flex flex-col items-center space-y-4 w-full">
              <div className="flex items-center space-x-2 bg-slate-800/80 px-4 py-1.5 rounded-full border border-slate-700">
                <Flame className={`w-4 h-4 ${fastingHours > 0 ? 'text-orange-500' : 'text-slate-500'}`} />
                <span className="text-sm font-medium text-slate-200">
                  {fastingHours > 16 ? "Autophagy Active" : fastingHours > 0 ? "Fat Burning" : "Ready to Fast"}
                </span>
              </div>

              {fastingHours > 0 ? (
                <button
                  onClick={async () => {
                    try {
                      await fastingApi.stop();
                      setFastingHours(0);
                      // Refresh other data if needed
                      window.location.reload(); // Simple reload to fetch fresh user stats
                    } catch (e) {
                      console.error("Failed to stop fast", e);
                    }
                  }}
                  className="bg-red-500/10 hover:bg-red-500/20 text-red-400 border border-red-500/30 px-8 py-2 rounded-full font-bold transition-all"
                >
                  Stop Fast
                </button>
              ) : (
                <button
                  onClick={async () => {
                    try {
                      await fastingApi.start("circadian", 16);
                      setFastingHours(0.1); // Optimistic update
                      window.location.reload();
                    } catch (e) {
                      console.error("Failed to start fast", e);
                    }
                  }}
                  className="bg-cyan-500 hover:bg-cyan-400 text-black px-8 py-2 rounded-full font-bold shadow-[0_0_15px_rgba(6,182,212,0.5)] transition-all"
                >
                  Start Fast
                </button>
              )}
            </div>
          </div>

          {/* Bio-Narrative Timeline */}
          <FastingJourney fastingHours={fastingHours} />

        </TabsContent>

        {/* TAB 2: TRIBE - Social Accountability */}
        <TabsContent value="tribe" className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
          <div className="bg-slate-900/50 rounded-xl p-6 border border-slate-800">
            <div className="flex justify-between items-center mb-6">
              <h3 className="font-bold text-white flex items-center text-lg">
                <Users className="w-5 h-5 mr-2 text-purple-400" /> Tribe Leaderboard
              </h3>
              <span className="text-xs bg-purple-500/20 text-purple-300 px-3 py-1 rounded-full border border-purple-500/30">Rank #42</span>
            </div>

            <div className="space-y-4">
              {[
                { name: 'Sarah K.', status: 'Fasting 18h', streak: 12, alert: false },
                { name: 'Mike R.', status: 'Eating Window', streak: 45, alert: false },
                { name: 'You', status: `Fasting ${fastingHours.toFixed(0)}h`, streak: 8, alert: false },
                { name: 'Dave L.', status: 'Danger Zone', streak: 2, alert: true },
              ].map((member, i) => (
                <div key={i} className="flex justify-between items-center p-3 hover:bg-slate-800 rounded-xl cursor-pointer transition-all border border-transparent hover:border-slate-700">
                  <div className="flex items-center">
                    <div className={`w-3 h-3 rounded-full mr-3 ${member.alert ? 'bg-red-500 animate-pulse shadow-[0_0_8px_rgba(239,68,68,0.6)]' : 'bg-green-500 shadow-[0_0_8px_rgba(34,197,94,0.4)]'}`}></div>
                    <span className="text-slate-200 font-medium">{member.name}</span>
                  </div>
                  <div className="flex items-center space-x-4">
                    <span className="text-slate-500 text-xs">{member.status}</span>
                    <span className="font-mono text-cyan-400 text-xs bg-cyan-950/30 px-2 py-1 rounded">Day {member.streak}</span>
                  </div>
                </div>
              ))}
            </div>

            <button className="w-full mt-6 bg-purple-600/10 hover:bg-purple-600/20 text-purple-300 border border-purple-500/30 py-3 rounded-xl text-sm flex justify-center items-center transition-all font-medium">
              <Zap className="w-4 h-4 mr-2" /> Nudge Dave (Cost: 0.5 pts)
            </button>
          </div>

          <div className="bg-slate-800/30 p-4 rounded-xl border border-slate-700/50 flex items-center justify-between">
            <div>
              <h4 className="font-bold text-white text-sm">Proof of Discipline</h4>
              <p className="text-xs text-slate-400">Share verification card to Instagram</p>
            </div>
            <button className="bg-pink-600 hover:bg-pink-500 text-white px-4 py-2 rounded-lg text-sm flex items-center transition-colors shadow-lg shadow-pink-900/20">
              <Share2 className="w-4 h-4 mr-2" /> Share (+5 pts)
            </button>
          </div>
        </TabsContent>

        {/* TAB 3: YOU - Metrics & Pricing */}
        <TabsContent value="you" className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
          <div className="grid grid-cols-2 gap-4">
            <MetricCard
              icon={Droplet}
              title="Ketones"
              value={ketoneLevel}
              unit="mmol/L"
              status="Optimal"
              trend="+0.2"
            />
            <div className="bg-slate-900 border border-slate-800 rounded-xl p-4 flex flex-col justify-between relative overflow-hidden">
              <div className="absolute top-0 right-0 p-2 opacity-10">
                <Target className="w-12 h-12 text-white" />
              </div>
              <div className="flex items-center space-x-2 mb-2">
                <Target className="w-4 h-4 text-emerald-400" />
                <span className="text-xs text-slate-400 uppercase tracking-wider font-bold">Discipline</span>
              </div>
              <div className="flex items-baseline space-x-1">
                <span className="text-2xl font-bold text-white">{disciplineScore}</span>
                <span className="text-xs text-slate-500">/ 100</span>
              </div>
              <div className="w-full bg-slate-800 rounded-full h-1.5 mt-3">
                <div className="bg-emerald-500 h-1.5 rounded-full" style={{ width: `${disciplineScore}%` }}></div>
              </div>
            </div>
          </div>

          <PricingMechanism currentPrice={price} disciplineScore={disciplineScore} />

          <div className="bg-yellow-900/10 border border-yellow-700/20 p-6 rounded-xl relative overflow-hidden group">
            <div className="absolute inset-0 bg-gradient-to-r from-yellow-500/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
            <h4 className="text-yellow-500 font-bold flex items-center mb-2 text-lg">
              <Lock className="w-5 h-5 mr-2" />
              Premium Protocols
            </h4>
            <p className="text-sm text-yellow-200/70 mb-4 leading-relaxed">
              Unlock "The Monk Fast" (36h) and "Bio-hacker Recipes" by reaching 80% Discipline Score.
            </p>
            <div className="w-full bg-slate-800 rounded-full h-2">
              <div className="bg-yellow-600 h-2 rounded-full w-[65%] shadow-[0_0_10px_rgba(202,138,4,0.4)]"></div>
            </div>
            <p className="text-xs text-right text-yellow-500/50 mt-2 font-mono">65% / 80%</p>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  );
};

export default Dashboard;
