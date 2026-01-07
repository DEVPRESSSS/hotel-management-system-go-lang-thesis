
//#region--Modal navigations
const headerTitle = document.getElementById('header-action');
const btnSubmit = document.getElementById('btn-submit');
let id = "";

function openModal(modalId) {
    const modal = document.getElementById(modalId);
    modal.classList.remove('hidden');
    modal.classList.add('flex');
}

function closeModal(modalId) {
    const modal = document.getElementById(modalId);
    modal.classList.add('hidden');
    modal.classList.remove('flex');
}

function createModal(modalId){
    headerTitle.innerText = "Create";
    btnSubmit.innerText = "Create";
    openModal(modalId);
}
 
//#endregion

//#region--Role based access management
fetch('/api/rbac')
    .then(response => response.json())
        .then(roleAccess => {

            const tbody = document.getElementById('role-access-body');

            roleAccess.forEach(rbac => {
                tbody.innerHTML += `
                
                    <tr>
                        <td>${rbac.AccessID}</td>
                        <td>${rbac.RoleID}</td>
                        <td>
                          
                            <button 
                                class="delete-btn text-red-600 hover:text-red-800 font-medium"
                                data-rbacid="${rbac.AccessID}">
                                Remove
                            </button>
                        </td>
                    </tr>
                
                `;               
            });
});
//#endregion



//#region--Update and insert role access
document.getElementById('upsertform').addEventListener('submit', function (e) {

    e.preventDefault();

    const roleId = document.getElementById('roleid').value;
    const accessId = document.getElementById('accessid').value;

    const formData = {
        RoleID: roleId,
        AccessID: accessId
    };

    fetch('/api/createrc', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(formData),
    })
    .then(async response => {
        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || "Request failed");
        }

        return data;
    })
    .then(data => {
        notification("success", data.success);
        closeModal('roleAccessModal'); 
    })
    .catch(err => {
        notification("error", err.message);
    });
});
//#endregion


//#region--Delete role access function

document.getElementById('role-access-body').addEventListener('click', function (e) {
    if (!e.target.classList.contains('delete-btn')) return;

    const rbacid = e.target.dataset.rbacid;
    if (!rbacid) return;

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
                    fetch(`/api/deleterc/${rbacid}`, {
            method: 'DELETE'
        })
        .then(res => {
            if (!res.ok) throw new Error('Delete failed');

            Swal.fire({
                title: "Deleted!",
                text: "Role access has been deleted.",
                icon: "success"
            });
            rbacid.value = "";
            e.target.closest('tr').remove();
          
        })
        .catch(err => {
            console.error(err);
            Swal.fire({
                title: "Error!",
                text: "Failed to delete role access.",
                icon: "error"
            });
        });


      
    }
    });

   
});

//#endregion


//#region--Access management 
fetch('/api/access')
    .then(response => response.json())
        .then(access => {

            const tbody = document.getElementById('Access-body');

            access.forEach(ac => {
                tbody.innerHTML += `
                
                    <tr>
                        <td>${ac.accessid}</td>
                        <td>${ac.accessname}</td>    
                        <td>
                            <button class="update-btn text-blue-600 hover:text-blue-800 font-medium"
                                data-accessid = "${ac.accessid}"
                                >Edit
                            
                            </button>
                            <button 
                                class="delete-btn text-red-600 hover:text-red-800 font-medium"
                                data-accessid="${ac.accessid}">
                                Delete
                            </button>
                        </td>
                    </tr>
                
                `;               
            });
});
//#endregion

//Generate Unique Id
function uuidv4() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'
    .replace(/[xy]/g, function (c) {
        const r = Math.random() * 16 | 0, 
            v = c == 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}
//#region --Update and insert access
document.getElementById('upsertcreateform').addEventListener('submit', function (e) {

    e.preventDefault();

    let uid = "";
    //Get the accessname 
    const accessId = document.getElementById('access-id').value;
    const accessName = document.getElementById('access-name').value;
    if( id === ""){
        uid = uuidv4();
    }else{
        uid = id;
    }
    const formData = {
        accessid: uid,
        accessname: accessName
    };

    if(accessId == ""){
        //Create access
        fetch('/api/createac', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(formData),
        })
        .then(async response => {
            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error || "Request failed");
            }

            return data;
        })
        .then(data => {
            notification("success", data.success);
            id = "";
            closeModal('access-modal'); 
        })
        .catch(err => {
            notification("error", err.message);
        });
    }else{
        fetch(`/api/updateac/${id}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(formData),
        })
        .then(async response => {
            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error || "Request failed");
            }

            return data;
        })
        .then(data => {
            notification("success", data.success);
            id = "";
            closeModal('access-modal'); 
        })
        .catch(err => {
            notification("error", err.message);
        });
    }

});


//#endregion
//#region--Update function
document.getElementById('Access-body').addEventListener('click' , function(e){

    if (!e.target.classList.contains('update-btn')) return;

    id = e.target.dataset.accessid;
    if (!id) return;

    const accessInputId = document.getElementById('access-id');
    accessInputId.value = id;
    
    btnSubmit.innerText = "Update";
    headerTitle.innerText = "Update access";  
        
    fetch(`/api/access/${id}`)
        .then(async res => {
            const data = await res.json();
            if (!res.ok) throw new Error(data.error || "Request failed");
            return data;
        })
        .then(data => {
            const accessName= data.success;
            document.getElementById('access-name').value = accessName.accessname;

        })
        .catch(err => {
            console.log(err);
        });

    openModal('access-modal');



}); 
//#endregion

//#region--Delete access function

document.getElementById('Access-body').addEventListener('click', function (e) {
    if (!e.target.classList.contains('delete-btn')) return;

    const accessid = e.target.dataset.accessid;
    if (!accessid) return;

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
                    fetch(`/api/deleteac/${accessid}`, {
            method: 'DELETE'
        })
        .then(res => {
            if (!res.ok) throw new Error('Delete failed');

            Swal.fire({
                title: "Deleted!",
                text: "Access has been deleted.",
                icon: "success"
            });
            accessid.value = "";
            e.target.closest('tr').remove();
          
        })
        .catch(err => {
            Swal.fire({
                title: "Error!",
                text: "Failed to delete role access.",
                icon: "error"
            });
        });

    }
    });

   
});