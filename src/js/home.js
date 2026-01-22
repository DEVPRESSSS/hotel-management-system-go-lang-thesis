document.addEventListener("DOMContentLoaded", () => {
    fetchRooms();
});
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
                <span class="absolute top-4 right-4 bg-blue-600 text-white px-3 py-1 rounded-full text-sm font-semibold">
                    ${room.price}/night
                </span>
            </div>

            <div class="p-6">
                <h3 class="text-xl font-bold text-gray-800 mb-2">
                    ${room.RoomType.roomtypename} Room
                </h3>

                <div class="flex items-center mb-3">
                    <span class="text-yellow-400">★★★★★</span>
                    <span class="text-gray-600 text-sm ml-2">(89 reviews)</span>
                </div>

                <p class="text-gray-600 text-sm mb-4">
                    ${room.RoomType.description || "Comfortable and well-furnished room."}
                </p>

                <div class="flex items-center text-gray-500 text-sm mb-4 space-x-4">
                    <span>👥 ${room.capacity}</span>
                    <span>🏢 ${room.Floor.floorname}</span>
                </div>

                <a id="book-btn" data-room_id = "${room.roomid}"
                    class="booknow-btn w-full bg-blue-600 text-white py-3 rounded-lg hover:bg-blue-700 transition duration-300 font-semibold p-2"
                    ${room.status !== "available" ? "disabled class='opacity-50 cursor-not-allowed'" : ""}
                >
                    ${room.status === "available" ? "Book Now" : "Unavailable"}
                </a>
            </div>
        `;

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

