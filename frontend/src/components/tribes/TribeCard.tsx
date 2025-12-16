import { Card } from '../ui/card';
import { Button } from '../ui/button';
import { Users, Clock, Target } from 'lucide-react';
import type { Tribe } from '../../api/types';
import { useNavigate } from 'react-router-dom';

interface TribeCardProps {
    tribe: Tribe;
}

export function TribeCard({ tribe }: TribeCardProps) {
    const navigate = useNavigate();

    const getScheduleDisplay = (schedule: string) => {
        const scheduleMap: Record<string, string> = {
            '16:8': '16:8 Intermittent',
            '18:6': '18:6 Extended',
            'omad': 'OMAD',
            'custom': 'Custom',
        };
        return scheduleMap[schedule] || schedule;
    };

    return (
        <Card
            className="overflow-hidden cursor-pointer hover:shadow-lg transition-all duration-300 hover:-translate-y-1"
            onClick={() => navigate(`/tribes/${tribe.id}`)}
        >
            {/* Cover Photo */}
            {tribe.cover_photo_url ? (
                <div className="h-32 bg-gradient-to-r from-purple-500 to-indigo-600 relative">
                    <img
                        src={tribe.cover_photo_url}
                        alt={tribe.name}
                        className="w-full h-full object-cover"
                    />
                </div>
            ) : (
                <div className="h-32 bg-gradient-to-r from-purple-500 to-indigo-600" />
            )}

            {/* Content */}
            <div className="p-4 relative">
                {/* Avatar overlaying cover */}
                {tribe.avatar_url && (
                    <div className="absolute -top-10 left-4">
                        <img
                            src={tribe.avatar_url}
                            alt={tribe.name}
                            className="w-20 h-20 rounded-full border-4 border-white object-cover"
                        />
                    </div>
                )}

                <div className={tribe.avatar_url ? 'mt-12' : ''}>
                    {/* Badge */}
                    <div className="flex items-center gap-2 mb-2">
                        <span
                            className={`text-xs px-2 py-1 rounded-full ${tribe.privacy === 'public'
                                    ? 'bg-green-100 text-green-700'
                                    : tribe.privacy === 'private'
                                        ? 'bg-orange-100 text-orange-700'
                                        : 'bg-purple-100 text-purple-700'
                                }`}
                        >
                            {tribe.privacy}
                        </span>
                        {tribe.is_joined && (
                            <span className="text-xs px-2 py-1 rounded-full bg-blue-100 text-blue-700">
                                Joined
                            </span>
                        )}
                    </div>

                    {/* Title */}
                    <h3 className="text-xl font-bold mb-2">{tribe.name}</h3>

                    {/* Description */}
                    <p className="text-sm text-gray-600 mb-4 line-clamp-2">{tribe.description}</p>

                    {/* Stats */}
                    <div className="grid grid-cols-3 gap-2 mb-4">
                        <div className="flex items-center gap-1.5 text-sm text-gray-600">
                            <Users className="w-4 h-4" />
                            <span>{tribe.member_count}</span>
                        </div>
                        <div className="flex items-center gap-1.5 text-sm text-gray-600">
                            <Clock className="w-4 h-4" />
                            <span>{getScheduleDisplay(tribe.fasting_schedule)}</span>
                        </div>
                        <div className="flex items-center gap-1.5 text-sm text-gray-600">
                            <Target className="w-4 h-4" />
                            <span className="capitalize text-xs">{tribe.primary_goal.replace('_', ' ')}</span>
                        </div>
                    </div>

                    {/* CTA */}
                    <Button
                        variant={tribe.is_joined ? 'outline' : 'default'}
                        className="w-full"
                        onClick={(e) => {
                            e.stopPropagation();
                            navigate(`/tribes/${tribe.id}`);
                        }}
                    >
                        {tribe.is_joined ? 'View Tribe' : 'Join Tribe'}
                    </Button>
                </div>
            </div>
        </Card>
    );
}
