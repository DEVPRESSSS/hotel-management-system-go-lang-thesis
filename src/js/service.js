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
    console.warn("Service page elements not found. JS skipped.");
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
    headerTitle.innerText = "Create Service";
    btnText.innerText = "Create";
    form.reset();
    openModal();
  };

  /* =====================
     FETCH SERVICES & INIT TABLE
  ====================== */
  fetch('/api/services')
    .then(res => {
      if (!res.ok) throw new Error("API failed");
      return res.json();
    })
    .then(services => {
      // Populate table body
      tbody.innerHTML = services.map(service => `
        <tr>
          <td class="px-4 py-3">${service.serviceid}</td>
          <td class="px-4 py-3">${service.servicename}</td>
          <td class="px-4 py-3">${formatTime(service.start_time)}</td>
          <td class="px-4 py-3">${formatTime(service.end_time)}</td>
          <td class="px-4 py-3">${new Date(service.created_at).toLocaleDateString()}</td>
          <td class="px-4 py-3 text-center">
            <button class="update-btn px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600 mr-2" data-id="${service.serviceid}">Edit</button>
            <button class="delete-btn px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600" data-id="${service.serviceid}">Delete</button>
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
     HELPER: Format time for input[type="time"]
  ====================== */
  function formatTimeForInput(dateString) {
    if (!dateString) return '';
    const date = new Date(dateString);
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    return `${hours}:${minutes}`;
  }

  /* =====================
     FORM SUBMIT
  ====================== */
  form.addEventListener('submit', e => {
    e.preventDefault();

    const serviceName = document.getElementById('servicename').value;
    const startTime = document.getElementById('start-time').value;
    const endTime = document.getElementById('end-time').value;

    const payload = {
      servicename: serviceName,
      start_time: startTime,
      end_time: endTime
    };

    const url = id ? `/api/updateservice/${id}` : '/api/createservice';
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
      fetch(`/api/service/${id}`)
        .then(res => res.json())
        .then(data => {
          const service = data.success;
          document.getElementById('servicename').value = service.servicename;
          document.getElementById('start-time').value = formatTimeForInput(service.start_time);
          document.getElementById('end-time').value = formatTimeForInput(service.end_time);
          
          headerTitle.innerText = "Update Service";
          btnText.innerText = "Update";
          openModal();
        })
        .catch(err => notification("error", err.message));
    }

    // Handle Delete Button
    if (e.target.classList.contains('delete-btn')) {
      const serviceid = e.target.dataset.id;
      
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
          fetch(`/api/deleteservice/${serviceid}`, { method: 'DELETE' })
            .then(res => {
              if (!res.ok) {
                throw new Error('Delete failed');
              }
              return;
            })
            .then(() => {
              Swal.fire("Deleted!", "Service has been deleted.", "success");
              setTimeout(() => location.reload(), 500);
            })
            .catch(err => notification("error", err.message));
        }
      });
    }
  });

});

function formatTime(time) {
  if (!time) return '';

  const [hours, minutes] = time.split(':');
  const date = new Date();
  date.setHours(hours, minutes);

  return date.toLocaleTimeString('en-US', {
    hour: 'numeric',
    minute: '2-digit',
    hour12: true
  });
}
