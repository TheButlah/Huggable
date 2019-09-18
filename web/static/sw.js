"use strict";

console.log('Executing sw.js');

// Use `me` instead of `self` as type of `self` is incorrect.
const me = (/** @type {ServiceWorkerGlobalScope} */ ((/** @type {unknown} */ (self))))

const cacheName = 'huggable';

// Essential components of app shell
const essentials = [
  // App Shell
  '/index.html',
  '/scripts/app.js',
  '/manifest.json',
  '/sw.js',

  // Icons
  '/icons/favicon.ico',
  '/icons/android-chrome-192x192.png',
  '/icons/android-chrome-512x512.png',
  '/icons/apple-touch-icon.png',
];

me.addEventListener("install", /** @param {FetchEvent} event */ event => {
  console.log("Installing service worker...")
  // Force all essential files to be in cache, in the event that we never get
  // to fetch them, such as when installing app but then losing internet.
  // Also ensure that we don't complete the event until the promise settles.
  event.waitUntil((async () => {
    console.log("Caching app shell...");
    try {
      const cache = await caches.open(cacheName);
      await cache.addAll(essentials);
      console.log("Done caching app shell");
    } catch (error) {
      console.error("Error while caching app shell: " + error);
    }
    // Don't wait for previous service worker to shutdown
    me.skipWaiting()
  })());
});

me.addEventListener("activate", event => {
  console.info("Activated service worker!");
});

me.addEventListener('fetch', /** @param {FetchEvent} event */ event => {
  // Only mess with GET requests
  if (event.request.method != 'GET') return;

  // Do not cache API calls
  if (new URL(event.request.url).pathname.startsWith("/api")) {
    return
  }

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
      if (!cachedResponse) {
        console.log("Cache doesn't contain an entry for " + event.request.url);
      }
      // Return whats in cache, otherwise use response from fetch
      return cachedResponse || fetchPromise;

    } catch (error) {
      console.error("Error during fetch event: " + error);
    }
  })());
});