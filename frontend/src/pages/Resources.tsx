import { useState, useMemo } from 'react';
import { motion } from 'framer-motion';
import { Search, SlidersHorizontal, X, ChefHat, Play, ChevronDown, BookOpen } from 'lucide-react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { recipes, dietFilters, mealTypeFilters, Recipe } from '@/data/recipes';
import { videos, categoryFilters, channels, Video } from '@/data/videos';
import RecipeCard from '@/components/recipes/RecipeCard';
import RecipeDetail from '@/components/recipes/RecipeDetail';
import VideoCard from '@/components/videos/VideoCard';
import VideoPlayer from '@/components/videos/VideoPlayer';
import { cn } from '@/lib/utils';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';

const Resources = () => {
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useSearchParams();
  const defaultTab = searchParams.get('tab') || 'recipes';

  // Recipe state
  const [recipeSearch, setRecipeSearch] = useState('');
  const [selectedDiet, setSelectedDiet] = useState<string>('all');
  const [selectedMealType, setSelectedMealType] = useState<string>('all');
  const [showFilters, setShowFilters] = useState(false);
  const [selectedRecipe, setSelectedRecipe] = useState<Recipe | null>(null);
  const [recipeDetailOpen, setRecipeDetailOpen] = useState(false);

  // Video state
  const [videoSearch, setVideoSearch] = useState('');
  const [selectedCategory, setSelectedCategory] = useState<string>('all');
  const [selectedVideo, setSelectedVideo] = useState<Video | null>(null);
  const [playerOpen, setPlayerOpen] = useState(false);

  // Filtered recipes
  const filteredRecipes = useMemo(() => {
    return recipes.filter(recipe => {
      const matchesSearch = recipeSearch === '' ||
        recipe.name.toLowerCase().includes(recipeSearch.toLowerCase()) ||
        recipe.description.toLowerCase().includes(recipeSearch.toLowerCase()) ||
        recipe.tags.some(tag => tag.toLowerCase().includes(recipeSearch.toLowerCase()));
      const matchesDiet = selectedDiet === 'all' || recipe.dietType === selectedDiet;
      const matchesMealType = selectedMealType === 'all' || recipe.mealType === selectedMealType;
      return matchesSearch && matchesDiet && matchesMealType;
    });
  }, [recipeSearch, selectedDiet, selectedMealType]);

  // Filtered videos
  const filteredVideos = useMemo(() => {
    return videos.filter(video => {
      const matchesSearch = videoSearch === '' ||
        video.title.toLowerCase().includes(videoSearch.toLowerCase()) ||
        video.description.toLowerCase().includes(videoSearch.toLowerCase()) ||
        video.channelName.toLowerCase().includes(videoSearch.toLowerCase()) ||
        video.tags.some(tag => tag.toLowerCase().includes(videoSearch.toLowerCase()));
      const matchesCategory = selectedCategory === 'all' || video.category === selectedCategory;
      return matchesSearch && matchesCategory;
    });
  }, [videoSearch, selectedCategory]);

  const featuredVideos = videos.filter(v => v.featured).slice(0, 3);
  const currentCategory = categoryFilters.find(c => c.id === selectedCategory);

  const hasRecipeFilters = recipeSearch || selectedDiet !== 'all' || selectedMealType !== 'all';
  const hasVideoFilters = videoSearch || selectedCategory !== 'all';

  const clearRecipeFilters = () => {
    setRecipeSearch('');
    setSelectedDiet('all');
    setSelectedMealType('all');
  };

  const clearVideoFilters = () => {
    setVideoSearch('');
    setSelectedCategory('all');
  };

  const handleTabChange = (value: string) => {
    setSearchParams({ tab: value });
  };

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="border-b border-border bg-card/50 backdrop-blur-sm sticky top-0 z-40">
        <div className="container mx-auto px-4 py-4 flex items-center justify-between">
          <button
            onClick={() => navigate('/dashboard')}
            className="flex items-center gap-2"
          >
            <img src="/fasthero.png" alt="FastingHero" className="w-8 h-8 rounded-lg" />
            <span className="font-bold text-lg text-foreground">Resources</span>
          </button>
          <Button variant="ghost" onClick={() => navigate('/dashboard')}>
            Back to Dashboard
          </Button>
        </div>
      </header>

      <main className="container mx-auto px-4 py-8 max-w-7xl">
        {/* Page Header */}
        <div className="text-center mb-8">
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-2">
            Knowledge Hub
          </h1>
          <p className="text-muted-foreground">
            Recipes, videos, and resources to support your fasting journey
          </p>
        </div>

        {/* Tabs */}
        <Tabs defaultValue={defaultTab} onValueChange={handleTabChange} className="w-full">
          <TabsList className="grid w-full max-w-md mx-auto grid-cols-2 mb-8">
            <TabsTrigger value="recipes" className="gap-2">
              <ChefHat className="w-4 h-4" />
              Recipes
            </TabsTrigger>
            <TabsTrigger value="videos" className="gap-2">
              <Play className="w-4 h-4" />
              Videos
            </TabsTrigger>
          </TabsList>

          {/* Recipes Tab */}
          <TabsContent value="recipes" className="mt-0">
            {/* Search & Filter Bar */}
            <div className="bg-background/95 backdrop-blur-sm py-4 border-b border-border mb-6">
              <div className="flex gap-3 mb-4">
                <div className="relative flex-1">
                  <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-muted-foreground" />
                  <Input
                    placeholder="Search recipes..."
                    value={recipeSearch}
                    onChange={(e) => setRecipeSearch(e.target.value)}
                    className="pl-10 bg-card border-border"
                  />
                  {recipeSearch && (
                    <button
                      onClick={() => setRecipeSearch('')}
                      className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
                    >
                      <X className="w-4 h-4" />
                    </button>
                  )}
                </div>
                <Button
                  variant="outline"
                  onClick={() => setShowFilters(!showFilters)}
                  className={cn(showFilters && 'border-secondary text-secondary')}
                >
                  <SlidersHorizontal className="w-4 h-4 mr-2" />
                  Filters
                </Button>
              </div>

              {/* Diet Filter Tabs */}
              <div className="flex gap-2 overflow-x-auto pb-2 scrollbar-hide">
                {dietFilters.map((filter) => (
                  <button
                    key={filter.id}
                    onClick={() => setSelectedDiet(filter.id)}
                    className={cn(
                      'px-4 py-2 rounded-full border whitespace-nowrap transition-all text-sm font-medium',
                      selectedDiet === filter.id
                        ? 'bg-secondary text-secondary-foreground border-secondary shadow-emerald-glow'
                        : 'bg-transparent border-border text-muted-foreground hover:border-secondary/50'
                    )}
                  >
                    {filter.emoji} {filter.label}
                  </button>
                ))}
              </div>

              {/* Advanced Filters */}
              {showFilters && (
                <motion.div
                  initial={{ height: 0, opacity: 0 }}
                  animate={{ height: 'auto', opacity: 1 }}
                  exit={{ height: 0, opacity: 0 }}
                  className="pt-4 border-t border-border mt-4"
                >
                  <div className="flex flex-wrap gap-4">
                    <div>
                      <label className="text-sm text-muted-foreground mb-2 block">Meal Type</label>
                      <div className="flex gap-2">
                        {mealTypeFilters.map((filter) => (
                          <button
                            key={filter.id}
                            onClick={() => setSelectedMealType(filter.id)}
                            className={cn(
                              'px-3 py-1.5 rounded-lg border text-sm transition-all',
                              selectedMealType === filter.id
                                ? 'bg-primary text-primary-foreground border-primary'
                                : 'bg-transparent border-border text-muted-foreground hover:border-primary/50'
                            )}
                          >
                            {filter.label}
                          </button>
                        ))}
                      </div>
                    </div>
                  </div>
                </motion.div>
              )}
            </div>

            {/* Results Count */}
            <div className="flex items-center justify-between mb-6">
              <p className="text-sm text-muted-foreground">
                {filteredRecipes.length} recipe{filteredRecipes.length !== 1 ? 's' : ''} found
              </p>
              {hasRecipeFilters && (
                <button onClick={clearRecipeFilters} className="text-sm text-secondary hover:underline">
                  Clear all filters
                </button>
              )}
            </div>

            {/* Recipe Grid */}
            {filteredRecipes.length > 0 ? (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {filteredRecipes.map((recipe) => (
                  <RecipeCard
                    key={recipe.id}
                    recipe={recipe}
                    onClick={() => {
                      setSelectedRecipe(recipe);
                      setRecipeDetailOpen(true);
                    }}
                  />
                ))}
              </div>
            ) : (
              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                className="text-center py-16"
              >
                <div className="w-24 h-24 rounded-full bg-muted mx-auto mb-4 flex items-center justify-center">
                  <ChefHat className="w-12 h-12 text-muted-foreground" />
                </div>
                <h3 className="text-xl font-semibold text-foreground mb-2">No recipes found</h3>
                <p className="text-muted-foreground mb-4">Try adjusting your filters or search terms</p>
                <Button onClick={clearRecipeFilters} variant="outline">Clear Filters</Button>
              </motion.div>
            )}
          </TabsContent>

          {/* Videos Tab */}
          <TabsContent value="videos" className="mt-0">
            {/* Search Bar */}
            <div className="max-w-2xl mx-auto mb-8">
              <div className="relative">
                <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-muted-foreground" />
                <Input
                  placeholder="Search videos..."
                  value={videoSearch}
                  onChange={(e) => setVideoSearch(e.target.value)}
                  className="pl-12 h-12 bg-card border-border text-lg"
                />
                {videoSearch && (
                  <button
                    onClick={() => setVideoSearch('')}
                    className="absolute right-4 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
                  >
                    <X className="w-5 h-5" />
                  </button>
                )}
              </div>
            </div>

            <div className="flex flex-col lg:flex-row gap-8">
              {/* Sidebar - Categories */}
              <aside className="lg:w-64 flex-shrink-0">
                {/* Mobile Dropdown */}
                <div className="lg:hidden mb-4">
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="outline" className="w-full justify-between">
                        {currentCategory?.emoji} {currentCategory?.label}
                        <ChevronDown className="w-4 h-4 ml-2" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent className="w-56">
                      {categoryFilters.map((category) => (
                        <DropdownMenuItem
                          key={category.id}
                          onClick={() => setSelectedCategory(category.id)}
                          className={cn(selectedCategory === category.id && 'bg-secondary/20')}
                        >
                          {category.emoji} {category.label}
                        </DropdownMenuItem>
                      ))}
                    </DropdownMenuContent>
                  </DropdownMenu>
                </div>

                {/* Desktop Sidebar */}
                <div className="hidden lg:block space-y-2 sticky top-24">
                  <h3 className="font-semibold text-foreground mb-3">Categories</h3>
                  {categoryFilters.map((category) => (
                    <button
                      key={category.id}
                      onClick={() => setSelectedCategory(category.id)}
                      className={cn(
                        'w-full flex items-center gap-2 px-3 py-2 rounded-lg text-sm font-medium transition-all text-left',
                        selectedCategory === category.id
                          ? 'bg-secondary text-secondary-foreground'
                          : 'text-muted-foreground hover:bg-muted hover:text-foreground'
                      )}
                    >
                      <span>{category.emoji}</span>
                      <span>{category.label}</span>
                    </button>
                  ))}

                  {/* Featured Channels */}
                  <div className="mt-8 pt-6 border-t border-border">
                    <h3 className="font-semibold text-foreground mb-3">Featured Experts</h3>
                    <div className="space-y-2">
                      {channels.slice(0, 5).map((channel) => (
                        <div key={channel.id} className="flex items-center gap-2 text-sm">
                          <div className="w-8 h-8 rounded-full bg-muted flex items-center justify-center">
                            <span className="text-xs font-bold text-muted-foreground">
                              {channel.name.charAt(0)}
                            </span>
                          </div>
                          <div className="flex-1 min-w-0">
                            <p className="text-foreground truncate">{channel.name}</p>
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                </div>
              </aside>

              {/* Main Content */}
              <div className="flex-1">
                {/* Featured Videos */}
                {!hasVideoFilters && (
                  <motion.div
                    className="mb-8"
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                  >
                    <h2 className="text-xl font-semibold text-foreground mb-4">Featured Videos</h2>
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                      {featuredVideos.map((video) => (
                        <VideoCard
                          key={video.id}
                          video={video}
                          onClick={() => {
                            setSelectedVideo(video);
                            setPlayerOpen(true);
                          }}
                        />
                      ))}
                    </div>
                  </motion.div>
                )}

                {/* Results Count */}
                <div className="flex items-center justify-between mb-4">
                  <p className="text-sm text-muted-foreground">
                    {filteredVideos.length} video{filteredVideos.length !== 1 ? 's' : ''} found
                  </p>
                  {hasVideoFilters && (
                    <button onClick={clearVideoFilters} className="text-sm text-secondary hover:underline">
                      Clear filters
                    </button>
                  )}
                </div>

                {/* Video Grid */}
                {filteredVideos.length > 0 ? (
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    {filteredVideos.map((video) => (
                      <VideoCard
                        key={video.id}
                        video={video}
                        onClick={() => {
                          setSelectedVideo(video);
                          setPlayerOpen(true);
                        }}
                      />
                    ))}
                  </div>
                ) : (
                  <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    className="text-center py-16"
                  >
                    <div className="w-24 h-24 rounded-full bg-muted mx-auto mb-4 flex items-center justify-center">
                      <Play className="w-12 h-12 text-muted-foreground" />
                    </div>
                    <h3 className="text-xl font-semibold text-foreground mb-2">No videos found</h3>
                    <p className="text-muted-foreground mb-4">Try adjusting your search or category filter</p>
                    <Button onClick={clearVideoFilters} variant="outline">Clear Filters</Button>
                  </motion.div>
                )}
              </div>
            </div>
          </TabsContent>
        </Tabs>
      </main>

      {/* Modals */}
      <RecipeDetail
        recipe={selectedRecipe}
        open={recipeDetailOpen}
        onClose={() => setRecipeDetailOpen(false)}
      />
      <VideoPlayer
        video={selectedVideo}
        open={playerOpen}
        onClose={() => setPlayerOpen(false)}
        onVideoSelect={(video) => setSelectedVideo(video)}
      />
    </div>
  );
};

export default Resources;