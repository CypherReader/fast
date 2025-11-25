import React, { useState, useEffect } from 'react';
import { Activity, Watch, Smartphone, RefreshCw, CheckCircle, XCircle, AlertTriangle } from 'lucide-react';
import { api } from '@/api/client';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

interface DeviceConnection {
    id: string;
    source: string;
    status: string;
    last_sync_at?: string;
}

const SUPPORTED_DEVICES = [
    { id: 'garmin', name: 'Garmin', icon: Watch, color: 'text-blue-400' },
    { id: 'apple_health', name: 'Apple Health', icon: Activity, color: 'text-pink-500' },
    { id: 'oura', name: 'Oura Ring', icon: CheckCircle, color: 'text-white' }, // Placeholder icon
];

export const TelemetryUplink = () => {
    const [connections, setConnections] = useState<DeviceConnection[]>([]);
    const [loading, setLoading] = useState(false);

    const fetchStatus = async () => {
        try {
            const res = await api.get<DeviceConnection[]>('/telemetry/status');
            setConnections(res.data || []);
        } catch (error) {
            console.error("Failed to fetch telemetry status", error);
        }
    };

    useEffect(() => {
        fetchStatus();
    }, []);

    const handleConnect = async (source: string) => {
        setLoading(true);
        try {
            await api.post('/telemetry/connect', { source });
            await fetchStatus();
        } catch (error) {
            alert("Failed to connect device");
        } finally {
            setLoading(false);
        }
    };

    const handleSync = async (source: string) => {
        setLoading(true);
        try {
            await api.post('/telemetry/sync', { source });
            await fetchStatus();
            alert("Sync successful");
        } catch (error) {
            alert("Sync failed");
        } finally {
            setLoading(false);
        }
    };

    return (
        <Card className="bg-slate-900 border-slate-800">
            <CardHeader>
                <CardTitle className="flex items-center text-white">
                    <Activity className="mr-2 h-5 w-5 text-cyan-400" />
                    Bio-Telemetry Uplink
                </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
                {SUPPORTED_DEVICES.map((device) => {
                    const connection = connections.find(c => c.source === device.id);
                    const isConnected = connection?.status === 'connected';
                    const Icon = device.icon;

                    return (
                        <div key={device.id} className="flex items-center justify-between p-3 bg-slate-800 rounded-lg border border-slate-700">
                            <div className="flex items-center space-x-3">
                                <div className={`p-2 rounded-full bg-slate-900 ${isConnected ? 'border-green-500/50 border' : 'border-slate-700 border'}`}>
                                    <Icon className={`h-5 w-5 ${device.color}`} />
                                </div>
                                <div>
                                    <h4 className="font-medium text-slate-200">{device.name}</h4>
                                    <div className="flex items-center text-xs">
                                        {isConnected ? (
                                            <>
                                                <span className="w-2 h-2 rounded-full bg-green-500 mr-2 animate-pulse"></span>
                                                <span className="text-green-400">Live Uplink</span>
                                            </>
                                        ) : (
                                            <>
                                                <span className="w-2 h-2 rounded-full bg-slate-600 mr-2"></span>
                                                <span className="text-slate-500">Disconnected</span>
                                            </>
                                        )}
                                    </div>
                                </div>
                            </div>

                            <div className="flex items-center space-x-2">
                                {isConnected ? (
                                    <>
                                        {connection?.last_sync_at && (
                                            <span className="text-xs text-slate-500 hidden sm:inline-block">
                                                Last: {new Date(connection.last_sync_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                                            </span>
                                        )}
                                        <Button
                                            size="sm"
                                            variant="outline"
                                            className="h-8 border-slate-600 text-slate-300 hover:bg-slate-700"
                                            onClick={() => handleSync(device.id)}
                                            disabled={loading}
                                        >
                                            <RefreshCw className={`h-3 w-3 ${loading ? 'animate-spin' : ''}`} />
                                        </Button>
                                    </>
                                ) : (
                                    <Button
                                        size="sm"
                                        className="h-8 bg-cyan-600 hover:bg-cyan-700 text-white"
                                        onClick={() => handleConnect(device.id)}
                                        disabled={loading}
                                    >
                                        Connect
                                    </Button>
                                )}
                            </div>
                        </div>
                    );
                })}

                <div className="mt-4 p-3 bg-yellow-900/10 border border-yellow-700/20 rounded-lg flex items-start">
                    <AlertTriangle className="h-4 w-4 text-yellow-500 mr-2 mt-0.5 flex-shrink-0" />
                    <p className="text-xs text-yellow-200/70">
                        <strong>Anti-Cheat Active:</strong> Manual entries carry a 50% Trust Score penalty. Connect a hardware device for full credit.
                    </p>
                </div>
            </CardContent>
        </Card>
    );
};
