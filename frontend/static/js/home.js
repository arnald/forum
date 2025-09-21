/**
 * Home Page Manager
 *
 * This class handles functionality specific to the home page including:
 * - Loading and displaying available categories
 * - Loading and displaying recent posts
 * - Providing an overview of forum activity
 */
class HomeManager {
    constructor() {
        // Initialize the home page functionality
        this.init();
    }

    /**
     * Initialize the home page manager
     * Loads categories and recent posts when page loads
     */
    init() {
        // Load all available categories to display to users
        this.loadCategories();
        // Load recent posts to show forum activity
        this.loadRecentPosts();
    }

    async loadCategories() {
        try {
            const response = await fetch('/api/v1/categories/all', {
                credentials: 'include'
            });
            const data = await response.json();

            if (response.ok && data.success) {
                this.displayCategories(data.categories || []);
            }
        } catch (error) {
            console.error('Error loading categories:', error);
        }
    }

    async loadRecentPosts() {
        try {
            const response = await fetch('/api/v1/posts/all?limit=5', {
                credentials: 'include'
            });
            const data = await response.json();

            if (response.ok && data.success) {
                this.displayRecentPosts(data.posts || []);
            }
        } catch (error) {
            console.error('Error loading recent posts:', error);
        }
    }

    displayCategories(categories) {
        const container = document.getElementById('categories-list');
        if (!container) return;

        if (categories.length === 0) {
            container.innerHTML = '<p>No categories available.</p>';
            return;
        }

        container.innerHTML = categories.map(category => `
            <div class="category-item">
                <div class="category-name">${this.escapeHtml(category.name)}</div>
                <div class="category-description">${this.escapeHtml(category.description || 'No description')}</div>
            </div>
        `).join('');
    }

    displayRecentPosts(posts) {
        const container = document.getElementById('recent-posts-list');
        if (!container) return;

        if (posts.length === 0) {
            container.innerHTML = '<p>No recent posts.</p>';
            return;
        }

        container.innerHTML = posts.map(post => `
            <div class="post">
                <h4 class="post-title">${this.escapeHtml(post.title)}</h4>
                <div class="post-meta">
                    By ${this.escapeHtml(post.username || 'Unknown')} • ${this.formatDate(post.created_at)}
                </div>
                <div class="post-content">
                    ${this.escapeHtml(post.content).substring(0, 150)}${post.content.length > 150 ? '...' : ''}
                </div>
                <div class="post-actions">
                    <button class="btn" onclick="window.location.href='/post/?id=${post.id}'">
                        Read More
                    </button>
                </div>
            </div>
        `).join('');
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
        return date.toLocaleDateString();
    }
}

// Initialize home manager when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.homeManager = new HomeManager();
});