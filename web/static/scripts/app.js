console.log('Executing /main/app.js');

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

    window.addEventListener('beforeinstallprompt', (event) => {
        console.log("install prompt fired");

        // Prevent Chrome 67 and earlier from automatically showing the prompt
        event.preventDefault();

        const deferredPrompt = event;
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
    });

    /*
    try {
        askNotificationPermission()
    } catch (error) {
        console.error("Error when asking for notification: " + error)
    }*/
}

/** Asks for permission to show notifications. */
function askNotificationPermission() {
    // The latest api uses a promise for `requestPermission()`, but the old API used a callback.
    // Solution: use deprecated API, and check to see if the function returns undefined. If so,
    return new Promise(function (resolve, reject) {
        const permissionResult = Notification.requestPermission(function (result) {
            resolve(result);
        });
        console.log(permissionResult);
        if (permissionResult) {
            console.log('notification API using promises');
            permissionResult.then(resolve, reject);
        } else {
            console.log('notification API using deprecated callbacks');
        }
    })
    .then(function (permissionResult) {
        if (permissionResult !== 'granted') {
            throw new Error('We weren\'t granted permission.');
        }
    });
}
