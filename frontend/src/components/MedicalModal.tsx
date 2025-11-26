import React from 'react';

interface MedicalModalProps {
    isOpen: boolean;
    onConfirm: () => void;
    onCancel: () => void;
    startTime: string;
    onStartTimeChange: (time: string) => void;
    goalHours: number;
    onGoalHoursChange: (hours: number) => void;
    showGoalInput: boolean;
}

const MedicalModal: React.FC<MedicalModalProps> = ({
    isOpen,
    onConfirm,
    onCancel,
    startTime,
    onStartTimeChange,
    goalHours,
    onGoalHoursChange,
    showGoalInput
}) => {
    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 bg-black bg-opacity-80 flex items-center justify-center z-50 p-4">
            <div className="bg-gray-800 border border-red-500 rounded-2xl max-w-md w-full p-6 shadow-2xl">
                <div className="flex items-center gap-3 mb-4 text-red-500">
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                    </svg>
                    <h2 className="text-xl font-bold">Confirm Fast</h2>
                </div>

                <p className="text-gray-300 mb-6 leading-relaxed">
                    Ready to start your fast?
                    <br /><br />
                    <strong>Consult your physician</strong> before attempting prolonged fasting.
                </p>

                {showGoalInput && (
                    <div className="mb-4">
                        <label className="block text-sm font-medium text-gray-400 mb-2">
                            Goal Duration (Hours)
                        </label>
                        <input
                            type="number"
                            min="1"
                            max="168"
                            value={goalHours}
                            onChange={(e) => onGoalHoursChange(parseInt(e.target.value) || 0)}
                            className="w-full bg-gray-900 border border-gray-700 rounded-lg p-2 text-white focus:ring-2 focus:ring-red-500 focus:outline-none"
                        />
                    </div>
                )}

                <div className="mb-6">
                    <label className="block text-sm font-medium text-gray-400 mb-2">
                        Did you already start? (Optional)
                    </label>
                    <input
                        type="datetime-local"
                        value={startTime}
                        onChange={(e) => onStartTimeChange(e.target.value)}
                        className="w-full bg-gray-900 border border-gray-700 rounded-lg p-2 text-white focus:ring-2 focus:ring-red-500 focus:outline-none"
                    />
                </div>

                <div className="flex gap-4 justify-end">
                    <button
                        onClick={onCancel}
                        className="px-4 py-2 rounded-lg text-gray-400 hover:text-white hover:bg-gray-700 transition-colors"
                    >
                        Cancel
                    </button>
                    <button
                        onClick={onConfirm}
                        className="px-6 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg font-bold shadow-lg transition-colors"
                    >
                        Start Fast
                    </button>
                </div>
            </div>
        </div>
    );
};

export default MedicalModal;
