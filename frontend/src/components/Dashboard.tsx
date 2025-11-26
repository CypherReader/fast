import React, { useState, useEffect } from 'react';
import { fastingApi } from '../api/client';
import { FastingSession } from '../api/types';
import MedicalModal from './MedicalModal';
import CortexWidget from './CortexWidget';
import FastingClock from './FastingClock';
import VaultDashboard from './VaultDashboard';
import VaultStatus from './VaultStatus';

const Dashboard: React.FC = () => {
    const [session, setSession] = useState<FastingSession | null>(null);
    const [elapsed, setElapsed] = useState(0);
    const [loading, setLoading] = useState(true);
    const [showModal, setShowModal] = useState(false);
    const [pendingFast, setPendingFast] = useState<{ plan: string; hours: number; startTime?: string } | null>(null);

    useEffect(() => {
        fetchCurrentSession();
    }, []);

    useEffect(() => {
        let interval: NodeJS.Timeout;
        if (session && session.status === 'active') {
            interval = setInterval(() => {
                const start = new Date(session.start_time).getTime();
                const now = new Date().getTime();
                setElapsed(Math.floor((now - start) / 1000));
            }, 1000);
        }
        return () => clearInterval(interval);
    }, [session]);

    const fetchCurrentSession = async () => {
        try {
            const res = await fastingApi.getCurrent();
            setSession(res.data);
        } catch (error) {
            console.error("No active session", error);
            setSession(null);
        } finally {
            setLoading(false);
        }
    };

    const handleStart = (plan: string, hours: number) => {
        setPendingFast({ plan, hours });
        setShowModal(true);
    };

    const confirmStart = async () => {
        if (pendingFast) {
            try {
                const res = await fastingApi.start(pendingFast.plan, pendingFast.hours, pendingFast.startTime);
                setSession(res.data);
            } catch (error) {
                alert("Failed to start fast");
            } finally {
                setShowModal(false);
                setPendingFast(null);
            }
        }
    };

    const handleStop = async () => {
        try {
            await fastingApi.stop();
            setSession(null);
            setElapsed(0);
        } catch (error) {
            alert("Failed to stop fast");
        }
    };

    const getFastingStage = (hours: number) => {
        if (hours < 12) return "Digestion Phase";
        if (hours < 18) return "Fat Burning & Ketosis";
        if (hours < 24) return "Autophagy Initiating";
        return "Deep Cellular Repair";
    };

    if (loading) return <div className="text-white">Loading...</div>;

    return (
        <div className="min-h-screen bg-gray-900 text-white p-6 flex flex-col items-center">
            <h1 className="text-3xl font-bold mb-8 text-emerald-400">FastingHero</h1>

            {/* Main Grid */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 w-full max-w-6xl">
                {/* Left Column: Fasting Clock & Vault */}
                <div className="lg:col-span-2 space-y-8">
                    <FastingClock
                        elapsedSeconds={elapsed}
                        goalHours={session?.goal_hours || 16}
                        isFasting={!!session}
                        stage={getFastingStage(elapsed / 3600)}
                    />

                    <VaultDashboard
                        deposit={20.00}
                        earned={session ? 2.00 : 0.00} // Mock data for now
                        discipline={85} // Mock data for now
                        dailyPotential={2.00}
                    />
                </div>

                {/* Right Column: Cortex & Stats */}
                <div className="space-y-8">
                    <CortexWidget />
                    <VaultStatus
                        deposit={20.00}
                        earned={session ? 5.50 : 0}
                        potentialRefund={session ? 5.50 : 0}
                    />
                </div>
            </div>

            {/* Controls */}
            {!session ? (
                <div className="grid grid-cols-2 gap-4 w-full max-w-md mt-8">
                    <button onClick={() => handleStart('16_8', 16)} className="bg-emerald-600 hover:bg-emerald-700 p-4 rounded-xl font-bold transition-all">
                        Start 16:8
                    </button>
                    <button onClick={() => handleStart('18_6', 18)} className="bg-emerald-600 hover:bg-emerald-700 p-4 rounded-xl font-bold transition-all">
                        Start 18:6
                    </button>
                    <button onClick={() => handleStart('omad', 23)} className="bg-emerald-600 hover:bg-emerald-700 p-4 rounded-xl font-bold transition-all">
                        OMAD (23h)
                    </button>
                    <button onClick={() => handleStart('custom', 16)} className="bg-emerald-600 hover:bg-emerald-700 p-4 rounded-xl font-bold transition-all">
                        Custom
                    </button>
                </div>
            ) : (
                <div className="mt-8 w-full max-w-md">
                    <button onClick={handleStop} className="bg-red-600 hover:bg-red-700 p-4 rounded-xl font-bold w-full transition-all">
                        Stop Fast
                    </button>
                </div>
            )}

            {/* Social Feed Teaser */}
            <div className="mt-8 w-full max-w-md">
                <h3 className="text-xl font-bold mb-4 text-purple-400">Community Pulse</h3>
                <div className="space-y-4">
                    <div className="bg-gray-800 p-4 rounded-xl border border-gray-700 flex items-center gap-3">
                        <div className="w-10 h-10 rounded-full bg-gradient-to-tr from-purple-500 to-pink-500 flex items-center justify-center font-bold">FK</div>
                        <div>
                            <div className="font-bold text-sm">FastingKing</div>
                            <div className="text-xs text-gray-400">Just hit 36 hours! Feeling amazing. 🔥</div>
                        </div>
                    </div>
                    <div className="bg-gray-800 p-4 rounded-xl border border-gray-700 flex items-center gap-3">
                        <div className="w-10 h-10 rounded-full bg-gradient-to-tr from-emerald-500 to-teal-500 flex items-center justify-center font-bold">KQ</div>
                        <div>
                            <div className="font-bold text-sm">KetoQueen</div>
                            <div className="text-xs text-gray-400">Broke my fast with avocado and eggs. 🥑</div>
                        </div>
                    </div>
                </div>
            </div>

            <MedicalModal
                isOpen={showModal}
                onConfirm={confirmStart}
                onCancel={() => setShowModal(false)}
                startTime={pendingFast?.startTime || ''}
                onStartTimeChange={(time) => setPendingFast(prev => prev ? { ...prev, startTime: time } : null)}
                goalHours={pendingFast?.hours || 16}
                onGoalHoursChange={(hours) => setPendingFast(prev => prev ? { ...prev, hours } : null)}
                showGoalInput={true}
            />
        </div>
    );
};

export default Dashboard;
