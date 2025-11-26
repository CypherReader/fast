import { useState, useEffect } from 'react';
import { Activity, Watch, RefreshCw, CheckCircle, AlertTriangle } from 'lucide-react';
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

interface TelemetryUplinkProps {
    onDataUpdate?: () => void;
}

export const TelemetryUplink = ({ onDataUpdate }: TelemetryUplinkProps) => {
    const [showManualModal, setShowManualModal] = useState(false);
    const [manualValue, setManualValue] = useState("");
    const [metricType, setMetricType] = useState<'steps' | 'weight'>('steps');
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

    const handleManualLog = async () => {
        if (!manualValue) return;
        setLoading(true);
        try {
            await api.post('/telemetry/manual', {
                type: metricType,
                value: parseFloat(manualValue),
                unit: metricType === 'steps' ? 'count' : 'kg'
            });
            setManualValue("");
            setShowManualModal(false);
            alert(`${metricType === 'steps' ? 'Steps' : 'Weight'} logged successfully (Trust Score: 0.5)`);
            if (onDataUpdate) onDataUpdate();
        } catch (error) {
            alert(`Failed to log ${metricType}`);
        } finally {
            setLoading(false);
        }
    };

    return (
        <Card className="bg-slate-900 border-slate-800">
            <CardHeader>
                <CardTitle className="flex items-center justify-between text-white">
                    <div className="flex items-center">
                        <Activity className="mr-2 h-5 w-5 text-cyan-400" />
                        Bio-Telemetry Uplink
                    </div>
                    <Button
                        size="sm"
                        variant="ghost"
                        className="text-xs text-slate-400 hover:text-white"
                        onClick={() => setShowManualModal(!showManualModal)}
                    >
                        + Manual Entry
                    </Button>
                </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
                {showManualModal && (
                    <div className="p-4 bg-slate-800 rounded-lg border border-slate-700 animate-in fade-in slide-in-from-top-2">
                        <h4 className="text-sm font-medium text-slate-200 mb-3">Log Manual Data</h4>

                        <div className="flex gap-2 mb-3">
                            <Button
                                size="sm"
                                variant={metricType === 'steps' ? 'default' : 'outline'}
                                onClick={() => setMetricType('steps')}
                                className={metricType === 'steps' ? 'bg-cyan-600' : 'border-slate-600 text-slate-300'}
                            >
                                Steps
                            </Button>
                            <Button
                                size="sm"
                                variant={metricType === 'weight' ? 'default' : 'outline'}
                                onClick={() => setMetricType('weight')}
                                className={metricType === 'weight' ? 'bg-cyan-600' : 'border-slate-600 text-slate-300'}
                            >
                                Weight
                            </Button>
                        </div>

                        <div className="flex gap-2">
                            <input
                                type="number"
                                placeholder={metricType === 'steps' ? "e.g. 5000" : "e.g. 75.5"}
                                value={manualValue}
                                onChange={(e) => setManualValue(e.target.value)}
                                className="flex-1 bg-slate-900 border border-slate-700 rounded-md px-3 py-2 text-sm text-white focus:outline-none focus:border-cyan-500"
                            />
                            <Button
                                size="sm"
                                className="bg-cyan-600 hover:bg-cyan-700"
                                onClick={handleManualLog}
                                disabled={loading || !manualValue}
                            >
                                Log
                            </Button>
                        </div>
                    </div>
                )}

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
