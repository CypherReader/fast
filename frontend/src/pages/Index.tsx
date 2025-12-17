import { Navbar } from "@/components/landing/Navbar";
import { HeroSection } from "@/components/landing/HeroSection";
import { LiveActivityTicker } from "@/components/landing/LiveActivityTicker";
import { ProblemSection } from "@/components/landing/ProblemSection";
import { HowItWorksSection } from "@/components/landing/HowItWorksSection";
import { TribesSection } from "@/components/landing/TribesSection";
import { TransformationShowcase } from "@/components/landing/TransformationShowcase";
import { ResourcesSection } from "@/components/landing/ResourcesSection";
import { TestimonialsSection } from "@/components/landing/TestimonialsSection";
import { PricingSection } from "@/components/landing/PricingSection";
import { FAQSection } from "@/components/landing/FAQSection";
import { FinalCTASection } from "@/components/landing/FinalCTASection";
import { Footer } from "@/components/landing/Footer";
import { ScrollProgress } from "@/components/landing/ScrollProgress";

const Index = () => {
  return (
    <div className="min-h-screen bg-background">
      <ScrollProgress />
      <Navbar />
      <main>
        <HeroSection />
        <LiveActivityTicker />
        <ProblemSection />
        <HowItWorksSection />
        <TribesSection />
        <TransformationShowcase />
        <ResourcesSection />
        <TestimonialsSection />
        <PricingSection />
        <FAQSection />
        <FinalCTASection />
      </main>
      <Footer />
    </div>
  );
};

export default Index;
