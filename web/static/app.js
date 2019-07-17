console.log('Hello World!')
if ('serviceWorker' in navigator) {
    navigator.serviceWorker.register('/sw.js')
    .then(function(registration) {
        console.info('Registration successful, scope is:', registration.scope);
    })
    .catch(function(error) {
        console.info('Service worker registration failed, error:', error);
    });
}

window.addEventListener("beforeinstallprompt", function(event) {
    console.info("install prompt")
    event.preventDefault()
    deferredPrompt = event
    //event.prompt()

})
