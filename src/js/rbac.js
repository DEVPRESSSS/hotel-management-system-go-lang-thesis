fetch('/api/rbac')
    .then(response => response.json())
        .then(roleAccess => {

            const tbody = document.getElementById('users-body');

            roleAccess.forEach(rbac => {
                tbody.innerHTML += `
                
                    <tr>
                        <td>${rbac.AccessID}</td>
                        <td>${rbac.Access.accessname}</td>    
                        <td>${rbac.Role.roleid}</td>
                        <td>${rbac.Role.rolename}</td>                 
                        <td>
                            <button class="update-btn text-blue-600 hover:text-blue-800 font-medium"
                                data-roleid = "${rbac.accessid}"
                                >Edit
                            
                            </button>
                            <button 
                                class="delete-btn text-red-600 hover:text-red-800 font-medium"
                                data-roleid="${rbac.accessid}">
                                Delete
                            </button>
                        </td>
                    </tr>
                
                `;               
            });
});