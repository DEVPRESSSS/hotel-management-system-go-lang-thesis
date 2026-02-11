
 /* =====================
     PASSWORD VALIDATION
  ====================== */
  const form= document.getElementById('register')
  function validatePassword() {
    const password = document.getElementById('password').value;
    const confirmPassword = document.getElementById('confirm-password').value;

    if (password !== confirmPassword) {
      errorMessage.textContent = "Passwords don't match";
      errorMessage.classList.remove("hidden");
      return false;
    }
    
    // errorMessage.textContent = "";
    // errorMessage.classList.add("hidden");
    return true;
  }

/* =====================
     FORM SUBMIT
  ====================== */
  form.addEventListener('submit', e => {
    e.preventDefault();

    // Validate passwords match
    if (!validatePassword()) {
      notification("error", "Passwords don't match");
      return;
    }

    const fullname = document.getElementById('fullname').value;
    const email = document.getElementById('email').value;
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    const payload = {
      fullname: fullname,
      email: email,
      username: username,
      password: password,
      locked: false,
    };


    fetch('/createaccount', {
      method:'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
    .then(res => {
        return res.json().then(data => {
        // Check if request was successful
        if (!res.ok) {
          // Throw error with server message
          throw new Error(data.error || 'Request failed');
        }
        return data;
      });
    })
    .then(data => {
      notification("success", data.success || data.message);
    
      setTimeout(() => {window.location.href= "/login"

      },3000)
      
    })
    .catch(err => {
      notification("error", err.message);
    });
  });
