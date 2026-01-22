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
    console.warn("Facility page elements not found. JS skipped.");
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
    headerTitle.innerText = "Create Facility";
    btnText.innerText = "Create";
    form.reset();
    openModal();
  };

  /* =====================
     HELPER: Format date for input[type="date"]
  ====================== */
  function formatDateForInput(dateString) {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toISOString().split('T')[0];
  }

  /* =====================
     HELPER: Safe date display
  ====================== */
  function formatDateDisplay(dateString) {
    if (!dateString || dateString === 'Invalid Date') return 'N/A';
    try {
      return new Date(dateString).toLocaleDateString();
    } catch {
      return 'N/A';
    }
  }

  /* =====================
     FETCH FACILITIES & INIT TABLE
  ====================== */
  fetch('/api/facility')
    .then(res => {
      if (!res.ok) throw new Error("API failed");
      return res.json();
    })
    .then(facilities => {
      // Populate table body
      tbody.innerHTML = facilities.map(facility => `
        <tr>
          <td class="px-4 py-3">${facility.facility_id}</td>
          <td class="px-4 py-3">${facility.facility_name}</td>
          <td class="px-4 py-3">${facility.status}</td>
          <td class="px-4 py-3">${formatDateDisplay(facility.maintenance_date)}</td>
          <td class="px-4 py-3">${formatDateDisplay(facility.created_at)}</td>
          <td class="px-4 py-3 text-center">
            <button class="update-btn px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600 mr-2" data-id="${facility.facility_id}">Edit</button>
            <button class="delete-btn px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600" data-id="${facility.facility_id}">Delete</button>
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

    const facilityName = document.getElementById('facilityname').value;
    const dateInput = document.getElementById('maintenance_date').value;

    // Convert date to ISO format if provided, otherwise null
    let isoDate = null;
    if (dateInput) {
      isoDate = new Date(dateInput + "T00:00:00").toISOString();
    }

    const payload = {
      facility_name: facilityName,
      maintenance_date: isoDate
    };

    const url = id ? `/api/updatefacility/${id}` : '/api/createfacility';
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
      fetch(`/api/facility/${id}`)
        .then(res => res.json())
        .then(data => {
          const facility = data.success;
          document.getElementById('facilityname').value = facility.facility_name;
          document.getElementById('maintenance_date').value = formatDateForInput(facility.maintenance_date);
          
          headerTitle.innerText = "Update Facility";
          btnText.innerText = "Update";
          openModal();
        })
        .catch(err => notification("error", err.message));
    }

    // Handle Delete Button
    if (e.target.classList.contains('delete-btn')) {
      const facilityId = e.target.dataset.id;
      
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
          fetch(`/api/deletefacility/${facilityId}`, { method: 'DELETE' })
            .then(res => {
              if (!res.ok) {
                throw new Error('Delete failed');
              }
              return;
            })
            .then(() => {
              Swal.fire("Deleted!", "Facility has been deleted.", "success");
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