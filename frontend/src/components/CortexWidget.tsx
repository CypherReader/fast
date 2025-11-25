import React, { useState, useRef, useEffect } from 'react';
import { cortexApi } from '../api/client';
import { Button } from './ui/button';
import { X, Send, Brain } from 'lucide-react';

const CortexWidget: React.FC = () => {
    const [isOpen, setIsOpen] = useState(false);
    const [messages, setMessages] = useState<{ role: 'user' | 'assistant'; content: string }[]>([
        { role: 'assistant', content: "I am Cortex. I know your discipline index. Don't disappoint me." }
    ]);
    const [input, setInput] = useState('');
    const [loading, setLoading] = useState(false);
    const messagesEndRef = useRef<HTMLDivElement>(null);

    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
    };

    useEffect(() => {
        scrollToBottom();
    }, [messages, isOpen]);

    const handleSend = async () => {
        if (!input.trim() || loading) return;

        const userMsg = input;
        setInput('');
        setMessages(prev => [...prev, { role: 'user', content: userMsg }]);
        setLoading(true);

        try {
            const res = await cortexApi.chat(userMsg);
            setMessages(prev => [...prev, { role: 'assistant', content: res.data.response }]);
        } catch (error: any) {
            console.error("Cortex API Error:", error);
            const errorMessage = error.response?.data?.error || error.message || "Unknown error";
            setMessages(prev => [...prev, { role: 'assistant', content: `Error: ${errorMessage}` }]);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="fixed bottom-6 right-6 z-50 flex flex-col items-end">
            {isOpen && (
                <div className="bg-gray-900 border border-cyan-500/50 rounded-2xl shadow-[0_0_30px_rgba(8,145,178,0.3)] w-80 h-96 mb-4 flex flex-col overflow-hidden animate-in slide-in-from-bottom-10 fade-in duration-300">
                    {/* Header */}
                    <div className="bg-gray-800 p-4 flex justify-between items-center border-b border-gray-700">
                        <div className="flex items-center gap-2">
                            <Brain className="w-5 h-5 text-cyan-400" />
                            <span className="font-bold text-cyan-400">Cortex AI</span>
                        </div>
                        <button onClick={() => setIsOpen(false)} className="text-gray-400 hover:text-white">
                            <X className="w-5 h-5" />
                        </button>
                    </div>

                    {/* Messages */}
                    <div className="flex-1 overflow-y-auto p-4 space-y-4">
                        {messages.map((msg, idx) => (
                            <div key={idx} className={`flex ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}>
                                <div className={`max-w-[80%] p-3 rounded-xl text-sm ${msg.role === 'user'
                                    ? 'bg-cyan-600 text-white rounded-tr-none'
                                    : 'bg-gray-800 text-gray-200 rounded-tl-none border border-gray-700'
                                    }`}>
                                    {msg.content}
                                </div>
                            </div>
                        ))}
                        {loading && (
                            <div className="flex justify-start">
                                <div className="bg-gray-800 text-cyan-400 text-xs p-2 rounded-lg animate-pulse">
                                    Analyzing neural patterns...
                                </div>
                            </div>
                        )}
                        <div ref={messagesEndRef} />
                    </div>

                    {/* Input */}
                    <div className="p-3 bg-gray-800 border-t border-gray-700 flex gap-2">
                        <input
                            type="text"
                            value={input}
                            onChange={(e) => setInput(e.target.value)}
                            onKeyDown={(e) => e.key === 'Enter' && handleSend()}
                            placeholder="Type a message..."
                            className="flex-1 bg-gray-900 border border-gray-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-cyan-500 transition-colors"
                        />
                        <Button onClick={handleSend} disabled={loading} size="icon" className="bg-cyan-600 hover:bg-cyan-700">
                            <Send className="w-4 h-4" />
                        </Button>
                    </div>
                </div>
            )}

            <button
                onClick={() => setIsOpen(!isOpen)}
                className={`bg-cyan-600 hover:bg-cyan-500 text-white p-4 rounded-full shadow-[0_0_20px_rgba(8,145,178,0.6)] transition-all hover:scale-110 ${isOpen ? 'rotate-180 bg-gray-700 hover:bg-gray-600' : ''}`}
            >
                {isOpen ? <X className="w-6 h-6" /> : <span className="text-2xl">🧠</span>}
            </button>
        </div>
    );
};

export default CortexWidget;
