/**
 * Posts Manager
 *
 * This class handles all post-related functionality including:
 * - Loading and displaying posts from the API
 * - Post filtering (all, by category, user's posts, liked posts)
 * - Post creation and form handling
 * - Voting on posts and comments
 * - Category management and selection
 * - UI updates and user interactions
 */
class PostsManager {
    constructor() {
        // Track current active filter (all, category, my-posts, liked)
        this.currentFilter = 'all';
        // Store loaded categories for filtering and post creation
        this.categories = [];
        // Store currently displayed posts
        this.posts = [];
        // Initialize the posts management system
        this.init();
    }

    /**
     * Initialize the posts manager
     * Sets up event listeners and loads initial data
     */
    init() {
        // Set up all event listeners for user interactions
        this.setupEventListeners();
        // Load available categories from API
        this.loadCategories();
        // Load and display initial posts
        this.loadPosts();
    }

    /**
     * Set up event listeners for post-related interactions
     * Handles filter buttons, category selection, and form submissions
     */
    setupEventListeners() {
        // Filter button click handlers (using event delegation)
        document.addEventListener('click', (e) => {
            if (e.target.classList.contains('filter-btn')) {
                this.handleFilterChange(e.target);
            }
        });

        // Category selection dropdown handler
        const categorySelect = document.getElementById('category-select');
        if (categorySelect) {
            categorySelect.addEventListener('change', (e) => {
                this.filterPostsByCategory(e.target.value);
            });
        }

        // Create post form submission handler
        const createPostForm = document.getElementById('createPostForm');
        if (createPostForm) {
            createPostForm.addEventListener('submit', (e) => this.handleCreatePost(e));
        }
    }

    /**
     * Handle filter button changes
     * Updates UI state and applies the selected filter
     * @param {HTMLElement} button - The clicked filter button
     */
    handleFilterChange(button) {
        // Update visual state of filter buttons
        document.querySelectorAll('.filter-btn').forEach(btn => {
            btn.classList.remove('active');
        });
        button.classList.add('active');

        // Get filter type from button data attribute
        const filter = button.dataset.filter;
        this.currentFilter = filter;

        // Show/hide category dropdown based on filter type
        const categoryFilter = document.getElementById('category-filter');
        if (categoryFilter) {
            categoryFilter.style.display = filter === 'category' ? 'block' : 'none';
        }

        // Apply the selected filter to posts
        this.applyFilter(filter);
    }

    /**
     * Apply the selected filter to load and display appropriate posts
     * Handles different filter types: all, my-posts, liked, category
     * @param {string} filter - The filter type to apply
     */
    async applyFilter(filter) {
        try {
            let endpoint = '/api/v1/posts/all';
            let requestOptions = {
                credentials: 'include' // Include cookies for authentication
            };

            // Determine API endpoint based on filter type
            switch (filter) {
                case 'all':
                    // Load all posts (default)
                    endpoint = '/api/v1/posts/all';
                    break;
                case 'my-posts':
                    // Load only user's own posts (requires authentication)
                    if (!window.authManager.isLoggedIn()) {
                        this.showMessage('Please login to view your posts', 'error');
                        return;
                    }
                    endpoint = `/api/v1/posts/user/${window.authManager.getUserId()}`;
                    break;
                case 'liked':
                    // Load posts user has liked (requires authentication)
                    if (!window.authManager.isLoggedIn()) {
                        this.showMessage('Please login to view liked posts', 'error');
                        return;
                    }
                    endpoint = `/api/v1/votes/user?user_id=${window.authManager.getUserId()}&target_type=post&vote_type=like`;
                    break;
                case 'category':
                    // Category filtering is handled separately by filterPostsByCategory\n                    return;
            }

            // Fetch posts from API
            const response = await fetch(endpoint, requestOptions);
            const data = await response.json();

            if (response.ok && data.success) {
                if (filter === 'liked') {
                    // For liked posts, we get vote records and need to fetch actual posts
                    await this.loadLikedPosts(data.votes);
                } else {
                    // For other filters, we get posts directly
                    this.posts = data.posts || [];
                    this.displayPosts(this.posts);
                }
            } else {
                this.showMessage(data.message || 'Failed to load posts', 'error');
            }
        } catch (error) {
            this.showMessage('Error loading posts', 'error');
        }
    }

    async loadLikedPosts(votes) {
        if (!votes || votes.length === 0) {
            this.displayPosts([]);
            return;
        }

        const likedPosts = [];
        for (const vote of votes) {
            try {
                const response = await fetch(`/api/v1/posts/get?id=${vote.target_id}`, {
                    credentials: 'include'
                });
                const data = await response.json();
                if (response.ok && data.success) {
                    likedPosts.push(data.post);
                }
            } catch (error) {
                console.error('Error loading liked post:', error);
            }
        }

        this.posts = likedPosts;
        this.displayPosts(this.posts);
    }

    async filterPostsByCategory(categoryId) {
        if (!categoryId) {
            this.loadPosts();
            return;
        }

        try {
            const response = await fetch(`/api/v1/posts/category?category_id=${categoryId}`, {
                credentials: 'include'
            });
            const data = await response.json();

            if (response.ok && data.success) {
                this.posts = data.posts || [];
                this.displayPosts(this.posts);
            } else {
                this.showMessage('Failed to load posts for category', 'error');
            }
        } catch (error) {
            this.showMessage('Error loading posts', 'error');
        }
    }

    async loadCategories() {
        try {
            const response = await fetch('/api/v1/categories/all', {
                credentials: 'include'
            });
            const data = await response.json();

            if (response.ok && data.success) {
                this.categories = data.categories || [];
                this.populateCategorySelect();
            }
        } catch (error) {
            console.error('Error loading categories:', error);
        }
    }

    populateCategorySelect() {
        const categorySelect = document.getElementById('category-select');
        if (!categorySelect) return;

        // Clear existing options except the first one
        categorySelect.innerHTML = '<option value="">Select Category</option>';

        this.categories.forEach(category => {
            const option = document.createElement('option');
            option.value = category.id;
            option.textContent = category.name;
            categorySelect.appendChild(option);
        });

        // Also populate create post categories if it exists
        const categoriesList = document.getElementById('categories-list');
        if (categoriesList) {
            categoriesList.innerHTML = this.categories.map(category => `
                <label>
                    <input type="checkbox" name="categories" value="${category.id}">
                    ${this.escapeHtml(category.name)}
                </label>
            `).join('');
        }
    }

    async loadPosts() {
        try {
            const response = await fetch('/api/v1/posts/all', {
                credentials: 'include'
            });
            const data = await response.json();

            if (response.ok && data.success) {
                this.posts = data.posts || [];
                this.displayPosts(this.posts);
            } else {
                this.showMessage('Failed to load posts', 'error');
            }
        } catch (error) {
            this.showMessage('Error loading posts', 'error');
        }
    }

    displayPosts(posts) {
        const postsContainer = document.getElementById('posts-list');
        if (!postsContainer) return;

        if (!posts || posts.length === 0) {
            postsContainer.innerHTML = '<p>No posts found.</p>';
            return;
        }

        postsContainer.innerHTML = posts.map(post => this.renderPost(post)).join('');
    }

    renderPost(post) {
        const categoriesHtml = post.categories && post.categories.length > 0
            ? `<div class="post-categories">
                ${post.categories.map(cat => `<span class="category-tag">${this.escapeHtml(cat.name || cat)}</span>`).join('')}
               </div>`
            : '';

        return `
            <div class="post" data-post-id="${post.id}">
                <div class="post-header">
                    <div>
                        <h3 class="post-title">${this.escapeHtml(post.title)}</h3>
                        <div class="post-meta">
                            By ${this.escapeHtml(post.username || 'Unknown')} • ${this.formatDate(post.created_at)}
                        </div>
                    </div>
                </div>
                <div class="post-content">
                    ${this.escapeHtml(post.content).substring(0, 300)}${post.content.length > 300 ? '...' : ''}
                </div>
                ${categoriesHtml}
                <div class="post-actions">
                    <div class="vote-buttons">
                        <button class="vote-btn" onclick="postsManager.vote('${post.id}', 'post', 'like')">
                            👍 Like (${post.like_count || 0})
                        </button>
                        <button class="vote-btn" onclick="postsManager.vote('${post.id}', 'post', 'dislike')">
                            👎 Dislike (${post.dislike_count || 0})
                        </button>
                    </div>
                    <div>
                        <button class="btn" onclick="window.location.href='/post/?id=${post.id}'">
                            View Details
                        </button>
                    </div>
                </div>
            </div>
        `;
    }

    async handleCreatePost(e) {
        e.preventDefault();

        if (!window.authManager.isLoggedIn()) {
            this.showMessage('Please login to create a post', 'error');
            window.location.href = '/login';
            return;
        }

        const title = document.getElementById('title').value;
        const content = document.getElementById('content').value;
        const selectedCategories = Array.from(document.querySelectorAll('input[name="categories"]:checked'))
            .map(cb => cb.value);

        try {
            const response = await fetch('/api/v1/posts', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
                body: JSON.stringify({
                    title: title,
                    content: content,
                    category_ids: selectedCategories
                })
            });

            const data = await response.json();

            if (response.ok && data.success) {
                this.showMessage('Post created successfully!', 'success');
                setTimeout(() => {
                    window.location.href = '/posts';
                }, 1000);
            } else {
                this.showMessage(data.message || 'Failed to create post', 'error');
            }
        } catch (error) {
            this.showMessage('Error creating post', 'error');
        }
    }

    async vote(targetId, targetType, voteType) {
        if (!window.authManager.isLoggedIn()) {
            this.showMessage('Please login to vote', 'error');
            return;
        }

        try {
            const response = await fetch('/api/v1/votes/cast', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
                body: JSON.stringify({
                    target_id: targetId,
                    target_type: targetType,
                    vote_type: voteType
                })
            });

            const data = await response.json();

            if (response.ok && data.success) {
                // Reload posts to update vote counts
                this.applyFilter(this.currentFilter);
            } else {
                this.showMessage(data.message || 'Failed to vote', 'error');
            }
        } catch (error) {
            this.showMessage('Error voting', 'error');
        }
    }

    showMessage(message, type) {
        const errorDiv = document.getElementById('error-message');
        const successDiv = document.getElementById('success-message');

        if (type === 'error' && errorDiv) {
            errorDiv.textContent = message;
            errorDiv.style.display = 'block';
            if (successDiv) successDiv.style.display = 'none';
        } else if (type === 'success' && successDiv) {
            successDiv.textContent = message;
            successDiv.style.display = 'block';
            if (errorDiv) errorDiv.style.display = 'none';
        }

        setTimeout(() => {
            if (errorDiv) errorDiv.style.display = 'none';
            if (successDiv) successDiv.style.display = 'none';
        }, 5000);
    }

    escapeHtml(text) {
        const map = {
            '&': '&amp;',
            '<': '&lt;',
            '>': '&gt;',
            '"': '&quot;',
            "'": '&#039;'
        };
        return text.replace(/[&<>"']/g, function(m) { return map[m]; });
    }

    formatDate(dateString) {
        const date = new Date(dateString);
        return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
    }
}

// Initialize posts manager when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.postsManager = new PostsManager();
});