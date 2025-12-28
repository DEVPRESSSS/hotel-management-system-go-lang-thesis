document.getElementById('loginform').addEventListener('submit', function (event) {
    event.preventDefault(); // Prevent form submission

    const formData = {
      username: document.getElementById('username').value,
      password: document.getElementById('password').value
    };


    // Send POST request to backend
    fetch('/api/auth', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
        // 'Access-Control-Allow-Headers': 'Content-Type, Authorization, X-Requested-With',
        // 'Access-Control-Expose-Headers': 'Authorization',
        // 'Access-Control-Allow-Origin': '*',
        // 'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, OPTIONS',
        // 'Access-Control-Allow-Credentials': 'true',
      },
      body: JSON.stringify(formData)
    })
      .then(response => response.json())
      .then(data => {
            
        notification("Success", data.message)
                
      })
      .catch(error => {
        console.log('Error:', error);

      });
  });