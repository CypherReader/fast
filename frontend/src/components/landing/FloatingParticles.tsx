import { useEffect, useState } from "react";
import { motion, AnimatePresence } from "framer-motion";

interface Particle {
  id: number;
  x: number;
  size: number;
  opacity: number;
  rotation: number;
  duration: number;
}

export const FloatingParticles = () => {
  const [particles, setParticles] = useState<Particle[]>([]);

  useEffect(() => {
    // Spawn initial particles
    const initialParticles: Particle[] = Array.from({ length: 8 }, (_, i) => ({
      id: i,
      x: Math.random() * 100,
      size: 24 + Math.random() * 16,
      opacity: 0.1 + Math.random() * 0.15,
      rotation: Math.random() * 30 - 15,
      duration: 15 + Math.random() * 10,
    }));
    setParticles(initialParticles);

    // Spawn new particles periodically
    const interval = setInterval(() => {
      const newParticle: Particle = {
        id: Date.now(),
        x: Math.random() * 100,
        size: 24 + Math.random() * 16,
        opacity: 0.1 + Math.random() * 0.15,
        rotation: Math.random() * 30 - 15,
        duration: 15 + Math.random() * 10,
      };

      setParticles((prev) => {
        // Keep max 15 particles
        const filtered = prev.length >= 15 ? prev.slice(1) : prev;
        return [...filtered, newParticle];
      });
    }, 2000);

    return () => clearInterval(interval);
  }, []);

  return (
    <div className="absolute inset-0 overflow-hidden pointer-events-none">
      <AnimatePresence>
        {particles.map((particle) => (
          <motion.div
            key={particle.id}
            initial={{ y: "110vh", x: `${particle.x}vw`, rotate: 0, opacity: 0 }}
            animate={{
              y: "-10vh",
              x: [
                `${particle.x}vw`,
                `${particle.x + 3}vw`,
                `${particle.x - 2}vw`,
                `${particle.x + 1}vw`,
              ],
              rotate: [0, particle.rotation, -particle.rotation, 0],
              opacity: [0, particle.opacity, particle.opacity, 0],
            }}
            exit={{ opacity: 0 }}
            transition={{
              duration: particle.duration,
              ease: "linear",
              x: {
                duration: particle.duration,
                ease: "easeInOut",
              },
            }}
            className="absolute text-primary font-bold"
            style={{ fontSize: particle.size }}
          >
            $
          </motion.div>
        ))}
      </AnimatePresence>
    </div>
  );
};
