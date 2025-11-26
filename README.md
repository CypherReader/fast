# FastingHero - UI/UX Design Documentation

> "Put Your Money Where Your Goals Are."

FastingHero is a high-stakes accountability application that gamifies discipline through a unique "Commitment Vault" model. Users deposit funds upfront and earn them back through disciplined actions like fasting, meal logging, and activity tracking.

---

## 🎨 Design System

### Color Palette

The application uses a **dark theme** with vibrant accent colors for a premium, modern feel:

#### Primary Colors

- **Primary Green**: `#10b981` (Emerald-500) - Success states, CTAs, fasting indicators
- **Secondary Purple**: `#a855f7` (Purple-500) - Premium features, special highlights
- **Background**: `#000000` to `#0f172a` (Black to Slate-900) - Dark, immersive base
- **Card Background**: `#1e293b` (Slate-800) - Elevated surfaces

#### Accent Colors

- **Emerald**: `#10b981` - Food logging, earnings, positive actions
- **Cyan**: `#06b6d4` - Telemetry data, metrics
- **Yellow/Gold**: `#fbbf24` - Badges, achievements, rewards
- **Red**: `#ef4444` - Warnings, penalties, discipline failures
- **Chart Green**: `#adfa1d` - Data visualization

#### Text Colors

- **Primary Text**: `#ffffff` - Headings, important content
- **Secondary Text**: `#94a3b8` (Slate-400) - Body text, labels
- **Muted Text**: `#64748b` (Slate-500) - Hints, placeholders

### Typography

**Font Family:**

- **Primary**: `Inter` - Clean, modern sans-serif for UI
- **Accent**: `Crimson Pro` - Serif for emphasis
- **Mono**: `Roboto Mono` - Code and data display

**Font Sizes (Tailwind):**

- `text-xs`: 12px - Small labels, captions
- `text-sm`: 14px - Body text, descriptions
- `text-base`: 16px - Standard text
- `text-lg`: 18px - Section headings  
- `text-xl`: 20px - Card titles
- `text-2xl`: 24px - Page headings
- `text-3xl`: 30px - Hero text
- `text-4xl`: 36px - Large metrics

### Spacing & Layout

**Container Widths:**

- Maximum content width: Responsive (no fixed max)
- Card padding: `p-4` to `p-6` (16px to 24px)
- Section spacing: `space-y-6` (24px vertical gap)

**Grid System:**

- Mobile: Single column
- Tablet: `md:grid-cols-2`
- Desktop: `lg:grid-cols-3` or `lg:grid-cols-4`

### Component Library

Built with **shadcn/ui** - A collection of re-usable components:

- **Card**: Main container component for sections
- **Button**: Multiple variants (default, outline, ghost, destructive)
- **Dialog/Modal**: Overlay content
- **Progress**: Linear progress bars
- **Badge**: Status indicators, tags
- **Tabs**: Navigation within pages
- **Input**: Form fields

---

## 📱 Application Structure

### Navigation System

**Bottom Navigation Bar** (Mobile-first, persistent):

- 🏠 Dashboard
- 📈 Progress
- 📊 Activity  
- 👥 Community
- 📖 Resources
- 👤 Profile

**Style:**

- Fixed bottom position
- Active state: Primary color with glow effect
- Inactive: Muted gray
- Icons from `lucide-react`

### Page Layout Pattern

All pages follow this structure:

```tsx
<div className="space-y-6">
  {/* Header */}
  <div className="animate-fade-in">
    <h1 className="text-2xl font-bold bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
      Page Title
    </h1>
    <p className="text-sm text-muted-foreground">Subtitle</p>
  </div>

  {/* Content Cards */}
  <Card className="border-primary/20 animate-fade-in-up">
    {/* Card content */}
  </Card>
</div>
```

---

## 🎬 Animations & Interactions

### Animation Patterns

**Entry Animations:**

- `animate-fade-in`: Fade in on page load
- `animate-fade-in-up`: Fade in with upward slide
- `animationDelay`: Staggered delays (0.1s, 0.2s, 0.3s, etc.)

**Hover Effects:**

- **Cards**: `hover:border-primary/40 hover:shadow-lg hover:scale-[1.02]`
- **Buttons**: `hover:scale-110 transition-transform`
- **Interactive elements**: Glow effect with `shadow-primary/10`

**Transitions:**

- Default: `transition-all duration-300`
- Fast: `duration-200`
- Smooth: `ease-in-out`

### Interactive States

**Button States:**

- Default: Solid color with border
- Hover: Slightly brighter, scale up
- Active/Pressed: Slightly darker
- Disabled: 50% opacity, no pointer events

**Form Inputs:**

- Default: Slate-900 background, slate-700 border
- Focus: Primary color border, no outline
- Error: Red border
- Disabled: Muted background

---

## 📄 Page-by-Page Breakdown

### 1. Dashboard (Home)

**Purpose:** Main fasting control center

**Sections:**

1. **Hero Section**
   - Current fasting status
   - Large timer display
   - Start/Stop buttons

2. **Quick Stats Grid** (4 columns on desktop)
   - Discipline score (0-100)
   - Current price/lazy tax
   - Fast duration
   - Vault balance

3. **Bio-Narrative Timeline**
   - Visual phases: Digestion → Glycogen Depletion → Ketosis → Autophagy
   - Progress indicator
   - Current phase description

4. **Focus Tab**
   - Fasting start/stop controls
   - Real-time insights

**Color Scheme:**

- Primary: Emerald green for active fasting
- Accent: Yellow for achievements
- Danger: Red for penalties

---

### 2. Progress Page

**Purpose:** Track health metrics and earn vault rewards

**Sections:** (Vertical stack, staggered animations)

1. **Weight Tracker** (animationDelay: 0.1s)
   - Current weight display
   - Week trend (-X lbs/kg)
   - Unit toggle (kg/lbs) in modal
   - Historical data list

2. **Hydration Tracker** (animationDelay: 0.2s)
   - 8 interactive water glasses
   - Visual fill states
   - Progress counter

3. **Food Logging** (animationDelay: 0.3s)
   - Meal counter (X / 3)
   - Earnings display (+$0.50 per meal)
   - Camera upload button
   - Recent meals carousel with images
   - Keto/Fake badges on meals

4. **Ketosis Level** (animationDelay: 0.4s)
   - Score display (0-100)
   - Progress bar
   - Premium upsell (Lock icon)

**Vault Earnings:**

- Emerald-500 color for earnings
- Real-time updates
- Toast notifications on success

---

### 3. Activity Page

**Purpose:** Steps, weight, and telemetry tracking

**Layout:**

1. **Telemetry Uplink** (Full width)
   - Tab selector: Steps / Weight
   - Manual input fields
   - Submit button (Cyan accent)
   - Trust score display

2. **Quick Stats** (4-column grid)
   - Total steps (with thousands separator)
   - Current weight (lbs + kg conversion)
   - Distance (km)
   - Calories burned

3. **Charts** (2-column on desktop)
   - Weekly steps bar chart (Lime-green bars)
   - Map placeholder

**Data Visualization:**

- Chart colors: `#adfa1d` (lime green)
- Grid: Subtle gray (`#888888`)
- Tooltips: Dark card with stats

---

### 4. Community Page

**Purpose:** Social features and engagement

**Tabs Layout:**

- Grid: `grid-cols-4` (Feed, Tribes, Leaderboard, My Progress)

#### Tab: Feed

- Activity cards with avatars
- User initials in circular badges
- Timestamp badges
- Like buttons (heart icon)
- Staggered card animations

#### Tab: Tribes

- Tribe cards with member counts
- Join/Leave buttons
- Admin indicators

#### Tab: Leaderboard

- Ranked list (1st, 2nd, 3rd with emoji medals)
- User highlight (primary background)
- Hours fasted metric

#### Tab: My Progress

- Current streak display
- Badge showcase (grid of earned badges)
- Achievement unlocks

---

### 5. Resources Page

**Purpose:** Recipes and knowledge hub

**Tabs:**

- Recipes / Knowledge

**Recipes Tab:**

- Diet filter badges (All, Vegan, Vegetarian, Normal)
- Recipe cards with images
- Badge indicators (Simple, diet icons)
- Nutritional info (calories, carbs)

**Knowledge Tab:**

- KnowledgeHub component
- FAQ-style content

---

### 6. Profile Page

**Purpose:** User settings and account management

**Sections:**

- User info display
- Referral code copy
- Vault balance
- Settings toggles

---

## 🎯 Design Patterns & Best Practices

### Card Design

**Standard Card:**

```tsx
<Card className="border-primary/20 hover:border-primary/40 transition-all">
  <CardHeader>
    <CardTitle className="flex items-center gap-2">
      <Icon className="h-5 w-5 text-primary" />
      Title
    </CardTitle>
  </CardHeader>
  <CardContent>
    {/* Content */}
  </CardContent>
</Card>
```

**Premium Card** (with lock):

```tsx
<Card className="border-secondary/20 relative overflow-hidden">
  <div className="absolute top-2 right-2">
    <Lock className="h-4 w-4 text-secondary" />
  </div>
  {/* Card content */}
</Card>
```

### Button Hierarchy

1. **Primary Action**: `bg-primary` (Emerald)
2. **Secondary Action**: `variant="outline"`
3. **Tertiary**: `variant="ghost"`
4. **Destructive**: `variant="destructive"` (Red)
5. **Premium**: Gradient `bg-gradient-to-r from-primary to-secondary`

### Modal/Dialog Pattern

- Dark background overlay (`DialogOverlay`)
- Centered content card
- Close button (X icon)
- Action buttons at bottom
- Descriptive title and subtitle

### Form Design

- Labels above inputs
- Placeholder text with `text-muted-foreground`
- Full-width inputs on mobile
- Submit button spans full width
- Validation errors in red below fields

---

## 📐 Responsive Design

### Breakpoints (Tailwind)

- **sm**: 640px - Small tablets
- **md**: 768px - Tablets
- **lg**: 1024px - Laptops
- **xl**: 1280px - Desktops
- **2xl**: 1536px - Large screens

### Mobile-First Approach

All designs start mobile (single column) and expand:

```tsx
<div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
```

### Touch Targets

- Minimum button height: `h-10` (40px)
- Minimum touch area: 44x44px
- Spacing between interactive elements: `gap-2` minimum

---

## 🔔 Notifications & Feedback

### Toast Notifications

- Success: Green background
- Error: Red background
- Info: Blue background
- Position: Top-right corner
- Auto-dismiss: 3-5 seconds

### Loading States

- Spinner components
- Skeleton placeholders
- Disabled buttons with reduced opacity
- Loading text: "Loading..."

### Empty States

- Centered content
- Icon illustration
- Descriptive text
- Call-to-action button

---

## 🎨 Special UI Components

### Fasting Timer

- Large circular display
- Real-time countdown
- Phase indicator
- Pulsing animation when active

### Progress Bars

- Standard height: `h-2`
- Rounded: `rounded-full`
- Custom colors per context
- Animated fill

### Badges

- Rounded: `rounded` or `rounded-full`
- Compact padding: `px-2 py-1`
- Small text: `text-xs`
- Icon + text combination

### Avatar System

- Circular: `rounded-full`
- Fallback: User initials
- Background: Primary color at 20% opacity

---

## 🚀 PWA Features

### Installability

- Manifest: `/manifest.json`
- App name: "FastingHero"
- Short name: "Fasting"
- Theme color: `#10b981`
- Icons: 192x192 and 512x512

### Offline Support

- Service worker caching
- Font caching
- Asset precaching

---

## 🛠️ Tech Stack (Frontend)

- **Framework**: React 18+ with TypeScript
- **Build Tool**: Vite
- **Styling**: Tailwind CSS 3+
- **Components**: shadcn/ui (Radix UI primitives)
- **Icons**: Lucide React
- **Charts**: Recharts
- **State**: React Hooks (useState, useEffect)
- **Routing**: React Router v6
- **HTTP**: Axios

---

## 📦 Component File Structure

```txt
frontend/src/
├── components/
│   ├── ui/              # shadcn/ui components
│   ├── bio/             # Telemetry components
│   ├── community/       # KnowledgeHub, etc.
│   └── Layout.tsx       # Main layout with nav
├── pages/
│   ├── Dashboard.tsx
│   ├── Progress.tsx
│   ├── Community.tsx
│   ├── Resources.tsx
│   └── Profile.tsx
├── context/
│   └── AuthContext.tsx  # Auth state
├── hooks/
│   └── useNotifications.ts
└── api/
    └── client.ts        # API calls
```

---

## 🎭 Design Philosophy

### Visual Hierarchy

1. **Most Important**: Large, bold, primary color
2. **Secondary**: Medium size, standard weight
3. **Tertiary**: Small, muted color

### Information Density

- **Desktop**: Rich, multi-column layouts
- **Mobile**: Simplified, single-column, progressive disclosure

### Feedback Loops

- Every action gets immediate visual feedback
- Animations confirm state changes
- Success states use green, errors use red

### Gamification

- Visual rewards (badges, streaks)
- Progress bars everywhere
- Positive reinforcement with earnings display
- Loss aversion with vault balance

---

## 🎨 Accessibility

- **Contrast**: WCAG AA compliant
- **Focus States**: Visible keyboard focus
- **ARIA Labels**: On interactive elements
- **Semantic HTML**: Proper heading hierarchy
- **Alt Text**: On all images

---

## 📖 Additional Resources

- **Component Reference**: `frontend/src/components/ui/`
- **Implementation Details**: See `walkthrough.md`
- **Configuration**: `.env.example`
- **Design Tokens**: `frontend/tailwind.config.ts`

---

**Last Updated:** Phase 4 Complete - PWA, Logging, Security
