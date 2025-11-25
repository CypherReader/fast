import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer } from 'recharts';
import { Activity } from 'lucide-react';
import { TelemetryUplink } from '@/components/bio/TelemetryUplink';

// Mock data for now, will replace with API call later
const mockStepsData = [
    { day: 'Mon', steps: 6000 },
    { day: 'Tue', steps: 8500 },
    { day: 'Wed', steps: 7200 },
    { day: 'Thu', steps: 10500 },
    { day: 'Fri', steps: 5000 },
    { day: 'Sat', steps: 12000 },
    { day: 'Sun', steps: 9000 },
];

export const ActivityDashboard = () => {
    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <h2 className="text-3xl font-bold tracking-tight text-white">Activity Tracking</h2>
            </div>

            {/* Telemetry Uplink Section */}
            <div className="grid gap-4 md:grid-cols-1">
                <TelemetryUplink />
            </div>

            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Total Steps Today</CardTitle>
                        <Activity className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">9,000</div>
                        <p className="text-xs text-muted-foreground">+20.1% from yesterday</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Distance</CardTitle>
                        <Activity className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">6.5 km</div>
                        <p className="text-xs text-muted-foreground">Daily goal: 5 km</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Calories Burned</CardTitle>
                        <Activity className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">450 kcal</div>
                        <p className="text-xs text-muted-foreground">Active calories</p>
                    </CardContent>
                </Card>
            </div>

            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
                <Card className="col-span-4">
                    <CardHeader>
                        <CardTitle>Weekly Steps</CardTitle>
                    </CardHeader>
                    <CardContent className="pl-2">
                        <ResponsiveContainer width="100%" height={350}>
                            <BarChart data={mockStepsData}>
                                <XAxis
                                    dataKey="day"
                                    stroke="#888888"
                                    fontSize={12}
                                    tickLine={false}
                                    axisLine={false}
                                />
                                <YAxis
                                    stroke="#888888"
                                    fontSize={12}
                                    tickLine={false}
                                    axisLine={false}
                                    tickFormatter={(value) => `${value}`}
                                />
                                <Tooltip
                                    cursor={{ fill: 'transparent' }}
                                    contentStyle={{ backgroundColor: '#1f2937', border: 'none', borderRadius: '8px', color: '#fff' }}
                                />
                                <Bar dataKey="steps" fill="#adfa1d" radius={[4, 4, 0, 0]} />
                            </BarChart>
                        </ResponsiveContainer>
                    </CardContent>
                </Card>
                <Card className="col-span-3">
                    <CardHeader>
                        <CardTitle>Recent Routes</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="flex items-center justify-center h-[300px] bg-muted/20 rounded-md border border-dashed">
                            <p className="text-muted-foreground">Map View Placeholder</p>
                            {/* We will add Leaflet map here later */}
                        </div>
                    </CardContent>
                </Card>
            </div>
        </div>
    );
};
