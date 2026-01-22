document.addEventListener("DOMContentLoaded", () => {

  /* =====================
     DOM ELEMENTS
  ====================== */
  const headerTitle = document.getElementById('header-action');
  const btnSubmit = document.getElementById('btn-submit');
  const btnText = document.getElementById('btn-text');
  const form = document.getElementById('upsertform');
  const tbody = document.getElementById('users-body');
  const tableElement = document.querySelector("#default-table");


  if (!tbody || !form || !headerTitle || !btnSubmit || !tableElement) {
    console.warn("User page elements not found. JS skipped.");
    console.log({tbody, form, headerTitle, btnSubmit, tableElement});
    return;
  }


  let dataTable = null;

  /* =====================
     FETCH USERS & INIT TABLE
  ====================== */
  fetch('/api/guest')
    .then(res => {
      if (!res.ok) throw new Error("API failed");
      return res.json();
    })
    .then(users => {
      // Populate table body
      tbody.innerHTML = users.map(user => `
        <tr>
          <td class="px-4 py-3">${user.userid}</td>
          <td class="px-4 py-3">${user.fullname || user.username}</td>
          <td class="px-4 py-3">${user.email}</td>
          <td class="px-4 py-3">${user.Role ? user.Role.rolename : 'N/A'}</td>
          <td class="px-4 py-3">
            ${user.locked === false 
              ? '<span class="px-2 py-1 text-xs font-semibold text-green-800 bg-green-100 rounded-full">Active</span>' 
              : '<span class="px-2 py-1 text-xs font-semibold text-red-800 bg-red-100 rounded-full">Locked</span>'}
          </td>
          <td class="px-4 py-3">${new Date(user.created_at).toLocaleDateString()}</td>
         
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

});
