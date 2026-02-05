document.getElementById('forgot-password-form').addEventListener('submit', function (event) {
  event.preventDefault();
  
  const formData = {
    email: document.getElementById('email').value,
  };

  // Show loading state (optional)
  const submitButton = this.querySelector('button[type="submit"]');
  const originalText = submitButton.textContent;
  submitButton.textContent = 'Sending...';
  submitButton.disabled = true;

  fetch('/forgotpassword', {  
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',  
    },
    body: JSON.stringify(formData)
  })
    .then(async response => {
      const data = await response.json();
      
      // Reset button state
      submitButton.textContent = originalText;
      submitButton.disabled = false;
      
      if (!response.ok) {
        throw new Error(data.message || 'Request failed');
      }
      return data;
    })
    .then(data => {
      // Success - show success message
      showNotification("success", data.message || "Password reset link has been sent to your email!");
      
      // Optional: Clear the form
      document.getElementById('forgot-form').reset();
    })
    .catch(error => {
      // Reset button state
    //   submitButton.textContent = originalText;
    //   submitButton.disabled = false;
      
      // Error - show error message
      showNotification("error", error.message || "Failed to send reset email. Please try again.");
    });
});