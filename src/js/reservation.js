document.addEventListener("DOMContentLoaded", () => {

  /* =====================
     DOM ELEMENTS
  ====================== */
//   const headerTitle = document.getElementById('header-action');
//   const btnSubmit = document.getElementById('btn-submit');
//   const btnText = document.getElementById('btn-text');
//   const form = document.getElementById('upsertform');
  const tbody = document.getElementById('reservation');
  const tableElement = document.querySelector("#default-table");
  //const userModal = document.getElementById('userModal');

  if (!tbody || !tableElement) {
    console.warn("Role page elements not found. JS skipped.");
    return;
  }

  let id = "";
  let dataTable = null;

  /* =====================
     MODAL FUNCTIONS
  ====================== */
//   function openModal() {
//     userModal.classList.remove('hidden');
//     userModal.classList.add('flex');
//   }

//   function closeModal() {
//     userModal.classList.add('hidden');
//     userModal.classList.remove('flex');
//   }

//   window.closeModal = closeModal;

//   window.createModal = function () {
//     id = "";
//     headerTitle.innerText = "Create Role";
//     btnText.innerText = "Create";
//     form.reset();
//     openModal();
//   };

  /* =====================
     FETCH ROLES & INIT TABLE
  ====================== */
  fetch('/api/reservations')
    .then(res => {
      if (!res.ok) throw new Error("API failed");
      return res.json();
    })
    .then(data => {
      
        const reservations = data.reservations;
      // Populate table body
      tbody.innerHTML = reservations.map(r => `
        <tr class="text-gray-700 dark:text-gray-400">
                      <td class="px-4 py-3">
                        <div class="flex items-center text-sm">
                        <!-- Avatar with inset shadow -->
                        <div
                            class="relative hidden w-8 h-8 mr-3 rounded-full md:block"
                          >
                            <img
                              class="object-cover w-full h-full rounded-full"
                              src="https://images.unsplash.com/flagged/photo-1570612861542-284f4c12e75f?ixlib=rb-1.2.1&q=80&fm=jpg&crop=entropy&cs=tinysrgb&w=200&fit=max&ixid=eyJhcHBfaWQiOjE3Nzg0fQ"
                              alt=""
                              loading="lazy"
                            />
                            <div
                              class="absolute inset-0 rounded-full shadow-inner"
                              aria-hidden="true"
                            ></div>
                          </div>
                          <div>
                            <p class="font-semibold">${r.User.fullname}</p>
                            
                          </div>
                        </div>
                      </td>
                     <td class="px-4 py-3 text-sm">
                        ${r.book_id}
                      </td>
                      <td class="px-4 py-3 text-sm">
                        ${r.room_number}
                      </td>
                      <td class="px-4 py-3 text-sm">
                         ${r.room_type}
                      </td>
                      <td class="px-4 py-3 text-sm">
                        
                         ${new Date(r.check_in_date).toLocaleDateString()}
                      </td>
                      <td class="px-4 py-3 text-sm">
                        ${new Date(r.check_out_date).toLocaleDateString()}

                      </td>
                      <td class="px-4 py-3 text-sm">
                         ${r.num_guests}
                      </td>
                      <td class="px-4 py-3 text-sm">
                         ${r.price_per_night}
                      </td>
                      <td class="px-4 py-3 text-sm">
                         ${r.total_price}
                      </td>
                      <td class="px-4 py-3 text-xs">
                        ${
                            r.payment_status === "Paid"
                            ? `
                                <span
                                class="px-2 py-1 font-semibold leading-tight text-green-700 bg-green-100 rounded-full dark:bg-green-700 dark:text-green-100">
                                ${r.payment_status}
                                </span>
                            `
                            : `
                                <span
                                class="px-2 py-1 font-semibold leading-tight text-orange-700 bg-orange-100 rounded-full dark:text-white dark:bg-orange-600">
                                ${r.payment_status}
                                </span>
                            `
                        }
                        </td>

                      <td class="px-4 py-3 text-sm">
                         ${r.special_requests}
                      </td>
                       <td class="px-4 py-3 text-sm">
                         ${new Date(r.created_at).toLocaleDateString()}
                      </td>
                     <td class="px-4 py-3 text-sm">
                       
                     </td>
        </tr>
      `).join("");

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
//   form.addEventListener('submit', e => {
//     e.preventDefault();

//     const roleName = document.getElementById('rolename').value;
//     let uid = "";

//         if (id === "" || id === null) {
//             uid = uuidv4(); 
//         } else {
//             uid = id;      
//     }
//     const payload = {
//       roleid: uid,
//       roleName: roleName
//     };

//     const url = id ? `/api/updaterole/${id}` : '/api/createrole';
//     const method = id ? 'PUT' : 'POST';

//     fetch(url, {
//       method,
//       headers: { 'Content-Type': 'application/json' },
//       body: JSON.stringify(payload)
//     })
//     .then(res => {
//       return res.json().then(data => {
//         // Check if request was successful
//         if (!res.ok) {
//           // Throw error with server message
//           throw new Error(data.error || 'Request failed');
//         }
//         return data;
//       });
//     })
//     .then(data => {
//       notification("success", data.success);
//       closeModal();
//       setTimeout(() => location.reload(), 500);
//     })
//     .catch(err => {
//       notification("error", err.message);
//     });
  });

