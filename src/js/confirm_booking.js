document.addEventListener("DOMContentLoaded", () => {
    // Get the booking info from session storage
    const bookingInfo = sessionStorage.getItem("bookingDraft");
    const bookingSummary = JSON.parse(bookingInfo);

    // Populate the summary values
    const roomId = bookingSummary.room_id;
    let price = document.getElementById("price");
    let total = document.getElementById("total");

    document.getElementById("room-type").textContent = bookingSummary.room_type;
    document.getElementById("check-in").textContent  = bookingSummary.check_in;
    document.getElementById("check-out").textContent = bookingSummary.check_out;
    document.getElementById("guest").textContent     = bookingSummary.guest;

    const guestNumber       = Number(bookingSummary.guest);
    const container         = document.getElementById('guest-container');

    // Store all guest data
    let guestData          = [];
    let currentGuestIndex  = 0;

    const shouldSkipFirstGuest = guestNumber >= 2;
    const totalFormsToShow     = shouldSkipFirstGuest ? guestNumber - 1 : 1;

    // ─── Guest-specific fields ───────────────────────────────────────────────
    function getGuestFields(guestNum) {
        return `
            <div class="flex flex-col sm:flex-row gap-4 mb-6">
                <div class="flex-1">
                    <label class="block text-sm font-medium text-gray-900 mb-2">
                        First Name <span class="text-red-600">*</span>
                    </label>
                    <input type="text" name="firstName" id="firstName-${guestNum}"
                        class="w-full px-4 py-3 bg-gray-100 border border-gray-100 rounded focus:outline-none focus:border-blue-500" required>
                </div>
            </div>

            <div class="flex flex-col sm:flex-row gap-4 mb-6">
                <div class="flex-1">
                    <label class="block text-sm font-medium text-gray-900 mb-2">
                        Last Name <span class="text-red-600">*</span>
                    </label>
                    <input type="text" name="lastName" id="lastName-${guestNum}"
                        class="w-full px-4 py-3 bg-gray-100 border border-gray-100 rounded focus:outline-none focus:border-blue-500" required>
                </div>
            </div>

            <div class="flex flex-col sm:flex-row gap-4 mb-6">
                <div class="flex-1">
                    <label class="block text-sm font-medium text-gray-900 mb-2">
                        Phone Number <span class="text-red-600">*</span>
                    </label>
                    <input type="tel" name="phone" id="phone-${guestNum}"
                        class="w-full px-4 py-3 bg-gray-100 border border-gray-100 rounded focus:outline-none focus:border-blue-500" required>
                </div>
            </div>
        `;
    }

    // ─── Common fields (special requests + terms) ────────────────────────────
    function getCommonFields() {
        return `
            <div class="mb-6">
                <label class="block text-sm font-medium text-gray-900 mb-2">
                    Special Requests (Optional)
                </label>
                <textarea rows="5" name="specialRequests" id="specialRequests"
                    class="w-full px-4 py-3 bg-gray-100 border border-gray-100 rounded focus:outline-none focus:border-blue-500 resize-none"
                    placeholder="Any special requests or requirements..."></textarea>
            </div>

            <div class="mb-6">
                <label class="flex items-start">
                    <input type="checkbox" name="terms" id="terms" class="mt-1 mr-3" required>
                    <span class="text-sm text-gray-900">
                        I agree to the terms and conditions and the cancellation policy
                        <span class="text-red-600">*</span>
                    </span>
                </label>
            </div>
        `;
    }

    // ─── Render form for current guest ───────────────────────────────────────
    function renderGuestForm() {
        const isLastGuest  = currentGuestIndex === totalFormsToShow - 1;
        const isFirstGuest = currentGuestIndex === 0;
        const actualGuestNum = shouldSkipFirstGuest ? currentGuestIndex + 2 : currentGuestIndex + 1;

        container.innerHTML = `
            <h1 class="text-3xl font-semibold text-gray-900 mb-2">
                ${guestNumber > 1
                    ? `Guest Information for Guest No. ${actualGuestNum}`
                    : 'Complete Your Booking'}
            </h1>
            <p class="text-gray-600 mb-8">
                ${guestNumber > 1
                    ? `Please provide details for guest ${actualGuestNum} of ${guestNumber}`
                    : 'Almost there! Just a few final details'}
            </p>

            ${guestNumber > 1 ? `
                <div class="mb-6">
                    <div class="flex items-center justify-between">
                        <span class="text-sm font-medium text-gray-600">Progress</span>
                        <span class="text-sm font-medium text-gray-900">${currentGuestIndex + 1} / ${totalFormsToShow}</span>
                    </div>
                    <div class="mt-2 w-full bg-gray-200 rounded-full h-2">
                        <div class="bg-gray-900 h-2 rounded-full transition-all duration-300"
                            style="width: ${((currentGuestIndex + 1) / totalFormsToShow) * 100}%">
                        </div>
                    </div>
                </div>
            ` : ''}

            <form id="guest-form">
                ${guestNumber >= 1 ? getGuestFields(actualGuestNum) : ''}
                ${isLastGuest || guestNumber === 1 ? getCommonFields() : ''}

                <div class="flex flex-col sm:flex-row gap-4">

                    <!-- Back Button -->
                    <button type="button" id="back-btn"
                        class="flex-1 px-6 py-3 bg-white border border-gray-300 rounded text-gray-700 font-medium hover:bg-gray-50 flex items-center justify-center ${isFirstGuest ? 'opacity-50 cursor-not-allowed' : ''}">
                        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
                        </svg>
                        Back
                    </button>

                    ${isLastGuest ? `
                        <!-- Card Payment -->
                        <button type="submit" id="stripe-btn"
                            class="stripe flex-1 px-6 py-3 bg-purple-600 text-white rounded font-medium hover:bg-gray-800 flex items-center justify-center">
                            Card
                            <svg class="w-5 h-5 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <rect x="2" y="5" width="20" height="14" rx="2" stroke-width="2"/>
                                <line x1="2" y1="10" x2="22" y2="10" stroke-width="2"/>
                            </svg>
                        </button>

                        <!-- E-Wallet Payment -->
                        <button type="submit" id="gcash-btn"
                            class="gcash flex-1 px-6 py-3 bg-purple-600 text-white rounded font-medium hover:bg-gray-800 flex items-center justify-center">
                            EWallet
                            <svg class="w-5 h-5 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <rect x="2" y="5" width="20" height="14" rx="2" stroke-width="2"/>
                                <line x1="2" y1="10" x2="22" y2="10" stroke-width="2"/>
                            </svg>
                        </button>
                    ` : `
                        <!-- Next Button -->
                        <button type="button" id="next-btn"
                            class="flex-1 px-6 py-3 bg-purple-600 text-white rounded font-medium hover:bg-gray-800 flex items-center justify-center">
                            Next
                            <svg class="w-5 h-5 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
                            </svg>
                        </button>
                    `}

                </div>
            </form>
        `;

        attachEventListeners();
    }

    // ─── Collect current form data ───────────────────────────────────────────
    function collectFormData() {
        const formData       = {};
        const actualGuestNum = shouldSkipFirstGuest ? currentGuestIndex + 2 : currentGuestIndex + 1;

        // Always try to collect guest fields if they exist in the DOM
        const firstNameEl  = document.getElementById(`firstName-${actualGuestNum}`);
        const lastNameEl   = document.getElementById(`lastName-${actualGuestNum}`);
        const phoneEl      = document.getElementById(`phone-${actualGuestNum}`);

        if (firstNameEl)  formData.firstName   = firstNameEl.value;
        if (lastNameEl)   formData.lastName    = lastNameEl.value;
        if (phoneEl)      formData.phoneNumber = phoneEl.value;

        // Collect common fields on last step
        if (currentGuestIndex === totalFormsToShow - 1 || guestNumber === 1) {
            formData.specialRequests = document.getElementById('specialRequests').value;
            formData.terms           = document.getElementById('terms').checked;
        }

        // DEBUG
        alert(`Collected Guest ${actualGuestNum}:\nFirst: ${formData.firstName}\nLast: ${formData.lastName}\nPhone: ${formData.phoneNumber}`);

        return formData;
    }

    // ─── Populate form with previously saved data ────────────────────────────
    function populateFormData() {
        if (!guestData[currentGuestIndex]) return;

        const data           = guestData[currentGuestIndex];
        const actualGuestNum = shouldSkipFirstGuest ? currentGuestIndex + 2 : currentGuestIndex + 1;

        const firstNameEl = document.getElementById(`firstName-${actualGuestNum}`);
        const lastNameEl  = document.getElementById(`lastName-${actualGuestNum}`);
        const phoneEl     = document.getElementById(`phone-${actualGuestNum}`);

        if (firstNameEl && data.firstName)   firstNameEl.value  = data.firstName;
        if (lastNameEl  && data.lastName)    lastNameEl.value   = data.lastName;
        if (phoneEl     && data.phoneNumber) phoneEl.value      = data.phoneNumber;

        if (currentGuestIndex === totalFormsToShow - 1 || guestNumber === 1) {
            const specialEl = document.getElementById('specialRequests');
            const termsEl   = document.getElementById('terms');
            if (specialEl && data.specialRequests) specialEl.value   = data.specialRequests;
            if (termsEl   && data.terms)           termsEl.checked   = data.terms;
        }
    }

    // ─── Handle payment ──────────────────────────────────────────────────────
    async function handleStripePayment() {
        try {
            const response = await fetch('/api/create-checkout-session', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    room_id:  bookingSummary.room_id,
                    check_in: bookingSummary.check_in,
                    check_out: bookingSummary.check_out,
                    guest:    bookingSummary.guest
                })
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.error || 'Failed to create payment session');
            }

            const data = await response.json();
            window.location.href = data.url;

        } catch (error) {
            alert(`Stripe Error: ${error.message}`);
        }
    }

    async function handleGcashPayment() {
        try {
            const response = await fetch('/paymongo/create/payment-intent', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    room_id:  bookingSummary.room_id,
                    check_in: bookingSummary.check_in,
                    check_out: bookingSummary.check_out,
                    guest:    bookingSummary.guest
                })
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.error || 'Failed to create payment session');
            }

            const data = await response.json();
            sessionStorage.setItem('payment_session_id', data.session_id);
            window.location.href = data.checkout_url;

        } catch (error) {
            alert(`GCash Error: ${error.message}`);
        }
    }

    // ─── Attach event listeners ──────────────────────────────────────────────
    function attachEventListeners() {
        const form    = document.getElementById('guest-form');
        const backBtn = document.getElementById('back-btn');
        const nextBtn = document.getElementById('next-btn');

        // Next button (not submit — just move forward)
        if (nextBtn) {
            nextBtn.addEventListener('click', () => {
                guestData[currentGuestIndex] = collectFormData();
                currentGuestIndex++;
                renderGuestForm();
                populateFormData();
            });
        }

        // Back button
        if (backBtn) {
            backBtn.addEventListener('click', () => {
                if (currentGuestIndex > 0) {
                    guestData[currentGuestIndex] = collectFormData();
                    currentGuestIndex--;
                    renderGuestForm();
                    populateFormData();
                }
            });
        }

        // Form submit (only on last step — Card or EWallet)
        if (form) {
            form.addEventListener('submit', async (e) => {
                e.preventDefault();

                // Save last guest data
                guestData[currentGuestIndex] = collectFormData();

                // Save all guest data to sessionStorage
                sessionStorage.setItem('guestData', JSON.stringify(guestData));

                // DEBUG — see all collected guests
                alert(`All guests saved:\n${JSON.stringify(guestData, null, 2)}`);

                // Determine which button was clicked
                const clickedBtn = e.submitter?.id;

                if (clickedBtn === 'stripe-btn') {
                    await handleStripePayment();
                } else if (clickedBtn === 'gcash-btn') {
                    await handleGcashPayment();
                }
            });
        }
    }

    // ─── Initialize ──────────────────────────────────────────────────────────
    renderGuestForm();
    populateFormData();

    // ─── Calculate price ─────────────────────────────────────────────────────
    fetch("/api/booking/calculate", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
            room_id:  roomId,
            check_in: bookingSummary.check_in,
            check_out: bookingSummary.check_out,
            guest:    Number(bookingSummary.guest)
        })
    })
    .then(response => {
        if (!response.ok) throw new Error("Calculation failed");
        return response.json();
    })
    .then(data => {
        price.textContent = "₱" + data.price_per_night;
        total.textContent = "₱" + data.total;
    })
    .catch(error => console.log(error));

});