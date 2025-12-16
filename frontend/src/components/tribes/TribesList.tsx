import { useState } from 'react';
import { useTribes } from '../../hooks/use-tribes';
import { TribeCard } from './TribeCard';
import { Input } from '../ui/input';
import { Button } from '../ui/button';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '../ui/select';
import { Search, Loader2, Plus } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { CreateTribeDialog } from './CreateTribeDialog';

export function TribesList() {
    const navigate = useNavigate();
    const [search, setSearch] = useState('');
    const [fastingSchedule, setFastingSchedule] = useState('');
    const [sortBy, setSortBy] = useState<'newest' | 'popular' | 'active' | 'members'>('popular');
    const [showCreateDialog, setShowCreateDialog] = useState(false);

    const { data, isLoading, error } = useTribes({
        search,
        fasting_schedule: fastingSchedule,
        sort_by: sortBy,
        limit: 20,
    });

    return (
        <div className="max-w-7xl mx-auto p-6 space-y-6">
            {/* Header */}
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold">Discover Tribes</h1>
                    <p className="text-gray-600 mt-1">Find your fasting community</p>
                </div>
                <Button onClick={() => setShowCreateDialog(true)}>
                    <Plus className="w-4 h-4 mr-2" />
                    Create Tribe
                </Button>
            </div>

            {/* Filters */}
            <div className="flex gap-4 items-center bg-white p-4 rounded-lg shadow-sm">
                {/* Search */}
                <div className="flex-1 relative">
                    <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
                    <Input
                        placeholder="Search tribes..."
                        value={search}
                        onChange={(e) => setSearch(e.target.value)}
                        className="pl-9"
                    />
                </div>

                {/* Fasting Schedule Filter */}
                <Select value={fastingSchedule} onValueChange={setFastingSchedule}>
                    <SelectTrigger className="w-48">
                        <SelectValue placeholder="All Schedules" />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectItem value="">All Schedules</SelectItem>
                        <SelectItem value="16:8">16:8 Intermittent</SelectItem>
                        <SelectItem value="18:6">18:6 Extended</SelectItem>
                        <SelectItem value="omad">OMAD</SelectItem>
                        <SelectItem value="custom">Custom</SelectItem>
                    </SelectContent>
                </Select>

                {/* Sort By */}
                <Select value={sortBy} onValueChange={(v) => setSortBy(v as any)}>
                    <SelectTrigger className="w-36">
                        <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectItem value="popular">Most Popular</SelectItem>
                        <SelectItem value="newest">Newest</SelectItem>
                        <SelectItem value="active">Most Active</SelectItem>
                        <SelectItem value="members">Most Members</SelectItem>
                    </SelectContent>
                </Select>
            </div>

            {/* Results */}
            {isLoading && (
                <div className="flex justify-center items-center py-12">
                    <Loader2 className="w-8 h-8 animate-spin text-purple-600" />
                </div>
            )}

            {error && (
                <div className="text-center py-12">
                    <p className="text-red-600">Failed to load tribes. Please try again.</p>
                </div>
            )}

            {data && data.tribes.length === 0 && (
                <div className="text-center py-12">
                    <p className="text-gray-600">No tribes found. Try adjusting your filters or create one!</p>
                </div>
            )}

            {data && data.tribes.length > 0 && (
                <>
                    <p className="text-sm text-gray-600">
                        Showing {data.tribes.length} of {data.total} tribes
                    </p>
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                        {data.tribes.map((tribe) => (
                            <TribeCard key={tribe.id} tribe={tribe} />
                        ))}
                    </div>
                </>
            )}

            {/* Create Dialog */}
            <CreateTribeDialog open={showCreateDialog} onOpenChange={setShowCreateDialog} />
        </div>
    );
}
