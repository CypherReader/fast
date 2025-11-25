import React from 'react';
import { useNavigate } from 'react-router-dom';
import './VaultIntro.css';

const VaultIntro = () => {
    const navigate = useNavigate();

    const handleStart = (e: React.MouseEvent) => {
        e.preventDefault();
        // Here you would typically trigger the subscription flow
        // For now, we'll navigate to the dashboard or a checkout page
        // Assuming we want to go back to the app after "starting"
        navigate('/');
    };

    return (
        <div className="vault-intro-wrapper">
            {/* Hero Section */}
            <section className="hero">
                <div className="container">
                    <h1>Put Your Money Where Your Goals Are</h1>
                    <p>The accountability app that pays you back for discipline</p>
                    <a href="#" onClick={handleStart} className="cta-button">Start Your Vault - It's Free to Try</a>
                    <div className="trust-badge">✓ 30-day money-back guarantee · ✓ Cancel anytime · ✓ No hidden fees</div>
                </div>
            </section>

            {/* How It Works */}
            <section className="how-it-works">
                <div className="container">
                    <h2 className="section-title">Here's How It Works</h2>
                    <div className="steps">
                        <div className="step">
                            <div className="step-number">1</div>
                            <h3>Commit $30/Month</h3>
                            <p>$10 keeps the app running. $20 goes into your personal Vault—money you can earn back.</p>
                        </div>
                        <div className="step">
                            <div className="step-number">2</div>
                            <h3>Take Action Daily</h3>
                            <p>Log your meals, hit your step goals, build streaks. Every action earns money from your Vault.</p>
                        </div>
                        <div className="step">
                            <div className="step-number">3</div>
                            <h3>Get Paid Back</h3>
                            <p>At month's end, we refund what you earned. Disciplined users pay as little as $10/month.</p>
                        </div>
                    </div>
                </div>
            </section>

            {/* Vault Visual */}
            <section className="vault-visual">
                <div className="container">
                    <div className="vault-box">
                        <h2 style={{ fontSize: '2.2em', marginBottom: '30px' }}>Your Monthly Breakdown</h2>
                        <div className="vault-breakdown">
                            <div className="vault-item">
                                <span className="vault-amount">$30</span>
                                <span className="vault-label">You Deposit</span>
                            </div>
                            <div className="vault-item">
                                <span className="vault-amount">$10</span>
                                <span className="vault-label">Base Fee (App Access)</span>
                            </div>
                            <div className="vault-item">
                                <span className="vault-amount">$20</span>
                                <span className="vault-label">Your Vault (Earn It Back)</span>
                            </div>
                        </div>
                        <p style={{ fontSize: '1.2em', marginTop: '30px', opacity: 0.95 }}>
                            <strong>Most users earn back $12-18 per month.</strong> That's a net cost of $12-18 for premium accountability.
                        </p>
                    </div>
                </div>
            </section>

            {/* Earning Examples */}
            <section className="earning-examples">
                <div className="container">
                    <h2 className="section-title">Real Examples: What You Could Earn</h2>
                    <div className="example-cards">
                        <div className="example-card">
                            <div className="persona">🔥 The Committed One</div>
                            <p style={{ color: '#718096', marginBottom: '20px' }}>Sarah logs every meal and walks 12k steps daily</p>
                            <ul className="activity-list">
                                <li><span>21 meals logged (3/day × 7 days)</span> <span className="earned">+$10.50</span></li>
                                <li><span>7 days of 10k+ steps</span> <span className="earned">+$3.50</span></li>
                                <li><span>Perfect week streak bonus</span> <span className="earned">+$5.00</span></li>
                                <li><span><strong>Monthly Refund</strong></span> <span className="earned">$19.00</span></li>
                            </ul>
                            <p style={{ textAlign: 'center', marginTop: '20px', fontSize: '1.3em', color: '#2d3748' }}>
                                <strong>Net Cost: $11/month</strong>
                            </p>
                        </div>

                        <div className="example-card">
                            <div className="persona">💪 The Improver</div>
                            <p style={{ color: '#718096', marginBottom: '20px' }}>Mike logs 2 meals/day and hits 8k steps on weekdays</p>
                            <ul className="activity-list">
                                <li><span>14 meals logged (2/day × 7 days)</span> <span className="earned">+$7.00</span></li>
                                <li><span>5 days of 5-10k steps</span> <span className="earned">+$1.00</span></li>
                                <li><span>2 days under 5k steps</span> <span className="earned">+$0.00</span></li>
                                <li><span><strong>Monthly Refund</strong></span> <span className="earned">$8.00</span></li>
                            </ul>
                            <p style={{ textAlign: 'center', marginTop: '20px', fontSize: '1.3em', color: '#2d3748' }}>
                                <strong>Net Cost: $22/month</strong>
                            </p>
                        </div>

                        <div className="example-card">
                            <div className="persona">😅 The Struggler</div>
                            <p style={{ color: '#718096', marginBottom: '20px' }}>Alex logs inconsistently and walks 3k steps most days</p>
                            <ul className="activity-list">
                                <li><span>8 meals logged (sporadic)</span> <span className="earned">+$4.00</span></li>
                                <li><span>Most days under 5k steps</span> <span className="earned">+$0.00</span></li>
                                <li><span>No streak bonuses</span> <span className="earned">+$0.00</span></li>
                                <li><span><strong>Monthly Refund</strong></span> <span className="earned">$4.00</span></li>
                            </ul>
                            <p style={{ textAlign: 'center', marginTop: '20px', fontSize: '1.3em', color: '#2d3748' }}>
                                <strong>Net Cost: $26/month</strong>
                            </p>
                        </div>
                    </div>
                </div>
            </section>

            {/* Social Proof */}
            <section className="social-proof">
                <div className="container">
                    <h2 className="section-title">What Our Users Say</h2>
                    <div className="testimonials">
                        <div className="testimonial">
                            <p className="quote">"I've tried every app and nothing stuck. Seeing real money on the line? Game changer. I'm down 18 lbs in 3 months."</p>
                            <p className="author">— Jennifer M., Beta User</p>
                        </div>
                        <div className="testimonial">
                            <p className="quote">"Finally, an app that doesn't let me lie to myself. I earned back $17 last month and felt amazing about it."</p>
                            <p className="author">— Carlos R., 6 Months In</p>
                        </div>
                        <div className="testimonial">
                            <p className="quote">"The streak bonuses are addictive in the best way. I haven't missed a meal log in 8 weeks."</p>
                            <p className="author">— Amanda K., Perfect Streaker</p>
                        </div>
                    </div>
                    <div className="stats">
                        <div className="stat">
                            <span className="stat-number">87%</span>
                            <span className="stat-label">Hit Their Goals Monthly</span>
                        </div>
                        <div className="stat">
                            <span className="stat-number">$16</span>
                            <span className="stat-label">Avg. Monthly Refund</span>
                        </div>
                        <div className="stat">
                            <span className="stat-number">4.8/5</span>
                            <span class="stat-label">User Rating</span>
                        </div>
                    </div>
                </div>
            </section>

            {/* FAQ */}
            <section className="faq">
                <div className="container">
                    <h2 className="section-title">Your Questions, Answered</h2>
                    <div className="faq-list">
                        <div className="faq-item">
                            <div className="faq-question">What if I don't earn anything back?</div>
                            <div className="faq-answer">Then you pay the full $30 that month. But that rarely happens—even small efforts add up. And if you're not seeing results after 30 days, we'll refund you completely.</div>
                        </div>
                        <div className="faq-item">
                            <div className="faq-question">Can I cancel anytime?</div>
                            <div className="faq-answer">Absolutely. No contracts, no penalties. Cancel with one click and you won't be charged next month.</div>
                        </div>
                        <div className="faq-item">
                            <div className="faq-question">How do you verify my meals and steps?</div>
                            <div className="faq-answer">Meals require a photo (we trust you, but accountability matters). Steps sync automatically from your phone's health app. We're working on AI verification for even better accuracy.</div>
                        </div>
                        <div className="faq-item">
                            <div className="faq-question">What happens to money I don't earn back?</div>
                            <div className="faq-answer">It funds the app's development, customer support, and new features. We're transparent: we want you to earn it back, but the stakes need to be real.</div>
                        </div>
                        <div className="faq-item">
                            <div className="faq-question">Is this really different from other apps?</div>
                            <div className="faq-answer">Yes. Other apps rely on motivation alone. We add financial accountability—real money, real stakes. Studies show financial commitment increases goal completion by 3-4x.</div>
                        </div>
                    </div>
                </div>
            </section>

            {/* Final CTA */}
            <section className="final-cta">
                <div className="container">
                    <h2>Ready to Stop Making Excuses?</h2>
                    <p>Join hundreds of people who've turned their goals into commitments—and earned money doing it.</p>
                    <a href="#" onClick={handleStart} className="cta-button">Start Your Vault Today</a>
                    <div className="guarantee">
                        🛡️ <strong>Risk-Free Guarantee:</strong> Try it for 30 days. If you don't see progress, we'll refund your entire deposit. No questions asked.
                    </div>
                </div>
            </section>
        </div>
    );
};

export default VaultIntro;
