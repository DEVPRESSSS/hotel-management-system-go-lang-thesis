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

            const tbody = document.getElementById('users-body');

            roleAccess.forEach(rbac => {
                tbody.innerHTML += `
                
                    <tr>
                        <td>${rbac.AccessID}</td>
                        <td>${rbac.Role.roleid}</td>
                        <td>
                            <button class="update-btn text-blue-600 hover:text-blue-800 font-medium"
                                data-roleid = "${rbac.AccessID}"
                                >Edit
                            
                            </button>
                            <button 
                                class="delete-btn text-red-600 hover:text-red-800 font-medium"
                                data-roleid="${rbac.AccessID}">
                                Delete
                            </button>
                        </td>
                    </tr>
                
                `;               
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
                                data-roleid = "${ac.accessid}"
                                >Edit
                            
                            </button>
                            <button 
                                class="delete-btn text-red-600 hover:text-red-800 font-medium"
                                data-roleid="${ac.accessid}">
                                Delete
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