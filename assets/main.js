// Function to check if at least one of the tags is selected when creating a new post
function validateForm() {
    const checkboxes = document.querySelectorAll('input[type="checkbox"]');
    let checkedCount = 0;
    for (let i = 0; i < checkboxes.length; i++) {
        if (checkboxes[i].checked) {
            checkedCount++;
        }
    }
    if (checkedCount === 0) {
        alert("Please select at least one tag.");
        return false; // Prevent form submission
    }
    return true; // Allow form submission
}

// Function to convert ISO 8601 date-time format to a readable format
function convertDateTime(dateTimeString) {
    const dateTime = new Date(dateTimeString);
    const options = { year: 'numeric', month: 'long', day: 'numeric', hour: 'numeric', minute: 'numeric', second: 'numeric' };
    return dateTime.toLocaleDateString('en-US', options);
}

// Function to update the date-time elements on page load
window.addEventListener('load', function () {
    const dateElements = document.querySelectorAll('.creationDate');
    dateElements.forEach(function (element) {
        const isoDateTime = element.textContent.trim();
        element.textContent = convertDateTime(isoDateTime);
    });
});
