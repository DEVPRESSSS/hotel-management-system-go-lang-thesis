

document.getElementById('reset-password-form').addEventListener('submit', async function (event) {
    event.preventDefault();

    const password = document.getElementById('password').value;
    const passwordConfirm = document.getElementById('passwordConfirm').value;
    const errorPassword = document.getElementById('validate-password');
    const errorConfirmPassword = document.getElementById('validate-confirm-password');
    
    if(!password || !passwordConfirm){
        errorPassword.textContent= "Password is required";
        errorConfirmPassword.textContent= "Confirm password is required";
        return;
    }
  
    // Client-side validation
    if (password !== passwordConfirm) {
        
        notification("error","Password don't match");
        errorPassword.textContent= "";
        errorConfirmPassword.textContent= "";
        return;
    }

    if (password.length < 8) {
        notification("error","Password is too short");
        errorPassword.textContent= "";
        errorConfirmPassword.textContent= "";
        return;
    }

    // Get the reset token from URL
    const pathParts = window.location.pathname.split('/');
    const resetToken = pathParts[pathParts.length - 1];

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

        notification("success","Password changed successfully");
        setTimeout(() => {
            window.location.href = '/login';
        }, 2000);

    } catch (error) {
        notification("error",`Failed to change password:${error.message}`);     
    }
});

