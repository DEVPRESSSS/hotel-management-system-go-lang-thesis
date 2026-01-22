document.addEventListener("DOMContentLoaded", () => {

  /* =====================
     DOM ELEMENTS
  ====================== */
  // Role Access Elements
  const roleAccessModal = document.getElementById('roleAccessModal');
  const roleAccessHeader = document.getElementById('role-access-header');
  const roleAccessBtnText = document.getElementById('role-access-btn-text');
  const roleAccessForm = document.getElementById('role-access-form');
  const roleAccessBody = document.getElementById('role-access-body');
  const roleAccessTable = document.querySelector("#role-access-table");

  // Access Elements
  const accessModal = document.getElementById('accessModal');
  const accessHeader = document.getElementById('access-header');
  const accessBtnText = document.getElementById('access-btn-text');
  const accessForm = document.getElementById('access-form');
  const accessBody = document.getElementById('access-body');
  const accessTable = document.querySelector("#access-table");

  let accessId = "";
  let roleAccessDataTable = null;
  let accessDataTable = null;

  /* =====================
     MODAL FUNCTIONS - ROLE ACCESS
  ====================== */
  function openRoleAccessModal() {
    roleAccessModal.classList.remove('hidden');
    roleAccessModal.classList.add('flex');
  }

  function closeRoleAccessModal() {
    roleAccessModal.classList.add('hidden');
    roleAccessModal.classList.remove('flex');
  }

  window.closeRoleAccessModal = closeRoleAccessModal;

  window.createRoleAccessModal = function () {
    roleAccessHeader.innerText = "Create Role Access";
    roleAccessBtnText.innerText = "Create";
    roleAccessForm.reset();
    openRoleAccessModal();
  };

  /* =====================
     MODAL FUNCTIONS - ACCESS
  ====================== */
  function openAccessModal() {
    accessModal.classList.remove('hidden');
    accessModal.classList.add('flex');
  }

  function closeAccessModal() {
    accessModal.classList.add('hidden');
    accessModal.classList.remove('flex');
  }

  window.closeAccessModal = closeAccessModal;

  window.createAccessModal = function () {
    accessId = "";
    accessHeader.innerText = "Create Access";
    accessBtnText.innerText = "Create";
    accessForm.reset();
    openAccessModal();
  };

  /* =====================
     LOAD ROLES & ACCESS FOR DROPDOWNS
  ====================== */
  function loadRolesDropdown() {
    fetch('/api/roles')
      .then(res => res.json())
      .then(roles => {
        const roleSelect = document.getElementById('roleid');
        roleSelect.innerHTML = '<option value="">-- Select Role --</option>';
        
        roles.forEach(role => {
          const option = document.createElement('option');
          option.value = role.roleid;
          option.textContent = role.rolename;
          roleSelect.appendChild(option);
        });
      })
      .catch(err => console.error('Failed to load roles:', err));
  }

  function loadAccessDropdown() {
    fetch('/api/access')
      .then(res => res.json())
      .then(accessList => {
        const accessSelect = document.getElementById('accessid');
        accessSelect.innerHTML = '<option value="">-- Select Access --</option>';
        
        accessList.forEach(access => {
          const option = document.createElement('option');
          option.value = access.accessid;
          option.textContent = access.accessname;
          accessSelect.appendChild(option);
        });
      })
      .catch(err => console.error('Failed to load access:', err));
  }

  // Load dropdowns on page load
  loadRolesDropdown();
  loadAccessDropdown();

  /* =====================
     FETCH ROLE ACCESS & INIT TABLE
  ====================== */
  fetch('/api/rbac')
    .then(res => {
      if (!res.ok) throw new Error("API failed");
      return res.json();
    })
    .then(roleAccessList => {
      // Populate table body
      roleAccessBody.innerHTML = roleAccessList.map(rbac => `
      <tr>
        <td class="px-4 py-3">${rbac.AccessID}</td>
        <td class="px-4 py-3">${rbac.RoleID}</td>
        <td class="px-4 py-3 text-center">
          ${
            rbac.Role.rolename !== "Admin"
              ? `<button class="delete-rbac-btn px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600"
                  data-id="${rbac.AccessID}">
                  Remove permission
                </button>`
              : `<p class="text-sm text-success">
                   <span
                          class="px-2 py-1 font-semibold leading-tight text-green-700 bg-green-100 rounded-full dark:bg-green-700 dark:text-green-100"
                        >
                          Default Admin Permission
                    </span>
                </p>`
          }
        </td>
      </tr>

      `).join("");

      // Initialize DataTable AFTER populating data
      if (window.simpleDatatables && roleAccessTable) {
        roleAccessDataTable = new simpleDatatables.DataTable(roleAccessTable, {
          searchable: true,
          paging: true,
          perPage: 10,
          perPageSelect: [5, 10, 20, 50],
          sortable: true,
        });
      }
    })
    .catch(console.error);

  /* =====================
     FETCH ACCESS & INIT TABLE
  ====================== */
  fetch('/api/access')
    .then(res => {
      if (!res.ok) throw new Error("API failed");
      return res.json();
    })
    .then(accessList => {
      // Populate table body
      accessBody.innerHTML = accessList.map(access => `
        <tr>
          <td class="px-4 py-3">${access.accessid}</td>
          <td class="px-4 py-3">
             <span class="px-2 py-1 font-semibold leading-tight text-green-700 bg-green-100 rounded-full dark:bg-green-700 dark:text-green-100">
                          ${access.accessname}
              </span>
                    
          </td>
          <td class="px-4 py-3 text-center">
            <button class="update-access-btn px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600 mr-2" data-id="${access.accessid}">Edit</button>
            <button class="delete-access-btn px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600" data-id="${access.accessid}">Delete</button>
          </td>
        </tr>
      `).join("");

      // Initialize DataTable AFTER populating data
      if (window.simpleDatatables && accessTable) {
        accessDataTable = new simpleDatatables.DataTable(accessTable, {
          searchable: true,
          paging: true,
          perPage: 10,
          perPageSelect: [5, 10, 20, 50],
          sortable: true,
        });
      }
    })
    .catch(console.error);

  /* =====================
     ROLE ACCESS FORM SUBMIT
  ====================== */
  roleAccessForm.addEventListener('submit', e => {
    e.preventDefault();

    const roleId = document.getElementById('roleid').value;
    const accessIdValue = document.getElementById('accessid').value;

    const payload = {
      RoleID: roleId,
      AccessID: accessIdValue
    };

    fetch('/api/createrc', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
    .then(res => {
      return res.json().then(data => {
        if (!res.ok) {
          throw new Error(data.error || 'Request failed');
        }
        return data;
      });
    })
    .then(data => {
      notification("success", data.success);
      closeRoleAccessModal();
      setTimeout(() => location.reload(), 500);
    })
    .catch(err => {
      notification("error", err.message);
    });
  });

  /* =====================
     ACCESS FORM SUBMIT
  ====================== */
  accessForm.addEventListener('submit', e => {
    e.preventDefault();

    const accessName = document.getElementById('access-name').value;

    const payload = {
      accessname: accessName
    };

    const url = accessId ? `/api/updateac/${accessId}` : '/api/createac';
    const method = accessId ? 'PUT' : 'POST';

    fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
    .then(res => {
      return res.json().then(data => {
        if (!res.ok) {
          throw new Error(data.error || 'Request failed');
        }
        return data;
      });
    })
    .then(data => {
      notification("success", data.success);
      closeAccessModal();
      setTimeout(() => location.reload(), 500);
    })
    .catch(err => {
      notification("error", err.message);
    });
  });

  /* =====================
     ROLE ACCESS TABLE CLICK HANDLER
  ====================== */
  roleAccessBody.addEventListener('click', e => {

    // Handle Delete Button
    if (e.target.classList.contains('delete-rbac-btn')) {
      const rbacId = e.target.dataset.id;
      
      Swal.fire({
        title: "Are you sure?",
        text: "This action cannot be undone!",
        icon: "warning",
        showCancelButton: true,
        confirmButtonColor: "#d33",
        cancelButtonColor: "#3085d6",
        confirmButtonText: "Yes, remove it!"
      }).then(result => {
        if (result.isConfirmed) {
          fetch(`/api/deleterc/${rbacId}`, { method: 'DELETE' })
            .then(res => {
              if (!res.ok) {
                throw new Error('Delete failed');
              }
              return;
            })
            .then(() => {
              Swal.fire("Deleted!", "Role access has been removed.", "success");
              setTimeout(() => location.reload(), 500);
            })
            .catch(err => notification("error", err.message));
        }
      });
    }
  });

  /* =====================
     ACCESS TABLE CLICK HANDLER
  ====================== */
  accessBody.addEventListener('click', e => {

    // Handle Update Button
    if (e.target.classList.contains('update-access-btn')) {
      accessId = e.target.dataset.id;
      
      fetch(`/api/access/${accessId}`)
        .then(res => res.json())
        .then(data => {
          const access = data.success;
          document.getElementById('access-name').value = access.accessname;
          
          accessHeader.innerText = "Update Access";
          accessBtnText.innerText = "Update";
          openAccessModal();
        })
        .catch(err => notification("error", err.message));
    }

    // Handle Delete Button
    if (e.target.classList.contains('delete-access-btn')) {
      const deleteAccessId = e.target.dataset.id;
      
      Swal.fire({
        title: "Are you sure?",
        text: "This action cannot be undone!",
        icon: "warning",
        showCancelButton: true,
        confirmButtonColor: "#d33",
        cancelButtonColor: "#3085d6",
        confirmButtonText: "Yes, delete it!"
      }).then(result => {
        if (result.isConfirmed) {
          fetch(`/api/deleteac/${deleteAccessId}`, { method: 'DELETE' })
            .then(res => {
              if (!res.ok) {
                throw new Error('Delete failed');
              }
              return;
            })
            .then(() => {
              Swal.fire("Deleted!", "Access has been deleted.", "success");
              setTimeout(() => location.reload(), 500);
            })
            .catch(err => notification("error", err.message));
        }
      });
    }
  });

});

/* =====================
   UUID GENERATOR (if needed)
====================== */
// function uuidv4() {
//   return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, c => {
//     const r = Math.random() * 16 | 0;
//     return (c === 'x' ? r : (r & 0x3 | 0x8)).toString(16);
//   });
// }