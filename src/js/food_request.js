// ─────────────────────────────────────────────────────────
//  STATE
// ─────────────────────────────────────────────────────────
let cart = JSON.parse(sessionStorage.getItem("foodCart") || "[]");

// modal state
let _modalFood = null; 
let _modalQty  = 1;     


// ─────────────────────────────────────────────────────────
//  INIT
// ─────────────────────────────────────────────────────────
document.addEventListener("DOMContentLoaded", () => {
    fetchFood();
    renderCart();

    // Cart toggle button
    document.getElementById("cartToggleBtn").addEventListener("click", openCart);

    // Close modal on backdrop click
    document.getElementById("foodModal").addEventListener("click", function (e) {
        if (e.target === this) closeModal();
    });
});


// ─────────────────────────────────────────────────────────
//  FETCH
// ─────────────────────────────────────────────────────────
function fetchFood() {
    fetch("/api/foodservices")
        .then(response => {
            if (!response.ok) throw new Error("Failed to fetch food services");
            return response.json();
        })
        .then(data => renderFood(data))
        .catch(error => console.error(error));
}



function renderFood(foods) {
    const grid = document.getElementById("foodGrid");
    grid.innerHTML = "";

    foods.forEach(food => {
        const isAvailable = food.status === "available";
        const imgSrc = food.image
            ? `/food_images/${food.image.split(/[\\/]/).pop()}`
            : '/src/placeholder.png';

        const card = document.createElement("div");
        card.className = "bg-white rounded-lg shadow-md overflow-hidden hover:shadow-xl hover:-translate-y-1 transition-all duration-300 flex flex-col";

        card.innerHTML = `
            <div class="relative overflow-hidden h-52 bg-gray-100">
                <img
                    src="${imgSrc}"
                    alt="${food.name}"
                    onerror="this.src='https://images.unsplash.com/photo-1504674900247-0877df9cc836?w=400&h=300&fit=crop'"
                    class="w-full h-full object-cover hover:scale-105 transition duration-500"
                >
                <span class="absolute top-3 left-3 bg-white/90 backdrop-blur-sm text-gray-700 text-xs font-semibold tracking-wider uppercase px-3 py-1 rounded-full shadow-sm">
                    ${food.FoodCategory?.name || "Food"}
                </span>
                <span class="absolute top-3 right-3 text-xs font-semibold px-3 py-1 rounded-full ${isAvailable ? "bg-green-100 text-green-700" : "bg-red-100 text-red-600"}">
                    ${isAvailable ? "● Available" : "● Unavailable"}
                </span>
            </div>

            <div class="p-5 flex flex-col flex-1">
                <div class="flex items-start justify-between gap-2 mb-2">
                    <h3 class="text-lg font-bold text-gray-800 leading-snug">${food.name}</h3>
                    <span class="bg-purple-600 text-white px-3 py-1 rounded-full text-sm font-semibold whitespace-nowrap">
                        ₱${parseFloat(food.price).toFixed(2)}
                    </span>
                </div>

                <p class="text-gray-500 text-sm mb-5 flex-1 leading-relaxed">
                    ${food.description || "Freshly prepared for your enjoyment."}
                </p>

                <button
                    data-food_id="${food.foodId}"
                    class="add-to-order-btn w-full py-3 rounded-lg text-sm font-semibold transition duration-300
                        ${isAvailable
                            ? "bg-purple-600 text-white hover:bg-purple-700 cursor-pointer"
                            : "bg-gray-100 text-gray-400 cursor-not-allowed opacity-60"}"
                    ${!isAvailable ? "disabled" : ""}
                >
                    ${isAvailable ? "Add to Cart" : "Unavailable"}
                </button>
            </div>
        `;

        grid.appendChild(card);
    });

    // ── Attach click → open modal (replaces the old navigation)
    grid.addEventListener("click", function (e) {
        const btn = e.target.closest(".add-to-order-btn");
        if (!btn || btn.disabled) return;

        const foodId = btn.dataset.food_id;
        if (!foodId) return;

        // Find the food object from the already-fetched list
        // We store them on window for easy access
        const food = (window._allFoods || []).find(f => String(f.foodId) === String(foodId));
        if (food) openModal(food);
    });

    // Cache foods for modal lookup
    window._allFoods = foods;
}


// ─────────────────────────────────────────────────────────
//  MODAL
// ─────────────────────────────────────────────────────────
function openModal(food) {
    _modalFood = food;
    _modalQty  = 1;

    const imgSrc = food.image
        ? `/food_images/${food.image.split(/[\\/]/).pop()}`
        : "https://images.unsplash.com/photo-1504674900247-0877df9cc836?w=800&h=600&fit=crop";

    document.getElementById("modalImg").src          = imgSrc;
    document.getElementById("modalImg").alt          = food.name;
    document.getElementById("modalImg").onerror      = () => {
        document.getElementById("modalImg").src = "https://images.unsplash.com/photo-1504674900247-0877df9cc836?w=800&h=600&fit=crop";
    };
    document.getElementById("modalCategory").textContent = food.FoodCategory?.name || "Food";
    document.getElementById("modalName").textContent      = food.name;
    document.getElementById("modalPrice").textContent     = `₱${parseFloat(food.price).toFixed(2)}`;
    document.getElementById("modalDesc").textContent      = food.description || "Freshly prepared for your enjoyment.";

    syncModalQty();

    // Show
    const modal = document.getElementById("foodModal");
    modal.classList.remove("modal-hidden");
    document.body.style.overflow = "hidden";
}

function closeModal() {
    document.getElementById("foodModal").classList.add("modal-hidden");
    document.body.style.overflow = "";
    _modalFood = null;
    _modalQty  = 1;
}

function changeModalQty(delta) {
    _modalQty = Math.max(1, _modalQty + delta);
    syncModalQty();
}

function syncModalQty() {
    document.getElementById("modalQty").textContent = _modalQty;
    document.getElementById("modalMinus").disabled  = _modalQty <= 1;
    document.getElementById("modalMinus").className = _modalQty <= 1
        ? "w-8 h-8 flex items-center justify-center rounded-lg text-gray-300 font-bold text-xl cursor-not-allowed"
        : "w-8 h-8 flex items-center justify-center rounded-lg text-gray-600 hover:bg-gray-200 hover:text-red-500 font-bold text-xl transition";

    if (_modalFood) {
        const lineTotal = parseFloat(_modalFood.price) * _modalQty;
        document.getElementById("modalLineTotal").textContent = `₱${lineTotal.toFixed(2)}`;
    }
}

function confirmAddToCart() {
    if (!_modalFood) return;

    const foodId = _modalFood.foodId; 

    addToCart(foodId, _modalFood.name, _modalFood.price, _modalQty);

    closeModal(); 

    const btn = document.querySelector(`.add-to-order-btn[data-food_id="${foodId}"]`);
    if (btn) {
        const original = btn.textContent;
        btn.textContent = "✓ Added!";
        btn.classList.replace("bg-purple-600", "bg-green-500");
        setTimeout(() => {
            btn.textContent = original;
            btn.classList.replace("bg-green-500", "bg-purple-600");
        }, 1200);
    }
}


// ─────────────────────────────────────────────────────────
//  CART SIDEBAR
// ─────────────────────────────────────────────────────────
function openCart() {
    const sidebar   = document.getElementById("cartSidebar");
    const backdrop  = document.getElementById("cartBackdrop");
    sidebar.classList.remove("cart-hidden");
    backdrop.classList.remove("hidden");
    document.body.style.overflow = "hidden";
}

function closeCart() {
    const sidebar   = document.getElementById("cartSidebar");
    const backdrop  = document.getElementById("cartBackdrop");
    sidebar.classList.add("cart-hidden");
    backdrop.classList.add("hidden");
    document.body.style.overflow = "";
}


// ─────────────────────────────────────────────────────────
//  CART LOGIC
// ─────────────────────────────────────────────────────────
function addToCart(id, name, price, qty = 1) {
    const existing = cart.find(c => String(c.id) === String(id));
    if (existing) {
        existing.qty += qty;
    } else {
        cart.push({ id: String(id), name, price: parseFloat(price), qty });
    }
    saveCart();
    renderCart();
}

function updateCartQty(id, delta) {
    const idx = cart.findIndex(c => String(c.id) === String(id));
    if (idx === -1) return;
    cart[idx].qty += delta;
    if (cart[idx].qty <= 0) cart.splice(idx, 1);
    saveCart();
    renderCart();
}

function clearCart() {
    cart = [];
    saveCart();
    renderCart();
}

function saveCart() {
    sessionStorage.setItem("foodCart", JSON.stringify(cart));
}

async function handleCheckout(event) {
   // event.preventDefault(); // add this
    if (!cart.length) return;

    const response = await fetch('/paymongo/food/checkout', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            items: cart.map(item => ({ id: item.id, qty: item.qty }))
        })
    });

    const data = await response.json();
    console.log("Checkout response:", data);

    if (!response.ok) {
        alert(data.error || "Checkout failed");
        return;
    }


    alert("Redirecting to: " + data.checkout_url); 

    window.location.href = data.checkout_url;
}



function renderCart() {
    const container  = document.getElementById("cartItemsContainer");
    const emptyEl    = document.getElementById("cartEmpty");
    const summaryEl  = document.getElementById("cartSummary");
    const badgeEl    = document.getElementById("cartBadge");
    const countEl    = document.getElementById("sidebarCount");

    const totalQty   = cart.reduce((s, i) => s + i.qty, 0);
    const totalPrice = cart.reduce((s, i) => s + i.price * i.qty, 0);

    // Badge
    if (totalQty > 0) {
        badgeEl.textContent = totalQty;   badgeEl.classList.remove("hidden");
        countEl.textContent = totalQty;   countEl.classList.remove("hidden");
    } else {
        badgeEl.classList.add("hidden");
        countEl.classList.add("hidden");
    }

    // Remove old rows (keep emptyEl)
    container.querySelectorAll(".cart-item").forEach(el => el.remove());

    if (!cart.length) {
        emptyEl.style.display  = "";
        summaryEl.classList.add("hidden");
        return;
    }

    emptyEl.style.display = "none";
    summaryEl.classList.remove("hidden");
    document.getElementById("cartSubtotal").textContent = `₱${totalPrice.toFixed(2)}`;
    document.getElementById("cartTotal").textContent    = `₱${totalPrice.toFixed(2)}`;

    cart.forEach(item => {
        const row = document.createElement("div");
        row.className = "cart-item flex items-center gap-3 bg-gray-50 rounded-xl p-3";
        row.innerHTML = `
            <div class="flex-1 min-w-0">
                <p class="text-sm font-semibold text-gray-800 truncate">${item.name}</p>
                <p class="text-xs text-purple-600 font-medium mt-0.5">₱${(item.price * item.qty).toFixed(2)}</p>
            </div>
            <div class="flex items-center gap-1.5 shrink-0">
                <button
                    onclick="updateCartQty('${item.id}', -1)"
                    class="w-7 h-7 p-2 bg-white border border-gray-200 rounded-lg text-gray-500 hover:text-red-500 hover:border-red-200 font-bold text-base flex items-center justify-center transition">
                    −
                </button>
                <span class="w-5 text-center text-sm font-bold text-gray-700">${item.qty}</span>
                <button
                    onclick="updateCartQty('${item.id}', 1)"
                    class="w-7 h-7 p-2 bg-purple-600 rounded-lg text-white font-bold text-base flex items-center justify-center hover:bg-purple-700 transition">
                    +
                </button>
            </div>
        `;
        container.appendChild(row);
    });
}