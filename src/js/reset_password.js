document.getElementById('reset-password-form').addEventListener('submit', async function (event) {
    event.preventDefault();

    const password = document.getElementById('password').value;
    const passwordConfirm = document.getElementById('passwordConfirm').value;
    const submitButton = document.getElementById('submit-btn');

    // Client-side validation
    // if (password !== passwordConfirm) {
    //     return;
    // }

    // if (password.length < 8) {
    //     return;
    // }

    // Get the reset token from URL
    const pathParts = window.location.pathname.split('/');
    const resetToken = pathParts[pathParts.length - 1];

    console.log('Reset Token:', resetToken); // For debugging

  
    try {
        const response = await fetch(`/api/resetpassword/${resetToken}`, {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                password: password,
                passwordConfirm: passwordConfirm
            })
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.message || 'Failed to reset password');
        }

        
        setTimeout(() => {
            window.location.href = '/login';
        }, 2000);

    } catch (error) {
        //alert(`${error.message}`);
      
    }
});
