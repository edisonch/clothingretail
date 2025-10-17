// Detect mode from script tag data attribute
const scriptTag = document.currentScript;
const mode = scriptTag && scriptTag.dataset.mode === 'edit' ? 'edit' : 'create';
const isEditMode = mode === 'edit';

// Get DOM elements
const form = document.getElementById('categoryForm');
const categoryIdInput = document.getElementById('categoryId');
const categoryNameInput = document.getElementById('categoryName');
const categoryNotesInput = document.getElementById('categoryNotes');
const nameCount = document.getElementById('nameCount');
const notesCount = document.getElementById('notesCount');
const messageDiv = document.getElementById('message');
const loadingDiv = document.getElementById('loading');
const submitBtn = document.getElementById('submitBtn');
const cancelBtn = document.getElementById('cancelBtn');
const loadingContainer = document.getElementById('loadingContainer');
const createdAtSpan = document.getElementById('createdAt');
const updatedAtSpan = document.getElementById('updatedAt');

// Load category data if in edit mode
if (isEditMode) {
    document.addEventListener('DOMContentLoaded', loadCategoryData);
}

// Load category data for editing
async function loadCategoryData() {
    // Get category ID from URL query parameter
    const urlParams = new URLSearchParams(window.location.search);
    const categoryId = urlParams.get('id');

    if (!categoryId) {
        showMessage('No category ID provided', 'error');
        loadingContainer.style.display = 'none';
        return;
    }

    try {
        const response = await fetch(`/api/categories/${categoryId}`);

        if (response.ok) {
            const category = await response.json();

            // Populate form fields
            categoryIdInput.value = category.id;
            categoryNameInput.value = category.clothes_cat_name;
            categoryNotesInput.value = category.clothes_notes || '';

            // Update character counts
            updateCharCount(categoryNameInput, nameCount, 32);
            updateCharCount(categoryNotesInput, notesCount, 256);

            // Update metadata
            if (createdAtSpan) {
                createdAtSpan.textContent = formatDateTime(category.created_at);
            }
            if (updatedAtSpan) {
                updatedAtSpan.textContent = formatDateTime(category.updated_at);
            }

            // Show form, hide loading
            loadingContainer.style.display = 'none';
            form.style.display = 'block';
        } else {
            const data = await response.json();
            showMessage(`Error: ${data.error || 'Category not found'}`, 'error');
            loadingContainer.style.display = 'none';
        }
    } catch (error) {
        showMessage(`Network error: ${error.message}`, 'error');
        loadingContainer.style.display = 'none';
    }
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
    } else if (length >= maxLength * 0.8) {
        counter.classList.add('warning');
        counter.classList.remove('error');
    } else {
        counter.classList.remove('warning', 'error');
    }
}

// Character counter for category name
categoryNameInput.addEventListener('input', function() {
    updateCharCount(this, nameCount, 32);
});

// Character counter for notes
categoryNotesInput.addEventListener('input', function() {
    updateCharCount(this, notesCount, 256);
});

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
    if (!categoryNameInput.value.trim()) {
        showMessage('Category name is required', 'error');
        return;
    }

    // Prepare data
    const formData = {
        clothes_cat_name: categoryNameInput.value.trim(),
        clothes_notes: categoryNotesInput.value.trim()
    };

    // Show loading
    loadingDiv.classList.add('show');
    submitBtn.disabled = true;
    messageDiv.style.display = 'none';

    try {
        let response;

        if (isEditMode) {
            // Update existing category
            const categoryId = categoryIdInput.value;
            response = await fetch(`/api/categories/${categoryId}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(formData)
            });
        } else {
            // Create new category
            response = await fetch('/api/categories', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(formData)
            });
        }

        const data = await response.json();

        if (response.ok) {
            const successMessage = isEditMode ? 'Category updated successfully!' : 'Category created successfully!';
            showMessage(successMessage, 'success');

            if (!isEditMode) {
                form.reset();
                nameCount.textContent = '0 / 32';
                notesCount.textContent = '0 / 256';
            }

            // Redirect after 2 seconds
            setTimeout(() => {
                window.location.href = '/';
            }, 2000);
        } else {
            const errorMessage = isEditMode ? 'Failed to update category' : 'Failed to create category';
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