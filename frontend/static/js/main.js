function filterTable() {

    // Declare variables
    let input, filter, table, tr, td, cell, i;
    input = document.getElementById("filter_leaderboard");
    filter = input.value.toUpperCase();
    table = document.getElementById("leaderboard");
    tr = table.getElementsByTagName("tr");

    // Loop through rows
    for (i = 1; i < tr.length; i++) {
        // Hide row initially
        tr[i].style.display = "none";

        td = tr[i].getElementsByTagName("td");
        cell = tr[i].getElementsByTagName("td")[1];
        if (cell) {
            if (cell.innerText.toUpperCase().indexOf(filter) > -1) {
                tr[i].style.display = "";
            }
        }
        cell = tr[i].getElementsByTagName("td")[2];
        if (cell) {
            if (cell.innerText.toUpperCase().indexOf(filter) > -1) {
                tr[i].style.display = "";
            }
        }
    }
}
