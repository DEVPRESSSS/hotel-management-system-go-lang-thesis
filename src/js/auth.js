document.getElementById('loginform').addEventListener('submit', function (event) {
  event.preventDefault();

  const formData = {
    username: document.getElementById('username').value,
    password: document.getElementById('password').value
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

       if(data.token){
          localStorage.setItem('token', data.token);
          notification("success", "Login successful");
          let url = "";
          switch (data.role) {
            case "Admin":
              url = "/api/dashboard";
              break;
            case "Guest":
              url= "/guest/dashboard";
              break;
            default:
              window.location.href = "/unauthorized";
          }
        
          setTimeout(() => {
             //Redirect to any dashboard depending on the role
              window.location.href = url;
              
           }, 100);
       }else{
            notification("error", `Incorect username or password`);

       }
    })
    .catch(error => {
      notification("error", error.message);
    });
});
