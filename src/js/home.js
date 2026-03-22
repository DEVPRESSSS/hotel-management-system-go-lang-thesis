document.addEventListener("DOMContentLoaded", () => {
    fetchRooms();

    //Fetch the rooms
    function fetchRooms() {
        fetch("/avail/rooms")
            .then(response => {
                if (!response.ok) {
                    throw new Error("Failed to fetch rooms");
                }
                return response.json();
            })
            .then(data => renderRooms(data))
            .catch(error => console.error(error));
    }

    //Populate the cards
    function renderRooms(rooms) {
        const grid = document.getElementById("roomsGrid");
        grid.innerHTML = "";

        rooms.forEach(room => {
            const card = document.createElement("div");
            card.className =
                "bg-white rounded-lg shadow-lg overflow-hidden hover:shadow-xl transition duration-300 gap-8";

         card.innerHTML = `
                <div class="relative overflow-hidden h-64 bg-gray-300">
                    <img 
                        src="https://images.unsplash.com/photo-1631049307264-da0ec9d70304?w=400&h=300&fit=crop"
                        alt="Room Image"
                        class="w-full h-full object-cover hover:scale-105 transition duration-300"
                    >
                </div>

                <div class="p-6">
                    <div class="flex items-center justify-between mb-2">
                        <h3 class="text-xl font-bold text-gray-800">
                            ${room.RoomType.roomtypename} Room
                        </h3>
                        <span class="bg-purple-600 text-white px-3 py-1 rounded-full text-sm font-semibold">
                            ${room.price}/night
                        </span>
                    </div>

                    <p class="text-gray-600 text-sm mb-4">
                        ${room.RoomType.description || "Comfortable and well-furnished room."}
                    </p>

                    <div class="flex items-center text-gray-500 text-sm mb-4 space-x-4">
                        <span>👥 ${room.capacity}</span>
                        <span>🏢 ${room.Floor.floorname}</span>
                    </div>

                    <a id="book-btn" data-room_id = "${room.roomid}"
                        class="booknow-btn w-full bg-purple-600 text-white py-3 rounded-lg hover:bg-blue-700 transition duration-300 font-semibold p-2"
                    
                    >
                    Book now
                    </a>
                </div>
            `;
            //   <div class="flex items-center mb-3">
            //             <span class="text-yellow-400">★★★★★</span>
            //             <span class="text-gray-600 text-sm ml-2">(89 reviews)</span>
            //     </div>

            grid.appendChild(card);
        });
    }

    //Get the selected room to populate the cards
    document.getElementById("roomsGrid").addEventListener("click", function (e) {
        if (!e.target.classList.contains("booknow-btn")) return;

        e.preventDefault();

        const id = e.target.dataset.room_id;
        if (!id) return;

        //Store the selected room to session storage
        sessionStorage.setItem("selectedRoomId", JSON.stringify(id));
        window.location.href = "/roomdetails";

        
    });

    function buildFilters() {
        const filters = document.getElementById('filter-wrapper');
        if (!filters) return;

        const options = ['All', 'Available Today', '1-2 guests','3-4 guests',"Deluxe","Standard"];

        options.forEach((label, i) => {
            const btn = document.createElement('button');
            btn.textContent = label;
            btn.dataset.filter = label.toLowerCase();
            btn.className = i === 0
                ? 'filter-btn bg-purple-600 text-white px-5 py-2 rounded-full text-sm font-semibold transition duration-200'
                : 'filter-btn bg-gray-100 text-gray-600 px-5 py-2 rounded-full text-sm font-semibold hover:bg-purple-100 hover:text-purple-600 transition duration-200';

            filters.appendChild(btn);
        });

        // Click handler
        filters.addEventListener('click', function (e) {
            if (!e.target.classList.contains('filter-btn')) return;
            
            // Reset all buttons
            document.querySelectorAll('.filter-btn').forEach(btn => {
                btn.className = 'filter-btn bg-gray-100 text-gray-600 px-5 py-2 rounded-full text-sm font-semibold hover:bg-purple-100 hover:text-purple-600 transition duration-200';
            });

            const selectedFilter = e.target.dataset.filter; 
            // Highlight active
            e.target.className = 'filter-btn bg-purple-600 text-white px-5 py-2 rounded-full text-sm font-semibold transition duration-200';

            if(selectedFilter == "all"){
                fetchRooms();
                return;
            }


            //Fetch
            fetch(`/avail/rooms/${selectedFilter}`)
                 .then(res =>{
                        if(!res.ok){
                            throw new Error("Failed to filter");
                        }
                        return res.json();
                 }).then(data => renderRooms(data))
                  .catch(error => console.error(error));

        });
    }
    buildFilters();

});