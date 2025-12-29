document.getElementById('loginform').addEventListener('submit', function (event) {
  event.preventDefault();

  const formData = {
    username: document.getElementById('username').value,
    password: document.getElementById('password').value
  };

   const token = localStorage.getItem('token');
  fetch('/api/auth', {
    method: 'POST',
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
       //notification("success", "Login successful");
       //window.location.href = "/dashboard";

       if(data.token){
          localStorage.setItem('token', data.token);
          window.location.href = '/dashboard';
          notification("success", "Login successful");

       }else{
            notification("error", `${data.token}`);

       }
    })
    .catch(err => {
      notification("error", err.message);
      //window.location.href = "/login"
    });
});
