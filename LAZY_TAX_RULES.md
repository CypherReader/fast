# Commitment Vault & Discipline Rules

The "Commitment Vault" is a deposit-based accountability system. Users commit a monthly sum upfront and earn it back through consistent discipline.

## 1. The Core Model (Tier 1)

**Structure:**

* **Monthly Charge:** $30.00 (charged upfront)
* **Base Fee:** $10.00 (Non-refundable, covers app costs)
* **Vault Deposit:** $20.00 (The "Earnable Refund Pool")

**Objective:** Earn back your $20 deposit by hitting daily goals.

## 2. Earning Rules

You can earn back money from your vault every day.

### Daily Actions

* **Food Logging:** $0.50 per meal photo (Breakfast, Lunch, Dinner).
* **Activity:** $0.50 for hitting 10,000+ steps.
* **Daily Cap:** Max earnings of **$2.00 per day**.

### Streak Bonuses

* **7-Day Perfect Streak:** +$5.00 bonus.
* **30-Day Perfect Month:** +$10.00 bonus (one-time).

### Refund Cap

* **Monthly Max Refund:** $20.00.
* *Note: You cannot earn more than your initial deposit.*

## 3. Tiers & Upgrades

### Tier 2: Accountability Plus (+$10/month)

* **Total Monthly:** $40.00
* **Features:**
  * Partner matching.
  * Shared goal tracking.
  * Group challenges (Winner takes 50% of losers' vault deposits).

### Tier 3: AI Coach (+$20/month)

* **Total Monthly:** $50.00
* **Features:**
  * AI meal analysis.
  * Personalized recommendations.
  * Weekly coaching calls.
  * Custom meal plans.

## 4. Payment & Refund Flow

1. **Start of Month:** User is charged the full monthly amount (e.g., $30).
2. **During Month:** System tracks `earned_refund` based on logs and steps.
3. **End of Month:**
    * If `earned_refund > 0`: A refund is processed to the original payment method.
    * Refund Amount = `MIN(earned_refund, vault_deposit)`.
    * Net Cost = `Monthly Charge - Refund Amount`.

## 5. Marketing Language (Do's & Don'ts)

* ❌ **Don't Say:** "Lazy Tax", "Penalties", "Fines".
* ✅ **Do Say:** "Commitment Vault", "Earn your money back", "Accountability that pays you back".
