let deltrue;
function confirmDeletion() {
    if (confirm("Are you sure to delete?") == true) {
        deltrue = true;
        return deltrue;
    } else {
        deltrue = false;
        return deltrue;
    }
}

