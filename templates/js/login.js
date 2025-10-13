document.addEventListener('DOMContentLoaded', function() {
    const loginForm = document.getElementById('loginForm');
    const usernameInput = document.getElementById('username');
    const pinInput = document.getElementById('pin');
    const usernameError = document.getElementById('username-error');
    const pinError = document.getElementById('pin-error');
    const errorMessage = document.getElementById('error-message');
    const loginBtn = loginForm.querySelector('button[type="submit"]');

    // Real-time validation for username
    usernameInput.addEventListener('input', function() {
        validateUsername();
    });

    // Real-time validation for PIN - only allow digits
    pinInput.addEventListener('input', function() {
        // Remove any non-digit characters
        this.value = this.value.replace(/\D/g, '');
        validatePin();
    });

    // Prevent paste of non-digit content in PIN field
    pinInput.addEventListener('paste', function(e) {
        e.preventDefault();
        const pastedText = (e.clipboardData || window.clipboardData).getData('text');
        const digitsOnly = pastedText.replace(/\D/g, '').substring(0, 6);
        this.value = digitsOnly;
        validatePin();
    });

    function validateUsername() {
        const username = usernameInput.value.trim();

        if (username === '') {
            usernameError.textContent = 'Username is required';
            usernameInput.classList.add('error');
            return false;
        }

        if (username.length > 32) {
            usernameError.textContent = 'Username must not exceed 32 characters';
            usernameInput.classList.add('error');
            return false;
        }

        usernameError.textContent = '';
        usernameInput.classList.remove('error');
        return true;
    }

    function validatePin() {
        const pin = pinInput.value;

        if (pin === '') {
            pinError.textContent = 'PIN is required';
            pinInput.classList.add('error');
            return false;
        }

        if (!/^\d+$/.test(pin)) {
            pinError.textContent = 'PIN must contain only digits';
            pinInput.classList.add('error');
            return false;
        }

        if (pin.length !== 6) {
            pinError.textContent = 'PIN must be exactly 6 digits';
            pinInput.classList.add('error');
            return false;
        }

        pinError.textContent = '';
        pinInput.classList.remove('error');
        return true;
    }

    function hideErrorMessage() {
        errorMessage.style.display = 'none';
        errorMessage.textContent = '';
    }

    function showErrorMessage(message) {
        errorMessage.textContent = message;
        errorMessage.style.display = 'block';
    }

    // Form submission
    loginForm.addEventListener('submit', async function(e) {
        e.preventDefault();

        hideErrorMessage();

        // Validate all fields
        const isUsernameValid = validateUsername();
        const isPinValid = validatePin();

        if (!isUsernameValid || !isPinValid) {
            return;
        }

        // Disable button during submission
        loginBtn.disabled = true;
        loginBtn.textContent = 'Logging in...';

        const formData = {
            username: usernameInput.value.trim(),
            pin: pinInput.value
        };

        try {
            const response = await fetch('/api/auth/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(formData)
            });

            const data = await response.json();

            if (response.ok) {
                // Success - redirect to index page
                window.location.href = '/';
            } else {
                // Show error message
                showErrorMessage(data.error || 'Failed to authenticate');
                loginBtn.disabled = false;
                loginBtn.textContent = 'Login';
            }
        } catch (error) {
            console.error('Login error:', error);
            showErrorMessage('An error occurred. Please try again.');
            loginBtn.disabled = false;
            loginBtn.textContent = 'Login';
        }
    });
});

// Hamburger menu toggle functionality
document.addEventListener('DOMContentLoaded', function() {
    const hamburgerBtn = document.getElementById('hamburgerBtn');
    const sidebar = document.getElementById('sidebar');
    const sidebarOverlay = document.getElementById('sidebarOverlay');
    const logoutBtn = document.getElementById('logoutBtn');

    // Ensure menu is closed by default
    sidebar.classList.remove('open');
    hamburgerBtn.classList.remove('active');
    sidebarOverlay.classList.remove('active');

    // Toggle menu when hamburger button is clicked
    function toggleMenu() {
        const isOpen = sidebar.classList.contains('open');

        if (isOpen) {
            closeMenu();
        } else {
            openMenu();
        }
    }

    // Open the menu
    function openMenu() {
        sidebar.classList.add('open');
        hamburgerBtn.classList.add('active');
        sidebarOverlay.classList.add('active');

        // Prevent body scroll when menu is open
        document.body.style.overflow = 'hidden';
    }

    // Close the menu
    function closeMenu() {
        sidebar.classList.remove('open');
        hamburgerBtn.classList.remove('active');
        sidebarOverlay.classList.remove('active');

        // Restore body scroll
        document.body.style.overflow = '';
    }

    // Event listeners
    hamburgerBtn.addEventListener('click', toggleMenu);
    sidebarOverlay.addEventListener('click', closeMenu);

    // Close menu when clicking on a nav link (mobile UX)
    const navLinks = document.querySelectorAll('.nav-link');
    navLinks.forEach(link => {
        link.addEventListener('click', function(e) {
            // Don't close for logout button - it has its own handler
            if (link.id !== 'logoutBtn' && !e.defaultPrevented) {
                closeMenu();
            }
        });
    });

    // Logout functionality
    logoutBtn.addEventListener('click', function(e) {
        e.preventDefault();

        // Call logout API
        fetch('/api/auth/logout', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'same-origin'
        })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    // Redirect to login page
                    window.location.href = '/login';
                } else {
                    console.error('Logout failed:', data.message);
                    // Redirect anyway to be safe
                    window.location.href = '/login';
                }
            })
            .catch(error => {
                console.error('Logout error:', error);
                // Redirect to login page even if there's an error
                window.location.href = '/login';
            });
    });

    // Close menu on escape key
    document.addEventListener('keydown', function(e) {
        if (e.key === 'Escape' && sidebar.classList.contains('open')) {
            closeMenu();
        }
    });

    // Handle window resize - close menu on larger screens if needed
    let resizeTimer;
    window.addEventListener('resize', function() {
        clearTimeout(resizeTimer);
        resizeTimer = setTimeout(function() {
            if (window.innerWidth >= 1024 && sidebar.classList.contains('open')) {
                closeMenu();
            }
        }, 250);
    });
});