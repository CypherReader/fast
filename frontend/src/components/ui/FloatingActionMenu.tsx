import { useState } from "react";
import { Plus, Scale, Droplet, Timer } from "lucide-react";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

interface FloatingActionMenuProps {
    onLogWeight: () => void;
    onLogWater: () => void;
    onStartFast: () => void;
}

export const FloatingActionMenu = ({ onLogWeight, onLogWater, onStartFast }: FloatingActionMenuProps) => {
    const [isOpen, setIsOpen] = useState(false);

    return (
        <div className="fixed bottom-24 right-6 z-50 flex flex-col items-end space-y-3">
            {isOpen && (
                <>
                    <Button
                        size="icon"
                        onClick={() => { onLogWeight(); setIsOpen(false); }}
                        className="h-12 w-12 rounded-full bg-slate-800 border border-slate-700 shadow-lg hover:bg-slate-700 animate-in slide-in-from-bottom-2 fade-in duration-200"
                    >
                        <Scale className="w-5 h-5 text-emerald-400" />
                        <span className="sr-only">Log Weight</span>
                    </Button>
                    <Button
                        size="icon"
                        onClick={() => { onLogWater(); setIsOpen(false); }}
                        className="h-12 w-12 rounded-full bg-slate-800 border border-slate-700 shadow-lg hover:bg-slate-700 animate-in slide-in-from-bottom-4 fade-in duration-300"
                    >
                        <Droplet className="w-5 h-5 text-blue-400" />
                        <span className="sr-only">Log Water</span>
                    </Button>
                    <Button
                        size="icon"
                        onClick={() => { onStartFast(); setIsOpen(false); }}
                        className="h-12 w-12 rounded-full bg-slate-800 border border-slate-700 shadow-lg hover:bg-slate-700 animate-in slide-in-from-bottom-6 fade-in duration-400"
                    >
                        <Timer className="w-5 h-5 text-purple-400" />
                        <span className="sr-only">Start Fast</span>
                    </Button>
                </>
            )}

            <Button
                size="icon"
                onClick={() => setIsOpen(!isOpen)}
                className={cn(
                    "h-14 w-14 rounded-full shadow-xl transition-all duration-300",
                    isOpen ? "bg-red-500 hover:bg-red-600 rotate-45" : "bg-cyan-600 hover:bg-cyan-700"
                )}
            >
                <Plus className="w-8 h-8 text-white" />
            </Button>
        </div>
    );
};
