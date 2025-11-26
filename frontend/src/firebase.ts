import { initializeApp } from 'firebase/app';
import { getMessaging, getToken, onMessage } from 'firebase/messaging';

const firebaseConfig = {
    apiKey: import.meta.env.VITE_FIREBASE_API_KEY,
    authDomain: import.meta.env.VITE_FIREBASE_AUTH_DOMAIN,
    projectId: import.meta.env.VITE_FIREBASE_PROJECT_ID,
    storageBucket: import.meta.env.VITE_FIREBASE_STORAGE_BUCKET,
    messagingSenderId: import.meta.env.VITE_FIREBASE_MESSAGING_SENDER_ID,
    appId: import.meta.env.VITE_FIREBASE_APP_ID,
};

// Initialize Firebase
const app = initializeApp(firebaseConfig);

// Initialize Firebase Cloud Messaging
let messaging: ReturnType<typeof getMessaging> | null = null;

try {
    messaging = getMessaging(app);
} catch (error) {
    console.warn('Firebase Messaging not supported in this environment:', error);
}

export { messaging };

export async function requestNotificationPermission(): Promise<string | null> {
    if (!messaging) {
        console.warn('Firebase Messaging not available');
        return null;
    }

    try {
        const permission = await Notification.requestPermission();
        if (permission === 'granted') {
            console.log('Notification permission granted');

            const vapidKey = import.meta.env.VITE_FIREBASE_VAPID_KEY;
            const token = await getToken(messaging, { vapidKey });

            if (token) {
                console.log('FCM Token:', token);
                return token;
            } else {
                console.log('No registration token available');
                return null;
            }
        } else {
            console.log('Notification permission denied');
            return null;
        }
    } catch (error) {
        console.error('Error getting notification permission:', error);
        return null;
    }
}

export function onMessageListener(): Promise<any> {
    if (!messaging) {
        return Promise.reject('Firebase Messaging not available');
    }

    return new Promise((resolve) => {
        onMessage(messaging!, (payload) => {
            console.log('Message received:', payload);
            resolve(payload);
        });
    });
}
