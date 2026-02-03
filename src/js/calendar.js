// import { Calendar } from "@fullcalendar/core";
// import dayGridPlugin from '@fullcalendar/daygrid'
// import interactionPlugin from '@fullcalendar/interaction'



document.addEventListener('DOMContentLoaded', () => {
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
    }
  
   
  })

  calendar.render()
})
