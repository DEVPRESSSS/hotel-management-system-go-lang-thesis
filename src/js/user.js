
//#region--Modal Upsert 
    function openModal() {
        document.getElementById('userModal').classList.remove('hidden');
        document.getElementById('userModal').classList.add('flex');
    }

    function closeModal() {
        document.getElementById('userModal').classList.add('hidden');
        document.getElementById('userModal').classList.remove('flex');
    }

    // Close modal on outside click
    document.getElementById('userModal').addEventListener('click', function(e) {
        if (e.target === this) {
            closeModal();
        }
    });

    //validate password first check if it is equal to confirm password
    function validatePassword(){
        const password = document.getElementById('password').value;
        const confirmPassword = document.getElementById('confirm-password').value;
        const errorMessage = document.getElementById('error-message');

        if(password !== confirmPassword){
            errorMessage.textContent = "Password don't match";
            errorMessage.classList.remove("hidden");
            return;
        }
        errorMessage.textContent = "";
        
    }

    //Send request to the backend

    document.getElementById('createUserForm').addEventListener('submit',function(event){
        event.preventDefault();

        const fullname = document.getElementById('fullname').value;
        const email = document.getElementById('email').value;
        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;
        const locked = false;
        const role = document.getElementById('roleid').value;
        //const role = "ROLE-101"
        const uid = uuidv4();
        //Assign the value
        const formData = {
            userid: uid,
            fullname: fullname,
            email: email,
            username: username,
            password: password,
            locked: locked,
            roleid: role
        };

        fetch('/userslist', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(formData)
        })
        .then(async response => {
            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error);
            }
            return data;
        })
        .then(data => {
            notification("success", data.message);

        })
        .catch(err => {
            //document.getElementById('error-message').textContent = err.message;
            notification("error", err.message);
        });
    });

//#endregion

//Create uuid for userid in usrs
function uuidv4() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'
    .replace(/[xy]/g, function (c) {
        const r = Math.random() * 16 | 0, 
            v = c == 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}

//Fetch all users
fetch('/api/users')
    .then(response => response.json())
        .then(users => {

            const tbody = document.getElementById('users-body');

            users.forEach(user => {
                tbody.innerHTML += `
                
                    <tr>
                        <td>${user.userid}</td>
                        <td>${user.fullname}</td>
                        <td>${user.email}</td>
                        <td>${user.role.rolename}</td>
                        <td>
                        ${
                            user.locked === false
                            ? '<p class="text-green-600 font-semibold">Active</p>'
                            : '<p class="text-red-600 font-semibold">Locked</p>'
                        }
                        </td>
                        <td>${new Date(user.created_at).toLocaleDateString()}</td>
                        <td>
                            <button class="text-blue-600 hover:text-blue-800 font-medium">Edit</button>
                            <button 
                                class="delete-btn text-red-600 hover:text-red-800 font-medium"
                                data-userid="${user.userid}">
                                Delete
                            </button>
                        </td>
                    </tr>
                
                `;
            
                
            });
});

document.getElementById('users-body').addEventListener('click', function (e) {
    if (!e.target.classList.contains('delete-btn')) return;

    const userid = e.target.dataset.userid;
    if (!userid) return;

    if (!confirm('Are you sure you want to delete this user?')) return;

    fetch(`/api/delete/${userid}`, {
        method: 'DELETE'
    })
    .then(res => {
        if (!res.ok) throw new Error('Delete failed');
        e.target.closest('tr').remove(); // remove row instantly
        alert('User deleted successfully');
    })
    .catch(err => {
        console.error(err);
        alert('Error deleting user');
    });
});


