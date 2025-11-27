import { X, CheckCircle, Share2, Bookmark, ExternalLink } from 'lucide-react';
import { Video, videos, categoryFilters } from '@/data/videos';
import { Dialog, DialogContent } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';

interface VideoPlayerProps {
  video: Video | null;
  open: boolean;
  onClose: () => void;
  onVideoSelect: (video: Video) => void;
}

const VideoPlayer = ({ video, open, onClose, onVideoSelect }: VideoPlayerProps) => {
  if (!video) return null;

  // Get related videos from same category
  const relatedVideos = videos
    .filter(v => v.category === video.category && v.id !== video.id)
    .slice(0, 4);

  const category = categoryFilters.find(c => c.id === video.category);

  return (
    <Dialog open={open} onOpenChange={onClose}>
      <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto p-0 bg-card border-border">
        {/* Video Player */}
        <div className="relative aspect-video bg-black">
          <iframe
            width="100%"
            height="100%"
            src={`https://www.youtube.com/embed/${video.youtubeId}?autoplay=1`}
            title={video.title}
            frameBorder="0"
            allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
            allowFullScreen
            className="absolute inset-0"
          />
          <button
            onClick={onClose}
            className="absolute top-4 right-4 w-10 h-10 rounded-full bg-black/60 flex items-center justify-center hover:bg-black/80 transition-colors z-10"
          >
            <X className="w-5 h-5 text-white" />
          </button>
        </div>

        <div className="p-6">
          {/* Title & Channel */}
          <div className="mb-4">
            <h2 className="text-xl font-bold text-foreground mb-3">{video.title}</h2>
            <div className="flex items-center justify-between flex-wrap gap-4">
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-full bg-muted flex items-center justify-center">
                  <span className="text-sm font-bold text-muted-foreground">
                    {video.channelName.charAt(0)}
                  </span>
                </div>
                <div>
                  <div className="flex items-center gap-1.5">
                    <span className="font-medium text-foreground">{video.channelName}</span>
                    {video.channelVerified && (
                      <CheckCircle className="w-4 h-4 text-secondary" />
                    )}
                  </div>
                  <span className="text-xs text-muted-foreground">
                    {video.viewCount} views • {new Date(video.uploadDate).toLocaleDateString()}
                  </span>
                </div>
              </div>
              <div className="flex gap-2">
                <Button variant="outline" size="sm" className="gap-2">
                  <Bookmark className="w-4 h-4" />
                  Save
                </Button>
                <Button variant="outline" size="sm" className="gap-2">
                  <Share2 className="w-4 h-4" />
                  Share
                </Button>
                <Button 
                  variant="outline" 
                  size="sm" 
                  className="gap-2"
                  onClick={() => window.open(`https://www.youtube.com/watch?v=${video.youtubeId}`, '_blank')}
                >
                  <ExternalLink className="w-4 h-4" />
                  YouTube
                </Button>
              </div>
            </div>
          </div>

          {/* Category Badge */}
          <div className="mb-4">
            <span className="px-3 py-1 bg-muted text-muted-foreground text-sm rounded-full">
              {category?.emoji} {category?.label}
            </span>
          </div>

          {/* Description */}
          <div className="bg-muted/50 rounded-xl p-4 mb-6">
            <p className="text-sm text-foreground whitespace-pre-line">
              {video.description}
            </p>
            {video.tags.length > 0 && (
              <div className="flex flex-wrap gap-2 mt-3 pt-3 border-t border-border">
                {video.tags.map(tag => (
                  <span key={tag} className="text-xs text-secondary">
                    #{tag}
                  </span>
                ))}
              </div>
            )}
          </div>

          {/* Related Videos */}
          {relatedVideos.length > 0 && (
            <div>
              <h3 className="font-semibold text-foreground mb-4">Related Videos</h3>
              <div className="space-y-3">
                {relatedVideos.map((relatedVideo) => (
                  <button
                    key={relatedVideo.id}
                    onClick={() => onVideoSelect(relatedVideo)}
                    className="w-full flex gap-3 p-2 rounded-lg hover:bg-muted/50 transition-colors text-left"
                  >
                    <div className="relative w-32 aspect-video rounded-lg overflow-hidden flex-shrink-0">
                      <img
                        src={relatedVideo.thumbnailUrl}
                        alt={relatedVideo.title}
                        className="w-full h-full object-cover"
                        onError={(e) => {
                          e.currentTarget.src = `https://img.youtube.com/vi/${relatedVideo.youtubeId}/hqdefault.jpg`;
                        }}
                      />
                      <div className="absolute bottom-1 right-1 px-1 py-0.5 bg-black/80 text-white text-xs rounded">
                        {relatedVideo.duration}
                      </div>
                    </div>
                    <div className="flex-1 min-w-0">
                      <h4 className="text-sm font-medium text-foreground line-clamp-2 mb-1">
                        {relatedVideo.title}
                      </h4>
                      <p className="text-xs text-muted-foreground">{relatedVideo.channelName}</p>
                      <p className="text-xs text-muted-foreground">{relatedVideo.viewCount} views</p>
                    </div>
                  </button>
                ))}
              </div>
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
};

export default VideoPlayer;
