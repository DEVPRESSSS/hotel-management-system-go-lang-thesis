document.addEventListener("DOMContentLoaded", () => {

 /* =====================
     FETCH ROOMS & INIT TABLE
  ====================== */
  const headerTitle = document.getElementById('header-action');
  const btnSubmit = document.getElementById('btn-submit');
  const btnText = document.getElementById('btn-text');
  const form = document.getElementById('upsertform');
  const tbody = document.getElementById('room-body');
  const tableElement = document.querySelector("#default-table");
  const userModal = document.getElementById('userModal');

  if (!tbody || !form || !headerTitle || !btnSubmit || !tableElement) {
    console.warn("Room page elements not found. JS skipped.");
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
    headerTitle.innerText = "Create Room";
    btnText.innerText = "Create";
    form.reset();
    openModal();
  };

  /* =====================
     FETCH ROOMS & INIT TABLE
  ====================== */
  fetch('/api/rooms')
    .then(res => {
      if (!res.ok) throw new Error("API failed");
      return res.json();
    })
    .then(rooms => {
      // Populate table body
      tbody.innerHTML = rooms.map(room => `
        <tr>
          <td class="px-4 py-3">${room.roomid}</td>
          <td class="px-4 py-3">${room.roomnumber}</td>
          <td class="px-4 py-3">${room.RoomType.roomtypename}</td>
          <td class="px-4 py-3">${room.Floor.floorname}</td>
          <td class="px-4 py-3">${room.capacity}</td>
          <td class="px-4 py-3">${room.price}</td>
          <td class="px-4 py-3">${new Date(room.created_at).toLocaleDateString()}</td>
          <td class="px-4 py-3">
            <button class="update-btn px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600 mr-2" data-id="${room.roomid}">Edit</button>
            <button class="delete-btn px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600" data-id="${room.roomid}">Delete</button>
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

      // =============================================
      // ATTACH EVENT LISTENER AFTER DATATABLE INIT
      // =============================================
      attachTableEventListeners();
    })
    .catch(console.error);

  /* =====================
     TABLE CLICK HANDLER - NOW A SEPARATE FUNCTION
  ====================== */
  function attachTableEventListeners() {
    // Use delegation on tableElement instead of tbody
    tableElement.addEventListener('click', e => {

      // Handle Update Button
      if (e.target.classList.contains('update-btn')) {
        id = e.target.dataset.id;
        
        fetch(`/api/room/${id}`)
          .then(res => res.json())
          .then(data => {
            const room = data.success;
            document.getElementById('roomnumber').value = room.roomnumber;
            document.getElementById('roomtypeid').value = room.roomtypeid;
            document.getElementById('floorid').value = room.floorid;
            document.getElementById('capacity').value = room.capacity;
            document.getElementById('price').value = room.price;
            
            headerTitle.innerText = "Update Room";
            btnText.innerText = "Update";
            openModal();
          })
          .catch(err => notification("error", err.message));
      }

      // Handle Delete Button
      if (e.target.classList.contains('delete-btn')) {
        const roomid = e.target.dataset.id;

        
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
            fetch(`/api/deleteroom/${roomid}`, { method: 'DELETE' })
              .then(res => {
                if (!res.ok) {
                  throw new Error('Delete failed');
                }
                return;
              })
              .then(() => {
                Swal.fire("Deleted!", "Room has been deleted.", "success");
                setTimeout(() => location.reload(), 500);
              })
              .catch(err => notification("error", err.message));
          }
        });
      }
    });
  }

  /* =====================
     FORM SUBMIT
  ====================== */
  form.addEventListener('submit', e => {
    e.preventDefault();

    const roomNumber = document.getElementById('roomnumber').value;
    const roomType = document.getElementById('roomtypeid').value;
    const floor = document.getElementById('floorid').value;
    const capacity = document.getElementById('capacity').value;
    const price = document.getElementById('price').value;

    const payload = {
      roomnumber: roomNumber,
      roomtypeid: roomType,
      floorid: floor,
      capacity: capacity,
      price: price
    };

    const url = id ? `/api/updateroom/${id}` : '/api/createroom';
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

});