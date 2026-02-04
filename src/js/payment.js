document.addEventListener("DOMContentLoaded", () => {
const bookingSummary = JSON.parse(sessionStorage.getItem("bookingDraft"));
const guestData = JSON.parse(sessionStorage.getItem("guestData"));

//Call immmediately the API after successfull payment
    fetch('/api/booking/confirmbooking',{
            method: 'POST',
            headers: {
            "Content-Type": "application/json"
            },
            body: JSON.stringify({
                room_id:bookingSummary.room_id,
                room_number: bookingSummary.room_number,
                room_type: bookingSummary.room_type,
                check_in_date: new Date(bookingSummary.check_in).toISOString(),
                check_out_date: new Date(bookingSummary.check_out).toISOString(),
                num_guests: Number(bookingSummary.guest),           
                special_requests: guestData[0]?.specialRequests || '',
                guests: guestData
            })
        })
        .then(response =>{
            if(!response.ok){
                throw new Error("Booking error");
            }
            return response.json();
        })
        .then(data=>{
            setTimeout(() => {
                window.location.href = "/guest/dashboard";
            },2000);
            
        })
        .catch(error =>{
            console.log(error);
        });
});