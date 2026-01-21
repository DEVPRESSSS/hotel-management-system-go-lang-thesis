document.addEventListener("DOMContentLoaded", () => {

  /* =====================
     DOM ELEMENTS
  ====================== */
  const headerTitle = document.getElementById('header-action');
  const btnSubmit = document.getElementById('btn-submit');
  const form = document.getElementById('upsertform');
  const tbody = document.getElementById('users-body');
  const tableElement = document.querySelector("#default-table");

  if (!tbody || !form || !headerTitle || !btnSubmit || !tableElement) {
    console.warn("Amenity page elements not found. JS skipped.");
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
    headerTitle.innerText = "Create Amenity";
    btnSubmit.innerText = "Create";
    form.reset();
    openModal();
  };

  /* =====================
     FETCH AMENITIES & INIT TABLE
  ====================== */
  fetch('/api/aminities')
    .then(res => {
      if (!res.ok) throw new Error("API failed");
      return res.json();
    })
    .then(amenities => {
      // Populate table body
      tbody.innerHTML = amenities.map(a => `
        <tr>
          <td class="px-4 py-3">${a.amenityid}</td>
          <td class="px-4 py-3">${a.amenityname}</td>
          <td class="px-4 py-3">
            <button class="update-btn px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600 mr-2" data-id="${a.amenityid}">Edit</button>
            <button class="delete-btn px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600" data-id="${a.amenityid}">Delete</button>
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

    const name = document.getElementById('aminity').value;
    const payload = {
      amenityid: id || uuidv4(),
      amenityname: name
    };

    const url = id ? `/api/updateaminity/${id}` : '/api/createaminity';
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
      notification("error", err)
    });
  });

  /* =====================
     TABLE CLICK HANDLER
  ====================== */
  tbody.addEventListener('click', e => {

    if (e.target.classList.contains('update-btn')) {
      id = e.target.dataset.id;
      fetch(`/api/aminity/${id}`)
        .then(res => res.json())
        .then(data => {
          document.getElementById('aminity').value = data.success.amenityname;
          headerTitle.innerText = "Update Amenity";
          btnSubmit.innerText = "Update";
          openModal();
        })
        .catch(err => notification("error", err.message));
    }

    if (e.target.classList.contains('delete-btn')) {
      const aid = e.target.dataset.id;
      console.log(aid);
      Swal.fire({
        title: "Are you sure?",
        text: "This action cannot be undone!",
        icon: "warning",
        showCancelButton: true,
        confirmButtonColor: "#d33",
        cancelButtonColor: "#3085d6",
        confirmButtonText: "Yes, delete it!"
      }).then(r => {
        if (r.isConfirmed) {
          fetch(`/api/deleteamenity/${aid}`, { method: 'DELETE' })
            .then(res => {
              if (!res.ok) {
                throw new Error('Delete failed');
           
              }
               return;
            })
            .then(() => {
              Swal.fire("Deleted!", "Amenity has been deleted.", "success");
              location.reload();
            })
            .catch(err => notification("error", err.message))
        }
      });
    }
  });

});

/* =====================
   UUID GENERATOR
====================== */
function uuidv4() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, c => {
    const r = Math.random() * 16 | 0;
    return (c === 'x' ? r : (r & 0x3 | 0x8)).toString(16);
  });
}