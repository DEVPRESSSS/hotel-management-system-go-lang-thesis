document.addEventListener("DOMContentLoaded", () => {

    const categoryModal = document.getElementById("categoryModal");
    if (!categoryModal) return;

    const btnText = document.getElementById("btn-text");
    const headerTitle = document.getElementById("header-action");

    function openModal() {
        categoryModal.classList.remove("hidden");
        categoryModal.classList.add("flex");
    }

    function closeModal() {
        categoryModal.classList.add("hidden");
        categoryModal.classList.remove("flex");
    }

    window.closeModal = closeModal;

    window.createModal = function () {
        headerTitle.innerText = "Create food category";
        btnText.innerText = "Create";
        openModal();
    };

    const name = document.getElementById("name");

});