document.addEventListener("DOMContentLoaded", () => {
    let id = "";

    const categoryModal = document.getElementById("categoryModal");
    if (!categoryModal) return;

    const btnText     = document.getElementById("btn-text");
    const headerTitle = document.getElementById("header-action");
    const tbody       = document.getElementById("category-body");
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
        id = "";
        headerTitle.innerText = "Create food category";
        btnText.innerText = "Create";
        form.reset();
        openModal();
    };

    window.editModal = function (categoryId, categoryName) {
        id = categoryId;
        inputId.value = categoryId;
        nameInput.value = categoryName;
        headerTitle.innerText = "Edit food category";
        btnText.innerText = "Update";
        openModal();
    };

    window.deleteCategory = function (categoryId) {
            Swal.fire({
                title: "Are you sure?",
                text: "This action cannot be undone!",
                icon: "warning",
                showCancelButton: true,
                confirmButtonColor: "#d33",
                cancelButtonColor: "#3085d6",
                confirmButtonText: "Yes, delete it!"
            }).then(result => {
                if (result.isConfirmed) {
                 fetch(`/api/food/category/${categoryId}`, {
                        method: "DELETE",
                        headers: { "Content-Type": "application/json" },
                    })
                    .then(res => res.json().then(data => 
                        ({ ok: res.ok, data })))
                    .then(({ ok, data }) => {
                            if (!ok) throw new Error(data.error || "Failed to delete");
                            notification("success", data.success || "Deleted successfully");
                            setTimeout(() => window.location.reload(), 800);
                        })
                    .catch(err => notification("error", err.message));
                }
            });

     
    };

    fetch("/api/food/category")
        .then(res => {
            if (!res.ok) throw new Error("Failed to load food categories");
            return res.json();
        })
        .then(foodcategory => {
            tbody.innerHTML = foodcategory.map(category => `
                <tr>
                    <td class="px-4 py-3">${category.foodcategoryId}</td>
                    <td class="px-4 py-3">${category.name}</td>
                    <td class="px-4 py-3 text-center">
                        <button onclick="editModal('${category.foodcategoryId}', '${category.name}')"
                            class="px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600 mr-2 text-xs">
                            Edit
                        </button>
                        <button onclick="deleteCategory('${category.foodcategoryId}')"
                            class="px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600 text-xs">
                            Delete
                        </button>
                    </td>
                </tr>
            `).join("");

            if (window.simpleDatatables && tableElement) {
                new simpleDatatables.DataTable(tableElement, {
                    searchable: true,
                    paging: true,
                    perPage: 10,
                    perPageSelect: [5, 10, 20, 50],
                    sortable: true,
                });
            }
        })
        .catch(err => console.error(err));

    // ── Submit ────────────────────────────────────────────────────────────────
    form.addEventListener("submit", e => {
        e.preventDefault();

        const payload = { name: nameInput.value.trim() };

        if (!payload.name) {
            notification("error", "Category name is required");
            return;
        }

        const url    = id ? `/api/food/category/${id}` : "/api/food/category";
        const method = id ? "PUT" : "POST";

        const btnSubmit    = document.getElementById("btn-submit");
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
                setTimeout(() => window.location.reload(), 800);
            })
            .catch(err => {
                notification("error", err.message);
                btnSubmit.disabled = false;
                btnText.innerText  = id ? "Update" : "Create";
            });
    });

    categoryModal.addEventListener("click", e => {
        if (e.target === categoryModal) closeModal();
    });

    document.addEventListener("keydown", e => {
        if (e.key === "Escape" && !categoryModal.classList.contains("hidden")) closeModal();
    });
});