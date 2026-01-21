import { DataTable } from "simple-datatables";
const dataTable = new DataTable("#default-table", {
    searchable:true,
    sensitivity: "base",
    searchQuerySeparator: " ",
    paging: true, 
    perPage: 10, 
    perPageSelect: [5, 10, 20, 50], 
    firstLast: true, 
    nextPrev: true, 
    sortable: true, // enable or disable sorting
    locale: "en-US", // set the locale for sorting
    numeric: true, // enable or disable numeric sorting
    caseFirst: "false", // set the case first for sorting (upper, lower)
    ignorePunctuation: true // enable or disable punctuation sorting
});


