//Access the header title to modify it later
const headerTitle = document.getElementById('header-action');
//Access the header text to modify its text inside
const btnSubmit = document.getElementById('btn-submit');
//Stors user id globally
let id = "";
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
function createModal(){
    headerTitle.innerText = "Create user";
    btnSubmit.innerText = "Create";
    openModal();
}

fetch('/api/roles')
  .then(res => res.json())
  .then(data => {
      AppState.roles = data;

      const select = document.getElementById('roleid');
      if (!select) return;

      // clear existing options
      select.innerHTML = '<option value="">Select role</option>';

      data.forEach(r => {
          const option = document.createElement('option');
          option.value = r.roleid;
          option.textContent = r.rolename; 
          select.appendChild(option);
      });
  })
  .catch(err => console.error('Failed to load roles:', err));


  //Upsert role
document.getElementById('upsertform').addEventListener('submit', function(e){

    e.preventDefault();
    //Generate unique uid
    const uid = uuidv4();
    //Get the input in role textbox
    const roleName = document.getElementById('rolename').value;

    const formData ={
        roleId: uid,
        roleName: roleName,

    };

    fetch('/api/create', {
        method:'POST',
        headers:{
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(formData),
         })
         .then(async response => {
                    const data = await response.json();

                    if (!response.ok) {
                        throw new Error(data.error);
                    }
                    return data;
        
         })
         .then(data=>{

            //Show notication success
            notification("success", data.success);

         })
         .catch(err=>{

            notification("error", err.message);

         });


});

function uuidv4() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'
    .replace(/[xy]/g, function (c) {
        const r = Math.random() * 16 | 0, 
            v = c == 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}

//Fetch all roles
fetch('/api/roles')
    .then(response => response.json())
        .then(roles => {

            const tbody = document.getElementById('users-body');

            roles.forEach(role => {
                tbody.innerHTML += `
                
                    <tr>
                        <td>${role.roleid}</td>
                        <td>${role.rolename}</td>                  
                        <td>
                            <button class="update-btn text-blue-600 hover:text-blue-800 font-medium"
                                data-roleid = "${role.roleid}"
                                >Edit
                            
                            </button>
                            <button 
                                class="delete-btn text-red-600 hover:text-red-800 font-medium"
                                data-roleid="${role.roleid}">
                                Delete
                            </button>
                        </td>
                    </tr>
                
                `;               
            });
});
//Edit click 
document.getElementById('users-body').addEventListener('click' , function(e){

    if (!e.target.classList.contains('update-btn')) return;

    id = e.target.dataset.roleid;
    if (!id) return;

    const roleInputId = document.getElementById('input-id');
    roleInputId.value = id;

    btnSubmit.innerText = "Update";
    headerTitle.innerText = "Update user";  
        
    
    fetch(`/api/roles/${id}`)
        .then(async res => {
            const data = await res.json();
            if (!res.ok) throw new Error(data.error || "Request failed");
            return data;
        })
        .then(data => {
            //notification("success", data.success || "Operation successful");           
            const userRole= data.success;
            document.getElementById('rolename').value = userRole.rolename;

        })
        .catch(err => {
            console.log(err);
            alert(`${err}`);
            //notification("error", err.err);
        });

    openModal();

}); 
