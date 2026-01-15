// Notification System
(function () {
  let eventSource = null;
  let unreadCount = 0;

  // DOM elements
  const bell = document.getElementById("notificationBell");
  const badge = document.getElementById("notificationBadge");
  const dropdown = document.getElementById("notificationDropdown");
  const notificationList = document.getElementById("notificationList");
  const markAllReadBtn = document.getElementById("markAllReadBtn");

  // Initialize
  function init() {
    if (!bell) return;

    // Toggle dropdown on bell click
    bell.addEventListener("click", (e) => {
      e.stopPropagation();

      const isHidden = getComputedStyle(dropdown).display === "none";
      dropdown.style.display = isHidden ? "flex" : "none";
      if (isHidden) loadNotifications();
    });

    // Close dropdown when clicking outside
    document.addEventListener("click", (e) => {
      if (!dropdown.contains(e.target) && !bell.contains(e.target)) {
        dropdown.style.display = "none";
      }
    });

    // Mark all as read
    markAllReadBtn.addEventListener("click", markAllAsRead);

    // Connect to SSE and load initial notifications
    connectSSE();
    loadNotifications();
  }

  // Connect to SSE stream
  function connectSSE() {
    eventSource = new EventSource("/api/notifications/stream");

    eventSource.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);

        if (data.type === "connected") {
          console.log("Connected to notification stream");
        } else if (data.type === "unread_count") {
          updateBadge(data.count);
        } else if (data.id) {
          // New notification
          unreadCount++;
          updateBadge(unreadCount);
          if (getComputedStyle(dropdown).display !== "none") {
            loadNotifications();
          }
        }
      } catch (error) {
        console.error("Error parsing notification:", error);
      }
    };

    eventSource.onerror = () => {
      console.log("SSE connection lost, reconnecting...");
      eventSource.close();
      setTimeout(connectSSE, 5000); // Reconnect after 5 seconds
    };
  }

  // Load notifications
  async function loadNotifications() {
    try {
      const response = await fetch("/api/notifications?limit=20");
      if (!response.ok) throw new Error("Failed to load notifications");

      const notifications = await response.json();
      renderNotifications(notifications);

      const unread = notifications.filter((n) => !n.isRead).length;
      updateBadge(unread);
    } catch (error) {
      console.error("Error loading notifications:", error);
    }
  }

  // Render notifications
  function renderNotifications(notifications) {
    if (!notifications || notifications.length === 0) {
      notificationList.innerHTML =
        '<p class="notification-empty">No notifications yet</p>';
      return;
    }

    notificationList.innerHTML = notifications
      .map((n) => {
        const icon =
          n.type === "like" ? "💚" : n.type === "dislike" ? "🤮" : n.type === "mention" ? "@" : "💬";
        const timeAgo = formatTimeAgo(new Date(n.createdAt));

        return `
        <div class="notification-item ${n.isRead ? "" : "unread"}" 
             data-id="${n.id}"
             data-read="${n.isRead}"
             data-related-type="${n.relatedType || ""}"
             data-related-id="${n.relatedId || ""}">
          <div class="notification-icon ${n.type}">${icon}</div>
          <div class="notification-content">
            <div class="notification-title">${escapeHtml(n.title)}</div>
            <div class="notification-message">${escapeHtml(n.message)}</div>
            <div class="notification-time">${timeAgo}</div>
          </div>
          ${!n.isRead ? '<div class="notification-unread-dot"></div>' : ""}
        </div>
      `;
      })
      .join("");

    // Add click handlers
    notificationList.querySelectorAll("notification-item").forEach((item) => {
      item.addEventListener("click", () => {
        const id = item.dataset.id;
        const isRead = item.dataset.read === "true";
        const relatedId = item.dataset.id;
        const relatedType = item.dataset.relatedType;

        if (!isRead) markAsRead(id, item);

        // Navigate to related content
        if (relatedType === "topic") {
          window.location.href = `/topic/${relatedId}`;
        }
      });
    });
  }

  // Mark as read
  async function markAsRead(id, element) {
    try {
      const response = await fetch(`/api/notifications/mark-read?id=${id}`, {
        method: "POST",
      });

      if (response.ok) {
        element.classList.remove("unread");
        element.dataset.read = "true";
        const dot = document.querySelector(".notification-unread-dot");
        if (dot) dot.remove();

        unreadCount = Math.max(0, unreadCount - 1);
        updateBadge(unreadCount);
      }
    } catch (error) {
      console.error("Error marking as read:", error);
    }
  }

  // Mark all as read
  async function markAllAsRead() {
    try {
      const response = await fetch("/api/notifications/mark-all-read", {
        method: "POST",
      });

      if (response.ok) {
        document
          .querySelectorAll(".notification-item.unread")
          .forEach((item) => {
            element.classList.remove("unread");
            element.dataset.read = "true";
            const dot = document.querySelector(".notification-unread-dot");
            if (dot) dot.remove();
          });
        updateBadge(0);
      }
    } catch (error) {
      console.error("Error marking all as read:", error);
    }
  }

  // Update badge
  function updateBadge(count) {
    unreadCount = count;
    if (count > 0) {
      badge.textContent = count > 99 ? "99+" : count;
      badge.style.display = "block";
    } else {
      badge.style.display = "none";
    }
  }

  // Format time ago
  function formatTimeAgo(date) {
    const seconds = Math.floor((new Date() - date) / 1000);
    if (seconds < 60) return "just now";
    if (seconds < 3600) return `${Math.floor(seconds) / 60}m ago`;
    if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`;
    if (seconds < 604800) return `${Math.floor(seconds / 86400)}d ago`;
    return date.toLocaleDateString();
  }

  // Escape HTML
  function escapeHtml(text) {
    const div = document.createElement("div");
    div.textContent = text;
    return div.innerHTML;
  }

  // Initialize
  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", init);
  }

  // Cleanup on unload
  document.addEventListener("beforeunload", () => {
    if (eventSource) eventSource.close();
  });
})();
