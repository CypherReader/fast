import { useParams, useNavigate } from 'react-router-dom';
import { useTribe, useJoinTribe, useLeaveTribe } from '../hooks/use-tribes';
import { Button } from '@/components/ui/button';
import { useAuth } from '@/contexts/AuthContext';
import { Loader2, Users, Clock, Target, Calendar, Shield, ArrowLeft } from 'lucide-react';
import { toast } from 'sonner';

export default function TribeDetails() {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const { user } = useAuth();
    const { data: tribe, isLoading, error } = useTribe(id);
    const joinTribe = useJoinTribe();
    const leaveTribe = useLeaveTribe();

    if (isLoading) {
        return (
            <div className="flex justify-center items-center h-screen">
                <Loader2 className="w-8 h-8 animate-spin text-purple-600" />
            </div>
        );
    }

    if (error || !tribe) {
        return (
            <div className="flex flex-col items-center justify-center h-screen gap-4">
                <p className="text-red-500">Failed to load tribe details.</p>
                <Button onClick={() => navigate('/tribes')}>Back to Tribes</Button>
            </div>
        );
    }

    const handleJoin = () => {
        if (!user) {
            navigate('/login');
            return;
        }
        joinTribe.mutate(tribe.id);
    };

    const handleLeave = () => {
        if (confirm('Are you sure you want to leave this tribe?')) {
            leaveTribe.mutate(tribe.id);
        }
    };

    return (
        <div className="min-h-screen bg-gray-50 pb-20">
            {/* Cover Photo */}
            <div className="h-64 relative bg-gray-900">
                {tribe.cover_photo_url ? (
                    <img
                        src={tribe.cover_photo_url}
                        alt="Cover"
                        className="w-full h-full object-cover opacity-80"
                    />
                ) : (
                    <div className="w-full h-full bg-gradient-to-r from-purple-600 to-indigo-700 opacity-80" />
                )}

                <Button
                    variant="ghost"
                    className="absolute top-4 left-4 text-white hover:bg-white/20"
                    onClick={() => navigate('/tribes')}
                >
                    <ArrowLeft className="w-5 h-5 mr-2" />
                    Back to Tribes
                </Button>
            </div>

            <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 -mt-20 relative z-10">
                <div className="bg-white rounded-xl shadow-lg p-6 sm:p-8">
                    {/* Header */}
                    <div className="flex flex-col sm:flex-row sm:items-end gap-6 mb-8">
                        {/* Avatar */}
                        <div className="flex-shrink-0">
                            {tribe.avatar_url ? (
                                <img
                                    src={tribe.avatar_url}
                                    alt={tribe.name}
                                    className="w-32 h-32 rounded-xl border-4 border-white shadow-md object-cover bg-gray-100"
                                />
                            ) : (
                                <div className="w-32 h-32 rounded-xl border-4 border-white shadow-md bg-purple-100 flex items-center justify-center">
                                    <Users className="w-12 h-12 text-purple-600" />
                                </div>
                            )}
                        </div>

                        {/* Info */}
                        <div className="flex-1">
                            <h1 className="text-3xl font-bold text-gray-900 mb-2">{tribe.name}</h1>
                            <div className="flex flex-wrap gap-2 mb-4">
                                <span className={`px-3 py-1 rounded-full text-sm font-medium ${tribe.privacy === 'public' ? 'bg-green-100 text-green-800' :
                                        tribe.privacy === 'private' ? 'bg-orange-100 text-orange-800' :
                                            'bg-gray-100 text-gray-800'
                                    }`}>
                                    {tribe.privacy.charAt(0).toUpperCase() + tribe.privacy.slice(1)}
                                </span>
                                {tribe.category && JSON.parse(String(tribe.category)).map((cat: string) => (
                                    <span key={cat} className="px-3 py-1 rounded-full text-sm bg-purple-50 text-purple-700">
                                        {cat}
                                    </span>
                                ))}
                            </div>
                        </div>

                        {/* Actions */}
                        <div className="flex-shrink-0">
                            {tribe.is_joined ? (
                                <div className="flex gap-2">
                                    <Button variant="outline" onClick={() => { }} disabled>Member</Button>
                                    <Button variant="ghost" className="text-red-600 hover:text-red-700 hover:bg-red-50" onClick={handleLeave}>
                                        Leave Tribe
                                    </Button>
                                </div>
                            ) : (
                                <Button size="lg" onClick={handleJoin} disabled={joinTribe.isPending}>
                                    {joinTribe.isPending ? (
                                        <>
                                            <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                                            Joining...
                                        </>
                                    ) : (
                                        'Join Tribe'
                                    )}
                                </Button>
                            )}
                        </div>
                    </div>

                    {/* Stats Grid */}
                    <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
                        <div className="bg-gray-50 p-4 rounded-lg">
                            <div className="flex items-center text-gray-500 mb-1">
                                <Users className="w-4 h-4 mr-2" />
                                <span className="text-sm">Members</span>
                            </div>
                            <span className="text-2xl font-bold text-gray-900">{tribe.member_count}</span>
                        </div>
                        <div className="bg-gray-50 p-4 rounded-lg">
                            <div className="flex items-center text-gray-500 mb-1">
                                <Clock className="w-4 h-4 mr-2" />
                                <span className="text-sm">Fasting Schedule</span>
                            </div>
                            <span className="text-xl font-semibold text-gray-900">{tribe.fasting_schedule}</span>
                        </div>
                        <div className="bg-gray-50 p-4 rounded-lg">
                            <div className="flex items-center text-gray-500 mb-1">
                                <Target className="w-4 h-4 mr-2" />
                                <span className="text-sm">Primary Goal</span>
                            </div>
                            <span className="text-lg font-medium text-gray-900 capitalize">{tribe.primary_goal.replace('_', ' ')}</span>
                        </div>
                        <div className="bg-gray-50 p-4 rounded-lg">
                            <div className="flex items-center text-gray-500 mb-1">
                                <Calendar className="w-4 h-4 mr-2" />
                                <span className="text-sm">Created</span>
                            </div>
                            <span className="text-lg font-medium text-gray-900">
                                {new Date(tribe.created_at).toLocaleDateString()}
                            </span>
                        </div>
                    </div>

                    {/* Description & Rules */}
                    <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                        <div className="lg:col-span-2 space-y-6">
                            <section>
                                <h2 className="text-xl font-bold text-gray-900 mb-3">About</h2>
                                <p className="text-gray-600 leading-relaxed whitespace-pre-wrap">
                                    {tribe.description}
                                </p>
                            </section>

                            {tribe.rules && (
                                <section>
                                    <h2 className="text-xl font-bold text-gray-900 mb-3 flex items-center">
                                        <Shield className="w-5 h-5 mr-2 text-gray-400" />
                                        Rules
                                    </h2>
                                    <div className="bg-yellow-50 border border-yellow-100 rounded-lg p-4">
                                        <p className="text-gray-700 whitespace-pre-wrap">{tribe.rules}</p>
                                    </div>
                                </section>
                            )}
                        </div>

                        {/* Recent Activity / Members Placeholder */}
                        <div className="space-y-6">
                            <div className="bg-gray-50 rounded-xl p-6">
                                <h3 className="font-semibold text-gray-900 mb-4">Active Members</h3>
                                <div className="space-y-4">
                                    <p className="text-sm text-gray-500 italic">Member list coming soon...</p>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
