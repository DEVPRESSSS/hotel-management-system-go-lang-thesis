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

function closeModal() {
        document.getElementById('userModal').classList.add('hidden');
        document.getElementById('userModal').classList.remove('flex');
}
function createModal(){
    headerTitle.innerText = "Create aminity";
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
    let uid = "";
    //Get the input in role textbox
    const aminityname = document.getElementById('aminity').value;
    if( id === ""){
        uid = uuidv4();
    }else{
        uid = id;
    }

    const formData ={
        aminityid: uid,
        aminityname: aminityname,

    };
    if(id === "" || id === null){
        fetch('/api/createaminity', {
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
    }else{
           //Update user
            const updateFormData = {
                  aminityid: uid,
                  aminityname: aminityname,       
            };
            fetch(`/api/updateaminity/${id}`, {
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

    //Close the modal
    closeModal();

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
fetch('/api/aminities')
    .then(response => response.json())
        .then(aminities => {

            const tbody = document.getElementById('users-body');

            aminities.forEach(aminity => {
                tbody.innerHTML += `
                
                    <tr>
                        <td>${aminity.aminityid}</td>
                        <td>${aminity.aminityname}</td>                  
                        <td>
                            <button class="update-btn text-blue-600 hover:text-blue-800 font-medium"
                                data-aminityid = "${aminity.aminityid}"
                                >Edit
                            
                            </button>
                            <button 
                                class="delete-btn text-red-600 hover:text-red-800 font-medium"
                                data-aminityid="${aminity.aminityid}">
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

    id = e.target.dataset.aminityid;
    if (!id) return;

    const roleInputId = document.getElementById('input-id');
    roleInputId.value = id;

    btnSubmit.innerText = "Update";
    headerTitle.innerText = "Update aminity";  
        
    
    fetch(`/api/aminity/${id}`)
        .then(async res => {
            const data = await res.json();
            if (!res.ok) throw new Error(data.error || "Request failed");
            return data;
        })
        .then(data => {
            //notification("success", data.success || "Operation successful");           
            const aminityname= data.success;
            document.getElementById('aminity').value = aminityname.aminityname;

        })
        .catch(err => {
            console.log(err);
            alert(`${err}`);
        });

    openModal();

}); 


//Delete user function
document.getElementById('users-body').addEventListener('click', function (e) {
    if (!e.target.classList.contains('delete-btn')) return;

    const aminityid = e.target.dataset.aminityid;
    if (!aminityid) return;

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
            fetch(`/api/deleteaminity/${aminityid}`, {
            method: 'DELETE'

            
        },  
            Swal.fire({
            title: "Deleted!",
            text: "Your file has been deleted.",
            icon: "success"
            })
        )
        .then(res => {
            if (!res.ok) throw new Error('Delete failed');
            e.target.closest('tr').remove(); 
            
        })
        .catch(err => {
            console.log(err);
            alert('Error deleting user');
        });

      
    }
    });

   
});