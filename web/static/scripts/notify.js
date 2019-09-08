console.log('Executing /scripts/notify.js');

/**
 * Asks for permission to show notifications.
 * @returns true if permission granted, otherwise false.
*/
function askNotificationPermission() {
  // The latest api uses a promise for `requestPermission()`, but the old API used a callback.
  // Solution: use deprecated API, and check to see if the function returns undefined, or a Promise
  if (Notification.permission === 'granted') {
      console.warn('Asking to grant notification perms, but they are already granted. Are you sure you want to do this?')
  }
  return new Promise(function (resolve, reject) {
      const permissionResult = Notification.requestPermission(function (result) {
          resolve(result)
      });
      if (permissionResult) {
          console.log('notification API using promises')
          permissionResult.then(resolve, reject)
      } else {
          console.log('notification API using deprecated callbacks')
      }
  })
  .then(function (permissionResult) {
      if (permissionResult === 'granted') {
          console.log('Notification permission granted')
          return true
      } else {
          console.log('Notification permission denied')
          return false
      }
  });
}

function showNotification(title, body) {
  const options = {
    image: "/icons/android-chrome-192x192.png",
    icon: "/icons/android-chrome-192x192.png",
    body: body,
  }
  const notification = new Notification(title, options)
}