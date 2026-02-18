document.addEventListener("DOMContentLoaded", () => {

  /* =====================
     DOM ELEMENTS
  ====================== */
  const headerTitle = document.getElementById('header-action');
  const btnSubmit = document.getElementById('btn-submit');
  const form = document.getElementById('upsertform');
  const tbody = document.getElementById('food-body');
  const tableElement = document.querySelector("#default-table");

  if (!tbody || !form || !headerTitle || !btnSubmit || !tableElement) {
    console.warn("Food service page elements not found. JS skipped.");
    console.log({tbody, form, headerTitle, btnSubmit, tableElement});
    return;
  }

  let id = "";
  let dataTable = null;

  /* =====================
     MODAL FUNCTIONS
  ====================== */
  function openModal() {
    userModal.classList.remove('hidden');
    userModal.classList.add('flex');
  }

  function closeModal() {
    userModal.classList.add('hidden');
    userModal.classList.remove('flex');
  }

  window.closeModal = closeModal;

  window.createModal = function () {
    id = "";
    headerTitle.innerText = "Create Service";
    btnSubmit.innerText = "Create";
    form.reset();
    openModal();
  };

  /* =====================
     FETCH FOOD SERVICES & INIT TABLE
  ====================== */
  fetch('/api/foodservices')
    .then(res => {
      if (!res.ok) throw new Error("API failed");
      return res.json();
    })
    .then(foodservices => {
      // Populate table body
      tbody.innerHTML = foodservices.map(a => {
        const imagePath = a.image 
            ? `/food_images/${a.image.split(/[\\/]/).pop()}`
            : '/src/placeholder.png'; 

        return `
            <tr>
            <td class="px-4 py-3">${a.foodId}</td>
            <td class="px-4 py-3">
                <img src="${imagePath}" class="rounded-sm w-16 h-16 object-cover"
                    onerror="this.onerror=null; this.src='/src/placeholder.png'"/>
            </td>
            <td class="px-4 py-3">${a.name}</td>
            <td class="px-4 py-3">${a.description}</td>
            <td class="px-4 py-3">${a.FoodCategory.name}</td>
            <td class="px-4 py-3">${a.price}</td>
            <td class="px-4 py-3">${a.status}</td>
            <td class="px-4 py-3">${new Date(a.created_at).toLocaleDateString()}</td>
            <td class="px-4 py-3">
                <button class="update-btn px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600 mr-2" data-id="${a.foodId}">Edit</button>
                <button class="delete-btn px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600" data-id="${a.foodId}">Delete</button>
            </td>
            </tr>
        `;
    }).join("");
      // Initialize DataTable AFTER populating data
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
    .catch(console.error);

  /* =====================
     FORM SUBMIT
  ====================== */
  form.addEventListener('submit', e => {
    e.preventDefault();

    const name = document.getElementById('name').value;
    const description = document.getElementById('description').value;
    const image = document.getElementById('image').value;
    const category = document.getElementById('food-category-id').value;
    const price = document.getElementById('price').value;
    const status = document.getElementById('status').value;
    const payload = {
      name: name,
      description: description,
      image: image,
      foodcategoryid: category,
      price: price,
      status: status,
    };

    const url = id ? `/api/update/foodservice/${id}` : '/api/create/foodservice';
    const method = id ? 'PUT' : 'POST';
    const formData = new FormData(document.getElementById('upsertform'));

    fetch(url, {
      method,
      body:formData
    })
    .then(res => {
        return res.json().then(data => {
          // Check if request was successful
          if (!res.ok) {
            // Throw error with server message
            throw new Error(data.error || 'Request failed');
          }
          return data;
       });
    })
    .then(data => {
      notification("success", data.success);
      closeModal();
      setTimeout(() => location.reload(), 500);
    })
    .catch(err => {
      notification("error", err)
    });
  });

  /* =====================
     TABLE CLICK HANDLER
  ====================== */
  tbody.addEventListener('click', e => {

    if (e.target.classList.contains('update-btn')) {
      id = e.target.dataset.id;
      fetch(`/api/foodservice/${id}`)
        .then(res => res.json())
        .then(data => {
          document.getElementById('name').value = data.success.name;
          document.getElementById('description').value = data.success.description;
        //   document.getElementById('image').value = data.success.image;
          document.getElementById('food-category-id').value = data.success.	foodcategoryid;
          document.getElementById('price').value = data.success.price;
          document.getElementById('status').value =  data.success.status;
          headerTitle.innerText = "Update Food service";
          btnSubmit.innerText = "Update";
          openModal();
        })
        .catch(err => notification("error", err.message));
    }

    if (e.target.classList.contains('delete-btn')) {
      const aid = e.target.dataset.id;
      console.log(aid);
      Swal.fire({
        title: "Are you sure you want to delete this record?",
        text: "This action cannot be undone!",
        icon: "warning",
        showCancelButton: true,
        confirmButtonColor: "#d33",
        cancelButtonColor: "#3085d6",
        confirmButtonText: "Yes, delete it!"
      }).then(r => {
        if (r.isConfirmed) {
          fetch(`/api/delete/foodservice/${aid}`, { method: 'DELETE' })
            .then(res => {
              if (!res.ok) {
                throw new Error('Delete failed');
           
              }
               return;
            })
            .then(() => {
              Swal.fire("Deleted!", "Food service has been deleted.", "success");
              location.reload();
            })
            .catch(err => notification("error", err.message))
        }
      });
    }
  });

});

