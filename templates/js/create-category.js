// Get DOM elements
const form = document.getElementById('categoryForm');
const categoryNameInput = document.getElementById('categoryName');
const categoryNotesInput = document.getElementById('categoryNotes');
const nameCount = document.getElementById('nameCount');
const notesCount = document.getElementById('notesCount');
const messageDiv = document.getElementById('message');
const loadingDiv = document.getElementById('loading');
const submitBtn = document.getElementById('submitBtn');
const cancelBtn = document.getElementById('cancelBtn');

// Character counter for category name
categoryNameInput.addEventListener('input', function() {
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

// Character counter for notes
categoryNotesInput.addEventListener('input', function() {
    const length = this.value.length;
    notesCount.textContent = `${length} / 256`;

    if (length >= 256) {
        notesCount.classList.add('error');
    } else if (length >= 230) {
        notesCount.classList.add('warning');
        notesCount.classList.remove('error');
    } else {
        notesCount.classList.remove('warning', 'error');
    }
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
        const response = await fetch('/api/categories', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(formData)
        });

        const data = await response.json();

        if (response.ok) {
            showMessage('Category created successfully!', 'success');
            form.reset();
            nameCount.textContent = '0 / 32';
            notesCount.textContent = '0 / 256';

            // Redirect after 2 seconds
            setTimeout(() => {
                window.location.href = '/';
            }, 2000);
        } else {
            showMessage(`Error: ${data.error || 'Failed to create category'}`, 'error');
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