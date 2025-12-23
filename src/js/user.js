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

        fetch('/userslist',{
            method:'POST',
            headers:{
                'Content-Type':'application/json',
                'Accept':'application/json',
                //'Authorization': 'Bearer' +
            },
            body:JSON.stringify(formData)
        })
        .then(response =>{

            if(response.ok){

                alert('User created successfully');
            }else{
                return response.json();
            }
        })
        .then(data =>{

            alert(data.error);
            console.log(data.error);
        })
         .catch(error =>{

            console.error('Error:', error);
            document.getElementById('message').textContent = 'An error occurred while processing your request.';
        });
        


    });


 //Create uuid
function uuidv4() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'
    .replace(/[xy]/g, function (c) {
        const r = Math.random() * 16 | 0, 
            v = c == 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}