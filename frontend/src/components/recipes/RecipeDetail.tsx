import { useState } from 'react';
import { motion } from 'framer-motion';
import { X, Clock, Users, Printer, Share2, Bookmark, Lightbulb } from 'lucide-react';
import { Recipe } from '@/data/recipes';
import { Dialog, DialogContent } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import { cn } from '@/lib/utils';

interface RecipeDetailProps {
  recipe: Recipe | null;
  open: boolean;
  onClose: () => void;
}

const dietColors: Record<string, string> = {
  keto: 'bg-secondary text-secondary-foreground',
  regular: 'bg-blue-500 text-white',
  vegetarian: 'bg-orange-500 text-white',
  vegan: 'bg-yellow-500 text-black',
};

const RecipeDetail = ({ recipe, open, onClose }: RecipeDetailProps) => {
  const [checkedIngredients, setCheckedIngredients] = useState<Set<number>>(new Set());

  if (!recipe) return null;

  const totalTime = recipe.prepTime + recipe.cookTime;
  const totalMacros = recipe.protein + recipe.carbs + recipe.fat;

  const toggleIngredient = (index: number) => {
    const newChecked = new Set(checkedIngredients);
    if (newChecked.has(index)) {
      newChecked.delete(index);
    } else {
      newChecked.add(index);
    }
    setCheckedIngredients(newChecked);
  };

  return (
    <Dialog open={open} onOpenChange={onClose}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto p-0 bg-card border-border">
        {/* Hero Image */}
        <div className="relative h-64 md:h-80">
          <img
            src={recipe.imageUrl}
            alt={recipe.name}
            className="w-full h-full object-cover"
          />
          <button
            onClick={onClose}
            className="absolute top-4 right-4 w-10 h-10 rounded-full bg-card/80 backdrop-blur-sm flex items-center justify-center hover:bg-card transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
          <div className="absolute bottom-4 left-4 flex gap-2">
            <span className={cn('px-3 py-1.5 text-sm font-semibold rounded-full', dietColors[recipe.dietType])}>
              {recipe.dietType.charAt(0).toUpperCase() + recipe.dietType.slice(1)}
            </span>
          </div>
        </div>

        <div className="p-6">
          {/* Header */}
          <div className="flex items-start justify-between mb-4">
            <div>
              <h2 className="text-2xl font-bold text-foreground mb-2">{recipe.name}</h2>
              <p className="text-muted-foreground">{recipe.description}</p>
            </div>
            <Button variant="ghost" size="icon" className="flex-shrink-0">
              <Bookmark className="w-5 h-5" />
            </Button>
          </div>

          {/* Quick Info */}
          <div className="grid grid-cols-4 gap-4 p-4 bg-muted/50 rounded-xl mb-6">
            <div className="text-center">
              <div className="text-lg font-bold text-foreground">{recipe.prepTime}</div>
              <div className="text-xs text-muted-foreground">Prep (min)</div>
            </div>
            <div className="text-center">
              <div className="text-lg font-bold text-foreground">{recipe.cookTime}</div>
              <div className="text-xs text-muted-foreground">Cook (min)</div>
            </div>
            <div className="text-center">
              <div className="text-lg font-bold text-foreground">{totalTime}</div>
              <div className="text-xs text-muted-foreground">Total (min)</div>
            </div>
            <div className="text-center">
              <div className="text-lg font-bold text-foreground">{recipe.servings}</div>
              <div className="text-xs text-muted-foreground">Servings</div>
            </div>
          </div>

          {/* Nutrition Facts */}
          <div className="bg-muted/50 rounded-xl p-4 mb-6">
            <h3 className="font-semibold text-foreground mb-3">Nutrition Facts</h3>
            <div className="grid grid-cols-3 gap-4 mb-4">
              <div className="text-center">
                <div className="text-2xl font-bold text-primary">{recipe.calories}</div>
                <div className="text-xs text-muted-foreground">Calories</div>
              </div>
              <div className="text-center">
                <div className="text-2xl font-bold text-secondary">{recipe.protein}g</div>
                <div className="text-xs text-muted-foreground">Protein</div>
              </div>
              <div className="text-center">
                <div className="text-2xl font-bold text-accent">{recipe.fat}g</div>
                <div className="text-xs text-muted-foreground">Fat</div>
              </div>
            </div>
            <div className="grid grid-cols-3 gap-4">
              <div className="text-center">
                <div className="text-lg font-semibold text-foreground">{recipe.carbs}g</div>
                <div className="text-xs text-muted-foreground">Total Carbs</div>
              </div>
              <div className="text-center">
                <div className="text-lg font-semibold text-foreground">{recipe.netCarbs}g</div>
                <div className="text-xs text-muted-foreground">Net Carbs</div>
              </div>
              <div className="text-center">
                <div className="text-lg font-semibold text-foreground">{recipe.fiber}g</div>
                <div className="text-xs text-muted-foreground">Fiber</div>
              </div>
            </div>
            {/* Macro Bar */}
            <div className="flex gap-1 mt-4 h-3 rounded-full overflow-hidden">
              <div 
                className="bg-secondary" 
                style={{ width: `${(recipe.protein / totalMacros) * 100}%` }}
              />
              <div 
                className="bg-primary" 
                style={{ width: `${(recipe.carbs / totalMacros) * 100}%` }}
              />
              <div 
                className="bg-accent" 
                style={{ width: `${(recipe.fat / totalMacros) * 100}%` }}
              />
            </div>
            <div className="flex justify-between text-xs text-muted-foreground mt-2">
              <span>Protein</span>
              <span>Carbs</span>
              <span>Fat</span>
            </div>
          </div>

          {/* Ingredients */}
          <div className="mb-6">
            <h3 className="font-semibold text-foreground mb-3">Ingredients</h3>
            <div className="space-y-2">
              {recipe.ingredients.map((ingredient, index) => (
                <div
                  key={index}
                  className={cn(
                    'flex items-center gap-3 p-2 rounded-lg transition-colors',
                    checkedIngredients.has(index) ? 'bg-secondary/10' : 'hover:bg-muted/50'
                  )}
                >
                  <Checkbox
                    checked={checkedIngredients.has(index)}
                    onCheckedChange={() => toggleIngredient(index)}
                  />
                  <span className={cn(
                    'text-sm',
                    checkedIngredients.has(index) ? 'text-muted-foreground line-through' : 'text-foreground'
                  )}>
                    <span className="font-medium">{ingredient.quantity}</span> {ingredient.item}
                    {ingredient.optional && <span className="text-muted-foreground ml-1">(optional)</span>}
                  </span>
                </div>
              ))}
            </div>
          </div>

          {/* Instructions */}
          <div className="mb-6">
            <h3 className="font-semibold text-foreground mb-3">Instructions</h3>
            <div className="space-y-4">
              {recipe.instructions.map((instruction, index) => (
                <motion.div
                  key={index}
                  className="flex gap-4"
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: 0.1 * index }}
                >
                  <div className="flex-shrink-0 w-8 h-8 rounded-full bg-primary/20 flex items-center justify-center">
                    <span className="text-sm font-bold text-primary">{index + 1}</span>
                  </div>
                  <p className="text-sm text-foreground pt-1">{instruction}</p>
                </motion.div>
              ))}
            </div>
          </div>

          {/* Tips */}
          {recipe.tips && recipe.tips.length > 0 && (
            <div className="bg-primary/10 border border-primary/20 rounded-xl p-4 mb-6">
              <div className="flex items-center gap-2 mb-3">
                <Lightbulb className="w-5 h-5 text-primary" />
                <h3 className="font-semibold text-foreground">Pro Tips</h3>
              </div>
              <ul className="space-y-2">
                {recipe.tips.map((tip, index) => (
                  <li key={index} className="text-sm text-muted-foreground flex items-start gap-2">
                    <span className="text-primary">•</span>
                    {tip}
                  </li>
                ))}
              </ul>
            </div>
          )}

          {/* Actions */}
          <div className="flex gap-3">
            <Button variant="outline" className="flex-1 gap-2">
              <Share2 className="w-4 h-4" />
              Share
            </Button>
            <Button variant="outline" className="flex-1 gap-2">
              <Printer className="w-4 h-4" />
              Print
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
};

export default RecipeDetail;
