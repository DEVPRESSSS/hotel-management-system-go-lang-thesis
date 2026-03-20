document.addEventListener("DOMContentLoaded", () => {
    let id = "";
    let dataTable = null;

    const categoryModal = document.getElementById("categoryModal");
    if (!categoryModal) return;

    const btnText     = document.getElementById("btn-text");
    const headerTitle = document.getElementById("header-action");
    const tableElement = document.querySelector("#default-table");
    const form        = document.getElementById("upsertform");
    const nameInput   = document.getElementById("name");
    const inputId     = document.getElementById("input-id");

    // ── Modal ─────────────────────────────────────────────────────────────────
    function openModal() {
        categoryModal.classList.remove("hidden");
        categoryModal.classList.add("flex");
    }

    function closeModal() {
        categoryModal.classList.add("hidden");
        categoryModal.classList.remove("flex");
        form.reset();
        id = "";
        inputId.value = "";
    }

    window.closeModal = closeModal;

    window.createModal = function () {
        headerTitle.innerText = "Create food category";
        btnText.innerText = "Create";
        id = "";
        openModal();
    };

    window.editModal = function (categoryId, categoryName) {
        headerTitle.innerText = "Edit food category";
        btnText.innerText = "Update";
        id = categoryId;
        inputId.value = categoryId;
        nameInput.value = categoryName;
        openModal();
    };

    // ── Table ─────────────────────────────────────────────────────────────────
    function refreshTable() {
    if (dataTable) {
        dataTable.destroy();
        dataTable = null;
    }

    fetch("/api/food/category")
        .then(res => {
            if (!res.ok) throw new Error("Failed to load food categories");
            return res.json();
        })
        .then(foodcategory => {
            // 2. Re-query tbody AFTER destroy — the original element is back
            const currentTbody = document.querySelector("#default-table tbody");
            if (!currentTbody) return console.error("tbody not found");

            currentTbody.innerHTML = foodcategory.map(category => `
                <tr>
                    <td class="px-4 py-3">${category.foodcategoryId}</td>
                    <td class="px-4 py-3">${category.name}</td>
                    <td class="px-4 py-3">
                         <button class="update-btn px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600 mr-2" onclick="editModal('${category.foodcategoryId}', '${category.name}')" data-id="${category.foodcategoryId}">Edit</button>
                         <button class="delete-btn px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600" data-id="${category.foodcategoryId}">Delete</button>
                    </td>
                </tr>
            `).join("");

            // 3. Init DataTable AFTER populating
            if (window.simpleDatatables && tableElement) {
                dataTable = new simpleDatatables.DataTable(tableElement, {
                    searchable: true,
                    paging: true,
                    perPage: 10,
                    perPageSelect: [5, 10, 20, 50],
                    sortable: true,
                });
            }
        })
        .catch(error => {
           // notification("error", error);
        });
    }
    refreshTable();

    // ── Delete ────────────────────────────────────────────────────────────────
    window.deleteCategory = function (categoryId) {
        if (!confirm("Are you sure you want to delete this category?")) return;

        fetch(`/api/delete/foodcategory/${categoryId}`, {
            method: "DELETE",
            headers: { "Content-Type": "application/json" },
        })
            .then(res => res.json().then(data => ({ ok: res.ok, data })))
            .then(({ ok, data }) => {
                if (!ok) throw new Error(data.error || "Failed to delete");
                notification("success", data.success || "Deleted successfully");
                refreshTable();
            })
            .catch(err => notification("error", err.message));
    };

    // ── Submit ────────────────────────────────────────────────────────────────
    form.addEventListener("submit", e => {
        e.preventDefault();

        const payload = {
            name: nameInput.value.trim(),  
        };

        if (!payload.name) {
            notification("error", "Category name is required");
            return;
        }

        const url    = id ? `/api/food/category/${id}` : "/api/food/category";
        const method = id ? "PUT" : "POST";

        const btnSubmit = document.getElementById("btn-submit");
        btnSubmit.disabled = true;
        btnText.innerText  = "Saving...";

        fetch(url, {
            method,
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload),
        })
            .then(res => res.json().then(data => ({ ok: res.ok, data })))
            .then(({ ok, data }) => {
                if (!ok) throw new Error(data.error || "Request failed");
                notification("success", data.success || "Saved successfully");
                closeModal();
                refreshTable();
            })
            .catch(err => notification("error", err.message))
            .finally(() => {
                btnSubmit.disabled = false;
                btnText.innerText  = id ? "Update" : "Create";
            });
    });

    // ── Close on backdrop click / Escape ──────────────────────────────────────
    categoryModal.addEventListener("click", e => {
        if (e.target === categoryModal) closeModal();
    });

    document.addEventListener("keydown", e => {
        if (e.key === "Escape" && !categoryModal.classList.contains("hidden")) {
            closeModal();
        }
    });
});