import { useState, useEffect } from "react";
import { Brain, Users, Zap, Lock, Share2, Activity, Target, Droplet } from "lucide-react";
import FastingJourney from "@/components/bio/FastingJourney";
import MetricCard from "@/components/bio/MetricCard";
import { fastingApi } from "@/api/client";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import FastingClock from "@/components/FastingClock";
import VaultStatus from "@/components/VaultStatus";
import { FastingSession } from "@/api/types";
import MedicalModal from "@/components/MedicalModal";
import { VaultHeroCard } from "@/features/vault/VaultHeroCard";
import { SocialProofBanner } from "@/components/trust/SocialProofBanner";
import { ActivityTicker } from "@/components/trust/ActivityTicker";
import { ChallengeCard } from "@/features/engagement/ChallengeCard";
import { StreakCounter } from "@/features/engagement/StreakCounter";
import { ReferralTeaser } from "@/features/referral/ReferralTeaser";
import { ReferralModal } from "@/features/referral/ReferralModal";
import { TribeDiscovery } from "@/features/tribe/TribeDiscovery";
import { FloatingActionMenu } from "@/components/ui/FloatingActionMenu";
import { RefundCountdown } from "@/features/vault/RefundCountdown";

const Dashboard = () => {
  const [session, setSession] = useState<FastingSession | null>(null);
  const [elapsed, setElapsed] = useState(0);
  const [ketoneLevel] = useState(1.8);
  const [disciplineScore] = useState(65);
  const [showModal, setShowModal] = useState(false);
  const [pendingFast, setPendingFast] = useState<{ plan: string; hours: number } | null>(null);
  const [manualStartTime, setManualStartTime] = useState("");
  const [showReferralModal, setShowReferralModal] = useState(false);
  const [isInTribe] = useState(false); // Mock state for demo

  // Fetch real fasting data
  useEffect(() => {
    const fetchData = async () => {
      try {
        const res = await fastingApi.getCurrent();
        setSession(res.data);
      } catch (e) {
        console.error("Failed to fetch fasting data", e);
        setSession(null);
      }
    };
    fetchData();
    const interval = setInterval(fetchData, 60000); // Update every minute
    return () => clearInterval(interval);
  }, []);

  // Timer logic
  useEffect(() => {
    let interval: NodeJS.Timeout;
    if (session && session.status === 'active') {
      interval = setInterval(() => {
        const start = new Date(session.start_time).getTime();
        const now = new Date().getTime();
        setElapsed(Math.floor((now - start) / 1000));
      }, 1000);
    } else {
      setElapsed(0);
    }
    return () => clearInterval(interval);
  }, [session]);

  const handleStart = (plan: string, hours: number) => {
    setPendingFast({ plan, hours });
    setManualStartTime(""); // Reset
    setShowModal(true);
  };

  const confirmStart = async () => {
    if (pendingFast) {
      try {
        // Format manualStartTime to RFC3339 if present
        let startTimeStr = undefined;
        if (manualStartTime) {
          startTimeStr = new Date(manualStartTime).toISOString();
        }

        const res = await fastingApi.start(pendingFast.plan, pendingFast.hours, startTimeStr);
        setSession(res.data);
      } catch (error: any) {
        console.error("Start fast error:", error);
        if (error.response?.data?.error) {
          const errorMessage = error.response.data.error;
          alert(`Failed to start fast: ${errorMessage}`);

          // If session already exists, try to sync
          if (errorMessage.includes("active fasting session")) {
            try {
              const res = await fastingApi.getCurrent();
              setSession(res.data);
            } catch (e) {
              console.error("Failed to recover session", e);
            }
          }
        } else {
          alert("Failed to start fast. Please check console for details.");
        }
      } finally {
        setShowModal(false);
        setPendingFast(null);
      }
    }
  };

  const handleStop = async () => {
    if (!session) return;

    const hoursElapsed = elapsed / 3600;
    const goal = session.goal_hours;
    const progress = (hoursElapsed / goal) * 100;

    let message = "Are you sure you want to end your fast?";
    if (progress < 100) {
      if (progress > 80) {
        message = `You are nearly there! ${progress.toFixed(0)}% complete. Just ${(goal - hoursElapsed).toFixed(1)} hours to go. Are you sure you want to stop now?`;
      } else {
        message = `You are ${progress.toFixed(0)}% of the way to your ${goal}h goal. Stopping now will result in a discipline penalty. Continue?`;
      }
    } else {
      message = `Goal achieved! You've fasted for ${hoursElapsed.toFixed(1)} hours. Ready to end?`;
    }

    if (window.confirm(message)) {
      try {
        await fastingApi.stop();
        setSession(null);
        setElapsed(0);
        alert("Fast ended. Great work!");
      } catch (error) {
        alert("Failed to stop fast");
      }
    }
  };

  const getFastingStage = (hours: number) => {
    if (hours < 12) return "Digestion Phase";
    if (hours < 18) return "Fat Burning & Ketosis";
    if (hours < 24) return "Autophagy Initiating";
    return "Deep Cellular Repair";
  };

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
            <RefundCountdown daysRemaining={12} />
            <StreakCounter streak={8} />
            {/* Vault moved to Profile */}
          </div>
        </div>
      </header>

      {/* Trust Signals */}
      <SocialProofBanner />
      <ActivityTicker />

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

          {/* Vault Hero Card - NEW Priority */}
          <VaultHeroCard
            earned={session ? 5.50 : 0}
            deposit={20.00}
            daysUntilRefund={12}
            fastsRemaining={4}
          />

          {/* Daily Challenge Card */}
          <ChallengeCard />

          {/* Referral Teaser */}
          <ReferralTeaser onInvite={() => setShowReferralModal(true)} />

          {/* Hero Timer Card - Replaced with FastingClock */}
          <div className="flex flex-col items-center justify-center">
            <FastingClock
              elapsedSeconds={elapsed}
              goalHours={session?.goal_hours || 16}
              isFasting={!!session}
              stage={getFastingStage(elapsed / 3600)}
            />
          </div>

          {/* Controls */}
          <div className="flex justify-center w-full">
            {!session ? (
              <div className="grid grid-cols-2 gap-4 w-full max-w-md">
                <button onClick={() => handleStart('16_8', 16)} className="bg-emerald-600 hover:bg-emerald-700 p-4 rounded-xl font-bold transition-all text-white">
                  Start 16:8
                </button>
                <button onClick={() => handleStart('18_6', 18)} className="bg-emerald-600 hover:bg-emerald-700 p-4 rounded-xl font-bold transition-all text-white">
                  Start 18:6
                </button>
                <button onClick={() => handleStart('omad', 23)} className="bg-emerald-600 hover:bg-emerald-700 p-4 rounded-xl font-bold transition-all text-white">
                  OMAD (23h)
                </button>
                <button onClick={() => handleStart('24h', 24)} className="bg-purple-600 hover:bg-purple-700 p-4 rounded-xl font-bold transition-all text-white">
                  24h Reset
                </button>
                <button onClick={() => handleStart('custom', 12)} className="col-span-2 bg-slate-700 hover:bg-slate-600 p-4 rounded-xl font-bold transition-all text-white border border-slate-600">
                  Custom / Manual Start
                </button>
              </div>
            ) : (
              <button onClick={handleStop} className="bg-red-500 hover:bg-red-600 px-8 py-3 rounded-full font-bold shadow-lg transition-all text-white">
                End Fast
              </button>
            )}
          </div>


          {/* Bio-Narrative Timeline */}
          <FastingJourney fastingHours={elapsed / 3600} />

        </TabsContent>

        {/* TAB 2: TRIBE - Social Accountability */}
        <TabsContent value="tribe" className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
          {!isInTribe ? (
            <TribeDiscovery />
          ) : (
            <>
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
                    { name: 'You', status: `Fasting ${(elapsed / 3600).toFixed(0)}h`, streak: 8, alert: false },
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
            </>
          )}
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

          {/* Replaced PricingMechanism with VaultStatus */}
          <div className="h-full">
            <VaultStatus
              deposit={20.00}
              earned={session ? 5.50 : 0}
              potentialRefund={session ? 5.50 : 0}
            />
          </div>

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

      <FloatingActionMenu
        onLogWeight={() => console.log("Log Weight")}
        onLogWater={() => console.log("Log Water")}
        onStartFast={() => setShowModal(true)}
      />

      <MedicalModal
        isOpen={showModal}
        onConfirm={confirmStart}
        onCancel={() => setShowModal(false)}
        startTime={manualStartTime}
        onStartTimeChange={setManualStartTime}
        goalHours={pendingFast?.hours || 12}
        onGoalHoursChange={(h) => setPendingFast(prev => prev ? { ...prev, hours: h } : null)}
        showGoalInput={pendingFast?.plan === 'custom'}
      />

      <ReferralModal
        open={showReferralModal}
        onOpenChange={setShowReferralModal}
      />
    </div>
  );
};

export default Dashboard;
