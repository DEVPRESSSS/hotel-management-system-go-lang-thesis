document.addEventListener("DOMContentLoaded", () => {
  /* =====================
     DOM ELEMENTS
  ====================== */

  /* =====================
     DOM ELEMENTS
  ====================== */
  const headerTitle = document.getElementById('header-action');
  const btnSubmit = document.getElementById('btn-submit');
  const btnText = document.getElementById('btn-text');
  const form = document.getElementById('upsertform');
  const userModal = document.getElementById('cleanerModal');

  if ( !form || !headerTitle || !btnSubmit ) {
    console.warn("Facility page elements not found. JS skipped.");
    return;
  }


  /* =====================
     MODAL FUNCTIONS
  ====================== */
  function openModal() {
    userModal.classList.remove('hidden');
    userModal.classList.add('flex');

    fetchAttendants();
  }

  function closeModal() {
    userModal.classList.add('hidden');
    userModal.classList.remove('flex');

    document.getElementById('upsertform').reset();
    updateSelectedCount();
  }

  window.closeModal = closeModal;

  window.createModal = function () {
    id = "";
    headerTitle.innerText = "Create Facility";
    btnText.innerText = "Create";
    form.reset();
    openModal();
  };

  const tbody = document.getElementById('reservation');
  const tableElement = document.querySelector("#default-table");

  if (!tbody || !tableElement) {
    console.warn("Role page elements not found. JS skipped.");
    return;
  }

  let dataTable = null;
  let reservationsData = []; // Store data for reference

  /* =====================
     FETCH ROLES & INIT TABLE
  ====================== */
  fetch('/api/reservations')
    .then(res => {
      if (!res.ok) throw new Error("API failed");
      return res.json();
    })
    .then(data => {
      reservationsData = data.reservations; // Store for later use
      
      // Populate table body
      tbody.innerHTML = reservationsData.map((r, index) => `
        <tr class="text-gray-700 dark:text-gray-400" data-index="${index}">
          <td class="px-4 py-3">
            <div class="flex items-center text-sm">
              <div class="relative hidden w-8 h-8 mr-3 rounded-full md:block">
                <img
                  class="object-cover w-full h-full rounded-full"
                  src="https://images.unsplash.com/flagged/photo-1570612861542-284f4c12e75f?ixlib=rb-1.2.1&q=80&fm=jpg&crop=entropy&cs=tinysrgb&w=200&fit=max&ixid=eyJhcHBfaWQiOjE3Nzg0fQ"
                  alt=""
                  loading="lazy"
                />
                <div class="absolute inset-0 rounded-full shadow-inner" aria-hidden="true"></div>
              </div>
              <div>
                <p class="font-semibold">${r.User.fullname}</p>
              </div>
            </div>
          </td>
          <td class="px-4 py-3 text-sm">${r.book_id}</td>
          <td class="px-4 py-3 text-sm">${r.room_number}</td>
          <td class="px-4 py-3 text-sm">${r.room_type}</td>
          <td class="px-4 py-3 text-sm">
            <span class="px-2 py-1 font-semibold leading-tight text-yellow-800 bg-yellow-100 rounded-full">                              
              ${new Date(r.check_in_date).toLocaleString('en-CA', {
                year: 'numeric',
                month: '2-digit',
                day: '2-digit',
                hour: '2-digit',
                minute: '2-digit',
                hour12: true
              }).replace(',', '')}
            </span>
          </td>
          <td class="px-4 py-3 text-sm">
            <span class="px-2 py-1 font-semibold leading-tight text-red-800 bg-red-100 rounded-full">
              ${new Date(r.check_out_date).toLocaleString('en-CA', {
                year: 'numeric',
                month: '2-digit',
                day: '2-digit',
                hour: '2-digit',
                minute: '2-digit',
                hour12: true
              }).replace(',', '')}
            </span>
          </td>
          <td class="px-4 py-3 text-sm">${r.num_guests}</td>
          <td class="px-4 py-3 text-sm">${r.price_per_night}</td>
          <td class="px-4 py-3 text-sm">${r.total_price}</td>
          <td class="px-4 py-3 text-xs">
            ${
              r.payment_status === "Paid"
                ? `<span class="px-2 py-1 font-semibold leading-tight text-green-700 bg-green-100 rounded-full dark:bg-green-700 dark:text-green-100">${r.payment_status}</span>`
                : `<span class="px-2 py-1 font-semibold leading-tight text-orange-700 bg-orange-100 rounded-full dark:text-white dark:bg-orange-600">${r.payment_status}</span>`
            }
          </td>
          <td class="px-4 py-3 text-sm">${r.status}</td>
          <td class="px-4 py-3 text-sm">${new Date(r.created_at).toLocaleDateString()}</td>
          <td class="px-4 py-3 text-sm">
            ${r.status === "check-out" ? 
              `<button class="action-btn px-3 py-1 bg-green-500 text-white rounded hover:bg-green-600 mr-2" 
                data-room-id="${r.room_id}" data-id="${r.book_id}" data-action="clean">Clean</button>` 
              : ''
            }
            ${(() => {
              const today = new Date();
              const checkIn = new Date(r.check_in_date);
              const isSameDay = today.toDateString() === checkIn.toDateString();
              const timeHasPassed = today.getTime() >= checkIn.getTime();
              return (isSameDay && timeHasPassed && r.status ==="pending") ? 
                `<button class="action-btn px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600 mr-2" data-action="checkin">Checkin</button>` 
                : '';
            })()}
            ${(() => {
              const now = new Date();
              const checkOut = new Date(r.check_out_date);
              const isSameDay = now.toDateString() === checkOut.toDateString();
              const twoHoursBefore = 2 * 60 * 60 * 1000; 
              const withinTimeWindow = now.getTime() >= (checkOut.getTime() - twoHoursBefore);
              return (isSameDay && withinTimeWindow && r.status === "check-in") ? 
                `<button class="action-btn px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600 mr-2" data-action="checkout">Checkout</button>` 
                : '';
            })()}
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

      // Attach event listener AFTER DataTable
      attachEventListeners();
    })
    .catch(console.error);

  /* =====================
     EVENT HANDLERS
  ====================== */
  function attachEventListeners() {
    // Use event delegation on table element (not tbody, as DataTable might replace it)
    tableElement.addEventListener("click", handleButtonClick, true);
  }


  let bookingId = "";
  let roomId = "";
  function handleButtonClick(e) {
    const btn = e.target;
    
    // Check if it's an action button
    if (!btn.classList.contains('action-btn')) return;

    e.preventDefault();
    e.stopPropagation();

    // Find the parent row
    const row = btn.closest('tr');
    if (!row) {
      console.error("Could not find parent row");
      return;
    }

    // Get the index from the row
    const index = parseInt(row.dataset.index);
    if (isNaN(index) || !reservationsData[index]) {
      console.error("Invalid row index:", row.dataset.index);
      return;
    }

    const reservation = reservationsData[index];
    const action = btn.dataset.action;

    if (!action) {
      console.error("Missing action on button");
      return;
    }

     bookingId = btn.dataset.id;
     roomId = btn.dataset.roomId;
     console.log(roomId);

    handleAction(action, reservation, btn);
  }

  function handleAction(action, reservation, buttonElement) {
    // Disable button to prevent double-clicks
    buttonElement.disabled = true;
    

    switch(action) {
      case 'clean':
        handleClean(reservation, buttonElement);
        break;
      case 'checkin':
        handleCheckin(reservation, buttonElement);
        break;
      case 'checkout':
        handleCheckout(reservation, buttonElement);
        break;
      default:
        console.warn(`Unknown action: ${action}`);
        buttonElement.disabled = false;
    }
  }

  function handleClean(reservation, btn) {
      openModal();
  }

  function handleCheckin(reservation, btn) {

      Swal.fire({
      title: "Are you sure the guest want to checkin?",
      text: "This action cannot be undone!",
      icon: "warning",
      showCancelButton: true,
      confirmButtonColor: "#d33",
      cancelButtonColor: "#3085d6",
      confirmButtonText: "Yes, mark as checkin"
    }).then(result => {
      if (result.isConfirmed) {

          fetch(`/api/reservations/checkin/${reservation.book_id}`, { 
            method: 'POST',
            headers: { 'Content-Type': 'application/json' }
          })
            .then(res => res.json())
            .then(data => {
              alert(`Guest checked in successfully!`);
            })
            .catch(err => {
              alert('Failed to check in guest');
              btn.disabled = false;
            });
            }
    });
    
   
  }

  function handleCheckout(reservation, btn) {
    console.log(`Checking out reservation ${reservation.id}`);
    
    // fetch(`/api/reservations/${reservation.id}/checkout`, { 
    //   method: 'POST',
    //   headers: { 'Content-Type': 'application/json' }
    // })
    //   .then(res => res.json())
    //   .then(data => {
    //     alert(`Guest ${reservation.User.fullname} checked out successfully!`);
    //     location.reload();
    //   })
    //   .catch(err => {
    //     console.error('Checkout failed:', err);
    //     alert('Failed to check out guest');
    //     btn.disabled = false;
    //   });
  }

  /* =====================
    FETCH ATTENDANTS FOR MODAL
  ====================== */
  function fetchAttendants() {
    const container = document.getElementById('attendant-checkboxes');
    
    // Show loading state
    container.innerHTML = '<p class="text-sm text-gray-500 text-center py-4">Loading attendants...</p>';
    
    fetch('/api/maintenances')
      .then(res => {
        if (!res.ok) throw new Error("Failed to fetch attendants");
        return res.json();
      })
      .then(maintenance => {
        populateAttendantCheckboxes(maintenance);
      })
      .catch(error => {
        console.error(error);
        container.innerHTML = '<p class="text-sm text-red-500 text-center py-4">Error loading attendants</p>';
      });
  }

  /* =====================
    POPULATE CHECKBOXES
  ====================== */
  function populateAttendantCheckboxes(maintenanceList) {
    const container = document.getElementById('attendant-checkboxes');
    if (!container) return;
    
    container.innerHTML = ''; // Clear existing
    
    if (maintenanceList.length === 0) {
      container.innerHTML = '<p class="text-sm text-gray-500 text-center py-4">No attendants available</p>';
      return;
    }
    
    maintenanceList.forEach((attendant) => {
      const label = document.createElement('label');
      label.className = 'flex items-center space-x-3 p-2 hover:bg-gray-50 rounded-lg cursor-pointer transition-colors';
      
      // REMOVED onchange attribute
      label.innerHTML = `
        <input type="checkbox" 
              name="attendants[]" 
              value="${attendant.id}" 
              data-name="${attendant.name}"
              id="attendant-${attendant.id}"
              class="w-4 h-4 text-purple-600 border-gray-300 rounded focus:ring-purple-500 focus:ring-2">
        <span class="text-sm text-gray-700">${attendant.name}</span>
      `;
      
      container.appendChild(label);
    });
    
    // Initialize count to 0
    updateSelectedCount();
  }

  /* =====================
    EVENT DELEGATION FOR CHECKBOXES
  ====================== */
  document.addEventListener('DOMContentLoaded', function() {
    const container = document.getElementById('attendant-checkboxes');
    
    if (container) {
      container.addEventListener('change', function(e) {
        if (e.target.name === 'attendants[]') {
          updateSelectedCount();
        }
      });
    }
  });

  /* =====================
    UPDATE SELECTED COUNT
  ====================== */
  function updateSelectedCount() {
    const checkboxes = document.querySelectorAll('input[name="attendants[]"]:checked');
    const countElement = document.getElementById('selected-count');
    if (countElement) {
      countElement.textContent = checkboxes.length;
    }
  }

  /* =====================
    GET SELECTED ATTENDANTS
  ====================== */
  function getSelectedAttendants() {
    const checkboxes = document.querySelectorAll('input[name="attendants[]"]:checked');
    return Array.from(checkboxes).map(cb => ({
      id: cb.value,
      name: cb.dataset.name
    }));
  }

  /* =====================
    FORM SUBMISSION
  ====================== */
  document.getElementById('upsertform').addEventListener('submit', function(e) {
    e.preventDefault();
    
    const selectedAttendants = getSelectedAttendants();
    
    if (selectedAttendants.length === 0) {
      alert('Please select at least one attendant');
      return;
    }
    
    console.log('Selected attendants:', selectedAttendants);
    console.log('Selected IDs:', selectedAttendants.map(a => a.id));
    
    const formData = {
      cleaner_id :selectedAttendants.map(a => a.id),
      room_id: roomId
    };
    fetch(`/api/reservations/assigncleaner/${bookingId}`,{
        method:'POST',
        credentials:'include',
        body: formData


    }).then(response =>{

       return response.json().then(data => {
        // Check if request was successful
        if (!response.ok) {
          // Throw error with server message
          throw new Error(data.error || 'Request failed');
        }
        return data;
      })
    }).then(data =>{

        alert(`Successfully assign a cleaner!!!!!`);
    }).catch(error =>{

        alert(`Failed to create!`);

    });
        
  });
});