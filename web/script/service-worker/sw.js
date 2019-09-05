console.log('Executing sw.js');

const cacheName = 'huggable';

// Essential components of app shell
const essentials = [
    // App Shell
    '/index.html',
    '/app.js',
    '/manifest.json',
    '/sw.js',

    // Icons
    '/icons/favicon.ico',
    '/icons/android-chrome-192x192.png',
    '/icons/android-chrome-512x512.png',
    '/icons/apple-touch-icon.png',
];

self.addEventListener("install", /** @param {FetchEvent} event */ event => {
    console.log("Installing service worker...")
    // Force all essential files to be in cache, in the event that we never get
    // to fetch them, such as when installing app but then losing internet.
    // Also ensure that we don't complete the event until the promise settles.
    event.waitUntil(async () => {
        console.log("Caching app shell...");
        try {
            const cache = await caches.open(cacheName);
            await cache.addAll(essentials);
            console.log("Done caching app shell");
        } catch (error) {
            console.error("Error while caching app shell: " + error);
        }
    });
});

self.addEventListener("activate", event => {
    console.info("Activated service worker!");
});

self.addEventListener('fetch', /** @param {FetchEvent} event */ event => {
    // Only mess with GET requests
    if (event.request.method != 'GET') return;

    // Prevent the default, and handle the request ourselves.
    event.respondWith((async () => {
        const cache = await caches.open(cacheName);

        // fetch in background even if in cache
        const fetchPromise = fetch(event.request)
        .then(response => {
            cache.put(event.request, response.clone());
            return response;
        }).catch(error => {
            console.error("Error while fetching the request: " + error);
            return Promise.reject(error);
        });

        try {
            const cachedResponse = await cache.match(event.request);
            console.log("cachedResponse: " + cachedResponse);
            if (!cachedResponse) {
                console.log("Cache doesn't contain an entry for " + event.request);
            }
            console.log("fetchPromise: " + fetchPromise);

            // Return whats in cache, otherwise use response from fetch
            return cachedResponse || fetchPromise;

        } catch (error) {
            console.error("Error during fetch event: " + error);
        }
    })());
});