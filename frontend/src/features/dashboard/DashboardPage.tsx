import { useEffect, useState } from 'react';
import { api } from '../../api/client';
import { FastingSession } from '../../api/types';
import { Button } from '../../components/ui/button';

interface DashboardPageProps {
    onLogout: () => void;
}

export function DashboardPage({ onLogout }: DashboardPageProps) {
    const [session, setSession] = useState<FastingSession | null>(null);
    const [loading, setLoading] = useState(true);
    const [elapsed, setElapsed] = useState('');

    const fetchSession = async () => {
        try {
            const res = await api.get('/fasting/current');
            setSession(res.data);
        } catch (err) {
            setSession(null);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchSession();
    }, []);

    useEffect(() => {
        if (!session || !session.start_time) return;

        const interval = setInterval(() => {
            const start = new Date(session.start_time).getTime();
            const now = new Date().getTime();
            const diff = now - start;

            const hours = Math.floor(diff / (1000 * 60 * 60));
            const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
            const seconds = Math.floor((diff % (1000 * 60)) / 1000);

            setElapsed(`${hours}h ${minutes}m ${seconds}s`);
        }, 1000);

        return () => clearInterval(interval);
    }, [session]);

    const handleStartFast = async () => {
        try {
            const res = await api.post('/fasting/start', {
                plan_type: '16:8',
                target_duration_hours: 16,
            });
            setSession(res.data);
        } catch (err) {
            console.error('Failed to start fast', err);
        }
    };

    const handleStopFast = async () => {
        try {
            await api.post('/fasting/stop');
            setSession(null);
            setElapsed('');
        } catch (err) {
            console.error('Failed to stop fast', err);
        }
    };

    if (loading) {
        return <div className="flex justify-center items-center h-full">Loading...</div>;
    }

    return (
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
            <div className="md:flex md:items-center md:justify-between mb-8">
                <div className="min-w-0 flex-1">
                    <h2 className="text-2xl font-bold leading-7 text-gray-900 sm:truncate sm:text-3xl sm:tracking-tight">
                        Fasting Dashboard
                    </h2>
                </div>
                <div className="mt-4 flex md:ml-4 md:mt-0">
                    <Button variant="outline" onClick={onLogout}>
                        Sign out
                    </Button>
                </div>
            </div>

            <div className="bg-white shadow sm:rounded-lg">
                <div className="px-4 py-5 sm:p-6 text-center">
                    {session ? (
                        <div className="space-y-6">
                            <h3 className="text-base font-semibold leading-6 text-gray-900">
                                You are currently fasting
                            </h3>
                            <div className="text-5xl font-bold tracking-tight text-indigo-600">
                                {elapsed}
                            </div>
                            <div className="text-sm text-gray-500">
                                Started at: {new Date(session.start_time).toLocaleString()}
                            </div>
                            <div className="mt-5">
                                <Button variant="secondary" onClick={handleStopFast}>
                                    End Fast
                                </Button>
                            </div>
                        </div>
                    ) : (
                        <div className="space-y-6">
                            <h3 className="text-base font-semibold leading-6 text-gray-900">
                                Ready to start your fast?
                            </h3>
                            <div className="mt-5">
                                <Button size="lg" onClick={handleStartFast}>
                                    Start 16:8 Fast
                                </Button>
                            </div>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}
