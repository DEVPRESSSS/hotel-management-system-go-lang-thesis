document.addEventListener('DOMContentLoaded', () => {

  const eventModal    = document.getElementById('eventModal');
  const headerAction  = document.getElementById('header-action');
  const checkIn       = document.getElementById('check-in');
  const checkOut      = document.getElementById('check-out');
  const duration      = document.getElementById('duration');
  const roomNumber    = document.getElementById('room-number');
  const roomType      = document.getElementById('room-type');
  const roomPrice     = document.getElementById('room-price');
  const paymentStatus = document.getElementById('payment-status');
  const guestCount    = document.getElementById('guest-count');
  const fullName      = document.getElementById('fullname');
  const badgeName     = document.getElementById('badge-name');
  const guestList     = document.getElementById('guest-list'); 

  const today     = new Date();
  const nextMonth = new Date(today);

  const calendarEl = document.getElementById('calendar');

  const calendar = new FullCalendar.Calendar(calendarEl, {
    initialView: 'dayGridMonth',
    headerToolbar: {
      left:   'prev,next today',
      center: 'title',
      right:  'dayGridMonth,timeGridWeek,timeGridDay',
    },
    validRange: {
      start: today,
      end:   nextMonth.setMonth(today.getMonth() + 1)
    },
    events: {
      url: '/api/reservations/events',
      method: 'GET',
      failure() {
        alert('Failed to load reservations');
      }
    },
    eventClick: function(info) {
      openModal(info.event.title);
    }
  });

  calendar.render();

  // ─── Helper: generate badge initials ────────────────────────────────────
  function getBadge(fullname) {
    const parts        = fullname.trim().split(" ");
    const firstName    = parts[0] || "";
    const lastName     = parts[1] || "";
    const firstLetter  = firstName[0] || "";
    const secondLetter = lastName[0]  || firstName[1] || "";
    return (firstLetter + secondLetter).toUpperCase();
  }

  // ─── Helper: render one guest row ───────────────────────────────────────
  function renderGuestRow(name) {
    return `
      <div class="flex items-center gap-2">
        <div class="w-7 h-7 rounded-full bg-purple-100 text-purple-600 text-xs font-bold flex items-center justify-center flex-shrink-0">
          ${getBadge(name)}
        </div>
        <div class="flex flex-col">
          <span class="text-sm text-gray-800 font-medium">${name}</span>
        </div>
      </div>
    `;
  }

  // ─── Open Modal ──────────────────────────────────────────────────────────
  function openModal(title) {
    eventModal.classList.remove('hidden');
    eventModal.classList.add('flex');

    const bookingId = title.trim().split(" ")[0];

    fetch(`/api/events/bookingInfo/${bookingId}`)
      .then(res => {
        if (!res.ok) throw new Error("Failed to fetch");
        return res.json();
      })
      .then(data => {
        const d = data.booking_details;

        // Header
        headerAction.textContent = d.book_id;

        // Dates
        const checkInDate  = new Date(d.check_in_date);
        const checkOutDate = new Date(d.check_out_date);
        const nights       = (checkOutDate - checkInDate) / (1000 * 60 * 60 * 24);

        checkIn.textContent   = checkInDate.toISOString().slice(0, 10);
        checkOut.textContent  = checkOutDate.toISOString().slice(0, 10);
        duration.textContent  = `${nights} ${nights > 1 ? 'nights' : 'night'}`;

        // Room
        roomNumber.textContent = d.room_number;
        roomType.textContent   = d.room_type;
        roomPrice.textContent  = `₱${d.total_price}`;

        // Payment status
        paymentStatus.textContent = d.payment_status;

        // Guest count
        guestCount.textContent = `( ${d.num_guests} )`;

        // ─── Guest List ───────────────────────────────────────────────────
        let guestHTML = renderGuestRow(d.fullname);

        if (d.guests && d.guests.length > 0) {
          d.guests.forEach(guest => {
            const fullname = `${guest.firstname} ${guest.lastname}`;
            guestHTML += renderGuestRow(fullname);
          });
        }

        guestList.innerHTML = guestHTML;
      })
      .catch(console.error);
  }

  // ─── Close Modal ─────────────────────────────────────────────────────────
  function closeModal() {
    eventModal.classList.add('hidden');
    eventModal.classList.remove('flex');
  }

  window.closeModal = closeModal;
});