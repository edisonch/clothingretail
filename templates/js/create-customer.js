// Get DOM elements
const form = document.getElementById('customerForm');
const customerNameInput = document.getElementById('customerName');
const customerPhoneInput = document.getElementById('customerPhone');
const customerEmailInput = document.getElementById('customerEmail');
const customerAddressInput = document.getElementById('customerAddress');
const customerCityInput = document.getElementById('customerCity');
const customerNotesInput = document.getElementById('customerNotes');
const nameCount = document.getElementById('nameCount');
const phoneCount = document.getElementById('phoneCount');
const emailCount = document.getElementById('emailCount');
const addressCount = document.getElementById('addressCount');
const cityCount = document.getElementById('cityCount');
const notesCount = document.getElementById('notesCount');
const messageDiv = document.getElementById('message');
const loadingDiv = document.getElementById('loading');
const submitBtn = document.getElementById('submitBtn');
const cancelBtn = document.getElementById('cancelBtn');

// Character counter function
function setupCharCounter(input, counter, maxLength) {
    input.addEventListener('input', function() {
        const length = this.value.length;
        counter.textContent = `${length} / ${maxLength}`;

        if (length >= maxLength) {
            counter.classList.add('error');
            counter.classList.remove('warning');
        } else if (length >= maxLength * 0.8) {
            counter.classList.add('warning');
            counter.classList.remove('error');
        } else {
            counter.classList.remove('warning', 'error');
        }
    });
}

// Setup character counters
setupCharCounter(customerNameInput, nameCount, 64);
setupCharCounter(customerPhoneInput, phoneCount, 16);
setupCharCounter(customerEmailInput, emailCount, 128);
setupCharCounter(customerAddressInput, addressCount, 256);
setupCharCounter(customerCityInput, cityCount, 64);
setupCharCounter(customerNotesInput, notesCount, 256);

// Email validation
function isValidEmail(email) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
}

// Phone validation (basic)
function isValidPhone(phone) {
    const phoneRegex = /^[\d\s\-\+\(\)]+$/;
    return phoneRegex.test(phone) && phone.replace(/\D/g, '').length >= 8;
}

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
    if (!customerNameInput.value.trim()) {
        showMessage('Customer name is required', 'error');
        return;
    }

    if (!customerPhoneInput.value.trim()) {
        showMessage('Phone number is required', 'error');
        return;
    }

    if (!isValidPhone(customerPhoneInput.value.trim())) {
        showMessage('Please enter a valid phone number', 'error');
        return;
    }

    if (!customerEmailInput.value.trim()) {
        showMessage('Email address is required', 'error');
        return;
    }

    if (!isValidEmail(customerEmailInput.value.trim())) {
        showMessage('Please enter a valid email address', 'error');
        return;
    }

    if (!customerAddressInput.value.trim()) {
        showMessage('Address is required', 'error');
        return;
    }

    if (!customerCityInput.value.trim()) {
        showMessage('City is required', 'error');
        return;
    }

    // Prepare data
    const formData = {
        cust_name: customerNameInput.value.trim(),
        cust_phone: customerPhoneInput.value.trim(),
        cust_email: customerEmailInput.value.trim(),
        cust_address: customerAddressInput.value.trim(),
        cust_city: customerCityInput.value.trim(),
        cust_notes: customerNotesInput.value.trim()
    };

    // Show loading
    loadingDiv.classList.add('show');
    submitBtn.disabled = true;
    messageDiv.style.display = 'none';

    try {
        const response = await fetch('/api/customers', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(formData)
        });

        const data = await response.json();

        if (response.ok) {
            showMessage('Customer created successfully!', 'success');
            form.reset();

            // Reset all character counters
            nameCount.textContent = '0 / 64';
            phoneCount.textContent = '0 / 16';
            emailCount.textContent = '0 / 128';
            addressCount.textContent = '0 / 256';
            cityCount.textContent = '0 / 64';
            notesCount.textContent = '0 / 256';

            // Redirect after 2 seconds
            setTimeout(() => {
                window.location.href = '/';
            }, 2000);
        } else {
            showMessage(`Error: ${data.error || 'Failed to create customer'}`, 'error');
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