//#region--Modal navigations
const headerTitle = document.getElementById('header-action');
const btnSubmit = document.getElementById('btn-submit');
let id = "";
function openModal() {
    document.getElementById('roleAccessModal').classList.remove('hidden');
    document.getElementById('roleAccessModal').classList.add('flex');
}

function closeModal() {
        document.getElementById('roleAccessModal').classList.add('hidden');
        document.getElementById('roleAccessModal').classList.remove('flex');
}
function createModal(){
    headerTitle.innerText = "Create user";
    btnSubmit.innerText = "Create";
    openModal();
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
                                data-rbacid="${rbac.RoleID}">
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
        closeModal(); 
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

            e.target.closest('tr').remove();
            rbacid.value = "";
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

//#region --CRUD Access Functions
// document.getElementById('upsertform').addEventListener('submit', function (e) {

//     e.preventDefault();

//     const accessId = document.getElementById('access-name').value;

//     const formData = {
//         RoleID: roleId,
//         AccessID: accessId
//     };

//     fetch('/api/createrc', {
//         method: 'POST',
//         headers: {
//             'Content-Type': 'application/json'
//         },
//         body: JSON.stringify(formData),
//     })
//     .then(async response => {
//         const data = await response.json();

//         if (!response.ok) {
//             throw new Error(data.error || "Request failed");
//         }

//         return data;
//     })
//     .then(data => {
//         notification("success", data.success);
//         closeModal(); 
//     })
//     .catch(err => {
//         notification("error", err.message);
//     });
// });


//#endregion