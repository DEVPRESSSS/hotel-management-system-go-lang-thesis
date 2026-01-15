
// Get the room data from sessionStorage
const id = sessionStorage.getItem("selectedRoomId");

if(!id){
  window.location.href = "/"
} 
const roomId = JSON.parse(id);

fetch(`/api/roomselected/${roomId}`)
        .then(response => {
            if (!response.ok) {
                throw new Error("Failed to retrieve room");
            }
            return response.json();
        })
        .then(data => {

          const room = data.room;
          console.log(room);
          //Populate the room title
          const roomType = document.getElementById('roomTitle');
          roomType.textContent = room.RoomType.roomtypename;

          //Populate the room number
          const roomNumber = document.getElementById('roomNumber');
          roomNumber.textContent = room.roomnumber;

          //Populate the room price
          const roomPrice = document.getElementById('roomPrice');
          roomPrice.textContent = room.price;

          //Populate the room capacity
          const roomCapacity = document.getElementById('roomCapacity');
          roomCapacity.textContent = room.capacity;

         //Populate the room floor
          const roomFloor = document.getElementById('roomFloor');
          roomFloor.textContent = room.Floor.floorname;

          //Populate the room description
          const description = document.getElementById('roomDescription');
          description.textContent = room.RoomType.description;

          //Populate the room description
          const amenitiesContainer = document.getElementById("amenities");
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