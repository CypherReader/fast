import { Lock, DollarSign } from "lucide-react";

interface VaultStatusProps {
    deposit: number;
    earned: number;
    potentialRefund: number;
}

const VaultStatus = ({ deposit, earned, potentialRefund }: VaultStatusProps) => {
    const progress = (earned / deposit) * 100;

    return (
        <div className="bg-slate-900 border border-slate-800 rounded-xl p-4 flex flex-col justify-between relative overflow-hidden h-full">
            <div className="absolute top-0 right-0 p-2 opacity-10">
                <Lock className="w-12 h-12 text-white" />
            </div>

            <div className="flex items-center space-x-2 mb-2">
                <DollarSign className="w-4 h-4 text-yellow-400" />
                <span className="text-xs text-slate-400 uppercase tracking-wider font-bold">Vault Balance</span>
            </div>

            <div className="flex items-baseline space-x-1">
                <span className="text-2xl font-bold text-white">${earned.toFixed(2)}</span>
                <span className="text-xs text-slate-500">/ ${deposit.toFixed(2)}</span>
            </div>

            <div className="w-full bg-slate-800 rounded-full h-1.5 mt-3">
                <div
                    className="bg-gradient-to-r from-yellow-600 to-yellow-400 h-1.5 rounded-full shadow-[0_0_10px_rgba(234,179,8,0.4)]"
                    style={{ width: `${progress}%` }}
                ></div>
            </div>

            <div className="mt-2 flex justify-between text-xs">
                <span className="text-slate-500">Refundable</span>
                <span className="text-yellow-400 font-bold">${potentialRefund.toFixed(2)}</span>
            </div>
        </div>
    );
};

export default VaultStatus;
