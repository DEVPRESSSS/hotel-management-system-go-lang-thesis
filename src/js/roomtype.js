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
    headerTitle.innerText = "Create user";
    btnSubmit.innerText = "Create";
    openModal();
}

//Upsert room type
document.getElementById('upsertform').addEventListener('submit', function(e){

    e.preventDefault();

    let uid = "";
    const roomType = document.getElementById('roomtype').value;
    const description = document.getElementById('description').value;
    if( id === ""){
        uid = uuidv4();
    }else{
        uid = id;
    }

    const formData ={
        roomtypeid: uid,
        roomtypename: roomType,
        description: description,

    };
    if(id === "" || id === null){
        fetch('/api/createroomtype', {
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

                notification("success", data.success);

            })
            .catch(err=>{

                notification("error", err.message);

            });
    }else{
           //Update user
           const updateFormData ={
                roomtypeid: uid,
                roomtypename: roomType,
                description: description,
            };
            fetch(`/api/updateroomtype/${id}`, {
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
fetch('/api/roomtypes')
    .then(response => response.json())
        .then(rt => {

            const tbody = document.getElementById('users-body');

            rt.forEach(r => {
                tbody.innerHTML += `
                
                    <tr>
                        <td>${r.roomtypeid}</td>
                        <td>${r.roomtypename}</td>                 
                        <td>${r.description}</td>                 
                        <td>
                            <button class="update-btn text-blue-600 hover:text-blue-800 font-medium"
                                data-roomtypeid = "${r.roomtypeid}"
                                >Edit
                            
                            </button>
                            <button 
                                class="delete-btn text-red-600 hover:text-red-800 font-medium"
                                data-roomtypeid="${r.roomtypeid}">
                                Delete
                            </button>
                        </td>
                    </tr>
                
                `;               
            });
});
//Edit click 
document.getElementById('users-body').addEventListener('click', function (e) {

    if (!e.target.classList.contains('update-btn')) return;

    id = e.target.dataset.roomtypeid;
    if (!id) return;

    const roomTypeId = document.getElementById('input-id');
    roomTypeId.value = id;

    btnSubmit.innerText = "Update";
    headerTitle.innerText = "Update";

    fetch(`/api/roomtype/${id}`)
        .then(async res => {
            if (!res.ok) {
                const text = await res.text();
                throw new Error(text);
            }
            return res.json();
        })
        .then(data => {
            const rt = data.success;
            document.getElementById('roomtype').value = rt.roomtypename;
            document.getElementById('description').value = rt.description;
        })
        .catch(err => {
            console.error(err);
            Swal.fire("Error", "Failed to load room type", "error");
        });

    openModal();
});



//Delete user function
document.getElementById('users-body').addEventListener('click', function (e) {
    if (!e.target.classList.contains('delete-btn')) return;

    const roomtypeid = e.target.dataset.roomtypeid;
    if (!roomtypeid) return;

    Swal.fire({
        title: "Are you sure?",
        text: "You won't be able to revert this!",
        icon: "warning",
        showCancelButton: true,
        confirmButtonColor: "#3085d6",
        cancelButtonColor: "#d33",
        confirmButtonText: "Yes, delete it!"
    }).then((result) => {
        if (!result.isConfirmed) return;

        fetch(`/api/deleteroomtype/${roomtypeid}`, {
            method: 'DELETE'
        })
        .then(async res => {
            if (!res.ok) {
                const data = await res.json();
                throw new Error(data.error || 'Delete failed');
            }

            // ✅ Show success ONLY here
            Swal.fire({
                title: "Deleted!",
                text: "Room type has been deleted.",
                icon: "success"
            });

            // Remove row from table
            e.target.closest('tr').remove();
        })
        .catch(err => {
            Swal.fire({
                title: "Error!",
                text: err.message,
                icon: "error"
            });
        });
    });
});
