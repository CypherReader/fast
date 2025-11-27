import { useScroll, useTransform, MotionValue } from "framer-motion";
import { useRef } from "react";

interface ParallaxOptions {
  speed?: number; // Positive = slower, Negative = faster than scroll
  offset?: ["start start" | "start end" | "end start" | "end end", "start start" | "start end" | "end start" | "end end"];
}

export const useParallax = (options: ParallaxOptions = {}) => {
  const { speed = 0.5, offset = ["start end", "end start"] } = options;
  const ref = useRef<HTMLDivElement>(null);
  
  const { scrollYProgress } = useScroll({
    target: ref,
    offset,
  });

  // Transform scroll progress to Y movement
  const y = useTransform(scrollYProgress, [0, 1], [0, speed * 200]);
  
  return { ref, y, scrollYProgress };
};

export const useMultiParallax = () => {
  const { scrollY } = useScroll();
  
  // Different speeds for different layers
  const slowY = useTransform(scrollY, [0, 3000], [0, -150]);
  const mediumY = useTransform(scrollY, [0, 3000], [0, -300]);
  const fastY = useTransform(scrollY, [0, 3000], [0, -500]);
  
  // Horizontal parallax
  const slowX = useTransform(scrollY, [0, 3000], [0, -100]);
  const fastX = useTransform(scrollY, [0, 3000], [0, 100]);
  
  // Scale effects
  const scaleUp = useTransform(scrollY, [0, 1000], [1, 1.2]);
  const scaleDown = useTransform(scrollY, [0, 1000], [1, 0.9]);
  
  // Rotation
  const rotate = useTransform(scrollY, [0, 3000], [0, 360]);
  const rotateReverse = useTransform(scrollY, [0, 3000], [0, -360]);
  
  // Opacity based on scroll
  const fadeOut = useTransform(scrollY, [0, 500], [1, 0]);
  const fadeIn = useTransform(scrollY, [0, 500], [0, 1]);

  return {
    slowY,
    mediumY,
    fastY,
    slowX,
    fastX,
    scaleUp,
    scaleDown,
    rotate,
    rotateReverse,
    fadeOut,
    fadeIn,
  };
};
