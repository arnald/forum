/**
 * Authentication Manager
 *
 * This class handles all user authentication functionality including:
 * - Login/logout operations
 * - Session management
 * - User registration
 * - UI updates based on authentication state
 * - Form validation and API communication
 */
class AuthManager {
    constructor() {
        // Store current user information (null if not logged in)
        this.user = null;
        // Initialize the authentication system
        this.init();
    }

    /**
     * Initialize the authentication manager
     * Called when the page loads to set up authentication state and event listeners
     */
    init() {
        // Check if user is already logged in (from existing session/cookie)
        this.checkAuthStatus();
        // Set up form event listeners for login/register/logout
        this.setupEventListeners();
    }

    /**
     * Check if user is currently authenticated
     * Attempts to validate existing session by calling protected endpoint
     */
    async checkAuthStatus() {
        try {
            // Try to get current user info from a protected endpoint
            // This validates the session cookie automatically
            const response = await fetch('/api/v1/user/profile', {
                credentials: 'include' // Include cookies for session validation
            });

            // If response is successful, user is authenticated
            if (response.ok) {
                const data = await response.json();
                if (data.success) {
                    // Store user information locally
                    this.setUser(data.user);
                }
            }
        } catch (error) {
            // Failed to authenticate - user is not logged in
            console.log('User not authenticated');
        }

        // Update UI elements based on authentication status
        this.updateUI();
    }

    /**
     * Store user information locally
     * @param {Object} user - User object containing id, username, email, etc.
     */
    setUser(user) {
        this.user = user;
        // Store in localStorage for persistence across page reloads
        localStorage.setItem('user', JSON.stringify(user));
    }

    /**
     * Clear stored user information (logout)
     * Removes user data from memory and localStorage
     */
    clearUser() {
        this.user = null;
        localStorage.removeItem('user');
    }

    /**
     * Check if user is currently logged in
     * @returns {boolean} True if user is authenticated, false otherwise
     */
    isLoggedIn() {
        return this.user !== null;
    }

    /**
     * Update UI elements based on current authentication status
     * Shows/hides user info, navigation elements, and authentication-specific features
     */
    updateUI() {
        const userInfo = document.getElementById('user-info');
        const username = document.getElementById('username');

        if (this.isLoggedIn() && userInfo) {
            // Show user info section in navigation
            userInfo.style.display = 'inline';
            if (username) {
                // Display username or email as fallback
                username.textContent = this.user.username || this.user.email;
            }

            // Show authenticated user filter buttons (on posts page)
            const myPostsBtn = document.querySelector('[data-filter="my-posts"]');
            const likedBtn = document.querySelector('[data-filter="liked"]');
            if (myPostsBtn) myPostsBtn.style.display = 'inline-block';
            if (likedBtn) likedBtn.style.display = 'inline-block';
        } else if (userInfo) {
            // Hide user info section
            userInfo.style.display = 'none';

            // Hide authenticated user filter buttons
            const myPostsBtn = document.querySelector('[data-filter="my-posts"]');
            const likedBtn = document.querySelector('[data-filter="liked"]');
            if (myPostsBtn) myPostsBtn.style.display = 'none';
            if (likedBtn) likedBtn.style.display = 'none';
        }
    }

    /**
     * Set up event listeners for authentication forms and buttons
     * Attaches handlers to login form, register form, and logout button
     */
    setupEventListeners() {
        // Login form submission handler
        const loginForm = document.getElementById('loginForm');
        if (loginForm) {
            loginForm.addEventListener('submit', (e) => this.handleLogin(e));
        }

        // Register form submission handler
        const registerForm = document.getElementById('registerForm');
        if (registerForm) {
            registerForm.addEventListener('submit', (e) => this.handleRegister(e));
        }

        // Logout button click handler (using event delegation for dynamic elements)
        document.addEventListener('click', (e) => {
            if (e.target.id === 'logout-btn') {
                this.handleLogout();
            }
        });
    }

    /**
     * Handle login form submission
     * Supports both email and username login methods
     * @param {Event} e - Form submission event
     */
    async handleLogin(e) {
        e.preventDefault(); // Prevent default form submission

        // Extract form data
        const identifier = document.getElementById('identifier').value;
        const password = document.getElementById('password').value;
        const loginType = document.querySelector('input[name="loginType"]:checked').value;

        // Choose API endpoint based on login type (email or username)
        const endpoint = loginType === 'email' ? '/api/v1/login/email' : '/api/v1/login/username';

        try {
            // Send login request to API
            const response = await fetch(endpoint, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include', // Include cookies for session management
                body: JSON.stringify({
                    [loginType]: identifier, // Dynamic property name based on login type
                    password: password
                })
            });

            const data = await response.json();

            if (response.ok && data.success) {
                // Login successful
                this.setUser(data.user);
                this.showMessage('Login successful!', 'success');
                // Redirect to posts page after brief delay
                setTimeout(() => {
                    window.location.href = '/posts';
                }, 1000);
            } else {
                // Login failed - show error message
                this.showMessage(data.message || 'Login failed', 'error');
            }
        } catch (error) {
            // Network or other error
            this.showMessage('Network error. Please try again.', 'error');
        }
    }

    /**
     * Handle user registration form submission
     * Validates password confirmation and creates new user account
     * @param {Event} e - Form submission event
     */
    async handleRegister(e) {
        e.preventDefault(); // Prevent default form submission

        // Extract form data
        const email = document.getElementById('email').value;
        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;
        const confirmPassword = document.getElementById('confirmPassword').value;

        // Client-side password validation
        if (password !== confirmPassword) {
            this.showMessage('Passwords do not match', 'error');
            return;
        }

        try {
            // Send registration request to API
            const response = await fetch('/api/v1/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include', // Include cookies for session management
                body: JSON.stringify({
                    email: email,
                    username: username,
                    password: password // Will be encrypted on server
                })
            });

            const data = await response.json();

            if (response.ok && data.success) {
                // Registration successful
                this.showMessage('Registration successful! Please login.', 'success');
                // Redirect to login page after delay
                setTimeout(() => {
                    window.location.href = '/login';
                }, 2000);
            } else {
                // Registration failed - show error (email/username taken, etc.)
                this.showMessage(data.message || 'Registration failed', 'error');
            }
        } catch (error) {
            // Network or other error
            this.showMessage('Network error. Please try again.', 'error');
        }
    }

    /**
     * Handle user logout
     * Clears session on server and local authentication state
     */
    async handleLogout() {
        try {
            // Send logout request to API to clear server session
            const response = await fetch('/api/v1/logout', {
                method: 'POST',
                credentials: 'include' // Include cookies for session management
            });

            // Clear local authentication state
            this.clearUser();
            this.updateUI();
            // Redirect to home page
            window.location.href = '/';
        } catch (error) {
            // Even if the request fails, clear local auth state
            // This ensures user is logged out locally even if server request fails
            this.clearUser();
            this.updateUI();
            window.location.href = '/';
        }
    }

    /**
     * Display success or error messages to the user
     * Automatically hides messages after 5 seconds
     * @param {string} message - The message to display
     * @param {string} type - Message type: 'error' or 'success'
     */
    showMessage(message, type) {
        const errorDiv = document.getElementById('error-message');
        const successDiv = document.getElementById('success-message');

        if (type === 'error' && errorDiv) {
            // Show error message and hide success message
            errorDiv.textContent = message;
            errorDiv.style.display = 'block';
            if (successDiv) successDiv.style.display = 'none';
        } else if (type === 'success' && successDiv) {
            // Show success message and hide error message
            successDiv.textContent = message;
            successDiv.style.display = 'block';
            if (errorDiv) errorDiv.style.display = 'none';
        }

        // Auto-hide messages after 5 seconds for better UX
        setTimeout(() => {
            if (errorDiv) errorDiv.style.display = 'none';
            if (successDiv) successDiv.style.display = 'none';
        }, 5000);
    }

    /**
     * Helper method to get current user ID for API calls
     * @returns {string|null} User ID if logged in, null otherwise
     */
    getUserId() {
        return this.user ? this.user.id : null;
    }
}

/**
 * Initialize the authentication manager when the DOM is fully loaded
 * Creates a global authManager instance that can be accessed by other scripts
 * This ensures authentication state is available across all pages
 */
document.addEventListener('DOMContentLoaded', () => {
    // Create global auth manager instance
    window.authManager = new AuthManager();
});