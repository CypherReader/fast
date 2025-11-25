import { useEffect, useState } from "react";

interface SplashScreenProps {
  onComplete: () => void;
}

const SplashScreen = ({ onComplete }: SplashScreenProps) => {
  const [fadeOut, setFadeOut] = useState(false);

  useEffect(() => {
    const video = document.getElementById("splash-video") as HTMLVideoElement;
    
    const handleVideoEnd = () => {
      setFadeOut(true);
      setTimeout(onComplete, 1000);
    };

    if (video) {
      video.play().catch(() => {
        // Fallback if autoplay fails
        setTimeout(() => {
          setFadeOut(true);
          setTimeout(onComplete, 1000);
        }, 3000);
      });
      video.addEventListener("ended", handleVideoEnd);
    }

    return () => {
      if (video) {
        video.removeEventListener("ended", handleVideoEnd);
      }
    };
  }, [onComplete]);

  return (
    <div
      className={`fixed inset-0 z-50 bg-background flex items-center justify-center transition-opacity duration-1000 ${
        fadeOut ? "opacity-0" : "opacity-100"
      }`}
    >
      <video
        id="splash-video"
        className="w-full h-full object-cover"
        muted
        playsInline
        preload="auto"
      >
        <source src="/splash-video.mp4" type="video/mp4" />
      </video>
      
      {/* Gradient overlay for cinematic effect */}
      <div className="absolute inset-0 bg-gradient-to-b from-transparent via-transparent to-background/80" />
      
      {/* App title fade in towards the end */}
      <div className="absolute inset-0 flex items-center justify-center">
        <h1 className="text-6xl font-bold text-shimmer animate-fade-in opacity-0 animation-delay-2000">
          Autophagy Arc
        </h1>
      </div>
    </div>
  );
};

export default SplashScreen;
