const cacheName = 'iglu-v1';
const staticAssets = [
    './',
    'index.html'
]

self.addEventListener('install', async e => {
    const cache = await caches.open(cacheName);
    await cache.addAll(staticAssets);
    return self.skipWaiting();
});

self.addEventListener('activate', e => {
    self.clients.claim();
})