document.addEventListener("DOMContentLoaded",() =>{

        //  <td class="px-4 py-3">${role.rolename}</td>
        //   <td class="px-4 py-3">${role.rolename}</td>
        //   <td class="px-4 py-3">${role.rolename}</td>
        //   <td class="px-4 py-3">${role.rolename}</td>
        //   <td class="px-4 py-3">${role.rolename}</td>
        //   <td class="px-4 py-3">${role.rolename}</td>
        //   <td class="px-4 py-3">${role.rolename}</td>
const tbody = document.getElementById('users-body');
const tableElement = document.querySelector("#default-table");
fetch('/api/getlogs/')
    .then(res => {
      if (!res.ok) throw new Error("API failed");
      return res.json();
    })
    .then(logs => {
      // Populate table body
      tbody.innerHTML = logs.map(log => `
        <tr>
          <td class="px-4 py-3">${log.id}</td>
          <td class="px-4 py-3">${log.EntityType}</td>
          <td class="px-4 py-3">${log.EntityID}</td>
          <td class="px-4 py-3">${log.Action}</td>
          <td class="px-4 py-3">${log.Description}</td>
          <td class="px-4 py-3">${log.PerformedBy}</td>
          <td class="px-4 py-3">${new Date(log.CreatedAt).toLocaleDateString()}</td>
     
  
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
