
import { Toaster } from "@/components/ui/toaster";
import { Toaster as Sonner } from "@/components/ui/sonner";
import { TooltipProvider } from "@/components/ui/tooltip";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import Layout from "./components/Layout";
import Dashboard from "./pages/Dashboard";
import Progress from "./pages/Progress";
import Community from "./pages/Community";
import Resources from "./pages/Resources";
import Profile from "./pages/Profile";
import NotFound from "./pages/NotFound";
import VaultIntro from "./pages/VaultIntro";
import Referrals from "./pages/Referrals";

import Contract from "./pages/Contract";
import Login from "./pages/Login";
import Register from "./pages/Register";

import Tribe from "./pages/Tribe";
import { ActivityDashboard } from "./features/activity/ActivityDashboard";

import { Elements } from "@stripe/react-stripe-js";
import { loadStripe } from "@stripe/stripe-js";

const queryClient = new QueryClient();
const stripePromise = loadStripe(import.meta.env.VITE_STRIPE_PUBLIC_KEY || "pk_test_placeholder");

import { AuthProvider } from "@/context/AuthContext";

const App = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <Elements stripe={stripePromise}>
        <TooltipProvider>
          <Toaster />
          <Sonner />
          <AuthProvider>
            <BrowserRouter>
              <Routes>
                <Route path="/login" element={<Login />} />
                <Route path="/register" element={<Register />} />
                <Route path="/contract" element={<Contract />} />
                <Route path="/" element={<Layout />}>
                  <Route index element={<Dashboard />} />
                  <Route path="progress" element={<Progress />} />
                  <Route path="community" element={<Community />} />
                  <Route path="resources" element={<Resources />} />
                  <Route path="tribe" element={<Tribe />} />
                  <Route path="activity" element={<ActivityDashboard />} />
                  <Route path="profile" element={<Profile />} />
                  <Route path="vault-intro" element={<VaultIntro />} />
                  <Route path="referrals" element={<Referrals />} />
                </Route>
                <Route path="*" element={<NotFound />} />
              </Routes>
            </BrowserRouter>
          </AuthProvider>
        </TooltipProvider>
      </Elements>
    </QueryClientProvider>
  );
};

export default App;
