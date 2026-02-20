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
    console.warn("Room attendant page elements not found. JS skipped.");
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
    headerTitle.innerText = "Create Attendant";
    btnText.innerText = "Create";
    form.reset();
    openModal();
  };

  /* =====================
     FETCH ROLES & INIT TABLE
  ====================== */
  fetch('/api/maintenances')
    .then(res => {
      if (!res.ok) throw new Error("API failed");
      return res.json();
    })
    .then(maintenance => {
     

      // Populate table body
      tbody.innerHTML = maintenance.map(m => `
        <tr>
          <td class="px-4 py-3">${m.id}</td>
          <td class="px-4 py-3">${m.name}</td>
          <td class="px-4 py-3 text-sm">${new Date(m.created_at).toLocaleDateString()}</td>
          <td class="px-4 py-3 text-center">
            <button class="update-btn px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600 mr-2" data-id="${m.id}">Edit</button>
            <button class="delete-btn px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600" data-id="${m.id}">Delete</button>
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

    const name = document.getElementById('attendantname').value;
    const payload = {
      Name: name
    };

    const url = id ? `/api/maintenances/${id}` : '/api/maintenances';
    const method = id ? 'PUT' : 'POST';

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
      fetch(`/api/maintenances/${id}`)
        .then(res => res.json())
        .then(data => {
          const attendant = data.success;
          document.getElementById('attendantname').value = attendant.name;
          
          headerTitle.innerText = "Update Attendant";
          btnText.innerText = "Update";
          openModal();
        })
        .catch(err => notification("error", err.message));
    }

    // Handle Delete Button
    if (e.target.classList.contains('delete-btn')) {
      const id = e.target.dataset.id;
      
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
          fetch(`/api/maintenances/${id}`, { method: 'DELETE' })
            .then(res => {
              if (!res.ok) {
                throw new Error('Delete failed');
              }
              return;
            })
            .then(() => {
              Swal.fire("Deleted!", "Attendant has been deleted.", "success");
              setTimeout(() => location.reload(), 500);
            })
            .catch(err => notification("error", err.message));
        }
      });
    }
  });

});
