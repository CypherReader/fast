export interface WeightLog {
    id: string;
    user_id: string;
    weight_lbs: number;
    weight_kg: number;
    logged_at: string;
    created_at: string;
}

export interface HydrationLog {
    id: string;
    user_id: string;
    glasses_count: number;
    logged_date: string;
    created_at: string;
}

export interface ActivityLog {
    id: string;
    user_id: string;
    steps: number;
    distance_km: number;
    calories_burned: number;
    logged_date: string;
    created_at: string;
}

export interface FriendNetwork {
    id: string;
    user_id: string;
    friend_id: string;
    status: 'pending' | 'accepted' | 'blocked';
    created_at: string;
}

export interface Tribe {
    id: string;
    creator_id: string;
    name: string;
    description: string;
    is_public: boolean;
    member_count: number;
    created_at: string;
}

export interface Notification {
    id: string;
    user_id: string;
    type: string;
    title: string;
    message: string;
    read: boolean;
    link?: string;
    created_at: string;
}

export interface NotificationHistoryResponse {
    notifications: Notification[];
}

export interface LeaderboardEntry {
    user_id: string;
    user_name: string;
    total_fasting_hours: number;
    discipline_score: number;
    rank: number;
}

export interface SocialEvent {
    id: string;
    user_id: string;
    user_name: string;
    event_type: 'fast_completed' | 'tribe_joined' | 'challenge_won';
    data: string;
    created_at: string;
    likes: number;
    comments: number;
}
