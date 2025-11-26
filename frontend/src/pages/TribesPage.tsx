import React, { useState } from 'react';
import TribeList from '../components/TribeList';
import { api } from '../services/api';

const TribesPage: React.FC = () => {
    const [showCreateModal, setShowCreateModal] = useState(false);
    const [newTribeName, setNewTribeName] = useState('');
    const [newTribeDesc, setNewTribeDesc] = useState('');

    const handleCreateTribe = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            await api.post('/tribes/', { name: newTribeName, description: newTribeDesc });
            setShowCreateModal(false);
            setNewTribeName('');
            setNewTribeDesc('');
            // Ideally refresh list here, but for MVP we can just reload or rely on TribeList internal refresh if we lift state
            window.location.reload();
        } catch (error) {
            console.error('Failed to create tribe', error);
            alert('Failed to create tribe');
        }
    };

    return (
        <div className="min-h-screen bg-gray-900 text-white p-8 pb-24">
            <div className="max-w-6xl mx-auto">
                <div className="flex justify-between items-center mb-8">
                    <div>
                        <h1 className="text-3xl font-bold mb-2">Tribes</h1>
                        <p className="text-gray-400">Join a community to boost your discipline.</p>
                    </div>
                    <button
                        onClick={() => setShowCreateModal(true)}
                        className="px-6 py-3 bg-emerald-600 rounded-xl font-bold hover:bg-emerald-500 transition-colors shadow-lg shadow-emerald-900/20"
                    >
                        + Create Tribe
                    </button>
                </div>

                <TribeList />

                {/* Create Tribe Modal */}
                {showCreateModal && (
                    <div className="fixed inset-0 bg-black/80 flex items-center justify-center z-50 p-4">
                        <div className="bg-gray-800 rounded-2xl p-8 max-w-md w-full border border-gray-700">
                            <h2 className="text-2xl font-bold mb-6">Create New Tribe</h2>
                            <form onSubmit={handleCreateTribe}>
                                <div className="mb-4">
                                    <label className="block text-gray-400 text-sm mb-2">Tribe Name</label>
                                    <input
                                        type="text"
                                        value={newTribeName}
                                        onChange={(e) => setNewTribeName(e.target.value)}
                                        className="w-full bg-gray-900 border border-gray-700 rounded-lg p-3 text-white focus:border-emerald-500 outline-none"
                                        placeholder="e.g., Keto Warriors"
                                        required
                                    />
                                </div>
                                <div className="mb-6">
                                    <label className="block text-gray-400 text-sm mb-2">Description</label>
                                    <textarea
                                        value={newTribeDesc}
                                        onChange={(e) => setNewTribeDesc(e.target.value)}
                                        className="w-full bg-gray-900 border border-gray-700 rounded-lg p-3 text-white focus:border-emerald-500 outline-none h-24"
                                        placeholder="What is this tribe about?"
                                        required
                                    />
                                </div>
                                <div className="flex gap-4">
                                    <button
                                        type="button"
                                        onClick={() => setShowCreateModal(false)}
                                        className="flex-1 py-3 rounded-xl border border-gray-600 text-gray-300 font-bold hover:bg-gray-700"
                                    >
                                        Cancel
                                    </button>
                                    <button
                                        type="submit"
                                        className="flex-1 py-3 rounded-xl bg-emerald-600 text-white font-bold hover:bg-emerald-500"
                                    >
                                        Create
                                    </button>
                                </div>
                            </form>
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
};

export default TribesPage;
