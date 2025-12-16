import { useState } from 'react';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from '../ui/dialog';
import { Button } from '../ui/button';
import { Input } from '../ui/input';
import { Textarea } from '../ui/textarea';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '../ui/select';
import { Label } from '../ui/label';
import { useCreateTribe } from '../../hooks/use-tribes';
import type { CreateTribeRequest } from '../../api/types';

interface CreateTribeDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
}

export function CreateTribeDialog({ open, onOpenChange }: CreateTribeDialogProps) {
    const createTribe = useCreateTribe();
    const [formData, setFormData] = useState<CreateTribeRequest>({
        name: '',
        description: '',
        fasting_schedule: '16:8',
        primary_goal: 'weight_loss',
        privacy: 'public',
        category: [],
    });

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            await createTribe.mutateAsync(formData);
            onOpenChange(false);
            // Reset form
            setFormData({
                name: '',
                description: '',
                fasting_schedule: '16:8',
                primary_goal: 'weight_loss',
                privacy: 'public',
                category: [],
            });
        } catch (error) {
            // Error is handled by the hook
        }
    };

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
                <DialogHeader>
                    <DialogTitle>Create a New Tribe</DialogTitle>
                    <DialogDescription>
                        Start your own fasting community and connect with like-minded people
                    </DialogDescription>
                </DialogHeader>

                <form onSubmit={handleSubmit} className="space-y-6">
                    {/* Tribe Name */}
                    <div className="space-y-2">
                        <Label htmlFor="name">Tribe Name *</Label>
                        <Input
                            id="name"
                            placeholder="e.g., 16:8 Warriors"
                            value={formData.name}
                            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                            required
                            minLength={3}
                            maxLength={50}
                        />
                    </div>

                    {/* Description */}
                    <div className="space-y-2">
                        <Label htmlFor="description">Description *</Label>
                        <Textarea
                            id="description"
                            placeholder="Tell people what your tribe is about..."
                            value={formData.description}
                            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                            required
                            minLength={10}
                            maxLength={500}
                            rows={4}
                        />
                    </div>

                    {/* Fasting Schedule */}
                    <div className="space-y-2">
                        <Label htmlFor="fasting_schedule">Fasting Schedule *</Label>
                        <Select
                            value={formData.fasting_schedule}
                            onValueChange={(value) => setFormData({ ...formData, fasting_schedule: value })}
                        >
                            <SelectTrigger>
                                <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="16:8">16:8 Intermittent Fasting</SelectItem>
                                <SelectItem value="18:6">18:6 Extended Fasting</SelectItem>
                                <SelectItem value="omad">OMAD (One Meal A Day)</SelectItem>
                                <SelectItem value="custom">Custom Schedule</SelectItem>
                            </SelectContent>
                        </Select>
                    </div>

                    {/* Primary Goal */}
                    <div className="space-y-2">
                        <Label htmlFor="primary_goal">Primary Goal *</Label>
                        <Select
                            value={formData.primary_goal}
                            onValueChange={(value) => setFormData({ ...formData, primary_goal: value })}
                        >
                            <SelectTrigger>
                                <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="weight_loss">Weight Loss</SelectItem>
                                <SelectItem value="metabolic_health">Metabolic Health</SelectItem>
                                <SelectItem value="longevity">Longevity</SelectItem>
                                <SelectItem value="mental_clarity">Mental Clarity</SelectItem>
                                <SelectItem value="autophagy">Autophagy</SelectItem>
                                <SelectItem value="discipline">Discipline Building</SelectItem>
                            </SelectContent>
                        </Select>
                    </div>

                    {/* Privacy */}
                    <div className="space-y-2">
                        <Label htmlFor="privacy">Privacy *</Label>
                        <Select
                            value={formData.privacy}
                            onValueChange={(value: any) => setFormData({ ...formData, privacy: value })}
                        >
                            <SelectTrigger>
                                <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="public">Public - Anyone can join</SelectItem>
                                <SelectItem value="private">Private - Requires approval</SelectItem>
                                <SelectItem value="invite_only">Invite Only</SelectItem>
                            </SelectContent>
                        </Select>
                    </div>

                    {/* Rules (Optional) */}
                    <div className="space-y-2">
                        <Label htmlFor="rules">Community Rules (Optional)</Label>
                        <Textarea
                            id="rules"
                            placeholder="Set some ground rules for your community..."
                            value={formData.rules || ''}
                            onChange={(e) => setFormData({ ...formData, rules: e.target.value })}
                            maxLength={1000}
                            rows={3}
                        />
                    </div>

                    {/* Submit Button */}
                    <div className="flex gap-3 justify-end pt-4">
                        <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
                            Cancel
                        </Button>
                        <Button type="submit" disabled={createTribe.isPending}>
                            {createTribe.isPending ? 'Creating...' : 'Create Tribe'}
                        </Button>
                    </div>
                </form>
            </DialogContent>
        </Dialog>
    );
}
