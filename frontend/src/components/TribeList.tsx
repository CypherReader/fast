import React, { useEffect, useState } from 'react';
import { api } from '../services/api';

interface Tribe {
    id: string;
    name: string;
    description: string;
    member_count: number;
    total_discipline: number;
}

const TribeList: React.FC = () => {
    const [tribes, setTribes] = useState<Tribe[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        fetchTribes();
    }, []);

    const fetchTribes = async () => {
        try {
            const response = await api.get('/tribes/');
            setTribes(response.data);
        } catch (error) {
            console.error('Failed to fetch tribes', error);
        } finally {
            setLoading(false);
        }
    };

    const handleJoin = async (tribeId: string) => {
        try {
            await api.post(`/tribes/${tribeId}/join`);
            alert('Joined tribe successfully!');
            fetchTribes(); // Refresh list
        } catch (error) {
            console.error('Failed to join tribe', error);
            alert('Failed to join tribe');
        }
    };

    if (loading) return <div className="text-white">Loading tribes...</div>;

    return (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {tribes.map((tribe) => (
                <div key={tribe.id} className="bg-gray-800 rounded-xl p-6 border border-gray-700 hover:border-emerald-500 transition-colors">
                    <h3 className="text-xl font-bold text-white mb-2">{tribe.name}</h3>
                    <p className="text-gray-400 text-sm mb-4 h-12 overflow-hidden">{tribe.description}</p>

                    <div className="flex justify-between items-center mb-6">
                        <div className="text-center">
                            <div className="text-2xl font-bold text-emerald-400">{tribe.member_count}</div>
                            <div className="text-xs text-gray-500 uppercase">Members</div>
                        </div>
                        <div className="text-center">
                            <div className="text-2xl font-bold text-yellow-400">{tribe.total_discipline.toFixed(0)}</div>
                            <div className="text-xs text-gray-500 uppercase">Discipline</div>
                        </div>
                    </div>

                    <button
                        onClick={() => handleJoin(tribe.id)}
                        className="w-full py-2 rounded-lg bg-emerald-600 text-white font-semibold hover:bg-emerald-500 transition-colors"
                    >
                        Join Tribe
                    </button>
                </div>
            ))}
        </div>
    );
};

export default TribeList;
