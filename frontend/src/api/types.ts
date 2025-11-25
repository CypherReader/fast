export interface User {
  id: string;
  email: string;
  is_premium: boolean;
  subscription_tier: 'free' | 'premium' | 'elite';
}

export interface FastingSession {
  id: string;
  user_id: string;
  start_time: string;
  end_time?: string;
  goal_hours: number;
  target_duration_hours?: number;
  plan_type: 'beginner' | '16_8' | '18_6' | 'omad' | '24h' | '36h' | 'extended';
  status: 'active' | 'completed' | 'cancelled';
}

export interface KetoEntry {
  id: string;
  user_id: string;
  logged_at: string;
  ketone_level?: number;
  acetone_level?: number;
  source: 'manual' | 'device';
}

export interface SocialPost {
  id: string;
  user_id: string;
  username?: string;
  content: string;
  image_url?: string;
  type: 'streak' | 'meal';
  likes?: number;
  created_at: string;
}
