// Used in leaderboard. It searches through the 'username' and 'topic' columns
// and hides all rows without a match, for real-time filtering.
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

// Used in phase 3 of quiz. Clones the element and stores it in the results
// table. It also removes the event-listener for activating this function.
function addToResults(element) {
    document.getElementById('results').append(element);
}

// Closes the flash message
$('.alert').alert();
