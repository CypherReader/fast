import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { toast } from "@/components/ui/use-toast";

const Contract = () => {
    const navigate = useNavigate();
    const [weight, setWeight] = useState("");
    const [reason, setReason] = useState("");
    const [agreed, setAgreed] = useState(false);

    const handleSign = () => {
        if (!weight || !reason || !agreed) {
            toast({
                title: "Contract Incomplete",
                description: "You must fill out all fields and agree to the terms.",
                variant: "destructive",
            });
            return;
        }

        // In a real app, we would send this to the backend
        // For MVP, we'll just simulate success and redirect
        toast({
            title: "Contract Signed",
            description: "Welcome to Neuro-Fast. The clock is ticking.",
        });
        navigate("/");
    };

    return (
        <div className="min-h-screen flex items-center justify-center bg-background p-4">
            <Card className="w-full max-w-lg border-primary/50 shadow-lg shadow-primary/20">
                <CardHeader>
                    <CardTitle className="text-3xl font-bold text-primary tracking-tighter">The Ulysses Pact</CardTitle>
                    <CardDescription className="text-lg">
                        "I hereby authorize Neuro-Fast to charge me <span className="text-destructive font-bold">$50/month</span>."
                    </CardDescription>
                </CardHeader>
                <CardContent className="space-y-6">
                    <div className="space-y-2">
                        <Label htmlFor="weight">Current Weight (lbs)</Label>
                        <Input
                            id="weight"
                            type="number"
                            placeholder="e.g. 200"
                            value={weight}
                            onChange={(e) => setWeight(e.target.value)}
                        />
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="reason">Why are you doing this? (Your Anchor)</Label>
                        <Input
                            id="reason"
                            placeholder="e.g. To see my kids grow up..."
                            value={reason}
                            onChange={(e) => setReason(e.target.value)}
                        />
                    </div>

                    <div className="bg-muted p-4 rounded-md text-sm text-muted-foreground">
                        <p>
                            By signing this contract, you agree to the <strong>Lazy Tax</strong> protocol.
                            Your base price is $50. Every day you fast, verify ketosis, and engage with your tribe,
                            your price drops. If you are disciplined, you pay $1. If you are lazy, you pay.
                        </p>
                    </div>

                    <div className="flex items-center space-x-2">
                        <Checkbox id="terms" checked={agreed} onCheckedChange={(c) => setAgreed(c as boolean)} />
                        <Label htmlFor="terms" className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                            I accept the financial consequences of my own actions.
                        </Label>
                    </div>

                    <Button className="w-full text-lg py-6" onClick={handleSign} disabled={!agreed}>
                        Sign Contract
                    </Button>
                </CardContent>
            </Card>
        </div>
    );
};

export default Contract;
