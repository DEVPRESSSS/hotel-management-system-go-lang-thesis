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
    console.warn("Room Amenity page elements not found. JS skipped.");
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
    headerTitle.innerText = "Assign Room Amenity";
    btnText.innerText = "Assign";
    form.reset();
    openModal();
  };

  /* =====================
     LOAD ROOMS DROPDOWN
  ====================== */
  function loadRoomsDropdown() {
    fetch('/api/rooms')
      .then(res => res.json())
      .then(rooms => {
        const roomSelect = document.getElementById('roomid');
        roomSelect.innerHTML = '<option value="">-- Select Room --</option>';
        
        rooms.forEach(room => {
          const option = document.createElement('option');
          option.value = room.roomid;
          option.textContent = `Room ${room.roomnumber}`;
          roomSelect.appendChild(option);
        });
      })
      .catch(err => console.error('Failed to load rooms:', err));
  }

  /* =====================
     LOAD AMENITIES DROPDOWN
  ====================== */
  function loadAmenitiesDropdown() {
    fetch('/api/aminities')
      .then(res => res.json())
      .then(amenities => {
        const amenitySelect = document.getElementById('aminityid');
        amenitySelect.innerHTML = '<option value="">-- Select Amenity --</option>';
        
        amenities.forEach(amenity => {
          const option = document.createElement('option');
          option.value = amenity.amenityid;
          option.textContent = amenity.amenityname;
          amenitySelect.appendChild(option);
        });
      })
      .catch(err => console.error('Failed to load amenities:', err));
  }

  // Load dropdowns on page load
  loadRoomsDropdown();
  loadAmenitiesDropdown();

  /* =====================
     FETCH ROOM AMENITIES & INIT TABLE
  ====================== */
  fetch('/api/roomaminities')
    .then(res => {
      if (!res.ok) throw new Error("API failed");
      return res.json();
    })
    .then(roomAmenities => {
      // Populate table body
      tbody.innerHTML = roomAmenities.map(amenity => `
        <tr>
          <td class="px-4 py-3">${amenity.Room ? amenity.Room.roomnumber : 'N/A'}</td>
          <td class="px-4 py-3">${amenity.Amenity ? amenity.Amenity.amenityname : 'N/A'}</td>
          <td class="px-4 py-3 text-center">
            <button class="delete-btn px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600" data-roomid="${amenity.RoomId}" data-amenityid="${amenity.AmenityId}">Remove</button>
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

    const roomId = document.getElementById('roomid').value;
    const amenityId = document.getElementById('aminityid').value;

    const payload = {
      RoomId: roomId,
      AmenityId: amenityId
    };

    const url = id ? `/api/updateroomaminity/${id}` : '/api/createroomaminity';
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

    // Handle Delete Button
    if (e.target.classList.contains('delete-btn')) {
      const roomId = e.target.dataset.roomid;
      const amenityId = e.target.dataset.amenityid;
      
      Swal.fire({
        title: "Are you sure?",
        text: "This will remove the amenity from the room!",
        icon: "warning",
        showCancelButton: true,
        confirmButtonColor: "#d33",
        cancelButtonColor: "#3085d6",
        confirmButtonText: "Yes, remove it!"
      }).then(result => {
        if (result.isConfirmed) {
          // Use roomId as the identifier for deletion
          fetch(`/api/deleteroomaminity/${roomId}`, { method: 'DELETE' })
            .then(res => {
              if (!res.ok) {
                throw new Error('Delete failed');
              }
              return;
            })
            .then(() => {
              Swal.fire("Removed!", "Room amenity has been removed.", "success");
              setTimeout(() => location.reload(), 500);
            })
            .catch(err => notification("error", err.message));
        }
      });
    }
  });

});

