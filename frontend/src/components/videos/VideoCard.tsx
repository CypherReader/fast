import { motion } from 'framer-motion';
import { Play, CheckCircle } from 'lucide-react';
import { Video, categoryFilters } from '@/data/videos';
import { cn } from '@/lib/utils';

interface VideoCardProps {
  video: Video;
  onClick: () => void;
}

const categoryColors: Record<string, string> = {
  science: 'bg-blue-500',
  benefits: 'bg-secondary',
  mental: 'bg-purple-500',
  'getting-started': 'bg-primary',
  'breaking-fast': 'bg-orange-500',
  'weight-loss': 'bg-pink-500',
  exercise: 'bg-red-500',
  faq: 'bg-cyan-500',
};

const VideoCard = ({ video, onClick }: VideoCardProps) => {
  const category = categoryFilters.find(c => c.id === video.category);

  return (
    <motion.button
      onClick={onClick}
      className="bg-card border border-border rounded-2xl overflow-hidden text-left w-full hover:border-secondary/50 transition-all group"
      whileHover={{ scale: 1.02 }}
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
    >
      {/* Thumbnail */}
      <div className="relative aspect-video overflow-hidden">
        <img
          src={video.thumbnailUrl}
          alt={video.title}
          className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
          loading="lazy"
          onError={(e) => {
            e.currentTarget.src = `https://img.youtube.com/vi/${video.youtubeId}/hqdefault.jpg`;
          }}
        />
        
        {/* Play Button Overlay */}
        <div className="absolute inset-0 flex items-center justify-center bg-black/20 opacity-0 group-hover:opacity-100 transition-opacity">
          <div className="w-16 h-16 rounded-full bg-primary/90 flex items-center justify-center transform group-hover:scale-110 transition-transform">
            <Play className="w-7 h-7 text-primary-foreground ml-1" fill="currentColor" />
          </div>
        </div>

        {/* Duration Badge */}
        <div className="absolute bottom-2 right-2 px-2 py-1 bg-black/80 text-white text-xs font-medium rounded">
          {video.duration}
        </div>

        {/* Category Badge */}
        <div className={cn(
          'absolute top-2 left-2 px-2 py-1 text-xs font-semibold rounded text-white',
          categoryColors[video.category]
        )}>
          {category?.emoji} {category?.label.split(' ')[0]}
        </div>
      </div>

      {/* Content */}
      <div className="p-4">
        <h3 className="font-semibold text-foreground mb-2 line-clamp-2 leading-tight">
          {video.title}
        </h3>
        
        {/* Channel */}
        <div className="flex items-center gap-1.5 mb-2">
          <span className="text-sm text-muted-foreground">{video.channelName}</span>
          {video.channelVerified && (
            <CheckCircle className="w-3.5 h-3.5 text-secondary" />
          )}
        </div>

        {/* Stats */}
        <div className="flex items-center gap-2 text-xs text-muted-foreground">
          <span>{video.viewCount} views</span>
          <span>•</span>
          <span>{new Date(video.uploadDate).toLocaleDateString('en-US', { 
            month: 'short', 
            year: 'numeric' 
          })}</span>
        </div>

        {/* Description Preview */}
        <p className="text-sm text-muted-foreground mt-2 line-clamp-2">
          {video.description}
        </p>
      </div>
    </motion.button>
  );
};

export default VideoCard;
