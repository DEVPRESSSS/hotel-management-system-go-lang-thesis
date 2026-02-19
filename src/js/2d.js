document.addEventListener("DOMContentLoaded", function () {

    fetch('/api/rooms')
        .then(res => res.json())
        .then(rooms => {

            rooms.forEach(r => {

                // Use correct property name
                const roomGroup = document.getElementById(r.roomnumber);
                if (!roomGroup) return;

                const rect = roomGroup.querySelector("rect");
                if (!rect) return;
                const path = roomGroup.querySelector("path");
                if (!path) return;

                path.setAttribute("fill", "white");
                rect.setAttribute(
                    "fill",
                    r.status === "available" ? "green" : "red"
                );

            });

        })
        .catch(console.error);

});
