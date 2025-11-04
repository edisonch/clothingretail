// Get DOM elements
const form = document.getElementById('rentalForm');
const customerSelect = document.getElementById('customer');
const categorySelect = document.getElementById('category');
const subcategorySelect = document.getElementById('subcategory');
const sizeSelect = document.getElementById('size');
const quantityInput = document.getElementById('quantity');
const rentDateBeginInput = document.getElementById('rentDateBegin');
const rentDateEndInput = document.getElementById('rentDateEnd');
const messageDiv = document.getElementById('message');
const loadingDiv = document.getElementById('loading');
const submitBtn = document.getElementById('submitBtn');
const cancelBtn = document.getElementById('cancelBtn');
const infoBox = document.getElementById('infoBox');

// Info display elements
const infoCustomer = document.getElementById('infoCustomer');
const infoItem = document.getElementById('infoItem');
const infoSize = document.getElementById('infoSize');
const infoQuantity = document.getElementById('infoQuantity');
const infoPeriod = document.getElementById('infoPeriod');
const infoDuration = document.getElementById('infoDuration');

// Store data for display
let customers = [];
let categories = [];
let subcategories = [];
let sizes = [];

// Load data on page load
document.addEventListener('DOMContentLoaded', function() {
    loadCustomers();
    loadCategories();
    setMinDateTime();
});

// Set minimum datetime to current time
function setMinDateTime() {
    const now = new Date();
    const offset = now.getTimezoneOffset() * 60000;
    const localISOTime = new Date(now - offset).toISOString().slice(0, 16);
    rentDateBeginInput.min = localISOTime;
    rentDateEndInput.min = localISOTime;
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

// Load categories from API
async function loadCategories() {
    try {
        const response = await fetch('/api/categories');
        if (response.ok) {
            categories = await response.json();
            categories.forEach(category => {
                const option = document.createElement('option');
                option.value = category.id;
                option.textContent = category.clothes_cat_name;
                categorySelect.appendChild(option);
            });
        } else {
            showMessage('Failed to load categories', 'error');
        }
    } catch (error) {
        showMessage(`Error loading categories: ${error.message}`, 'error');
    }
}

// Load subcategories when category is selected
categorySelect.addEventListener('change', async function() {
    subcategorySelect.innerHTML = '<option value="">Select a subcategory...</option>';
    subcategorySelect.disabled = true;
    sizeSelect.innerHTML = '<option value="">Select a subcategory first...</option>';
    sizeSelect.disabled = true;

    if (!this.value) return;

    try {
        const response = await fetch('/api/categories-sub');
        if (response.ok) {
            const allSubcategories = await response.json();
            subcategories = allSubcategories.filter(sub => sub.id_clothing_category === parseInt(this.value));

            if (subcategories.length === 0) {
                subcategorySelect.innerHTML = '<option value="">No subcategories available</option>';
                return;
            }

            subcategories.forEach(subcategory => {
                const option = document.createElement('option');
                option.value = subcategory.id;
                option.textContent = `${subcategory.clothes_cat_name_sub} - ${subcategory.clothes_cat_location_sub}`;
                subcategorySelect.appendChild(option);
            });

            subcategorySelect.disabled = false;
        } else {
            showMessage('Failed to load subcategories', 'error');
        }
    } catch (error) {
        showMessage(`Error loading subcategories: ${error.message}`, 'error');
    }
});

// Load sizes when subcategory is selected
subcategorySelect.addEventListener('change', async function() {
    sizeSelect.innerHTML = '<option value="">Select a size...</option>';
    sizeSelect.disabled = true;

    if (!this.value) return;

    // Note: You may need to implement a sizes API endpoint
    // For now, we'll create a placeholder
    // In a real implementation, you would fetch sizes from /api/sizes?subcategory_id=X

    // Placeholder implementation - assuming sizes are stored somewhere
    // You may need to modify this based on your actual API structure
    const placeholderSizes = [
        { id: 1, name: 'XS' },
        { id: 2, name: 'S' },
        { id: 3, name: 'M' },
        { id: 4, name: 'L' },
        { id: 5, name: 'XL' },
        { id: 6, name: 'XXL' }
    ];

    placeholderSizes.forEach(size => {
        const option = document.createElement('option');
        option.value = size.id;
        option.textContent = size.name;
        sizeSelect.appendChild(option);
    });

    sizeSelect.disabled = false;
});

// Update info box
function updateInfoBox() {
    if (!customerSelect.value || !subcategorySelect.value || !sizeSelect.value) {
        infoBox.style.display = 'none';
        return;
    }

    const customer = customers.find(c => c.id === parseInt(customerSelect.value));
    const subcategory = subcategories.find(s => s.id === parseInt(subcategorySelect.value));
    const size = sizeSelect.options[sizeSelect.selectedIndex].text;

    infoCustomer.textContent = customer ? customer.cust_name : '-';
    infoItem.textContent = subcategory ? subcategory.clothes_cat_name_sub : '-';
    infoSize.textContent = size;
    infoQuantity.textContent = quantityInput.value;

    if (rentDateBeginInput.value && rentDateEndInput.value) {
        const begin = new Date(rentDateBeginInput.value);
        const end = new Date(rentDateEndInput.value);

        infoPeriod.textContent = `${formatDate(begin)} to ${formatDate(end)}`;

        const duration = Math.ceil((end - begin) / (1000 * 60 * 60 * 24));
        infoDuration.textContent = duration > 0 ? `${duration} day(s)` : 'Invalid';
    } else {
        infoPeriod.textContent = '-';
        infoDuration.textContent = '-';
    }

    infoBox.style.display = 'block';
}

// Format date for display
function formatDate(date) {
    const day = String(date.getDate()).padStart(2, '0');
    const month = date.toLocaleDateString('en-US', { month: 'short' });
    const year = date.getFullYear();
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');

    return `${day}-${month}-${year} ${hours}:${minutes}`;
}

// Listen for changes to update info box
customerSelect.addEventListener('change', updateInfoBox);
subcategorySelect.addEventListener('change', updateInfoBox);
sizeSelect.addEventListener('change', updateInfoBox);
quantityInput.addEventListener('input', updateInfoBox);
rentDateBeginInput.addEventListener('change', function() {
    // Set end date minimum to start date
    rentDateEndInput.min = this.value;
    updateInfoBox();
});
rentDateEndInput.addEventListener('change', updateInfoBox);

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

    if (!subcategorySelect.value) {
        showMessage('Please select a subcategory', 'error');
        return;
    }

    if (!sizeSelect.value) {
        showMessage('Please select a size', 'error');
        return;
    }

    if (!quantityInput.value || quantityInput.value < 1) {
        showMessage('Please enter a valid quantity', 'error');
        return;
    }

    if (!rentDateBeginInput.value || !rentDateEndInput.value) {
        showMessage('Please select rental dates', 'error');
        return;
    }

    // Validate dates
    const beginDate = new Date(rentDateBeginInput.value);
    const endDate = new Date(rentDateEndInput.value);

    if (endDate <= beginDate) {
        showMessage('End date must be after start date', 'error');
        return;
    }

    // Prepare data - format dates as ISO strings
    const formData = {
        id_clothing_category_sub: parseInt(subcategorySelect.value),
        id_clothing_size: parseInt(sizeSelect.value),
        id_clothing_customer: parseInt(customerSelect.value),
        clothes_qty_rent: parseInt(quantityInput.value),
        rent_date_begin: beginDate.toISOString(),
        rent_date_end: endDate.toISOString()
    };

    // Show loading
    loadingDiv.classList.add('show');
    submitBtn.disabled = true;
    messageDiv.style.display = 'none';

    try {
        const response = await fetch('/api/rentals', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(formData)
        });

        const data = await response.json();

        if (response.ok) {
            showMessage('Rental created successfully!', 'success');
            form.reset();
            infoBox.style.display = 'none';
            subcategorySelect.disabled = true;
            sizeSelect.disabled = true;
            setMinDateTime();

            // Redirect after 2 seconds
            setTimeout(() => {
                window.location.href = '/';
            }, 2000);
        } else {
            showMessage(`Error: ${data.error || 'Failed to create rental'}`, 'error');
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