// Topic voting functionality with proper toggle logic
document.addEventListener("DOMContentLoaded", function () {
  initializeVoteButtons();
  setInitialVoteStates();
});

function initializeVoteButtons() {
  const voteButtons = document.querySelectorAll(".like-btn, .dislike-btn");

  voteButtons.forEach((button) => {
    button.addEventListener("click", async function (e) {
      e.preventDefault();

      // Prevent multiple clicks
      if (this.disabled) return;

      const isLike = this.classList.contains("like-btn");
      const commentContent = this.closest(".comment-content");
      const isComment = commentContent !== null;

      let targetId, targetType;

      if (isComment) {
        targetId = parseInt(commentContent.dataset.commentId);
        targetType = "comment";
      } else {
        // It's a topic vote
        const urlParts = window.location.pathname.split("/");
        targetId = parseInt(urlParts[urlParts.length - 1]);
        targetType = "topic";
      }

      const reactionType = isLike ? 1 : -1;

      try {
        // Disable buttons during request
        const reactionsContainer = this.closest(".reactions");
        disableVoteButtons(reactionsContainer, true);

        // If clicking the same button again, we're toggling off
        await castVote(targetId, targetType, reactionType, this);
      } catch (error) {
        console.error("Error casting vote:", error);
        alert("Failed to cast vote. Please try again.");
      } finally {
        const reactionsContainer = this.closest(".reactions");
        disableVoteButtons(reactionsContainer, false);
      }
    });
  });
}

function disableVoteButtons(container, disabled) {
  const buttons = container.querySelectorAll(".like-btn, .dislike-btn");
  buttons.forEach((btn) => {
    btn.disabled = disabled;
    if (disabled) {
      container.classList.add("loading");
    } else {
      container.classList.remove("loading");
    }
  });
}

async function castVote(targetId, targetType, reactionType, buttonElement) {
  const payload = {
    reactionType: reactionType,
  };

  if (targetType === "comment") {
    payload.commentId = targetId;
  } else {
    payload.topicId = targetId;
  }

  const response = await fetch("/api/vote/cast", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
    body: JSON.stringify(payload),
  });

  if (response.status === 401) {
    window.location.href = "/login";
    return;
  }

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    throw new Error(errorData.message || "Failed to cast vote");
  }

  // Update the UI with new counts
  await updateVoteUI(targetId, targetType, buttonElement);
}

async function updateVoteUI(targetId, targetType, clickedButton) {
  const paramName = targetType === "comment" ? "comment_id" : "topic_id";
  // Get current counts from the backend
  const countsResponse = await fetch(
    `/api/vote/counts?${paramName}=${targetId}`,
    {
      method: "GET",
      credentials: "include",
    }
  );

  if (!countsResponse.ok) {
    throw new Error("Failed to get vote counts");
  }

  const countsData = await countsResponse.json();
  const counts = countsData.data;

  // Find the reaction container
  const reactionsContainer = clickedButton.closest(".reactions");
  const likeCount = reactionsContainer.querySelector(".like-count");
  const dislikeCount = reactionsContainer.querySelector(".dislike-count");
  const likeBtn = reactionsContainer.querySelector(".like-btn");
  const dislikeBtn = reactionsContainer.querySelector(".dislike-btn");

  // Add animation class
  likeCount.classList.add("updating");
  dislikeCount.classList.add("updating");

  // Update counts
  likeCount.textContent = counts.upvotes;
  dislikeCount.textContent = counts.downvotes;

  // Update vote score if it exists (for topics)
  const voteScoreElement = reactionsContainer
    .closest(".topic-body-container")
    ?.querySelector(".views-count");
  if (voteScoreElement) {
    voteScoreElement.textContent = counts.score;
  }

  // Toggle button states - if clicking active button, toggle it off
  // Otherwise, activate clicked button and deactivate the other
  const wasActive = clickedButton.classList.contains("active");

  likeBtn.classList.remove("active");
  dislikeBtn.classList.remove("active");

  if (!wasActive) {
    clickedButton.classList.add("active");
  }

  // Remove animation class after animation completes
  setTimeout(() => {
    likeCount.classList.remove("updating");
    dislikeCount.classList.remove("updating");
  }, 300);
}

// Set initial vote states based on UserVote data
function setInitialVoteStates() {
  // For topic vote
  const topicContainer = document.querySelector(".topic-body-container");
  if (topicContainer) {
    const userVote = parseInt(topicContainer.dataset.userVote);
    if (!isNaN(userVote)) {
      const topicReactions = topicContainer.querySelector(".reactions");
      if (userVote === 1) {
        topicReactions.querySelector(".like-btn")?.classList.add("active");
      } else if (userVote === -1) {
        topicReactions.querySelector(".dislike-btn")?.classList.add("active");
      }
    }
  }

  // For comment votes
  const comments = document.querySelectorAll(".comment-content");
  comments.forEach((comment) => {
    const userVote = parseInt(comment.dataset.userVote);
    if (!isNaN(userVote)) {
      const reactions = comment.querySelector(".reactions");
      if (userVote === 1) {
        reactions.querySelector(".like-btn")?.classList.add("active");
      } else if (userVote === -1) {
        reactions.querySelector(".dislike-btn")?.classList.add("active");
      }
    }
  });
}
