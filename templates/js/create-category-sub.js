// Get DOM elements
const form = document.getElementById('categorySubForm');
const parentCategorySelect = document.getElementById('parentCategory');
const subCategoryNameInput = document.getElementById('subCategoryName');
const subCategoryLocationInput = document.getElementById('subCategoryLocation');
const nameCount = document.getElementById('nameCount');
const locationCount = document.getElementById('locationCount');
const messageDiv = document.getElementById('message');
const loadingDiv = document.getElementById('loading');
const submitBtn = document.getElementById('submitBtn');
const cancelBtn = document.getElementById('cancelBtn');

// Store base64 images
const pictureData = {
    picture1: null,
    picture2: null,
    picture3: null,
    picture4: null,
    picture5: null
};

// Load categories on page load
document.addEventListener('DOMContentLoaded', loadCategories);

// Load categories from API
async function loadCategories() {
    try {
        const response = await fetch('/api/categories');
        if (response.ok) {
            const categories = await response.json();

            categories.forEach(category => {
                const option = document.createElement('option');
                option.value = category.id;
                option.textContent = category.clothes_cat_name;
                parentCategorySelect.appendChild(option);
            });
        } else {
            showMessage('Failed to load categories', 'error');
        }
    } catch (error) {
        showMessage(`Error loading categories: ${error.message}`, 'error');
    }
}

// Character counter for subcategory name
subCategoryNameInput.addEventListener('input', function() {
    const length = this.value.length;
    nameCount.textContent = `${length} / 32`;

    if (length >= 32) {
        nameCount.classList.add('error');
    } else if (length >= 25) {
        nameCount.classList.add('warning');
        nameCount.classList.remove('error');
    } else {
        nameCount.classList.remove('warning', 'error');
    }
});

// Character counter for location
subCategoryLocationInput.addEventListener('input', function() {
    const length = this.value.length;
    locationCount.textContent = `${length} / 64`;

    if (length >= 64) {
        locationCount.classList.add('error');
    } else if (length >= 55) {
        locationCount.classList.add('warning');
        locationCount.classList.remove('error');
    } else {
        locationCount.classList.remove('warning', 'error');
    }
});

// Handle picture uploads
for (let i = 1; i <= 5; i++) {
    const input = document.getElementById(`picture${i}`);
    const preview = document.getElementById(`preview${i}`);
    const removeBtn = document.querySelector(`.remove-picture[data-index="${i}"]`);

    input.addEventListener('change', function(e) {
        const file = e.target.files[0];
        if (file) {
            // Validate file type
            if (!file.type.startsWith('image/')) {
                showMessage('Please select a valid image file', 'error');
                return;
            }

            // Validate file size (max 5MB)
            if (file.size > 5 * 1024 * 1024) {
                showMessage('Image size must be less than 5MB', 'error');
                return;
            }

            const reader = new FileReader();
            reader.onload = function(event) {
                const base64 = event.target.result;
                pictureData[`picture${i}`] = base64;
                preview.style.backgroundImage = `url(${base64})`;
                preview.classList.add('active');
            };
            reader.readAsDataURL(file);
        }
    });

    removeBtn.addEventListener('click', function() {
        input.value = '';
        pictureData[`picture${i}`] = null;
        preview.style.backgroundImage = '';
        preview.classList.remove('active');
    });
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
    if (!parentCategorySelect.value) {
        showMessage('Please select a parent category', 'error');
        return;
    }

    if (!subCategoryNameInput.value.trim()) {
        showMessage('Subcategory name is required', 'error');
        return;
    }

    if (!subCategoryLocationInput.value.trim()) {
        showMessage('Location is required', 'error');
        return;
    }

    // Prepare data
    const formData = {
        id_clothing_category: parseInt(parentCategorySelect.value),
        clothes_cat_name_sub: subCategoryNameInput.value.trim(),
        clothes_cat_location_sub: subCategoryLocationInput.value.trim(),
        clothes_picture_1: pictureData.picture1 || '',
        clothes_picture_2: pictureData.picture2 || '',
        clothes_picture_3: pictureData.picture3 || '',
        clothes_picture_4: pictureData.picture4 || '',
        clothes_picture_5: pictureData.picture5 || ''
    };

    // Show loading
    loadingDiv.classList.add('show');
    submitBtn.disabled = true;
    messageDiv.style.display = 'none';

    try {
        const response = await fetch('/api/categories-sub', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(formData)
        });

        const data = await response.json();

        if (response.ok) {
            showMessage('Subcategory created successfully!', 'success');
            form.reset();
            nameCount.textContent = '0 / 32';
            locationCount.textContent = '0 / 64';

            // Clear all picture previews
            for (let i = 1; i <= 5; i++) {
                const preview = document.getElementById(`preview${i}`);
                preview.style.backgroundImage = '';
                preview.classList.remove('active');
                pictureData[`picture${i}`] = null;
            }

            // Redirect after 2 seconds
            setTimeout(() => {
                window.location.href = '/';
            }, 2000);
        } else {
            showMessage(`Error: ${data.error || 'Failed to create subcategory'}`, 'error');
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