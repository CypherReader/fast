import { motion } from 'framer-motion';
import { Clock, ChefHat, Users } from 'lucide-react';
import { Recipe } from '@/data/recipes';
import { cn } from '@/lib/utils';

interface RecipeCardProps {
  recipe: Recipe;
  onClick: () => void;
}

const dietColors: Record<string, string> = {
  keto: 'bg-secondary text-secondary-foreground',
  regular: 'bg-blue-500 text-white',
  vegetarian: 'bg-orange-500 text-white',
  vegan: 'bg-yellow-500 text-black',
};

const difficultyLabels: Record<string, string> = {
  easy: 'Easy',
  medium: 'Medium',
  hard: 'Hard',
};

const RecipeCard = ({ recipe, onClick }: RecipeCardProps) => {
  const totalTime = recipe.prepTime + recipe.cookTime;

  return (
    <motion.button
      onClick={onClick}
      className="bg-card border border-border rounded-2xl overflow-hidden text-left w-full hover:border-secondary/50 transition-all group"
      whileHover={{ scale: 1.02 }}
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
    >
      {/* Image */}
      <div className="relative h-48 overflow-hidden">
        <img
          src={recipe.imageUrl}
          alt={recipe.name}
          className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
          loading="lazy"
        />
        {/* Badges */}
        <div className="absolute top-3 left-3 flex gap-2">
          <span className={cn('px-2 py-1 text-xs font-semibold rounded-full', dietColors[recipe.dietType])}>
            {recipe.dietType.charAt(0).toUpperCase() + recipe.dietType.slice(1)}
          </span>
          <span className="px-2 py-1 text-xs font-semibold rounded-full bg-card/80 backdrop-blur-sm text-foreground">
            {difficultyLabels[recipe.difficulty]}
          </span>
        </div>
      </div>

      {/* Content */}
      <div className="p-5">
        <h3 className="font-semibold text-lg text-foreground mb-2 line-clamp-1">
          {recipe.name}
        </h3>
        <p className="text-sm text-muted-foreground line-clamp-2 mb-4">
          {recipe.description}
        </p>

        {/* Macros Bar */}
        <div className="flex gap-1 mb-4 h-2 rounded-full overflow-hidden bg-muted">
          <div 
            className="bg-secondary" 
            style={{ width: `${(recipe.protein / (recipe.protein + recipe.carbs + recipe.fat)) * 100}%` }}
            title={`Protein: ${recipe.protein}g`}
          />
          <div 
            className="bg-primary" 
            style={{ width: `${(recipe.carbs / (recipe.protein + recipe.carbs + recipe.fat)) * 100}%` }}
            title={`Carbs: ${recipe.carbs}g`}
          />
          <div 
            className="bg-accent" 
            style={{ width: `${(recipe.fat / (recipe.protein + recipe.carbs + recipe.fat)) * 100}%` }}
            title={`Fat: ${recipe.fat}g`}
          />
        </div>

        {/* Macro Labels */}
        <div className="flex justify-between text-xs text-muted-foreground mb-4">
          <span><span className="text-secondary font-medium">{recipe.protein}g</span> protein</span>
          <span><span className="text-primary font-medium">{recipe.carbs}g</span> carbs</span>
          <span><span className="text-accent font-medium">{recipe.fat}g</span> fat</span>
        </div>

        {/* Prep Info */}
        <div className="flex items-center justify-between text-xs text-muted-foreground">
          <div className="flex items-center gap-1">
            <Clock className="w-3.5 h-3.5" />
            <span>{totalTime} min</span>
          </div>
          <div className="flex items-center gap-1">
            <ChefHat className="w-3.5 h-3.5" />
            <span>{difficultyLabels[recipe.difficulty]}</span>
          </div>
          <div className="flex items-center gap-1">
            <Users className="w-3.5 h-3.5" />
            <span>{recipe.servings} servings</span>
          </div>
        </div>
      </div>
    </motion.button>
  );
};

export default RecipeCard;
