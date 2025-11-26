// firebase-messaging-sw.js

// Give the service worker access to Firebase Messaging.
// Note that you can only use Firebase Messaging here. Other Firebase libraries
// are not available in the service worker.
importScripts('https://www.gstatic.com/firebasejs/10.7.0/firebase-app-compat.js');
importScripts('https://www.gstatic.com/firebasejs/10.7.0/firebase-messaging-compat.js');

// Initialize the Firebase app in the service worker by passing in
// your app's Firebase config object.
firebase.initializeApp({
    apiKey: "AIzaSyCe9Gvi8I9mWL1w7PCJfy1nhwbZ3qq0z4E",
    authDomain: "fastinghero.firebaseapp.com",
    projectId: "fastinghero",
    storageBucket: "fastinghero.firebasestorage.app",
    messagingSenderId: "707237131980",
    appId: "1:707237131980:web:9584970e382b16afe23469"
});

// Retrieve an instance of Firebase Messaging so that it can handle background
// messages.
const messaging = firebase.messaging();

messaging.onBackgroundMessage((payload) => {
    console.log('[firebase-messaging-sw.js] Received background message ', payload);

    const notificationTitle = payload.notification.title;
    const notificationOptions = {
        body: payload.notification.body,
        icon: '/icon-192x192.png', // Make sure you have this icon
        data: payload.data
    };

    self.registration.showNotification(notificationTitle, notificationOptions);
});
