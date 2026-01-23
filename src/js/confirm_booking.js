document.addEventListener("DOMContentLoaded", () => {

    //Get the booking info from session storage since it is draft
    const bookingInfo = sessionStorage.getItem("bookingDraft");
    //Parse the booking info
    const bookingSummary = JSON.parse(bookingInfo);
    console.log(bookingSummary.checkOut);

    //Populate the value
    const roomId = bookingSummary.room_id;
    document.getElementById("room-type").textContent= bookingSummary.room_type;
    document.getElementById("check-in").textContent= bookingSummary.check_in;
    document.getElementById("check-out").textContent= bookingSummary.check_out;
    document.getElementById("guest").textContent= bookingSummary.guests;

    fetch("/api/booking/calculate", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            room_id: roomId,
            check_in: bookingSummary.check_in,
            check_out: bookingSummary.check_out
        })
        })
        .then(response => {
            if (!response.ok) {
                throw new Error("Calculation failed");
            }
            return response.json();
        })
        .then(data => {
            document.getElementById("price").textContent = data.price_per_night;
            document.getElementById("total").textContent = data.total;
        })
        .catch(error => console.error(error));

});