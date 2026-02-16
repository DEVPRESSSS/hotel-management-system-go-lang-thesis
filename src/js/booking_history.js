document.addEventListener("DOMContentLoaded", function () {
    const container = document.getElementById("booking-container");

    // ── Status badge config ──────────────────────────────────────────
    function getStatusBadge(status) {
        const badges = {
            "check-in": `
                <span class="inline-flex items-center bg-green-600 gap-1.5 text-xs font-semibold text-white bg-emerald-50 border border-emerald-200 rounded-full px-3 py-0.5 mb-3">
                    <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
                    Active Stay
                </span>`,
            "pending": `
                <span class="inline-flex items-center gap-1.5 text-xs font-semibold text-amber-600 bg-amber-50 border border-amber-200 rounded-full px-3 py-0.5 mb-3">
                    <span class="w-1.5 h-1.5 rounded-full bg-amber-400"></span>
                    Upcoming
                </span>`,
            "completed": `
                <span class="inline-flex items-center gap-1.5 text-xs font-semibold text-gray-500 bg-gray-100 border border-gray-200 rounded-full px-3 py-0.5 mb-3">
                    <span class="w-1.5 h-1.5 rounded-full bg-gray-400"></span>
                    Completed
                </span>`,
            "cancelled": `
                <span class="inline-flex items-center gap-1.5 text-xs font-semibold text-red-500 bg-red-50 border border-red-200 rounded-full px-3 py-0.5 mb-3">
                    <span class="w-1.5 h-1.5 rounded-full bg-red-400"></span>
                    Cancelled
                </span>`,
        };
        return badges[status] ?? `
            <span class="inline-flex items-center gap-1.5 text-xs font-semibold text-gray-500 bg-gray-100 border border-gray-200 rounded-full px-3 py-0.5 mb-3">
                <span class="w-1.5 h-1.5 rounded-full bg-gray-400"></span>
                ${status}
            </span>`;`1`
    }

    // ── Action buttons per status ────────────────────────────────────
    function getActionButtons(status, bookId) {
        switch (status) {
            case "check-in":
                return `
                    <button class="text-xs border border-gray-300 text-gray-500 rounded p-2 hover:border-gray-500 hover:text-gray-700 transition">Receipt</button>
                    <button class="text-xs bg-purple-700 text-white rounded p-2 hover:bg-purple-600 transition">Request Service</button>`;
            case "pending":
                return `
                    <button class="text-xs border border-gray-300 text-gray-500 rounded p-2 hover:border-gray-500 hover:text-gray-700 transition">Modify</button>
                    <button class="text-xs bg-purple-700 text-white rounded p-2 hover:bg-purple-700 transition">View</button>`;
            case "completed":
                return `
                    <button class="text-xs border border-gray-300 text-gray-500 rounded p-2 hover:border-gray-500 hover:text-gray-700 transition">Receipt</button>
                    <button class="text-xs bg-gray-800 text-white rounded px-3 py-1.5 hover:bg-gray-900 transition">Book Again</button>`;
            case "cancelled":
                return `
                    <button class="text-xs border border-gray-300 text-gray-500 rounded p-2 hover:border-gray-500 hover:text-gray-700 transition">View</button>`;
            default:
                return `
                    <button class="text-xs border border-gray-300 text-gray-500 rounded p-2 hover:border-gray-500 hover:text-gray-700 transition">View</button>`;
        }
    }

    // ── Progress bar (only for active check-in) ──────────────────────
    function getProgressBar(res) {
        if (res.status !== "check-in") return "";

        const checkIn  = new Date(res.check_in_date);
        const checkOut = new Date(res.check_out_date);
        const today    = new Date();

        const totalNights   = Math.round((checkOut - checkIn) / (1000 * 60 * 60 * 24));
        const nightsElapsed = Math.round((today - checkIn)   / (1000 * 60 * 60 * 24));
        const nightsLeft    = totalNights - nightsElapsed;
        const pct           = Math.min(100, Math.round((nightsElapsed / totalNights) * 100));

        return `
            <div class="mt-4">
                <div class="flex justify-between text-xs text-gray-400 mb-1">
                    <span>Day ${nightsElapsed + 1} of ${totalNights}</span>
                    <span>${nightsLeft} night${nightsLeft !== 1 ? "s" : ""} remaining</span>
                </div>
                <div class="h-1.5 bg-gray-100 rounded-full overflow-hidden">
                    <div class="h-full bg-emerald-400 rounded-full" style="width: ${pct}%"></div>
                </div>
            </div>`;
    }

    // ── Format date ──────────────────────────────────────────────────
    function fmt(dateStr) {
        return new Date(dateStr).toLocaleDateString("en-US", {
            month: "short", day: "numeric", year: "numeric"
        });
    }

    // ── Fetch & render ───────────────────────────────────────────────
    fetch("/booking/history")
        .then(res => {
            if (!res.ok) throw new Error("Failed to fetch bookings");
            return res.json();
        })
        .then(data => {
            const reservations = data.reservations;

            if (!Array.isArray(reservations)) {
                console.error("reservations is not an array");
                return;
            }

            const fragment = document.createDocumentFragment();

            reservations.forEach(res => {
                const wrapper = document.createElement("div");
                wrapper.className = "max-w-full w-full lg:flex mb-6 shadow-md rounded-lg overflow-hidden";

                wrapper.innerHTML = `
                    <!-- Room image -->
                    <div
                        class="h-48 lg:h-auto lg:w-48 flex-none bg-cover bg-center"
                        style="background-image: url('https://images.unsplash.com/photo-1611892440504-42a792e24d32?w=400&q=80')"
                    ></div>

                    <!-- Card content -->
                    <div class="border-r border-b border-l border-gray-300 lg:border-l-0 lg:border-t lg:border-gray-300 bg-white lg:rounded-r p-5 flex flex-col justify-between leading-normal w-full">
                        <div class="mb-4">

                            ${getStatusBadge(res.status)}

                            <div class="text-gray-900 font-bold text-xl mb-1">
                                ${res.room_type} — ${res.room_number}
                            </div>

                            <p class="text-xs text-gray-400 mb-3">
                                Booking Ref: <span class="font-medium text-gray-600">${res.book_id}</span>
                            </p>

                            <div class="flex flex-wrap gap-x-6 gap-y-2 text-sm text-gray-600">
                                <div>
                                    <span class="text-xs uppercase tracking-wide text-gray-400 block">Check-in</span>
                                    <span class="font-medium text-gray-800">${fmt(res.check_in_date)}</span>
                                </div>
                                <div>
                                    <span class="text-xs uppercase tracking-wide text-gray-400 block">Check-out</span>
                                    <span class="font-medium text-gray-800">${fmt(res.check_out_date)}</span>
                                </div>
                                <div>
                                    <span class="text-xs uppercase tracking-wide text-gray-400 block">Guests</span>
                                    <span class="font-medium  p-2 text-gray-800">${res.num_guests}</span>
                                </div>
                                <div>
                                    <span class="text-xs uppercase tracking-wide text-gray-400 block">Total</span>
                                    <span class="font-bold text-gray-900">₱${res.total_price}</span>
                                </div>
                            </div>

                            ${getProgressBar(res)}

                        </div>

                        <!-- Footer: actions -->
                        <div class="flex items-center justify-end flex-wrap gap-2">
                            ${getActionButtons(res.status, res.book_id)}
                        </div>
                    </div>
                `;

                fragment.appendChild(wrapper);
            });

            container.innerHTML = "";
            container.appendChild(fragment);
        })
        .catch(error => console.error("Booking fetch error:", error));
});