import React from 'react';

interface MedicalModalProps {
    isOpen: boolean;
    onConfirm: () => void;
    onCancel: () => void;
}

const MedicalModal: React.FC<MedicalModalProps> = ({ isOpen, onConfirm, onCancel }) => {
    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 bg-black bg-opacity-80 flex items-center justify-center z-50 p-4">
            <div className="bg-gray-800 border border-red-500 rounded-2xl max-w-md w-full p-6 shadow-2xl">
                <div className="flex items-center gap-3 mb-4 text-red-500">
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                    </svg>
                    <h2 className="text-xl font-bold">Medical Disclaimer</h2>
                </div>

                <p className="text-gray-300 mb-6 leading-relaxed">
                    You are about to start an extended fasting protocol (&gt;24 hours) or access advanced metabolic tracking.
                    <br /><br />
                    <strong>Consult your physician</strong> before attempting prolonged fasting, especially if you have diabetes, heart conditions, or are taking medication.
                    <br /><br />
                    Do you acknowledge the risks and wish to proceed?
                </p>

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
                        I Understand & Proceed
                    </button>
                </div>
            </div>
        </div>
    );
};

export default MedicalModal;
