document.addEventListener("DOMContentLoaded", () => {
    const book_id = sessionStorage.getItem("bookId");
    async function insertOrderAfterPayment() {
        const cart = JSON.parse(sessionStorage.getItem("foodCart") || "[]");
      
        if (!cart.length) return;

        const response = await fetch('/paymongo/food/insert-order', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                book_id:book_id,
                items: cart.map(item => ({ id: item.id, qty: item.qty }))
            })
        });

        if (response.ok) {
            sessionStorage.removeItem("foodCart"); 
            sessionStorage.removeItem("bookId"); 
        }
    }
   

    // Call on success page load
    insertOrderAfterPayment();
});