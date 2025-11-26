import React from 'react';

interface PricingTableProps {
    onSelectVault: () => void;
}

const PricingTable: React.FC<PricingTableProps> = ({ onSelectVault }) => {
    return (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8 max-w-4xl mx-auto mt-12">
            {/* Free Tier */}
            <div className="bg-gray-800 rounded-2xl p-8 border border-gray-700 flex flex-col">
                <h3 className="text-2xl font-bold text-white mb-2">Hero</h3>
                <p className="text-gray-400 mb-6">For those just starting their journey.</p>
                <div className="text-4xl font-bold text-white mb-6">$0<span className="text-lg text-gray-500 font-normal">/mo</span></div>

                <ul className="space-y-4 mb-8 flex-grow">
                    <li className="flex items-center text-gray-300">
                        <span className="text-green-500 mr-2">✓</span> Basic Fasting Timer
                    </li>
                    <li className="flex items-center text-gray-300">
                        <span className="text-green-500 mr-2">✓</span> Weight Logging
                    </li>
                    <li className="flex items-center text-gray-300">
                        <span className="text-green-500 mr-2">✓</span> 7-Day History
                    </li>
                </ul>

                <button className="w-full py-3 rounded-xl border border-gray-600 text-white font-bold hover:bg-gray-700 transition-all">
                    Current Plan
                </button>
            </div>

            {/* Vault Tier */}
            <div className="bg-gradient-to-b from-emerald-900 to-gray-900 rounded-2xl p-8 border-2 border-emerald-500 flex flex-col relative overflow-hidden shadow-2xl">
                <div className="absolute top-0 right-0 bg-emerald-500 text-white text-xs font-bold px-3 py-1 rounded-bl-lg">RECOMMENDED</div>
                <h3 className="text-2xl font-bold text-white mb-2">Vault Member</h3>
                <p className="text-emerald-200 mb-6">For those serious about results.</p>
                <div className="text-4xl font-bold text-white mb-6">$30<span className="text-lg text-emerald-200 font-normal">/mo</span></div>

                <ul className="space-y-4 mb-8 flex-grow">
                    <li className="flex items-center text-white">
                        <span className="text-emerald-400 mr-2">✓</span> <strong>Money at Risk</strong> (Earn back $20)
                    </li>
                    <li className="flex items-center text-white">
                        <span className="text-emerald-400 mr-2">✓</span> Unlimited History & Stats
                    </li>
                    <li className="flex items-center text-white">
                        <span className="text-emerald-400 mr-2">✓</span> Cortex AI Coach Access
                    </li>
                    <li className="flex items-center text-white">
                        <span className="text-emerald-400 mr-2">✓</span> Community Challenges
                    </li>
                </ul>

                <button
                    onClick={onSelectVault}
                    className="w-full py-3 rounded-xl bg-emerald-500 text-white font-bold hover:bg-emerald-400 shadow-lg hover:shadow-emerald-500/20 transition-all transform hover:-translate-y-1"
                >
                    Join the Vault
                </button>
                <p className="text-center text-xs text-emerald-300 mt-3">30-day money-back guarantee</p>
            </div>
        </div>
    );
};

export default PricingTable;
