var cacheName = 'huggable'

self.addEventListener("install", function (event) {
    console.info("Installing service worker...")

    /*
    // Force all essential files to be in cache
    event.waitUntil(
        caches.open(cacheName).then(cache => {
            cache.addAll([
                '/index.html',
                '/app.js',
                '/manifest.json',
                '/sw.js',
            ])
        })
    )*/
})

self.addEventListener("activate", event => {
    console.info("Activated service worker!")
})

self.addEventListener('fetch', function (event) {    
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