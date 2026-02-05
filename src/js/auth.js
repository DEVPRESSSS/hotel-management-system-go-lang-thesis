document.getElementById('loginform').addEventListener('submit', function (event) {
  event.preventDefault();

  //Get the value of the inputs
  const username = document.getElementById('username').value.trim();
  const password = document.getElementById('password').value.trim();

  const usernameError = document.getElementById('validation-username');
  const passwordError = document.getElementById('validation-password');

  let hasError = false;

  // Reset previous errors
  usernameError.textContent = "";
  passwordError.textContent = "";

  // Username validation
  if (!username) {
    usernameError.textContent = "Username is required";
    hasError = true;
  }

  // Password validation
  if (!password) {
    passwordError.textContent = "Password is required";
    hasError = true;
  }

  // Stop if any validation failed
  if (hasError) return;

  const formData = {
    username,
    password
  };

  const token = localStorage.getItem('token');
  fetch('/api/auth', {
    method: 'POST',
    credentials: 'include', 
    headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
        'Access-Control-Allow-Headers': 'Content-Type, Authorization, X-Requested-With',
        'Access-Control-Expose-Headers': 'Authorization',
        'Access-Control-Allow-Origin': '*',
        'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, OPTIONS',
        'Access-Control-Allow-Credentials': 'true',
        'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify(formData)
  })
    .then(async response => {
      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || 'Login failed');
      }

      return data;
    })
    .then(data => {

      if (data.token) {
      
          notification("success", "Login successful");
          const routes = {
              Admin: "/api/dashboard",
              FrontDesk: "/api/dashboard",
              Guest: "/guest/dashboard",
          };

          setTimeout(() => {
              window.location.href = routes[data.role] ;
          }, 100);
      } else {
          notification("error", "Incorrect username or password");
      }

    })
    .catch(error => {
      notification("error", error.message);
    });
});
