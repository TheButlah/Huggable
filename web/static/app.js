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

let deferredPrompt
window.addEventListener('beforeinstallprompt', function(event) {
    console.info("install prompt")

    // Prevent Chrome 67 and earlier from automatically showing the prompt
    event.preventDefault()

    deferredPrompt = event
    document.addEventListener('click', function(event) {
        deferredPrompt.prompt()
    })

    deferredPrompt.userChoice.then(choice => {
        if (choice.outcome === 'accepted') {
            console.info('User accepted the installation')
        } else {
            console.info('User declined the installation')
        }
        deferredPrompt = null
    })
})