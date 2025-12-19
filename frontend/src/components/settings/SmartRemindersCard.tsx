import { useState } from 'react';
import { Clock, Droplets, Sparkles, Timer, Info } from 'lucide-react';
import { Card } from '@/components/ui/card';
import { Switch } from '@/components/ui/switch';
import { Label } from '@/components/ui/label';
import { Slider } from '@/components/ui/slider';
import { Button } from '@/components/ui/button';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip';
import { useReminderSettings, useOptimalFastingWindow, formatHour } from '@/hooks/use-smart-reminders';

export function SmartRemindersCard() {
    const { settings, isLoading, updateSettings, isUpdating } = useReminderSettings();
    const { data: optimalWindow, isLoading: loadingWindow } = useOptimalFastingWindow();
    const [showAISuggestion, setShowAISuggestion] = useState(false);

    if (isLoading || !settings) {
        return (
            <Card className="p-6 animate-pulse">
                <div className="h-6 bg-muted rounded w-1/3 mb-4"></div>
                <div className="space-y-4">
                    <div className="h-16 bg-muted rounded"></div>
                    <div className="h-16 bg-muted rounded"></div>
                </div>
            </Card>
        );
    }

    return (
        <Card className="p-6">
            <h3 className="font-semibold text-lg mb-4 flex items-center gap-2">
                <Timer className="w-5 h-5 text-primary" />
                Smart Reminders
                <span className="text-xs bg-primary/20 text-primary px-2 py-0.5 rounded-full ml-2">AI-Powered</span>
            </h3>
            <div className="space-y-4">
                {/* Fast Start Reminder */}
                <div className="flex items-center justify-between py-3 border-b border-border">
                    <div>
                        <Label htmlFor="fast-start-reminder" className="flex items-center gap-2">
                            <Clock className="w-4 h-4 text-muted-foreground" />
                            Fast Start Reminder
                        </Label>
                        <p className="text-sm text-muted-foreground">
                            Get reminded at {formatHour(settings.preferred_fast_start_hour)} to start your fast
                        </p>
                    </div>
                    <Switch
                        id="fast-start-reminder"
                        checked={settings.reminder_fast_start}
                        onCheckedChange={(checked) => updateSettings({ ...settings, reminder_fast_start: checked })}
                        disabled={isUpdating}
                    />
                </div>

                {/* Preferred Start Time (show only if fast start is enabled) */}
                {settings.reminder_fast_start && (
                    <div className="py-3 border-b border-border">
                        <div className="flex items-center gap-2 mb-2">
                            <Label className="text-sm text-muted-foreground">Preferred Start Time</Label>
                            <TooltipProvider>
                                <Tooltip>
                                    <TooltipTrigger asChild>
                                        <Info className="w-4 h-4 text-muted-foreground cursor-help" />
                                    </TooltipTrigger>
                                    <TooltipContent className="max-w-[250px]">
                                        <p>Set the time you typically want to start your daily fast. You'll receive a reminder at this time.</p>
                                    </TooltipContent>
                                </Tooltip>
                            </TooltipProvider>
                        </div>
                        <div className="flex items-center gap-4">
                            <Slider
                                value={[settings.preferred_fast_start_hour]}
                                min={0}
                                max={23}
                                step={1}
                                onValueChange={([value]) => updateSettings({ ...settings, preferred_fast_start_hour: value })}
                                className="flex-1"
                                disabled={isUpdating}
                            />
                            <span className="min-w-[80px] text-right font-medium">
                                {formatHour(settings.preferred_fast_start_hour)}
                            </span>
                        </div>
                    </div>
                )}

                {/* Fast End Reminder */}
                <div className="flex items-center justify-between py-3 border-b border-border">
                    <div>
                        <Label htmlFor="fast-end-reminder" className="flex items-center gap-2">
                            <Clock className="w-4 h-4 text-muted-foreground" />
                            Fast End Reminder
                        </Label>
                        <p className="text-sm text-muted-foreground">Get notified 30 minutes before your fast ends</p>
                    </div>
                    <Switch
                        id="fast-end-reminder"
                        checked={settings.reminder_fast_end}
                        onCheckedChange={(checked) => updateSettings({ ...settings, reminder_fast_end: checked })}
                        disabled={isUpdating}
                    />
                </div>

                {/* Hydration Reminder */}
                <div className="flex items-center justify-between py-3 border-b border-border">
                    <div>
                        <Label htmlFor="hydration-reminder" className="flex items-center gap-2">
                            <Droplets className="w-4 h-4 text-blue-400" />
                            Hydration Reminders
                        </Label>
                        <p className="text-sm text-muted-foreground">
                            Remind me every {settings.hydration_interval_minutes} minutes during fasting
                        </p>
                    </div>
                    <Switch
                        id="hydration-reminder"
                        checked={settings.reminder_hydration}
                        onCheckedChange={(checked) => updateSettings({ ...settings, reminder_hydration: checked })}
                        disabled={isUpdating}
                    />
                </div>

                {/* Hydration Interval (show only if hydration is enabled) */}
                {settings.reminder_hydration && (
                    <div className="py-3 border-b border-border">
                        <Label className="text-sm text-muted-foreground mb-2 block">Hydration Interval</Label>
                        <div className="flex items-center gap-4">
                            <Slider
                                value={[settings.hydration_interval_minutes]}
                                min={30}
                                max={120}
                                step={15}
                                onValueChange={([value]) => updateSettings({ ...settings, hydration_interval_minutes: value })}
                                className="flex-1"
                                disabled={isUpdating}
                            />
                            <span className="min-w-[80px] text-right font-medium">
                                {settings.hydration_interval_minutes} min
                            </span>
                        </div>
                    </div>
                )}

                {/* AI Optimal Window Suggestion */}
                <div className="pt-3">
                    <Button
                        variant="outline"
                        className="w-full justify-start gap-2"
                        onClick={() => setShowAISuggestion(!showAISuggestion)}
                        disabled={loadingWindow}
                    >
                        <Sparkles className="w-4 h-4 text-amber-400" />
                        {showAISuggestion ? 'Hide' : 'Show'} AI Fasting Suggestion
                    </Button>

                    {showAISuggestion && optimalWindow && (
                        <div className="mt-4 p-4 bg-primary/5 rounded-lg border border-primary/20">
                            <div className="flex items-center gap-2 mb-2">
                                <Sparkles className="w-4 h-4 text-amber-400" />
                                <span className="font-medium text-sm">Your Optimal Fasting Window</span>
                            </div>
                            <p className="text-xl font-bold text-primary">
                                {optimalWindow.suggested_duration_hours}h fasting window
                            </p>
                            <p className="text-sm text-muted-foreground mt-2">
                                {optimalWindow.reasoning}
                            </p>
                            <div className="mt-3 flex gap-2">
                                <span className="text-xs bg-muted px-2 py-1 rounded">
                                    Confidence: {Math.round(optimalWindow.confidence_score * 100)}%
                                </span>
                            </div>
                        </div>
                    )}
                </div>
            </div>
        </Card>
    );
}
