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
    name: string;
    slug: string;
    description: string;
    avatar_url?: string;
    cover_photo_url?: string;
    creator_id: string;
    fasting_schedule: string; // "16:8", "18:6", "omad", "custom"
    primary_goal: string; // "weight_loss", "metabolic_health", etc.
    category: string[]; // Array of category tags
    privacy: 'public' | 'private' | 'invite_only';
    rules?: string;
    member_count: number;
    active_member_count: number;
    created_at: string;
    updated_at: string;
    // Computed fields
    is_joined?: boolean; // Whether current user is a member
    user_role?: 'creator' | 'moderator' | 'member'; // Current user's role if member
}

export interface TribeMember {
    id: string;
    tribe_id: string;
    user_id: string;
    role: 'creator' | 'moderator' | 'member';
    status: 'active' | 'pending' | 'left';
    joined_at: string;
    left_at?: string;
    notifications_enabled: boolean;
    user_name: string;
    user_avatar: string;
    user_streak: number;
}

export interface CreateTribeRequest {
    name: string;
    description: string;
    fasting_schedule: string;
    primary_goal: string;
    category?: string[];
    privacy: 'public' | 'private' | 'invite_only';
    rules?: string;
    avatar_url?: string;
    cover_photo_url?: string;
}

export interface TribeStats {
    tribe_id: string;
    total_fasts: number;
    total_fasting_hours: number;
    average_member_streak: number;
    weekly_growth_percent: number;
    active_members_percent: number;
}

export interface ListTribesResponse {
    tribes: Tribe[];
    total: number;
    limit: number;
    offset: number;
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
