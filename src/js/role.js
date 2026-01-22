document.addEventListener("DOMContentLoaded", () => {

  /* =====================
     DOM ELEMENTS
  ====================== */
  const headerTitle = document.getElementById('header-action');
  const btnSubmit = document.getElementById('btn-submit');
  const btnText = document.getElementById('btn-text');
  const form = document.getElementById('upsertform');
  const tbody = document.getElementById('users-body');
  const tableElement = document.querySelector("#default-table");
  const userModal = document.getElementById('userModal');

  if (!tbody || !form || !headerTitle || !btnSubmit || !tableElement) {
    console.warn("Role page elements not found. JS skipped.");
    console.log({tbody, form, headerTitle, btnSubmit, tableElement});
    return;
  }

  let id = "";
  let dataTable = null;

  /* =====================
     MODAL FUNCTIONS
  ====================== */
  function openModal() {
    userModal.classList.remove('hidden');
    userModal.classList.add('flex');
  }

  function closeModal() {
    userModal.classList.add('hidden');
    userModal.classList.remove('flex');
  }

  window.closeModal = closeModal;

  window.createModal = function () {
    id = "";
    headerTitle.innerText = "Create Role";
    btnText.innerText = "Create";
    form.reset();
    openModal();
  };

  /* =====================
     FETCH ROLES & INIT TABLE
  ====================== */
  fetch('/api/roles')
    .then(res => {
      if (!res.ok) throw new Error("API failed");
      return res.json();
    })
    .then(roles => {
      // Store roles in global state if needed
      if (window.AppState) {
        window.AppState.roles = roles;
      }

      // Populate table body
      tbody.innerHTML = roles.map(role => `
        <tr>
          <td class="px-4 py-3">${role.roleid}</td>
          <td class="px-4 py-3">${role.rolename}</td>
          <td class="px-4 py-3 text-center">
            <button class="update-btn px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600 mr-2" data-id="${role.roleid}">Edit</button>
            <button class="delete-btn px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600" data-id="${role.roleid}">Delete</button>
          </td>
        </tr>
      `).join("");

      // Initialize DataTable AFTER populating data
      if (window.simpleDatatables && tableElement) {
        dataTable = new simpleDatatables.DataTable(tableElement, {
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
     FORM SUBMIT
  ====================== */
  form.addEventListener('submit', e => {
    e.preventDefault();

    const roleName = document.getElementById('rolename').value;
    let uid = "";

        if (id === "" || id === null) {
            uid = uuidv4(); 
        } else {
            uid = id;      
    }
    const payload = {
      roleid: uid,
      roleName: roleName
    };

    const url = id ? `/api/updaterole/${id}` : '/api/createrole';
    const method = id ? 'PUT' : 'POST';

    fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
    .then(res => {
      return res.json().then(data => {
        // Check if request was successful
        if (!res.ok) {
          // Throw error with server message
          throw new Error(data.error || 'Request failed');
        }
        return data;
      });
    })
    .then(data => {
      notification("success", data.success);
      closeModal();
      setTimeout(() => location.reload(), 500);
    })
    .catch(err => {
      notification("error", err.message);
    });
  });

  /* =====================
     TABLE CLICK HANDLER
  ====================== */
  tbody.addEventListener('click', e => {

    // Handle Update Button
    if (e.target.classList.contains('update-btn')) {
      id = e.target.dataset.id;
      fetch(`/api/roles/${id}`)
        .then(res => res.json())
        .then(data => {
          const role = data.success;
          document.getElementById('rolename').value = role.rolename;
          
          headerTitle.innerText = "Update Role";
          btnText.innerText = "Update";
          openModal();
        })
        .catch(err => notification("error", err.message));
    }

    // Handle Delete Button
    if (e.target.classList.contains('delete-btn')) {
      const roleid = e.target.dataset.id;
      
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
          fetch(`/api/deleterole/${roleid}`, { method: 'DELETE' })
            .then(res => {
              if (!res.ok) {
                throw new Error('Delete failed');
              }
              return;
            })
            .then(() => {
              Swal.fire("Deleted!", "Role has been deleted.", "success");
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
function uuidv4() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, c => {
    const r = Math.random() * 16 | 0;
    return (c === 'x' ? r : (r & 0x3 | 0x8)).toString(16);
  });
}