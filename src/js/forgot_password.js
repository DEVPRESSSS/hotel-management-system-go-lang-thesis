document.getElementById('forgot-password-form').addEventListener('submit', function (event) {
  event.preventDefault();
  

  const email = document.getElementById('email').value;
  const error = document.getElementById('validation-email');
  if(!email){
    //notification("error","Email is required")
    error.textContent = "Email is required";
    return;
  }
  const formData = {
    email: email
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
      
        notification("success", "Password reset link has been sent to your email");
        document.getElementById('forgot-password-form').reset();
    })
    .catch(error => {
        notification("error", `${error.message}`);
    });
});