document.addEventListener('DOMContentLoaded', function(){

//#region--Display room details 
// Get the room data from sessionStorage
const id = sessionStorage.getItem("selectedRoomId");

if(!id){
  window.location.href = "/"
} 
const roomId = JSON.parse(id);
const roomType = document.getElementById('roomTitle');
const roomPrice = document.getElementById('roomPrice');
const roomNumber = document.getElementById('roomNumber');
const roomCapacity = document.getElementById('roomCapacity');
const roomFloor = document.getElementById('roomFloor');
const description = document.getElementById('roomDescription');
const amenitiesContainer = document.getElementById("amenities");
//Book now button


fetch(`/api/roomselected/${roomId}`)
        .then(response => {
            if (!response.ok) {
                throw new Error("Failed to retrieve room");
            }
            return response.json();
        })
        .then(data => {

          const room = data.room;
       
          //Populate the room title
          roomType.textContent = room.RoomType.roomtypename;

          //Populate the room number
          roomNumber.textContent = room.roomnumber;

          //Populate the room price
          roomPrice.textContent = room.price;

          //Populate the room capacity
          roomCapacity.textContent = room.capacity;

         //Populate the room floor
          roomFloor.textContent = room.Floor.floorname;

          //Populate the room description
          description.textContent = room.RoomType.description;

          //Populate the room description
          amenitiesContainer.innerHTML = ""; 

          room.amenities.forEach(amenity => {
            const div = document.createElement("div");
            div.className = "flex items-center gap-2 text-sm text-slate-600";

            div.innerHTML = `
              <svg class="w-5 h-5 text-green-700" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/>
              </svg>
              <span>${amenity.amenityname}</span>
            `;

            amenitiesContainer.appendChild(div);
           });


        })
        .catch(error => console.log(error));

//Store the information of the booking if the user click the book now btn
const bookNowBtn = document.getElementById("book-now");

bookNowBtn.addEventListener("click", function (e) {
  e.preventDefault(); 

  // const checkIn = document.getElementById("checkin").value;
  // const checkOut = document.getElementById("check-out").value;
  const checkIn = document.querySelector(".checkin").value
  const checkOut = document.querySelector(".checkout").value;
  const numberOfGuest = document.getElementById("guest").value;

  console.log(checkIn);
  console.log(checkOut);

  if (!checkIn || !checkOut) {
    alert("Please select check-in and check-out dates");
    return;
  }

  if(checkIn == checkOut){
    alert("Same day checkout is not supported");
    return;
  }
  // store draft booking
  sessionStorage.setItem(
    "bookingDraft",
    JSON.stringify({
      room_id:roomId,
      room_number:roomNumber.textContent,
      room_type:roomType.textContent,
      check_in: checkIn,
      check_out: checkOut,
      guest: Number(numberOfGuest)
    })
  );


  // redirect to next page
  window.location.href = "/booking/summary";
});
//#endregion


//#region--Display dates status in datetime picker
const currentDate = new Date();
console.log(currentDate);
const nextMonth = new Date(currentDate);

fetch(`/api/calendar/${roomId}`)
  .then(response => {
    if (!response.ok) throw new Error("Failed to load");
    return response.json();
  })
  .then(data => {
    const booked = data.books;
    
    // Create an array to store all booked date ranges
    const bookedRanges = [];
    
    if (booked.length > 0) {
      booked.forEach(booking => {
        bookedRanges.push({
          from: booking.check_in_date.split("T")[0],
          to: booking.check_out_date.split("T")[0]
        });
      });
    }

    console.log("Booked ranges:", bookedRanges);

    flatpickr("#date-picker", {
      enableTime: true,
      time_24hr: false,
      dateFormat: "Y-m-d h:i K",
      minDate: currentDate,
      maxDate: nextMonth.setMonth(currentDate.getMonth() + 1),
      
      // Disable all booked ranges
      disable: bookedRanges,

      onDayCreate: function(dObj, dStr, fp, dayElem) {
          const date = dayElem.dateObj;
          
          // Use local date string instead of UTC to avoid timezone issues
          const dateStr = date.getFullYear() + '-' + 
                          String(date.getMonth() + 1).padStart(2, '0') + '-' + 
                          String(date.getDate()).padStart(2, '0');
          
          // Check if date is within ANY booked range
          let isBooked = false;
          bookedRanges.forEach(range => {
              if (dateStr >= range.from && dateStr <= range.to) {
                  isBooked = true;
              }
          });
          
          // Apply styling based on booking status
          if (isBooked) {
              dayElem.setAttribute('style', 'background-color: #ff4444 !important; color: white !important;');
              dayElem.classList.add('red-date');
          } else if (date >= currentDate) {
              dayElem.style.backgroundColor = "#44ff44";
              dayElem.style.color = "white";
          }
      }
    });
  })
  .catch(error => console.log(error));
//#endregion
});
