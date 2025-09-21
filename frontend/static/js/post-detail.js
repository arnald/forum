/**
 * Post Detail Page Manager
 *
 * This class handles functionality for viewing individual posts including:
 * - Loading and displaying full post content
 * - Loading and displaying threaded comments
 * - Creating new comments and replies
 * - Voting on posts and comments
 * - Managing user interactions on the post detail page
 */
class PostDetailManager {
    constructor() {
        // Store the current post ID from URL
        this.postId = null;
        // Store loaded comments for the post
        this.comments = [];
        // Initialize the post detail functionality
        this.init();
    }

    /**
     * Initialize the post detail manager
     * Extracts post ID from URL and loads post data
     */
    init() {
        // Get post ID from URL parameters (?id=...)
        this.postId = this.getPostIdFromUrl();
        if (this.postId) {
            // Load the post content
            this.loadPost();
            // Load comments for this post
            this.loadComments();
        }
        // Set up event listeners for forms and interactions
        this.setupEventListeners();
    }

    getPostIdFromUrl() {
        const urlParams = new URLSearchParams(window.location.search);
        return urlParams.get('id');
    }

    setupEventListeners() {
        const commentForm = document.getElementById('commentForm');
        if (commentForm) {
            commentForm.addEventListener('submit', (e) => this.handleCreateComment(e));
        }
    }

    async loadPost() {
        try {
            const response = await fetch(`/api/v1/posts/get?id=${this.postId}`, {
                credentials: 'include'
            });
            const data = await response.json();

            if (response.ok && data.success) {
                this.displayPost(data.post);
            } else {
                this.showMessage('Post not found', 'error');
            }
        } catch (error) {
            this.showMessage('Error loading post', 'error');
        }
    }

    async loadComments() {
        try {
            const response = await fetch(`/api/v1/comments/tree?post_id=${this.postId}`, {
                credentials: 'include'
            });
            const data = await response.json();

            if (response.ok && data.success) {
                this.comments = data.commentTree || [];
                this.displayComments();
            }
        } catch (error) {
            console.error('Error loading comments:', error);
        }
    }

    displayPost(post) {
        const container = document.getElementById('post-container');
        if (!container) return;

        const categoriesHtml = post.categories && post.categories.length > 0
            ? `<div class="post-categories">
                ${post.categories.map(cat => `<span class="category-tag">${this.escapeHtml(cat.name || cat)}</span>`).join('')}
               </div>`
            : '';

        container.innerHTML = `
            <article class="post">
                <header class="post-header">
                    <h1 class="post-title">${this.escapeHtml(post.title)}</h1>
                    <div class="post-meta">
                        By ${this.escapeHtml(post.username || 'Unknown')} • ${this.formatDate(post.created_at)}
                    </div>
                </header>
                <div class="post-content">
                    ${this.escapeHtml(post.content).replace(/\n/g, '<br>')}
                </div>
                ${categoriesHtml}
                <div class="post-actions">
                    <div class="vote-buttons">
                        <button class="vote-btn" onclick="postDetailManager.vote('${post.id}', 'post', 'like')">
                            👍 Like (${post.like_count || 0})
                        </button>
                        <button class="vote-btn" onclick="postDetailManager.vote('${post.id}', 'post', 'dislike')">
                            👎 Dislike (${post.dislike_count || 0})
                        </button>
                    </div>
                </div>
            </article>
        `;

        // Show comment form if user is logged in
        const commentForm = document.getElementById('commentForm');
        if (commentForm && window.authManager && window.authManager.isLoggedIn()) {
            commentForm.style.display = 'block';
        }
    }

    displayComments() {
        const container = document.getElementById('comments-list');
        if (!container) return;

        if (this.comments.length === 0) {
            container.innerHTML = '<p>No comments yet. Be the first to comment!</p>';
            return;
        }

        container.innerHTML = this.comments.map(comment => this.renderComment(comment)).join('');
    }

    renderComment(commentNode, level = 0) {
        const comment = commentNode.Comment || commentNode;
        const replies = commentNode.Replies || [];

        const marginLeft = level * 20;

        let html = `
            <div class="comment" data-comment-id="${comment.id}" style="margin-left: ${marginLeft}px;">
                <div class="comment-header">
                    <span class="comment-author">${this.escapeHtml(comment.username || 'Anonymous')}</span>
                    <span class="comment-date">${this.formatDate(comment.created_at)}</span>
                </div>
                <div class="comment-content">
                    ${this.escapeHtml(comment.content).replace(/\n/g, '<br>')}
                </div>
                <div class="comment-actions">
                    <div class="vote-buttons">
                        <button class="vote-btn" onclick="postDetailManager.vote('${comment.id}', 'comment', 'like')">
                            👍 (${comment.like_count || 0})
                        </button>
                        <button class="vote-btn" onclick="postDetailManager.vote('${comment.id}', 'comment', 'dislike')">
                            👎 (${comment.dislike_count || 0})
                        </button>
                    </div>
                    ${window.authManager && window.authManager.isLoggedIn()
                        ? `<button class="btn" onclick="postDetailManager.showReplyForm('${comment.id}')">Reply</button>`
                        : ''
                    }
                </div>
                <div id="reply-form-${comment.id}" class="reply-form" style="display: none;">
                    <textarea id="reply-content-${comment.id}" placeholder="Write your reply..." rows="3"></textarea>
                    <div class="reply-actions">
                        <button class="btn" onclick="postDetailManager.submitReply('${comment.id}')">Post Reply</button>
                        <button class="btn btn-secondary" onclick="postDetailManager.hideReplyForm('${comment.id}')">Cancel</button>
                    </div>
                </div>
            </div>
        `;

        // Add replies recursively
        if (replies && replies.length > 0) {
            html += replies.map(reply => this.renderComment(reply, level + 1)).join('');
        }

        return html;
    }

    async handleCreateComment(e) {
        e.preventDefault();

        if (!window.authManager.isLoggedIn()) {
            this.showMessage('Please login to comment', 'error');
            return;
        }

        const content = document.getElementById('comment-content').value.trim();
        if (!content) {
            this.showMessage('Please enter a comment', 'error');
            return;
        }

        try {
            const response = await fetch('/api/v1/comments', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
                body: JSON.stringify({
                    content: content,
                    post_id: this.postId
                })
            });

            const data = await response.json();

            if (response.ok && data.success) {
                document.getElementById('comment-content').value = '';
                this.loadComments(); // Reload comments
                this.showMessage('Comment added successfully!', 'success');
            } else {
                this.showMessage(data.message || 'Failed to add comment', 'error');
            }
        } catch (error) {
            this.showMessage('Error adding comment', 'error');
        }
    }

    showReplyForm(commentId) {
        const form = document.getElementById(`reply-form-${commentId}`);
        if (form) {
            form.style.display = 'block';
        }
    }

    hideReplyForm(commentId) {
        const form = document.getElementById(`reply-form-${commentId}`);
        if (form) {
            form.style.display = 'none';
            document.getElementById(`reply-content-${commentId}`).value = '';
        }
    }

    async submitReply(parentId) {
        if (!window.authManager.isLoggedIn()) {
            this.showMessage('Please login to reply', 'error');
            return;
        }

        const content = document.getElementById(`reply-content-${parentId}`).value.trim();
        if (!content) {
            this.showMessage('Please enter a reply', 'error');
            return;
        }

        try {
            const response = await fetch('/api/v1/comments', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
                body: JSON.stringify({
                    content: content,
                    post_id: this.postId,
                    parent_id: parentId
                })
            });

            const data = await response.json();

            if (response.ok && data.success) {
                this.hideReplyForm(parentId);
                this.loadComments(); // Reload comments
                this.showMessage('Reply added successfully!', 'success');
            } else {
                this.showMessage(data.message || 'Failed to add reply', 'error');
            }
        } catch (error) {
            this.showMessage('Error adding reply', 'error');
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
                // Reload post and comments to update vote counts
                this.loadPost();
                this.loadComments();
            } else {
                this.showMessage(data.message || 'Failed to vote', 'error');
            }
        } catch (error) {
            this.showMessage('Error voting', 'error');
        }
    }

    showMessage(message, type) {
        const errorDiv = document.getElementById('error-message');

        if (errorDiv) {
            errorDiv.textContent = message;
            errorDiv.style.display = 'block';
            errorDiv.style.color = type === 'error' ? 'red' : 'green';

            setTimeout(() => {
                errorDiv.style.display = 'none';
            }, 5000);
        }
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

// Initialize post detail manager when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.postDetailManager = new PostDetailManager();
});