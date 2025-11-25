import React, { useState, useEffect } from 'react';
import { Dna, Clock } from 'lucide-react';
import { cortexApi } from '../../api/client';

const FASTING_STAGES = [
    { hours: 0, title: "Digestion", desc: "Blood sugar rises. Insulin is active. Energy is drawn from food.", color: "text-slate-400" },
    { hours: 4, title: "Glycogen Depletion", desc: "Insulin drops. The body begins searching for stored energy in the liver.", color: "text-blue-400" },
    { hours: 12, title: "Metabolic Switch", desc: "You are now entering Ketosis. Your liver is converting fat into ketones. Brain fog begins to clear.", color: "text-purple-400" },
    { hours: 16, title: "Autophagy", desc: "Cellular cleanup mode engaged. Your body is recycling damaged cells. HGH levels spike to preserve muscle.", color: "text-green-400" },
    { hours: 24, title: "Stem Cell Regen", desc: "Immune system reboot. Inflammation drops significantly. Deep healing active.", color: "text-yellow-400" }
];

const FALLBACK_INSIGHTS: Record<string, string> = {
    "Digestion": "Glucose levels are stabilizing. Your body is processing your last meal.",
    "Glycogen Depletion": "Liver glycogen stores are being utilized. You may feel a slight dip in energy.",
    "Metabolic Switch": "Ketone production has started. Your brain is switching to a more efficient fuel source.",
    "Autophagy": "Intracellular repair mechanisms are at peak efficiency. Old proteins are being recycled.",
    "Stem Cell Regen": "System-wide regeneration is occurring. Immune function is being reset."
};

interface FastingJourneyProps {
    fastingHours: number;
}

const FastingJourney: React.FC<FastingJourneyProps> = ({ fastingHours }) => {
    const [insight, setInsight] = useState<string>("");
    const [loadingInsight, setLoadingInsight] = useState(false);

    // Determine current stage based on hours
    // findLastIndex polyfill equivalent
    let currentStageIndex = 0;
    for (let i = FASTING_STAGES.length - 1; i >= 0; i--) {
        if (fastingHours >= FASTING_STAGES[i].hours) {
            currentStageIndex = i;
            break;
        }
    }

    const currentStage = FASTING_STAGES[currentStageIndex] || FASTING_STAGES[0];
    const nextStage = FASTING_STAGES[currentStageIndex + 1];

    useEffect(() => {
        const fetchInsight = async () => {
            setLoadingInsight(true);
            try {
                const res = await cortexApi.getInsight(fastingHours);
                setInsight(res.data.insight);
            } catch (error) {
                console.error("Failed to fetch insight", error);
                setInsight(FALLBACK_INSIGHTS[currentStage.title] || "Biological data unavailable. Neural link unstable.");
            } finally {
                setLoadingInsight(false);
            }
        };

        if (fastingHours > 0) {
            fetchInsight();
        }
    }, [fastingHours]);

    return (
        <div className="bg-slate-900 border border-slate-700 rounded-xl overflow-hidden flex flex-col h-full min-h-[400px]">
            <div className="bg-slate-800 p-3 border-b border-slate-700 flex justify-between items-center">
                <span className="font-semibold text-cyan-400 flex items-center">
                    <Dna className="w-4 h-4 mr-2" /> Biological Timeline
                </span>
                <span className="text-xs text-slate-500 uppercase flex items-center">
                    <Clock className="w-3 h-3 mr-1" /> T+{fastingHours.toFixed(1)} Hours
                </span>
            </div>

            <div className="p-6 flex-1 flex flex-col relative overflow-y-auto max-h-[400px]">
                {/* Timeline Visualization */}
                <div className="absolute left-8 top-6 bottom-6 w-0.5 bg-slate-800"></div>

                <div className="space-y-8 z-10">
                    {FASTING_STAGES.map((stage, idx) => {
                        const isActive = idx === currentStageIndex;
                        const isPast = idx < currentStageIndex;

                        return (
                            <div key={idx} className={`flex items-start pl-2 relative transition-all duration-500 ${isActive ? 'opacity-100 scale-100' : 'opacity-40 scale-95'}`}>
                                <div className={`w-4 h-4 rounded-full border-2 mt-1 mr-4 flex-shrink-0 z-20 bg-slate-900 ${isActive ? 'border-cyan-400 shadow-[0_0_10px_rgba(34,211,238,0.5)]' :
                                    isPast ? 'border-slate-600 bg-slate-800' : 'border-slate-800'
                                    }`}></div>

                                <div className={`${isActive ? 'bg-slate-800/80 p-4 rounded-lg border border-slate-700 w-full' : ''}`}>
                                    <div className="flex justify-between items-center mb-1">
                                        <span className={`font-bold text-sm ${isActive ? stage.color : 'text-slate-500'}`}>
                                            {stage.title}
                                        </span>
                                        <span className="text-xs font-mono text-slate-600">{stage.hours}h</span>
                                    </div>

                                    {isActive && (
                                        <div className="mt-2 space-y-3">
                                            <p className="text-sm text-slate-300 leading-relaxed">
                                                {stage.desc}
                                            </p>

                                            {/* DeepSeek "Voice" Box */}
                                            <div className="bg-cyan-950/30 border-l-2 border-cyan-500 p-3 rounded-r text-xs text-cyan-200/80 italic">
                                                {loadingInsight ? "Analyzing biological telemetry..." : `"${insight}"`}
                                            </div>
                                        </div>
                                    )}
                                </div>
                            </div>
                        );
                    })}
                </div>
            </div>

            {/* Next Milestone Footer */}
            {nextStage && (
                <div className="p-3 bg-slate-800 border-t border-slate-700">
                    <div className="flex justify-between items-center text-xs">
                        <span className="text-slate-400">Next Phase: {nextStage.title}</span>
                        <span className="text-white font-mono">in {(nextStage.hours - fastingHours).toFixed(1)} hrs</span>
                    </div>
                    <div className="w-full bg-slate-900 rounded-full h-1 mt-2">
                        <div
                            className="bg-slate-600 h-1 rounded-full transition-all duration-1000"
                            style={{ width: `${((fastingHours - currentStage.hours) / (nextStage.hours - currentStage.hours)) * 100}%` }}
                        ></div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default FastingJourney;
