//Access the header title to modify it later
const headerTitle = document.getElementById('header-action');
//Access the header text to modify its text inside
const btnSubmit = document.getElementById('btn-submit');

//Stors user id globally
let id = "";
//#region--Modal Upsert 
function openModal() {
    document.getElementById('userModal').classList.remove('hidden');
    document.getElementById('userModal').classList.add('flex');
}

function createModal(){
    headerTitle.innerText = "Create user";
    btnSubmit.innerText = "Create";
    openModal();
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


document.getElementById('upsertform').addEventListener('submit',function(event){
        event.preventDefault();

        const fullname = document.getElementById('fullname').value;
        const email = document.getElementById('email').value;
        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;
        const locked = false;
        const role = document.getElementById('roleid').value;
        let uid = "";

        if (id === "" || id === null) {
            uid = uuidv4(); 
        } else {
            uid = id;      
        }

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

        if(id === "" || id === null){
            //Create user
            console.log(id);
            
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
                    closeModal();

                })
                .catch(err => {
                    notification("error", err.message);
            });

            return;
        }else{
            //Update user
            const updateFormData = {
                userid: uid,
                fullname: fullname,
                email: email,
                username: username,
                password: password,
                locked: locked,
                roleid: role
            };
            fetch(`/api/update/${id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(updateFormData)
                })
                .then(async response => {
                    const data = await response.json();

                    if (!response.ok) {
                        throw new Error(data.error);
                    }
                    return data;
                })
                .then(data => {
                    notification("success", data.success);

                })
                .catch(err => {
                    notification("error", err.error);
            });
        }
       
        closeModal();
        
        

        
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
                        <td>${user.username}</td>
                        <td>${user.email}</td>
                        <td>${user.Role.rolename}</td>
                        <td>
                        ${
                            user.locked === false
                            ? '<p class="text-green-600 font-semibold">Active</p>'
                            : '<p class="text-red-600 font-semibold">Locked</p>'
                        }
                        </td>
                        <td>${new Date(user.created_at).toLocaleDateString()}</td>
                        <td>
                            <button class="update-btn text-blue-600 hover:text-blue-800 font-medium"
                                data-userid = "${user.userid}"
                                >Edit
                            
                            </button>
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


//Delete user function
document.getElementById('users-body').addEventListener('click', function (e) {
    if (!e.target.classList.contains('delete-btn')) return;

    const userid = e.target.dataset.userid;
    if (!userid) return;

    Swal.fire({
        title: "Are you sure?",
        text: "You won't be able to revert this!",
        icon: "warning",
        showCancelButton: true,
        confirmButtonColor: "#3085d6",
        cancelButtonColor: "#d33",
        confirmButtonText: "Yes, delete it!"
    }).then((result) => {
        if (result.isConfirmed) {
            fetch(`/api/delete/${userid}`, {
            method: 'DELETE'
        })
        .then(res => {
            if (!res.ok) throw new Error('Delete failed');
            e.target.closest('tr').remove(); 
            
        })
        .catch(err => {
            console.error(err);
            alert('Error deleting user');
        });

        Swal.fire({
        title: "Deleted!",
        text: "Your file has been deleted.",
        icon: "success"
        });
    }
    });

   
});


//Edit click 
document.getElementById('users-body').addEventListener('click' , function(e){

    if (!e.target.classList.contains('update-btn')) return;

    id = e.target.dataset.userid;
    if (!id) return;

    const userInputId = document.getElementById('input-id');
    userInputId.value = id;

    btnSubmit.innerText = "Update";
    headerTitle.innerText = "Update user";  
        
    
    fetch(`/api/user/${id}`)
        .then(async res => {
            const data = await res.json();
            if (!res.ok) throw new Error(data.error || "Request failed");
            return data;
        })
        .then(data => {
            //notification("success", data.message || "Operation successful");
            const user = data.success;           
            document.getElementById('fullname').value = user.fullname;
            document.getElementById('email').value = user.email;
            document.getElementById('username').value = user.username;
            document.getElementById('roleid').value = user.roleid;
            document.getElementById('password').value = user.value;
            document.getElementById('confirm-password').value = user.value;

        })
        .catch(err => {
            notification("error", err.message);
        });

    openModal();

}); 





