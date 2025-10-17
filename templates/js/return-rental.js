// Get DOM elements
const form = document.getElementById('returnForm');
const customerSelect = document.getElementById('customer');
const rentalSelect = document.getElementById('rental');
const quantityReturnInput = document.getElementById('quantityReturn');
const returnDateInput = document.getElementById('returnDate');
const messageDiv = document.getElementById('message');
const loadingDiv = document.getElementById('loading');
const submitBtn = document.getElementById('submitBtn');
const cancelBtn = document.getElementById('cancelBtn');
const rentalDetails = document.getElementById('rentalDetails');
const returnSummary = document.getElementById('returnSummary');
const maxQuantitySpan = document.getElementById('maxQuantity');

// Detail elements
const detailItem = document.getElementById('detailItem');
const detailSize = document.getElementById('detailSize');
const detailQtyRent = document.getElementById('detailQtyRent');
const detailQtyReturned = document.getElementById('detailQtyReturned');
const detailQtyRemaining = document.getElementById('detailQtyRemaining');
const detailDateBegin = document.getElementById('detailDateBegin');
const detailDateEnd = document.getElementById('detailDateEnd');
const detailStatus = document.getElementById('detailStatus');

// Summary elements
const summaryCustomer = document.getElementById('summaryCustomer');
const summaryQuantity = document.getElementById('summaryQuantity');
const summaryReturnDate = document.getElementById('summaryReturnDate');
const summaryDaysRented = document.getElementById('summaryDaysRented');
const summaryDaysLate = document.getElementById('summaryDaysLate');
const lateWarning = document.getElementById('lateWarning');

// Store data
let customers = [];
let rentals = [];
let currentRental = null;

// Load data on page load
document.addEventListener('DOMContentLoaded', function() {
    loadCustomers();
    setCurrentDateTime();
});

// Set current datetime as default
function setCurrentDateTime() {
    const now = new Date();
    const offset = now.getTimezoneOffset() * 60000;
    const localISOTime = new Date(now - offset).toISOString().slice(0, 16);
    returnDateInput.value = localISOTime;
}

// Load customers from API
async function loadCustomers() {
    try {
        const response = await fetch('/api/customers');
        if (response.ok) {
            customers = await response.json();
            customers.forEach(customer => {
                const option = document.createElement('option');
                option.value = customer.id;
                option.textContent = `${customer.cust_name} - ${customer.cust_phone}`;
                customerSelect.appendChild(option);
            });
        } else {
            showMessage('Failed to load customers', 'error');
        }
    } catch (error) {
        showMessage(`Error loading customers: ${error.message}`, 'error');
    }
}

// Load active rentals when customer is selected
customerSelect.addEventListener('change', async function() {
    rentalSelect.innerHTML = '<option value="">Loading...</option>';
    rentalSelect.disabled = true;
    rentalDetails.style.display = 'none';
    returnSummary.style.display = 'none';
    quantityReturnInput.disabled = true;
    returnDateInput.disabled = true;
    submitBtn.disabled = true;

    if (!this.value) {
        rentalSelect.innerHTML = '<option value="">Select a customer first...</option>';
        return;
    }

    try {
        const response = await fetch('/api/rentals');
        if (response.ok) {
            const allRentals = await response.json();
            // Filter for active rentals (status 1 = rent) for selected customer
            rentals = allRentals.filter(rental =>
                rental.id_clothing_customer === parseInt(this.value) &&
                rental.clothes_rent_status === 1 &&
                rental.clothes_qty_rent > rental.clothes_qty_return
            );

            if (rentals.length === 0) {
                rentalSelect.innerHTML = '<option value="">No active rentals found</option>';
                return;
            }

            rentalSelect.innerHTML = '<option value="">Select a rental...</option>';
            rentals.forEach(rental => {
                const option = document.createElement('option');
                option.value = rental.id;
                const remaining = rental.clothes_qty_rent - rental.clothes_qty_return;
                option.textContent = `Rental #${rental.id} - Qty: ${remaining} remaining`;
                rentalSelect.appendChild(option);
            });

            rentalSelect.disabled = false;
        } else {
            showMessage('Failed to load rentals', 'error');
        }
    } catch (error) {
        showMessage(`Error loading rentals: ${error.message}`, 'error');
    }
});

// Show rental details when rental is selected
rentalSelect.addEventListener('change', function() {
    if (!this.value) {
        rentalDetails.style.display = 'none';
        returnSummary.style.display = 'none';
        quantityReturnInput.disabled = true;
        returnDateInput.disabled = true;
        submitBtn.disabled = true;
        return;
    }

    currentRental = rentals.find(r => r.id === parseInt(this.value));
    if (currentRental) {
        displayRentalDetails(currentRental);
        quantityReturnInput.disabled = false;
        returnDateInput.disabled = false;

        const remaining = currentRental.clothes_qty_rent - currentRental.clothes_qty_return;
        quantityReturnInput.max = remaining;
        quantityReturnInput.value = Math.min(1, remaining);
        maxQuantitySpan.textContent = remaining;
    }
});

// Display rental details
function displayRentalDetails(rental) {
    const remaining = rental.clothes_qty_rent - rental.clothes_qty_return;

    // Note: Item and size names would need to be fetched or joined from the rental data
    detailItem.textContent = `Subcategory ID: ${rental.id_clothing_category_sub}`;
    detailSize.textContent = `Size ID: ${rental.id_clothing_size}`;
    detailQtyRent.textContent = rental.clothes_qty_rent;
    detailQtyReturned.textContent = rental.clothes_qty_return;
    detailQtyRemaining.textContent = remaining;
    detailDateBegin.textContent = formatDate(new Date(rental.clothes_rent_date_begin));
    detailDateEnd.textContent = formatDate(new Date(rental.clothes_rent_date_end));

    const statusMap = {
        1: 'Active',
        2: 'Returned',
        3: 'Cancelled',
        4: 'Not Returned',
        5: 'Loss'
    };
    detailStatus.textContent = statusMap[rental.clothes_rent_status] || 'Unknown';

    rentalDetails.style.display = 'block';
    updateReturnSummary();
}

// Update return summary
function updateReturnSummary() {
    if (!currentRental || !quantityReturnInput.value || !returnDateInput.value) {
        returnSummary.style.display = 'none';
        submitBtn.disabled = true;
        return;
    }

    const customer = customers.find(c => c.id === currentRental.id_clothing_customer);
    const returnDate = new Date(returnDateInput.value);
    const beginDate = new Date(currentRental.clothes_rent_date_begin);
    const endDate = new Date(currentRental.clothes_rent_date_end);

    summaryCustomer.textContent = customer ? customer.cust_name : 'Unknown';
    summaryQuantity.textContent = `${quantityReturnInput.value} item(s)`;
    summaryReturnDate.textContent = formatDate(returnDate);

    const daysRented = Math.ceil((returnDate - beginDate) / (1000 * 60 * 60 * 24));
    summaryDaysRented.textContent = `${daysRented} day(s)`;

    // Check if late
    if (returnDate > endDate) {
        const daysLate = Math.ceil((returnDate - endDate) / (1000 * 60 * 60 * 24));
        summaryDaysLate.textContent = `${daysLate} day(s) late`;
        lateWarning.style.display = 'block';
    } else {
        lateWarning.style.display = 'none';
    }

    returnSummary.style.display = 'block';
    submitBtn.disabled = false;
}

// Format date for display
function formatDate(date) {
    return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });
}

// Listen for changes to update summary
quantityReturnInput.addEventListener('input', function() {
    if (currentRental) {
        const remaining = currentRental.clothes_qty_rent - currentRental.clothes_qty_return;
        if (parseInt(this.value) > remaining) {
            this.value = remaining;
        }
        updateReturnSummary();
    }
});
returnDateInput.addEventListener('change', updateReturnSummary);

// Cancel button
cancelBtn.addEventListener('click', function() {
    if (confirm('Are you sure you want to cancel? Any unsaved changes will be lost.')) {
        window.location.href = '/';
    }
});

// Form submission
form.addEventListener('submit', async function(e) {
    e.preventDefault();

    // Validate form
    if (!customerSelect.value) {
        showMessage('Please select a customer', 'error');
        return;
    }

    if (!rentalSelect.value) {
        showMessage('Please select a rental', 'error');
        return;
    }

    if (!quantityReturnInput.value || quantityReturnInput.value < 1) {
        showMessage('Please enter a valid quantity', 'error');
        return;
    }

    const remaining = currentRental.clothes_qty_rent - currentRental.clothes_qty_return;
    if (parseInt(quantityReturnInput.value) > remaining) {
        showMessage(`Cannot return more than ${remaining} item(s)`, 'error');
        return;
    }

    if (!returnDateInput.value) {
        showMessage('Please select a return date', 'error');
        return;
    }

    // Validate return date is not before rental start
    const returnDate = new Date(returnDateInput.value);
    const beginDate = new Date(currentRental.clothes_rent_date_begin);

    if (returnDate < beginDate) {
        showMessage('Return date cannot be before rental start date', 'error');
        return;
    }

    // Prepare data
    const formData = {
        rental_id: parseInt(rentalSelect.value),
        clothes_qty_return: parseInt(quantityReturnInput.value),
        actual_return_date: returnDate.toISOString()
    };

    // Show loading
    loadingDiv.classList.add('show');
    submitBtn.disabled = true;
    messageDiv.style.display = 'none';

    try {
        const response = await fetch('/api/rentals/return', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(formData)
        });

        const data = await response.json();

        if (response.ok) {
            showMessage('Return processed successfully!', 'success');
            form.reset();
            rentalDetails.style.display = 'none';
            returnSummary.style.display = 'none';
            rentalSelect.disabled = true;
            quantityReturnInput.disabled = true;
            returnDateInput.disabled = true;
            setCurrentDateTime();

            // Redirect after 2 seconds
            setTimeout(() => {
                window.location.href = '/';
            }, 2000);
        } else {
            showMessage(`Error: ${data.error || 'Failed to process return'}`, 'error');
        }
    } catch (error) {
        showMessage(`Network error: ${error.message}`, 'error');
    } finally {
        loadingDiv.classList.remove('show');
        submitBtn.disabled = false;
    }
});

// Helper function to show messages
function showMessage(text, type) {
    messageDiv.textContent = text;
    messageDiv.className = `message ${type}`;
    messageDiv.style.display = 'block';

    if (type === 'error') {
        // Auto-hide error messages after 5 seconds
        setTimeout(() => {
            messageDiv.style.display = 'none';
        }, 5000);
    }
}