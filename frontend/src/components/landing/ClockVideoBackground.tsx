import { useEffect, useState } from "react";
import { motion, useScroll, useTransform } from "framer-motion";

export const ClockVideoBackground = () => {
  const [seconds, setSeconds] = useState(0);
  const [minutes, setMinutes] = useState(42);
  const [hours, setHours] = useState(16);

  // Parallax scroll effects
  const { scrollY } = useScroll();
  const clockY = useTransform(scrollY, [0, 500], [0, 150]);
  const clockScale = useTransform(scrollY, [0, 500], [1, 0.8]);
  const clockOpacity = useTransform(scrollY, [0, 400], [1, 0.3]);
  const particlesY = useTransform(scrollY, [0, 500], [0, -100]);
  const ringsY = useTransform(scrollY, [0, 500], [0, 200]);
  const sweepRotate = useTransform(scrollY, [0, 1000], [0, 180]);

  useEffect(() => {
    const interval = setInterval(() => {
      setSeconds((prev) => {
        if (prev >= 59) {
          setMinutes((m) => {
            if (m >= 59) {
              setHours((h) => (h + 1) % 24);
              return 0;
            }
            return m + 1;
          });
          return 0;
        }
        return prev + 1;
      });
    }, 1000);
    return () => clearInterval(interval);
  }, []);

  const secondRotation = seconds * 6;
  const minuteRotation = minutes * 6 + seconds * 0.1;
  const hourRotation = hours * 30 + minutes * 0.5;

  return (
    <div className="absolute inset-0 overflow-hidden">
      {/* Animated gradient background */}
      <motion.div
        className="absolute inset-0"
        animate={{
          background: [
            "radial-gradient(ellipse at 30% 20%, hsl(var(--primary) / 0.15) 0%, transparent 50%), radial-gradient(ellipse at 70% 80%, hsl(var(--secondary) / 0.1) 0%, transparent 50%)",
            "radial-gradient(ellipse at 70% 30%, hsl(var(--primary) / 0.2) 0%, transparent 50%), radial-gradient(ellipse at 30% 70%, hsl(var(--secondary) / 0.15) 0%, transparent 50%)",
            "radial-gradient(ellipse at 50% 50%, hsl(var(--primary) / 0.15) 0%, transparent 50%), radial-gradient(ellipse at 50% 50%, hsl(var(--secondary) / 0.1) 0%, transparent 50%)",
            "radial-gradient(ellipse at 30% 20%, hsl(var(--primary) / 0.15) 0%, transparent 50%), radial-gradient(ellipse at 70% 80%, hsl(var(--secondary) / 0.1) 0%, transparent 50%)",
          ],
        }}
        transition={{ duration: 10, repeat: Infinity, ease: "easeInOut" }}
      />

      {/* Giant clock face with parallax */}
      <motion.div 
        className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2"
        style={{ y: clockY, scale: clockScale, opacity: clockOpacity }}
      >
        <motion.div
          initial={{ scale: 0.8, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          transition={{ duration: 1.5, ease: "easeOut" }}
          className="relative"
        >
          {/* Outer glow ring */}
          <motion.div
            className="absolute -inset-20 rounded-full"
            animate={{
              boxShadow: [
                "0 0 100px 50px hsl(var(--primary) / 0.1)",
                "0 0 150px 80px hsl(var(--primary) / 0.2)",
                "0 0 100px 50px hsl(var(--primary) / 0.1)",
              ],
            }}
            transition={{ duration: 4, repeat: Infinity, ease: "easeInOut" }}
          />

          <svg width="700" height="700" viewBox="0 0 700 700" className="opacity-30">
            {/* Outer ring with tick marks */}
            <circle
              cx="350"
              cy="350"
              r="320"
              fill="none"
              stroke="hsl(var(--border))"
              strokeWidth="2"
              opacity="0.5"
            />

            {/* Hour markers */}
            {[...Array(12)].map((_, i) => {
              const angle = (i * 30 - 90) * (Math.PI / 180);
              const x1 = 350 + 280 * Math.cos(angle);
              const y1 = 350 + 280 * Math.sin(angle);
              const x2 = 350 + 310 * Math.cos(angle);
              const y2 = 350 + 310 * Math.sin(angle);
              return (
                <motion.line
                  key={i}
                  x1={x1}
                  y1={y1}
                  x2={x2}
                  y2={y2}
                  stroke="hsl(var(--primary))"
                  strokeWidth="4"
                  strokeLinecap="round"
                  initial={{ opacity: 0 }}
                  animate={{ opacity: [0.3, 0.8, 0.3] }}
                  transition={{ duration: 2, delay: i * 0.1, repeat: Infinity }}
                />
              );
            })}

            {/* Minute markers */}
            {[...Array(60)].map((_, i) => {
              if (i % 5 === 0) return null;
              const angle = (i * 6 - 90) * (Math.PI / 180);
              const x1 = 350 + 295 * Math.cos(angle);
              const y1 = 350 + 295 * Math.sin(angle);
              const x2 = 350 + 310 * Math.cos(angle);
              const y2 = 350 + 310 * Math.sin(angle);
              return (
                <line
                  key={i}
                  x1={x1}
                  y1={y1}
                  x2={x2}
                  y2={y2}
                  stroke="hsl(var(--muted-foreground))"
                  strokeWidth="1"
                  opacity="0.3"
                />
              );
            })}

            {/* Hour hand */}
            <motion.line
              x1="350"
              y1="350"
              x2="350"
              y2="180"
              stroke="hsl(var(--foreground))"
              strokeWidth="8"
              strokeLinecap="round"
              style={{
                transformOrigin: "350px 350px",
                transform: `rotate(${hourRotation}deg)`,
              }}
              opacity="0.6"
            />

            {/* Minute hand */}
            <motion.line
              x1="350"
              y1="350"
              x2="350"
              y2="120"
              stroke="hsl(var(--foreground))"
              strokeWidth="5"
              strokeLinecap="round"
              style={{
                transformOrigin: "350px 350px",
                transform: `rotate(${minuteRotation}deg)`,
              }}
              opacity="0.5"
            />

            {/* Second hand */}
            <motion.line
              x1="350"
              y1="350"
              x2="350"
              y2="80"
              stroke="hsl(var(--primary))"
              strokeWidth="2"
              strokeLinecap="round"
              style={{
                transformOrigin: "350px 350px",
                transform: `rotate(${secondRotation}deg)`,
              }}
              className="drop-shadow-[0_0_10px_hsl(var(--primary))]"
            />

            {/* Center dot */}
            <circle cx="350" cy="350" r="12" fill="hsl(var(--primary))" />
            <circle cx="350" cy="350" r="6" fill="hsl(var(--background))" />

            {/* Gradient definitions */}
            <defs>
              <radialGradient id="clockGlow" cx="50%" cy="50%" r="50%">
                <stop offset="0%" stopColor="hsl(var(--primary))" stopOpacity="0.3" />
                <stop offset="100%" stopColor="hsl(var(--primary))" stopOpacity="0" />
              </radialGradient>
            </defs>
          </svg>

          {/* Digital time overlay */}
          <div className="absolute top-1/2 left-1/2 -translate-x-1/2 translate-y-16 text-center">
            <motion.div
              className="font-mono text-6xl font-bold text-foreground/20 tabular-nums tracking-wider"
              animate={{ opacity: [0.15, 0.25, 0.15] }}
              transition={{ duration: 3, repeat: Infinity }}
            >
              {String(hours).padStart(2, "0")}:{String(minutes).padStart(2, "0")}:
              {String(seconds).padStart(2, "0")}
            </motion.div>
          </div>
        </motion.div>
      </motion.div>

      {/* Sweeping light effect with parallax */}
      <motion.div
        className="absolute top-1/2 left-1/2 w-[800px] h-[800px] -translate-x-1/2 -translate-y-1/2"
        style={{
          background:
            "conic-gradient(from 0deg, transparent 0deg, hsl(var(--primary) / 0.1) 30deg, transparent 60deg)",
          rotate: sweepRotate,
        }}
        animate={{ rotate: 360 }}
        transition={{ duration: 60, repeat: Infinity, ease: "linear" }}
      />

      {/* Floating time particles with parallax */}
      <motion.div style={{ y: particlesY }}>
        {[...Array(15)].map((_, i) => (
          <motion.div
            key={i}
            className="absolute text-primary/20 font-mono text-lg font-bold"
            initial={{
              x: Math.random() * (typeof window !== "undefined" ? window.innerWidth : 1200),
              y: Math.random() * (typeof window !== "undefined" ? window.innerHeight : 800),
            }}
            animate={{
              y: [null, -100],
              opacity: [0, 0.4, 0],
            }}
            transition={{
              duration: 8 + Math.random() * 4,
              repeat: Infinity,
              delay: Math.random() * 5,
            }}
          >
            {["16h", "18h", "24h", "12h", "20h", "⏰", "🔥", "⚡"][i % 8]}
          </motion.div>
        ))}
      </motion.div>

      {/* Pulse rings with parallax */}
      <motion.div style={{ y: ringsY }}>
        {[...Array(3)].map((_, i) => (
          <motion.div
            key={i}
            className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 rounded-full border border-primary/20"
            initial={{ width: 100, height: 100, opacity: 0.5 }}
            animate={{
              width: [100, 800],
              height: [100, 800],
              opacity: [0.3, 0],
            }}
            transition={{
              duration: 4,
              repeat: Infinity,
              delay: i * 1.3,
              ease: "easeOut",
            }}
          />
        ))}
      </motion.div>
    </div>
  );
};
