import React from 'react';

interface VaultDashboardProps {
    deposit: number;
    earned: number;
    discipline: number;
    dailyPotential: number;
}

const VaultDashboard: React.FC<VaultDashboardProps> = ({ deposit, earned, discipline, dailyPotential }) => {
    const earnedPercentage = Math.min((earned / deposit) * 100, 100);
    const disciplineColor = discipline >= 80 ? 'text-green-500' : discipline >= 50 ? 'text-yellow-500' : 'text-red-500';

    return (
        <div className="bg-gray-800 rounded-2xl p-6 shadow-xl border border-gray-700">
            <h2 className="text-2xl font-bold text-white mb-6 flex items-center gap-2">
                <span className="text-3xl">🔒</span> The Vault
            </h2>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {/* Money at Risk */}
                <div className="bg-gray-900 rounded-xl p-5 border border-gray-700 relative overflow-hidden">
                    <div className="absolute top-0 right-0 p-2 opacity-10">
                        <svg xmlns="http://www.w3.org/2000/svg" className="h-24 w-24 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                    </div>
                    <p className="text-gray-400 text-sm font-medium uppercase tracking-wider">Money at Risk</p>
                    <p className="text-4xl font-bold text-white mt-2">${deposit.toFixed(2)}</p>
                    <p className="text-xs text-gray-500 mt-2">Held in secure escrow</p>
                </div>

                {/* Daily Potential */}
                <div className="bg-gray-900 rounded-xl p-5 border border-gray-700 relative overflow-hidden">
                    <div className="absolute top-0 right-0 p-2 opacity-10">
                        <svg xmlns="http://www.w3.org/2000/svg" className="h-24 w-24 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
                        </svg>
                    </div>
                    <p className="text-gray-400 text-sm font-medium uppercase tracking-wider">Daily Potential</p>
                    <p className="text-4xl font-bold text-green-500 mt-2">+${dailyPotential.toFixed(2)}</p>
                    <p className="text-xs text-gray-500 mt-2">Earn back by staying disciplined</p>
                </div>
            </div>

            {/* Earned Refund Progress */}
            <div className="mt-8">
                <div className="flex justify-between items-end mb-2">
                    <p className="text-gray-300 font-medium">Earned Refund</p>
                    <p className="text-white font-bold">${earned.toFixed(2)} <span className="text-gray-500 text-sm font-normal">/ ${deposit.toFixed(2)}</span></p>
                </div>
                <div className="w-full bg-gray-700 rounded-full h-4 overflow-hidden">
                    <div
                        className="bg-gradient-to-r from-green-500 to-emerald-400 h-4 rounded-full transition-all duration-1000 ease-out"
                        style={{ width: `${earnedPercentage}%` }}
                    ></div>
                </div>
                <p className="text-xs text-gray-500 mt-2 text-right">{earnedPercentage.toFixed(0)}% Recovered</p>
            </div>

            {/* Discipline Index */}
            <div className="mt-8 bg-gray-900 rounded-xl p-6 border border-gray-700 flex items-center justify-between">
                <div>
                    <p className="text-gray-400 text-sm font-medium uppercase tracking-wider">Discipline Index</p>
                    <p className={`text-5xl font-bold mt-2 ${disciplineColor}`}>{discipline}</p>
                </div>
                <div className="h-20 w-20 rounded-full border-4 border-gray-700 flex items-center justify-center relative">
                    <svg className="h-full w-full transform -rotate-90" viewBox="0 0 100 100">
                        <circle
                            className="text-gray-700 stroke-current"
                            strokeWidth="8"
                            cx="50"
                            cy="50"
                            r="40"
                            fill="transparent"
                        ></circle>
                        <circle
                            className={`${disciplineColor} stroke-current transition-all duration-1000 ease-out`}
                            strokeWidth="8"
                            strokeLinecap="round"
                            cx="50"
                            cy="50"
                            r="40"
                            fill="transparent"
                            strokeDasharray={`${discipline * 2.51} 251.2`}
                        ></circle>
                    </svg>
                </div>
            </div>
        </div>
    );
};

export default VaultDashboard;
