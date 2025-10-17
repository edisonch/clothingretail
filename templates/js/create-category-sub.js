// Detect mode from script tag data attribute
const scriptTag = document.currentScript;
const mode = scriptTag && scriptTag.dataset.mode === 'edit' ? 'edit' : 'create';
const isEditMode = mode === 'edit';

// Get DOM elements
const form = document.getElementById('categorySubForm');
const subcategoryIdInput = document.getElementById('subcategoryId');
const parentCategorySelect = document.getElementById('parentCategory');
const subCategoryNameInput = document.getElementById('subCategoryName');
const subCategoryLocationInput = document.getElementById('subCategoryLocation');
const nameCount = document.getElementById('nameCount');
const locationCount = document.getElementById('locationCount');
const messageDiv = document.getElementById('message');
const loadingDiv = document.getElementById('loading');
const submitBtn = document.getElementById('submitBtn');
const cancelBtn = document.getElementById('cancelBtn');
const loadingContainer = document.getElementById('loadingContainer');
const createdAtSpan = document.getElementById('createdAt');
const updatedAtSpan = document.getElementById('updatedAt');

// Store base64 images
const pictureData = {
    picture1: null,
    picture2: null,
    picture3: null,
    picture4: null,
    picture5: null
};

// Store data for display
let categories = [];

// Load data on page load
document.addEventListener('DOMContentLoaded', function() {
    loadCategories();
    if (isEditMode) {
        loadSubcategoryData();
    }
});

// Load subcategory data for editing
async function loadSubcategoryData() {
    // Get subcategory ID from URL query parameter
    const urlParams = new URLSearchParams(window.location.search);
    const subcategoryId = urlParams.get('id');

    if (!subcategoryId) {
        showMessage('No subcategory ID provided', 'error');
        if (loadingContainer) loadingContainer.style.display = 'none';
        return;
    }

    try {
        const response = await fetch(`/api/categories-sub/${subcategoryId}`);

        if (response.ok) {
            const subcategory = await response.json();

            // Wait for categories to load first
            await waitForCategories();

            // Populate form fields
            subcategoryIdInput.value = subcategory.id;
            parentCategorySelect.value = subcategory.id_clothing_category;
            subCategoryNameInput.value = subcategory.clothes_cat_name_sub;
            subCategoryLocationInput.value = subcategory.clothes_cat_location_sub;

            // Load existing pictures
            for (let i = 1; i <= 5; i++) {
                const pictureKey = `picture${i}`;
                const base64Data = subcategory[`clothes_picture_${i}`];
                if (base64Data) {
                    pictureData[pictureKey] = base64Data;
                    const preview = document.getElementById(`preview${i}`);
                    preview.style.backgroundImage = `url(${base64Data})`;
                    preview.classList.add('active');
                }
            }

            // Update character counts
            updateCharCount(subCategoryNameInput, nameCount, 32);
            updateCharCount(subCategoryLocationInput, locationCount, 64);

            // Update metadata
            if (createdAtSpan) {
                createdAtSpan.textContent = formatDateTime(subcategory.created_at);
            }
            if (updatedAtSpan) {
                updatedAtSpan.textContent = formatDateTime(subcategory.updated_at);
            }

            // Show form, hide loading
            if (loadingContainer) loadingContainer.style.display = 'none';
            form.style.display = 'block';
        } else {
            const data = await response.json();
            showMessage(`Error: ${data.error || 'Subcategory not found'}`, 'error');
            if (loadingContainer) loadingContainer.style.display = 'none';
        }
    } catch (error) {
        showMessage(`Network error: ${error.message}`, 'error');
        if (loadingContainer) loadingContainer.style.display = 'none';
    }
}

// Wait for categories to load
function waitForCategories() {
    return new Promise((resolve) => {
        const checkInterval = setInterval(() => {
            if (categories.length > 0) {
                clearInterval(checkInterval);
                resolve();
            }
        }, 100);

        // Timeout after 5 seconds
        setTimeout(() => {
            clearInterval(checkInterval);
            resolve();
        }, 5000);
    });
}

// Format datetime for display
function formatDateTime(dateString) {
    const date = new Date(dateString);
    return date.toLocaleString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });
}

// Update character count
function updateCharCount(input, counter, maxLength) {
    const length = input.value.length;
    counter.textContent = `${length} / ${maxLength}`;

    if (length >= maxLength) {
        counter.classList.add('error');
        counter.classList.remove('warning');
    } else if (length >= maxLength * 0.85) {
        counter.classList.add('warning');
        counter.classList.remove('error');
    } else {
        counter.classList.remove('warning', 'error');
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
                parentCategorySelect.appendChild(option);
            });

            // Show form if in create mode
            if (!isEditMode) {
                if (loadingContainer) loadingContainer.style.display = 'none';
                form.style.display = 'block';
            }
        } else {
            showMessage('Failed to load categories', 'error');
        }
    } catch (error) {
        showMessage(`Error loading categories: ${error.message}`, 'error');
    }
}

// Character counter for subcategory name
subCategoryNameInput.addEventListener('input', function() {
    updateCharCount(this, nameCount, 32);
});

// Character counter for location
subCategoryLocationInput.addEventListener('input', function() {
    updateCharCount(this, locationCount, 64);
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
        let response;

        if (isEditMode) {
            // Update existing subcategory
            const subcategoryId = subcategoryIdInput.value;
            response = await fetch(`/api/categories-sub/${subcategoryId}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(formData)
            });
        } else {
            // Create new subcategory
            response = await fetch('/api/categories-sub', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(formData)
            });
        }

        const data = await response.json();

        if (response.ok) {
            const successMessage = isEditMode ? 'Subcategory updated successfully!' : 'Subcategory created successfully!';
            showMessage(successMessage, 'success');

            if (!isEditMode) {
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
            }

            // Redirect after 2 seconds
            setTimeout(() => {
                window.location.href = '/';
            }, 2000);
        } else {
            const errorMessage = isEditMode ? 'Failed to update subcategory' : 'Failed to create subcategory';
            showMessage(`Error: ${data.error || errorMessage}`, 'error');
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