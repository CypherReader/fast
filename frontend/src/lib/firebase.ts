// Mock Firebase configuration for development
// In production, replace with actual Firebase initialization

export const messaging = {
    // Mock messaging object
};

export const getToken = async (messaging: any, options?: any) => {
    console.log("Mock getToken called with options:", options);
    // Simulate network delay
    await new Promise(resolve => setTimeout(resolve, 500));
    return "mock-fcm-token-" + Date.now();
};

export const onMessage = (messaging: any, callback: (payload: any) => void) => {
    console.log("Mock onMessage listener registered");
    // Mock receiving a message after 10 seconds
    setTimeout(() => {
        callback({
            notification: {
                title: "Mock Notification",
                body: "This is a test notification from the mock Firebase setup."
            }
        });
    }, 10000);
};
