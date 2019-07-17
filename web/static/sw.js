var cacheName = 'huggable'

// Essential components of app shell
var essentials = [
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
]

self.addEventListener("install", function (event) {
    console.info("Installing service worker...")

    // Force all essential files to be in cache, in the event that we never get
    // to fetch them, such as when installing app but then losing internet
    event.waitUntil(
        caches.open(cacheName).then(cache => {
            cache.addAll(essentials)
        })
    )
})

self.addEventListener("activate", event => {
    console.info("Activated service worker!")
})

self.addEventListener('fetch', function (event) {    
    // Uses 
    event.respondWith(
        caches.open(cacheName).then(cache => {
            return cache.match(event.request).then(response => {
                // fetch in background even if in cache
                var fetchPromise = fetch(event.request).then(networkResponse => {
                    cache.put(event.request, networkResponse.clone());
                    return networkResponse;
                })

                // Return whats in cache, otherwise wait till fetch finishes and
                // return that.
                return response || fetchPromise;
            })
        })
    )
})