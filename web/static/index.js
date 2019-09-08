console.log('Executing /index.js');

// Enable UI For granting notifications if they aren't already granted
if (Notification.permission !== "granted") {
  console.log('Notification permission not granted. Enabling button UI.')
  setNotificationUIActive(true)
}

function setNotificationUIActive(setActive) {
  const div = document.getElementById('notify-button')
  const button = div.getElementsByTagName('button')[0]
  if (setActive) {
    div.style.display = "block"
    button.onclick = () => {
      askNotificationPermission()
      .then(granted => {
        // Regardless of if user accepts or denies, we can't prompt them again till the browser lets us.
        // So just hide the button for now.
        setNotificationUIActive(false)
      })
    }
  } else {
    div.style.display = "none"
    button.onclick = null
  }
}