import { useState } from "react";
import { Card, CardHeader } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Badge } from "@/components/ui/badge";
import { KnowledgeHub } from "@/components/community/KnowledgeHub";
import { recipeApi, Recipe } from "@/api/client";
import { Leaf, Beef, CheckCircle2 } from "lucide-react";

const Resources = () => {
    const [activeTab, setActiveTab] = useState("recipes");
    const [recipes, setRecipes] = useState<Recipe[]>([]);
    const [dietFilter, setDietFilter] = useState<string>("all");

    const fetchRecipes = async (diet?: string) => {
        try {
            const res = await recipeApi.list(diet === "all" ? undefined : diet);
            setRecipes(res.data);
        } catch (e) {
            console.error("Failed to fetch recipes", e);
        }
    };

    // Fetch recipes when tab changes to recipes or filter changes
    if (activeTab === "recipes" && recipes.length === 0) {
        fetchRecipes(dietFilter);
    }

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="animate-fade-in">
                <h1 className="text-2xl font-bold bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
                    Resources
                </h1>
                <p className="text-sm text-muted-foreground">Recipes and knowledge for your fasting journey</p>
            </div>

            {/* Tabs */}
            <Tabs value={activeTab} onValueChange={setActiveTab}>
                <TabsList className="grid w-full grid-cols-2">
                    <TabsTrigger value="recipes">Recipes</TabsTrigger>
                    <TabsTrigger value="knowledge">Knowledge</TabsTrigger>
                </TabsList>

                <TabsContent value="recipes" className="mt-6 space-y-6">
                    <div className="flex gap-2 overflow-x-auto pb-2">
                        {["all", "vegan", "vegetarian", "normal"].map((diet) => (
                            <Badge
                                key={diet}
                                variant={dietFilter === diet ? "default" : "outline"}
                                className="cursor-pointer capitalize px-4 py-2"
                                onClick={() => {
                                    setDietFilter(diet);
                                    fetchRecipes(diet);
                                }}
                            >
                                {diet}
                            </Badge>
                        ))}
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        {recipes.map((recipe) => (
                            <Card key={recipe.id} className="overflow-hidden hover:shadow-md transition-shadow">
                                <div className="relative h-48">
                                    <img src={recipe.image} alt={recipe.title} className="w-full h-full object-cover" />
                                    {recipe.is_simple && (
                                        <Badge className="absolute top-2 right-2 bg-green-500 hover:bg-green-600">
                                            <CheckCircle2 className="w-3 h-3 mr-1" />
                                            Simple
                                        </Badge>
                                    )}
                                </div>
                                <CardHeader className="p-4">
                                    <div className="flex justify-between items-start mb-2">
                                        <h3 className="font-bold text-lg">{recipe.title}</h3>
                                        <Badge variant="secondary" className="capitalize text-xs">
                                            {recipe.diet === 'vegan' && <Leaf className="w-3 h-3 mr-1 text-green-500" />}
                                            {recipe.diet === 'vegetarian' && <Leaf className="w-3 h-3 mr-1 text-yellow-500" />}
                                            {recipe.diet === 'normal' && <Beef className="w-3 h-3 mr-1 text-red-500" />}
                                            {recipe.diet}
                                        </Badge>
                                    </div>
                                    <p className="text-sm text-muted-foreground line-clamp-2">{recipe.description}</p>
                                    <div className="flex gap-4 mt-4 text-sm font-medium">
                                        <div>🔥 {recipe.calories} kcal</div>
                                        <div>🍞 {recipe.carbs}g net carbs</div>
                                    </div>
                                </CardHeader>
                            </Card>
                        ))}
                    </div>
                </TabsContent>

                <TabsContent value="knowledge" className="mt-6">
                    <KnowledgeHub />
                </TabsContent>
            </Tabs>
        </div>
    );
};

export default Resources;
