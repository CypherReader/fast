import React from 'react';
import { TrendingDown, LucideIcon } from 'lucide-react';

interface MetricCardProps {
    icon: LucideIcon;
    title: string;
    value: string | number;
    unit: string;
    status?: string;
    trend?: string;
}

const MetricCard: React.FC<MetricCardProps> = ({ icon: Icon, title, value, unit, status, trend }) => (
    <div className="bg-slate-800/50 rounded-xl p-5 border border-slate-700/50 backdrop-blur-sm">
        <div className="flex justify-between items-start mb-2">
            <div className="p-2 bg-slate-700/50 rounded-lg">
                <Icon className="w-5 h-5 text-cyan-400" />
            </div>
            {status && (
                <span className={`text-xs px-2 py-1 rounded-full ${status === 'Optimal' ? 'bg-green-500/20 text-green-400' : 'bg-yellow-500/20 text-yellow-400'
                    }`}>
                    {status}
                </span>
            )}
        </div>
        <div className="mt-2">
            <span className="text-2xl font-bold text-white">{value}</span>
            <span className="text-sm text-slate-400 ml-1">{unit}</span>
        </div>
        <div className="mt-1 text-xs text-slate-500">
            {title}
        </div>
        {trend && (
            <div className="mt-3 text-xs flex items-center text-emerald-400">
                <TrendingDown className="w-3 h-3 mr-1" />
                {trend} vs last week
            </div>
        )}
    </div>
);

export default MetricCard;
