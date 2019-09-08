console.log('Executing /scripts/app.js');

/** Registers a ServiceWorker and returns the registration promise. */
async function registerSW() {
    try {
        // Note that the jsconfig.json moves the service worker to the root web dir.
        // This is because the maximum scope for a service worker is its current directory.
        const registration = await navigator.serviceWorker.register('/sw.js');
        console.log('Service worker registration successful, scope is:', registration.scope);
        return registration;
    } catch (error) {
        console.error('Service worker registration failed, error: ' + error);
    }
}

if ('serviceWorker' in navigator) {
    const swRegistration = registerSW();  // This is a promise

    // Only supported on Chrome
    try {
        window.addEventListener('beforeinstallprompt', event => {
            console.log("install prompt fired");

            // Prevent Chrome 67 and earlier from automatically showing the prompt
            event.preventDefault();

            const deferredPrompt = event;
            /*
            Disable auto prompt for installing on google chrome
            until we can figure out how to standardize it across all browsers
            document.addEventListener('click', () => {
                deferredPrompt.prompt();
            });

            deferredPrompt.userChoice.then(choice => {
                if (choice.outcome === 'accepted') {
                    console.info('User accepted the installation');
                } else {
                    console.info('User declined the installation');
                }
            });
            */
        });
    } catch (error) {
        console.error("Error while adding listener for 'beforeinstallevent' error: " + error)
    }
}