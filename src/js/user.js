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
  const errorMessage = document.getElementById('error-message');

  if (!tbody || !form || !headerTitle || !btnSubmit || !tableElement) {
    console.warn("User page elements not found. JS skipped.");
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
    headerTitle.innerText = "Create User";
    btnText.innerText = "Create";
    form.reset();
    errorMessage.classList.add('hidden');
    openModal();
  };

  // Close modal on outside click
  userModal.addEventListener('click', function(e) {
    if (e.target === userModal) {
      closeModal();
    }
  });

  /* =====================
     PASSWORD VALIDATION
  ====================== */
  function validatePassword() {
    const password = document.getElementById('password').value;
    const confirmPassword = document.getElementById('confirm-password').value;

    if (password !== confirmPassword) {
      errorMessage.textContent = "Passwords don't match";
      errorMessage.classList.remove("hidden");
      return false;
    }
    
    errorMessage.textContent = "";
    errorMessage.classList.add("hidden");
    return true;
  }

  // Attach validation to confirm password input
  document.getElementById('confirm-password').addEventListener('input', validatePassword);

  /* =====================
     LOAD ROLES DROPDOWN
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

  // Load roles on page load
  loadRolesDropdown();

  /* =====================
     FETCH USERS & INIT TABLE
  ====================== */
  fetch('/api/users')
    .then(res => {
      if (!res.ok) throw new Error("API failed");
      return res.json();
    })
    .then(users => {
      // Populate table body
      tbody.innerHTML = users.map(user => `
        <tr>
          <td class="px-4 py-3">${user.userid}</td>
          <td class="px-4 py-3">${user.fullname || user.username}</td>
          <td class="px-4 py-3">${user.email}</td>
          <td class="px-4 py-3">${user.Role ? user.Role.rolename : 'N/A'}</td>
          <td class="px-4 py-3">
            ${user.locked === false 
              ? '<span class="px-2 py-1 text-xs font-semibold text-green-800 bg-green-100 rounded-full">Active</span>' 
              : '<span class="px-2 py-1 text-xs font-semibold text-red-800 bg-red-100 rounded-full">Locked</span>'}
          </td>
          <td class="px-4 py-3">${new Date(user.created_at).toLocaleDateString()}</td>
          <td class="px-4 py-3 text-center">
            <button class="update-btn px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600 mr-2" data-id="${user.userid}">Edit</button>
            <button class="delete-btn px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600" data-id="${user.userid}">Delete</button>
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

    // Validate passwords match
    if (!validatePassword()) {
      notification("error", "Passwords don't match");
      return;
    }

    const fullname = document.getElementById('fullname').value;
    const email = document.getElementById('email').value;
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    const roleid = document.getElementById('roleid').value;


     let uid = "";

        if (id === "" || id === null) {
            uid = uuidv4(); 
        } else {
            uid = id;      
    }
    const payload = {
      userid: uid,
      fullname: fullname,
      email: email,
      username: username,
      password: password,
      locked: false,
      roleid: roleid
    };

    const url = id ? `/api/update/${id}` : '/userslist';
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
      notification("success", data.success || data.message);
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
      fetch(`/api/user/${id}`)
        .then(res => res.json())
        .then(data => {
          const user = data.success;
          document.getElementById('fullname').value = user.fullname || '';
          document.getElementById('email').value = user.email || '';
          document.getElementById('username').value = user.username || '';
          document.getElementById('roleid').value = user.roleid || '';
          
          // For update, make password optional by clearing required
          document.getElementById('password').removeAttribute('required');
          document.getElementById('confirm-password').removeAttribute('required');
          document.getElementById('password').value = '';
          document.getElementById('confirm-password').value = '';
          
          headerTitle.innerText = "Update User";
          btnText.innerText = "Update";
          errorMessage.classList.add('hidden');
          openModal();
        })
        .catch(err => notification("error", err.message));
    }

    // Handle Delete Button
    if (e.target.classList.contains('delete-btn')) {
      const userid = e.target.dataset.id;
      
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
          fetch(`/api/delete/${userid}`, { method: 'DELETE' })
            .then(res => {
              if (!res.ok) {
                throw new Error('Delete failed');
              }
              return;
            })
            .then(() => {
              Swal.fire("Deleted!", "User has been deleted.", "success");
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