document.addEventListener("DOMContentLoaded", () => {
  const table = document.querySelector("#default-table");

  if (!table) return;

  new simpleDatatables.DataTable(table, {
    searchable: true,
    paging: true,
    perPage: 10,
    perPageSelect: [5, 10, 20, 50],
    sortable: true
  });
});
