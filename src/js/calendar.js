
document.addEventListener('DOMContentLoaded', () => {

  const eventModal = document.getElementById('eventModal');
  const headerAction = document.getElementById('header-action');
  const roomNo = document.getElementById('room-number');

  const today = new Date()
  const nextMonth = new Date(today)

  const calendarEl = document.getElementById('calendar')

  const calendar = new FullCalendar.Calendar(calendarEl, {
    initialView: 'dayGridMonth',
    headerToolbar: {
      left: 'prev,next today',
      center: 'title',
      right: 'dayGridMonth,timeGridWeek,timeGridDay',
    },
   
    validRange: {
      start: today,
      end: nextMonth.setMonth(today.getMonth() + 1)
    },
     events: {
      url: '/api/reservations/events', 
      method: 'GET',
      failure() {
        alert('Failed to load reservations')
      }
    },
    eventClick: function(info){
     
       openModal(info.event.title);
    }
  
   
  })

  calendar.render()

  function openModal(title) {
    eventModal.classList.remove('hidden');
    eventModal.classList.add('flex');

    const bookingInfo = title.trim();
    const bookingId = bookingInfo.split(" ")[0];
    const roomNumber = bookingInfo.split(" ")[1];
    //Render the booking information
    headerAction.textContent = bookingId;
    roomNo.textContent = roomNumber;
  }

  function closeModal() {
    eventModal.classList.add('hidden');
    eventModal.classList.remove('flex');
  }

  window.closeModal = closeModal;
  
});


//Call the modal for viewing of event
