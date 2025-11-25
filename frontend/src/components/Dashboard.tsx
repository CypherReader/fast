import React, { useState, useEffect } from 'react';
import { fastingApi } from '../api/client';
import { FastingSession } from '../api/types';
import MedicalModal from './MedicalModal';
import CortexWidget from './CortexWidget';

const Dashboard: React.FC = () => {
    const [session, setSession] = useState<FastingSession | null>(null);
    const [elapsed, setElapsed] = useState(0);
    const [loading, setLoading] = useState(true);
    const [showModal, setShowModal] = useState(false);
    const [pendingFast, setPendingFast] = useState<{ plan: string; hours: number } | null>(null);

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
                const res = await fastingApi.start(pendingFast.plan, pendingFast.hours);
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

    const formatTime = (seconds: number) => {
        const h = Math.floor(seconds / 3600);
        const m = Math.floor((seconds % 3600) / 60);
        const s = seconds % 60;
        return `${h}h ${m}m ${s}s`;
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

            {/* Timer Circle */}
            <div className="relative w-64 h-64 rounded-full border-4 border-emerald-500 flex items-center justify-center mb-8 shadow-[0_0_20px_rgba(16,185,129,0.5)]">
                <div className="text-center">
                    {session ? (
                        <>
                            <div className="text-4xl font-mono font-bold">{formatTime(elapsed)}</div>
                            <div className="text-sm text-gray-400 mt-2">Goal: {session.goal_hours}h</div>
                            <div className="text-xs text-emerald-300 mt-1 animate-pulse">
                                {getFastingStage(elapsed / 3600)}
                            </div>
                        </>
                    ) : (
                        <div className="text-xl text-gray-400">Ready to Fast?</div>
                    )}
                </div>
            </div>

            {/* Controls */}
            {!session ? (
                <div className="grid grid-cols-2 gap-4 w-full max-w-md">
                    <button onClick={() => handleStart('16_8', 16)} className="bg-emerald-600 hover:bg-emerald-700 p-4 rounded-xl font-bold transition-all">
                        Start 16:8
                    </button>
                    <button onClick={() => handleStart('18_6', 18)} className="bg-emerald-600 hover:bg-emerald-700 p-4 rounded-xl font-bold transition-all">
                        Start 18:6
                    </button>
                    <button onClick={() => handleStart('omad', 23)} className="bg-emerald-600 hover:bg-emerald-700 p-4 rounded-xl font-bold transition-all">
                        OMAD (23h)
                    </button>
                    <button onClick={() => handleStart('24h', 24)} className="bg-purple-600 hover:bg-purple-700 p-4 rounded-xl font-bold transition-all">
                        24h Reset
                    </button>
                </div>
            ) : (
                <button onClick={handleStop} className="bg-red-500 hover:bg-red-600 px-8 py-3 rounded-full font-bold shadow-lg transition-all">
                    End Fast
                </button>
            )}

            {/* Price Ticker - The Lazy Tax */}
            <div className="mt-8 bg-gray-900 border border-red-900/50 p-6 rounded-2xl w-full max-w-md shadow-[0_0_20px_rgba(220,38,38,0.2)]">
                <div className="flex justify-between items-center mb-2">
                    <div>
                        <h3 className="text-lg font-bold text-red-500">The Lazy Tax</h3>
                        <div className="text-xs text-gray-400">Projected bill for next month</div>
                    </div>
                    <div className="text-4xl font-mono font-bold text-white">$48.50</div>
                </div>
                <div className="w-full bg-gray-800 h-3 rounded-full overflow-hidden">
                    <div className="bg-red-600 h-full w-[97%]"></div>
                </div>
                <div className="text-xs text-right text-red-400 mt-1">You are paying for your lack of discipline.</div>
            </div>

            {/* Cortex AI Widget */}
            <CortexWidget />

            <MedicalModal
                isOpen={showModal}
                onConfirm={confirmStart}
                onCancel={() => setShowModal(false)}
            />

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
        </div>
    );
};

export default Dashboard;
