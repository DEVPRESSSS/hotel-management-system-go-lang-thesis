document.addEventListener("DOMContentLoaded", () => {

  /* =====================================================================
     DOM ELEMENTS
  ===================================================================== */
  const headerTitle = document.getElementById('header-action');
  const btnSubmit   = document.getElementById('btn-submit');
  const btnText     = document.getElementById('btn-text');
  const form        = document.getElementById('upsertform');
  const userModal   = document.getElementById('cleanerModal');

  const reservationModal = document.getElementById('reservationModal');

  if (!form || !headerTitle || !btnSubmit) {
    console.warn("Facility page elements not found. JS skipped.");
    return;
  }

  const tbody        = document.getElementById('reservation');
  const tableElement = document.querySelector("#default-table");

  if (!tbody || !tableElement) {
    console.warn("Role page elements not found. JS skipped.");
    return;
  }

  let dataTable        = null;
  let reservationsData = [];
  let bookingId        = "";
  let roomId           = "";


  /* =====================================================================
     MODAL — CLEANER
  ===================================================================== */
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
    bookingId = "";
    headerTitle.innerText = "Create Facility";
    btnText.innerText     = "Create";
    form.reset();
    openModal();
  };


  /* =====================================================================
     MODAL — RESERVATION (CHECK-IN)
  ===================================================================== */
  function openReservationModal() {
    reservationModal.classList.remove('hidden');
    reservationModal.classList.add('flex');
  }

  function closeReservationModal() {
    reservationModal.classList.add('hidden');
    reservationModal.classList.remove('flex');
  }

  window.openReservationModal  = openReservationModal;
  window.closeReservationModal = closeReservationModal;

  window.populateReservationModal = function (index) {
    const r = reservationsData[index];
    if (!r) return;

    // Store for confirm action
    bookingId = r.book_id;
    roomId    = r.room_id;

    // Calculate duration in nights
    const checkIn  = new Date(r.check_in_date);
    const checkOut = new Date(r.check_out_date);
    const nights   = Math.round((checkOut - checkIn) / (1000 * 60 * 60 * 24));

    // Guest info
    document.getElementById('name').textContent       = r.name          || '—';
    document.getElementById('email').textContent      = r.email         || '—';
    document.getElementById('contact-no').textContent    = r.contact       || '—';

    // Booking details
    document.getElementById('booking-no').textContent    = r.book_id;
    document.getElementById('check-in-date').textContent      = formatDateTime(r.check_in_date);
    document.getElementById('check-out-date').textContent     = formatDateTime(r.check_out_date);
    document.getElementById('duration').textContent      = `${nights} night${nights !== 1 ? 's' : ''}`;
    document.getElementById('guest-count').textContent   = `${r.num_guests} guest${r.num_guests !== 1 ? 's' : ''}`;
    document.getElementById('amount').textContent         = `₱${Number(r.total_price).toLocaleString()}`;
    document.getElementById('payment-status').textContent = r.payment_status || '—';

    //Populate all the booking guest
    
    const guestRows = document.getElementById('guest-rows');
    guestRows.innerHTML = "";
    if(r.guests && r.guests.length > 0){
         r.guests.forEach(g =>{
          guestRows.innerHTML += `
            <div class="grid grid-cols-3 text-sm text-gray-700 py-2 border-b border-gray-100">
              <span>${g.lastname}, ${g.firstname}</span>
              <span>${g.phonenumber}</span>
              <span class="text-right">${g.guest_number}</span>
            </div>
          `;
      });
    } else {
      guestRows.innerHTML = `
        <div class="py-4 text-sm text-gray-400 text-center">
          No additional guests
        </div>
      `;
    }
 
    // Special requests
    const specialReqEl = document.getElementById('special-request');
    if (specialReqEl) {
      specialReqEl.textContent = r.special_requests?.trim() || 'None';
    }

    openReservationModal();

  };

  
  window.confirmCheckin = function () {
    const reservation = reservationsData.find(r => r.book_id === bookingId);
    if (!reservation) return;
    //closeReservationModal();
    handleCheckin(reservation, document.getElementById('btn-confirm-checkin'));
  };


  /* =====================================================================
     FETCH RESERVATIONS & INIT TABLE
  ===================================================================== */
  fetch('/api/reservations')
    .then(res => {
      if (!res.ok) throw new Error("API failed");
      return res.json();
    })
    .then(data => {
      reservationsData = data.reservations;
      tbody.innerHTML  = reservationsData.map((r, index) => buildRow(r, index)).join("");

      if (window.simpleDatatables && tableElement) {
        dataTable = new simpleDatatables.DataTable(tableElement, {
          searchable:    true,
          paging:        true,
          perPage:       10,
          perPageSelect: [5, 10, 20, 50],
          sortable:      true,
        });
      }

      attachEventListeners();
    })
    .catch(console.error);


  /* =====================================================================
     BUILD TABLE ROW
  ===================================================================== */
  function buildRow(r, index) {
    return `
      <tr class="text-gray-700 dark:text-gray-400" data-index="${index}">

        <!-- Guest Name -->
        <td class="px-4 py-3">
          <div class="flex items-center text-sm">
            <div class="relative hidden w-8 h-8 mr-3 rounded-full md:block">
              <img
                class="object-cover w-full h-full rounded-full"
                src="https://images.unsplash.com/flagged/photo-1570612861542-284f4c12e75f?ixlib=rb-1.2.1&q=80&fm=jpg&crop=entropy&cs=tinysrgb&w=200&fit=max&ixid=eyJhcHBfaWQiOjE3Nzg0fQ"
                alt="" loading="lazy"
              />
              <div class="absolute inset-0 rounded-full shadow-inner" aria-hidden="true"></div>
            </div>
            <p class="font-semibold">${r.name}</p>
          </div>
        </td>

        <!-- Booking Info -->
        <td class="px-4 py-3 text-sm">${r.book_id}</td>
        <td class="px-4 py-3 text-sm">${r.room_number}</td>
        <td class="px-4 py-3 text-sm">${r.room_type}</td>

        <!-- Check-in Date -->
        <td class="px-4 py-3 text-sm">
          <span class="px-2 py-1 font-semibold leading-tight text-yellow-800 bg-yellow-100 rounded-full">
            ${formatDateTime(r.check_in_date)}
          </span>
        </td>

        <!-- Check-out Date -->
        <td class="px-4 py-3 text-sm">
          <span class="px-2 py-1 font-semibold leading-tight text-red-800 bg-red-100 rounded-full">
            ${formatDateTime(r.check_out_date)}
          </span>
        </td>

        <!-- Payment Status -->
        <td class="px-4 py-3 text-xs">
          ${r.payment_status === "Paid"
            ? `<span class="px-2 py-1 font-semibold leading-tight text-green-700 bg-green-100 rounded-full dark:bg-green-700 dark:text-green-100">${r.payment_status}</span>`
            : `<span class="px-2 py-1 font-semibold leading-tight text-orange-700 bg-orange-100 rounded-full dark:text-white dark:bg-orange-600">${r.payment_status}</span>`
          }
        </td>

        <!-- Status -->
        <td class="px-4 py-3 text-xs">
          <span class="capitalize px-2 py-1 font-semibold leading-tight rounded-full
          ${r.status === "pending"
              ? "text-blue-700 bg-blue-100 dark:bg-blue-700 dark:text-blue-100"
              : "text-green-700 bg-green-100 dark:bg-green-700 dark:text-green-100"}">
          ${r.status}
          </span>
        </td>

        <!-- Created -->
        <td class="px-4 py-3 text-sm">${new Date(r.created_at).toISOString().slice(0, 10)}</td>

        <!-- Action Buttons -->
        <td class="px-4 py-3 text-xs">
          ${buildActionButtons(r, index)}
        </td>

      </tr>
    `;
  }


  /* =====================================================================
     HELPERS
  ===================================================================== */
  function formatDateTime(dateStr) {
    return new Date(dateStr).toLocaleString('en-CA', {
      year:   'numeric',
      month:  '2-digit',
      day:    '2-digit',
      hour:   '2-digit',
      minute: '2-digit',
      hour12: true,
    }).replace(',', '');
  }

  function buildActionButtons(r, index) {
    let buttons = '';
    const now      = new Date();
    const checkIn  = new Date(r.check_in_date);
    const checkOut = new Date(r.check_out_date);

    // Clean button — shown after checkout
    if (r.status === "check-out") {
      buttons += `
        <button class="action-btn px-3 py-1 bg-green-500 text-white rounded hover:bg-green-600 mr-2"
          data-room-id="${r.room_id}" data-id="${r.book_id}" data-action="clean">
          Clean
        </button>`;
    }

    // Done button — shown while cleaning
    if (r.status === "cleaning") {
      buttons += `
        <button class="action-btn px-3 py-1 bg-purple-500 text-white rounded hover:bg-purple-600 mr-2"
          data-room-id="${r.room_id}" data-id="${r.book_id}" data-action="done">
          Done?
        </button>`;
    }

    // Check-in button — same day and time has passed
    const isSameDayCheckIn  = now.toDateString() === checkIn.toDateString();
    const checkInTimePassed = now.getTime() >= checkIn.getTime();

    if (isSameDayCheckIn && checkInTimePassed && r.status === "pending") {
      buttons += `
        <button class="px-3 py-1 bg-blue-500 text-white rounded-full hover:bg-blue-600 mr-2"
          onclick="populateReservationModal(${index})">
          CheckIn?
        </button>`;
    }

    // Check-out button — same day and within 2 hours before checkout
    const isSameDayCheckOut = now.toDateString() === checkOut.toDateString();
    const twoHoursBefore    = 2 * 60 * 60 * 1000;
    const withinWindow      = now.getTime() >= (checkOut.getTime() - twoHoursBefore);

    if (isSameDayCheckOut && withinWindow && r.status === "check-in") {
      buttons += `
        <button class="action-btn px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600 mr-2"
          data-room-id="${r.room_id}" data-id="${r.book_id}" data-action="checkout">
          Checkout
        </button>`;
    }

    return buttons;
  }


  /* =====================================================================
     EVENT LISTENERS
  ===================================================================== */
  function attachEventListeners() {
    tableElement.addEventListener("click", handleButtonClick, true);
  }

  function handleButtonClick(e) {
    const btn = e.target;
    if (!btn.classList.contains('action-btn')) return;

    e.preventDefault();
    e.stopPropagation();

    const row   = btn.closest('tr');
    const index = parseInt(row?.dataset.index);

    if (!row || isNaN(index) || !reservationsData[index]) {
      console.error("Invalid row or index:", row?.dataset.index);
      return;
    }

    const reservation = reservationsData[index];
    const action      = btn.dataset.action;

    if (!action) {
      console.error("Missing action on button");
      return;
    }

    bookingId = btn.dataset.id;
    roomId    = btn.dataset.roomId;

    btn.disabled = true;
    handleAction(action, reservation, btn);
  }

  function handleAction(action, reservation, btn) {
    switch (action) {
      case 'clean':    handleClean(reservation, btn);    break;
      case 'checkin':  handleCheckin(reservation, btn);  break;
      case 'checkout': handleCheckout(reservation, btn); break;
      case 'done':     handleCompleted(reservation, btn);break;
      default:
        console.warn(`Unknown action: ${action}`);
        btn.disabled = false;
    }
  }


  /* =====================================================================
     ACTION HANDLERS
  ===================================================================== */
  function handleClean(reservation, btn) {
    openModal();
  }

  function handleCheckin(reservation, btn) {
    Swal.fire({
      title:             "Are you sure the guest wants to check in?",
      text:              "This action cannot be undone!",
      icon:              "warning",
      showCancelButton:  true,
      confirmButtonColor:"#d33",
      cancelButtonColor: "#3085d6",
      confirmButtonText: "Yes, mark as check-in",
    }).then(result => {
      if (!result.isConfirmed) {
        if (btn) btn.disabled = false;
        return;
      }

      fetch(`/api/reservations/checkin/${reservation.book_id}`, {
        method:  'POST',
        headers: { 'Content-Type': 'application/json' },
      })
        .then(res => res.json())
        .then(() => {
          notification("success", "Guest successfully checked in");
          window.location.reload();
        })
        .catch(() => {
          alert("Failed to check in guest");
          if (btn) btn.disabled = false;
        });
    });
  }

  function handleCheckout(reservation, btn) {
    Swal.fire({
      title:             "Are you sure the guest wants to check out?",
      text:              "This action cannot be undone!",
      icon:              "warning",
      showCancelButton:  true,
      confirmButtonColor:"#d33",
      cancelButtonColor: "#3085d6",
      confirmButtonText: "Yes, mark as check-out",
    }).then(result => {
      if (!result.isConfirmed) {
        if (btn) btn.disabled = false;
        return;
      }

      fetch(`/api/reservations/checkout/${reservation.book_id}`, {
        method:  'POST',
        headers: { 'Content-Type': 'application/json' },
      })
        .then(res => res.json())
        .then(() => {
          notification("success", "Guest successfully checked out");
          window.location.reload();
        })
        .catch(() => {
          alert("Failed to check out guest");
          if (btn) btn.disabled = false;
        });
    });
  }

  function handleCompleted(reservation, btn) {
    Swal.fire({
      title:             "Is the room fully cleaned?",
      text:              "This action cannot be undone!",
      icon:              "warning",
      showCancelButton:  true,
      confirmButtonColor:"#d33",
      cancelButtonColor: "#3085d6",
      confirmButtonText: "Yes, mark as clean",
    }).then(result => {
      if (!result.isConfirmed) {
        if (btn) btn.disabled = false;
        return;
      }

      fetch(`/api/reservations/completed/${reservation.book_id}`, {
        method:  'POST',
        headers: { 'Content-Type': 'application/json' },
        body:    JSON.stringify({ room_id: roomId }),
      })
        .then(res => res.json())
        .then(() => {
          notification("success", "Room marked as clean");
          window.location.reload();
        })
        .catch(() => {
          alert("Failed to mark room as clean");
          if (btn) btn.disabled = false;
        });
    });
  }


  /* =====================================================================
     FETCH ATTENDANTS
  ===================================================================== */
  function fetchAttendants() {
    const container     = document.getElementById('attendant-checkboxes');
    container.innerHTML = '<p class="text-sm text-gray-500 text-center py-4">Loading attendants...</p>';

    fetch('/api/maintenances')
      .then(res => {
        if (!res.ok) throw new Error("Failed to fetch attendants");
        return res.json();
      })
      .then(maintenance => populateAttendantCheckboxes(maintenance))
      .catch(error => {
        console.error(error);
        container.innerHTML = '<p class="text-sm text-red-500 text-center py-4">Error loading attendants</p>';
      });
  }


  /* =====================================================================
     POPULATE ATTENDANT CHECKBOXES
  ===================================================================== */
  function populateAttendantCheckboxes(maintenanceList) {
    const container = document.getElementById('attendant-checkboxes');
    if (!container) return;

    container.innerHTML = '';

    if (maintenanceList.length === 0) {
      container.innerHTML = '<p class="text-sm text-gray-500 text-center py-4">No attendants available</p>';
      return;
    }

    maintenanceList.forEach(attendant => {
      const label     = document.createElement('label');
      label.className = 'flex items-center space-x-3 p-2 hover:bg-gray-50 rounded-lg cursor-pointer transition-colors';
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

    container.addEventListener('change', function (e) {
      if (e.target.name === 'attendants[]') updateSelectedCount();
    });

    updateSelectedCount();
  }


  /* =====================================================================
     SELECTED COUNT
  ===================================================================== */
  function updateSelectedCount() {
    const checkboxes   = document.querySelectorAll('input[name="attendants[]"]:checked');
    const countElement = document.getElementById('selected-count');
    if (countElement) countElement.textContent = checkboxes.length;
  }

  function getSelectedAttendants() {
    const checkboxes = document.querySelectorAll('input[name="attendants[]"]:checked');
    return Array.from(checkboxes).map(cb => ({ id: cb.value, name: cb.dataset.name }));
  }


  /* =====================================================================
     FORM SUBMISSION (assign cleaner)
  ===================================================================== */
  form.addEventListener('submit', function (e) {
    e.preventDefault();

    const selectedAttendants = getSelectedAttendants();

    if (selectedAttendants.length === 0) {
      alert('Please select at least one attendant');
      return;
    }

    const formData = {
      cleaner_id: selectedAttendants.map(a => a.id),
      room_id:    roomId,
    };

    fetch(`/api/reservations/assigncleaner/${bookingId}`, {
      method:      'POST',
      credentials: 'include',
      headers:     { 'Content-Type': 'application/json' },
      body:        JSON.stringify(formData),
    })
      .then(response =>
        response.json().then(data => {
          if (!response.ok) throw new Error(data.error || 'Request failed');
          return data;
        })
      )
      .then(() => {
        alert('Successfully assigned cleaner(s)!');
        closeModal();
      })
      .catch(error => {
        console.error('Error:', error);
        alert(`Failed to assign cleaner: ${error.message}`);
      });
  });

});