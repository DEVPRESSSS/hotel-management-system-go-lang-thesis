document.addEventListener("DOMContentLoaded", () => {
    //Get the booking info from session storage since it is draft
    const bookingInfo = sessionStorage.getItem("bookingDraft");
    //Parse the booking info
    const bookingSummary = JSON.parse(bookingInfo);

    //Populate the value
    const roomId = bookingSummary.room_id;
    let price = document.getElementById("price");
    let total = document.getElementById("total");
    
    document.getElementById("room-type").textContent= bookingSummary.room_type;
    document.getElementById("check-in").textContent= bookingSummary.check_in;
    document.getElementById("check-out").textContent= bookingSummary.check_out;
    document.getElementById("guest").textContent= bookingSummary.guest;

    const guestNumber = Number(bookingSummary.guest);
    const container = document.getElementById('guest-container');

    // Store all guest data
    let guestData = [];
    let currentGuestIndex = 0;

    const shouldSkipFirstGuest = guestNumber >= 2;
    const totalFormsToShow = shouldSkipFirstGuest ? guestNumber - 1 : 1;

    // Function to generate guest-specific fields
    function getGuestFields(guestNum) {
        return `
            <!-- Name Fields -->
            <div class="flex flex-col sm:flex-row gap-4 mb-6">
                <div class="flex-1">
                    <label class="block text-sm font-medium text-gray-900 mb-2">
                        First Name <span class="text-red-600">*</span>
                    </label>
                    <input type="text" name="firstName" id="firstName-${guestNum}" class="w-full px-4 py-3 bg-gray-100 border border-gray-100 rounded focus:outline-none focus:border-blue-500" required>
                </div>
            </div>

            <div class="flex flex-col sm:flex-row gap-4 mb-6">
                <div class="flex-1">
                    <label class="block text-sm font-medium text-gray-900 mb-2">
                        Last Name <span class="text-red-600">*</span>
                    </label>
                    <input type="text" name="lastName" id="lastName-${guestNum}" class="w-full px-4 py-3 bg-gray-100 border border-gray-100 rounded focus:outline-none focus:border-blue-500" required>
                </div>
            </div>

            <!-- Contact Fields -->
            <div class="flex flex-col sm:flex-row gap-4 mb-6">
                <div class="flex-1">
                    <label class="block text-sm font-medium text-gray-900 mb-2">
                        Phone Number <span class="text-red-600">*</span>
                    </label>
                    <input type="tel" name="phone" id="phone-${guestNum}" class="w-full px-4 py-3 bg-gray-100 border border-gray-100 rounded focus:outline-none focus:border-blue-500" required>
                </div>
            </div>
        `;
    }

    // Common fields (shown on last guest or single guest)
    function getCommonFields() {
        return `
            <!-- Special Requests -->
            <div class="mb-6">
                <label class="block text-sm font-medium text-gray-900 mb-2">
                    Special Requests (Optional)
                </label>
                <textarea rows="5" name="specialRequests" id="specialRequests" class="w-full px-4 py-3 bg-gray-100 border border-gray-100 rounded focus:outline-none focus:border-blue-500 resize-none" placeholder="Any special requests or requirements..."></textarea>
            </div>

            <!-- Terms Checkbox -->
            <div class="mb-6">
                <label class="flex items-start">
                    <input type="checkbox" name="terms" id="terms" class="mt-1 mr-3" required>
                    <span class="text-sm text-gray-900">
                        I agree to the terms and conditions and the cancellation policy <span class="text-red-600">*</span>
                    </span>
                </label>
            </div>
        `;
    }

    // Function to render form for current guest
    function renderGuestForm() {
       const isLastGuest = currentGuestIndex === totalFormsToShow - 1;
       const isFirstGuest = currentGuestIndex === 0;
        
       const actualGuestNum = shouldSkipFirstGuest ? currentGuestIndex + 2 : currentGuestIndex + 1;

        container.innerHTML = `
            <h1 class="text-3xl font-semibold text-gray-900 mb-2">
                ${guestNumber > 1 ? `Guest Information for Guest No. ${actualGuestNum}` : 'Complete Your Booking'}
            </h1>
            <p class="text-gray-600 mb-8">
                ${guestNumber > 1 ? `Please provide details for guest ${actualGuestNum} of ${guestNumber}` : 'Almost there! Just a few final details'}
            </p>
            
            ${guestNumber > 1 ? `
                <div class="mb-6">
                    <div class="flex items-center justify-between">
                        <span class="text-sm font-medium text-gray-600">Progress</span>
                        <span class="text-sm font-medium text-gray-900">${currentGuestIndex + 1} / ${guestNumber}</span>
                    </div>
                    <div class="mt-2 w-full bg-gray-200 rounded-full h-2">
                        <div class="bg-gray-900 h-2 rounded-full transition-all duration-300" style="width: ${((currentGuestIndex + 1) / guestNumber) * 100}%"></div>
                    </div>
                </div>
            ` : ''}

             ${totalFormsToShow > 1 ? `
                <div class="progress">
                    Progress ${currentGuestIndex + 1} / ${totalFormsToShow}
                </div>
              ` : ''}

            
            <form id="guest-form">
                ${guestNumber > 1 ? getGuestFields(currentGuestIndex + 1) : ''}
                ${isLastGuest || guestNumber === 1 ? getCommonFields() : ''}
                
                <!-- Buttons -->
                <div class="flex flex-col sm:flex-row gap-4">
                    <button type="button" id="back-btn" class="flex-1 px-6 py-3 bg-white border border-gray-300 rounded text-gray-700 font-medium hover:bg-gray-50 flex items-center justify-center ${isFirstGuest ? 'opacity-50 cursor-not-allowed' : ''}">
                        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
                        </svg>
                        Back
                    </button>
             
                      ${isLastGuest ? `
                      <!-- Stripe/Credit Card Icon -->
                        <button type="submit" data-gateway="stripe"
                            class="stripe flex-1 px-6 py-3 bg-purple-600 text-white rounded font-medium hover:bg-purple-700 flex items-center justify-center">
                            Credit Card
                            <svg class="w-5 h-5 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                <rect x="1" y="4" width="22" height="16" rx="2" ry="2"/>
                                <line x1="1" y1="10" x2="23" y2="10"/>
                            </svg>
                        </button>

                        <!-- GCash/Mobile Wallet Icon -->
                        <button type="submit" data-gateway="paymongo"
                            class="gcash flex-1 px-6 py-3 bg-blue-600 text-white rounded font-medium hover:bg-blue-700 flex items-center justify-center">
                            GCash
                            <svg class="w-5 h-5 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                <rect x="5" y="2" width="14" height="20" rx="2" ry="2"/>
                                <line x1="12" y1="18" x2="12.01" y2="18"/>
                            </svg>
                        </button>

                        <!-- Cash/Money Icon -->
                        <button type="button" id="cash-btn" 
                            class="cash flex-1 px-6 py-3 bg-green-600 text-white rounded font-medium hover:bg-green-700 flex items-center justify-center">
                            Cash
                            <svg class="w-5 h-5 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                <line x1="12" y1="1" x2="12" y2="23"/>
                                <path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/>
                            </svg>
                        </button>
                        ` : `
                            <button  id="next-btn" class="flex-1 px-6 py-3 bg-purple-600 text-white rounded font-medium hover:bg-gray-800 flex items-center justify-center ">             
                                Next
                                <svg class="w-5 h-5 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
                                </svg>
                           </button>`}
                      
                       
                </div>
            </form>
        `;
        
        attachEventListeners();
    }

    // Function to collect current form data
    function collectFormData() {
        const formData = {};
        const actualGuestNum = shouldSkipFirstGuest ? currentGuestIndex + 2 : currentGuestIndex + 1;
        
        if (totalFormsToShow > 1) {
            formData.firstName = document.getElementById(`firstName-${actualGuestNum}`).value;
            formData.lastName = document.getElementById(`lastName-${actualGuestNum}`).value;
            formData.phoneNumber = document.getElementById(`phone-${actualGuestNum}`).value;
        }
        
        // Collect common fields if on last guest
        if (currentGuestIndex === totalFormsToShow - 1 || guestNumber === 1) {
            formData.specialRequests = document.getElementById('specialRequests').value;
            formData.terms = document.getElementById('terms').checked;
        }
        
        return formData;
    }

    // Function to populate form with saved data
    function populateFormData() {
        if (guestData[currentGuestIndex]) {
            const data = guestData[currentGuestIndex];
            const actualGuestNum = shouldSkipFirstGuest ? currentGuestIndex + 2 : currentGuestIndex + 1;
            
            if (totalFormsToShow > 1) {
                if (data.firstName) document.getElementById(`firstName-${actualGuestNum}`).value = data.firstName;
                if (data.lastName) document.getElementById(`lastName-${actualGuestNum}`).value = data.lastName;
                if (data.phoneNumber) document.getElementById(`phone-${actualGuestNum}`).value = data.phoneNumber;
            }
            
            if (currentGuestIndex === totalFormsToShow - 1 || guestNumber === 1) {
                if (data.specialRequests) document.getElementById('specialRequests').value = data.specialRequests;
                if (data.terms) document.getElementById('terms').checked = data.terms;
            }
        }
    }

    // Attach event listeners
    function attachEventListeners() {
        const form = document.getElementById('guest-form');
        const backBtn = document.getElementById('back-btn');
        
        // Form submission
        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            
            // Save current guest data
            guestData[currentGuestIndex] = collectFormData();
         
            
            // If last guest, submit all data
            if (currentGuestIndex === totalFormsToShow - 1) {
                // Store guest data in sessionStorage for after payment
                sessionStorage.setItem('guestData', JSON.stringify(guestData));
                const bookingSummary = JSON.parse(sessionStorage.getItem('bookingDraft'));

                try {
                    //Stripe payment prcocess
                    document.querySelector('.stripe').addEventListener('click', async (event) => {
                        event.preventDefault();

                            const response = await fetch('/api/create-checkout-session', {
                                method: 'POST',
                                headers: {
                                    'Content-Type': 'application/json',
                                },
                                body: JSON.stringify({
                                    room_id: bookingSummary.room_id,
                                    check_in: bookingSummary.check_in,
                                    check_out: bookingSummary.check_out,
                                    guest: bookingSummary.guest
                                })
                            });
                            
                            if (!response.ok) {
                                const errorData = await response.json();
                                throw new Error(errorData.error || 'Failed to create payment session');
                            }
                            
                            const data = await response.json();
                            const { url } = data;
                            window.location.href = url;
                            return;
                     
                    });

                     //Paymongo payment
                    document.querySelector('.gcash').addEventListener('click', async (event) => {
                        event.preventDefault();
                     
                        const bookingSummary = JSON.parse(sessionStorage.getItem('bookingDraft'));         
                        const response = await fetch('/paymongo/create/payment-intent', {
                            method: 'POST',
                            headers: {
                                'Content-Type': 'application/json',
                            },
                            body: JSON.stringify({
                                room_id: bookingSummary.room_id,
                                check_in: bookingSummary.check_in,
                                check_out: bookingSummary.check_out,
                                guest: bookingSummary.guest,
                            })
                        });
                        
                        if (!response.ok) {
                            const errorData = await response.json();
                            throw new Error(errorData.error || 'Failed to create payment session');
                        }

                        const data = await response.json();
                        
                        sessionStorage.setItem('payment_session_id', data.session_id);
                        
                        window.location.href = data.checkout_url;
                    
                    });
                  
                } catch (error) {
                    alert(`Payment Error: ${error.message}`);
                }

            } else {
                // Move to next guest
                currentGuestIndex++;
                renderGuestForm();
                populateFormData();
            }
        });
        
        // Back button
        backBtn.addEventListener('click', () => {
            if (currentGuestIndex > 0) {
                // Save current data before going back
                guestData[currentGuestIndex] = collectFormData();
                currentGuestIndex--;
                renderGuestForm();
                populateFormData();
            }
        });
    }

    // Initialize form
    renderGuestForm();
    populateFormData();



//Calculate the price in the backend
fetch("/api/booking/calculate", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            room_id: roomId,
            check_in: bookingSummary.check_in,
            check_out: bookingSummary.check_out,
            guest: Number(bookingSummary.guest)
        })
        })
        .then(response => {
            if (!response.ok) {
                throw new Error("Calculation failed");
            }
            return response.json();
        })
        .then(data => {
            price.textContent= "₱" + data.price_per_night;
            total.textContent=  "₱" + data.total;
        })
        .catch(error => console.log(error));

});