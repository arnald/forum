////// Error handling utilities //////
function showError(errorElementId, message) {
  const errorElement = document.getElementById(errorElementId);
  if (errorElement) {
    errorElement.textContent = message;
  }
}

function clearError(errorElementId) {
  const errorElement = document.getElementById(errorElementId);
  if (errorElement) {
    errorElement.textContent = "";
  }
}

function clearAllErrors() {
  document.querySelectorAll(".field-error").forEach((el) => {
    el.textContent = "";
  });
}

////// Toggle add comment form //////
const addCommentBtn = document.querySelector(".btn-comment");
const addCommentForm = document.querySelector(".add-comment");
const closeCommentBtn = document.querySelector(".close-comment-form");

if (addCommentBtn && addCommentForm) {
  addCommentBtn.addEventListener("click", () => {
    closeAllEditForms();
    addCommentForm.classList.toggle("active");
  });
}

if (closeCommentBtn && addCommentForm) {
  closeCommentBtn.addEventListener("click", () => {
    addCommentForm.classList.remove("active");
    clearAllErrors();
  });
}

////// Edit button handlers - Show edit forms //////
document.addEventListener("click", (e) => {
  const target = e.target;

  if (target.classList.contains("btn-edit")) {
    const dataType = target.getAttribute("data-type");

    closeAllEditForms();
    clearAllErrors();

    if (addCommentForm) {
      addCommentForm.classList.remove("active");
    }

    if (dataType === "topic") {
      const topicEditForm = document.querySelector(".edit-topic-form");
      if (topicEditForm) topicEditForm.style.display = "block";
    }

    if (dataType === "comment") {
      const commentId = target.getAttribute("data-comment-id");
      const commentContent = document.querySelector(
        `.comment-content[data-comment-id="${commentId}"]`
      );
      if (commentContent) {
        const editForm = commentContent.querySelector(".edit-comment-form");
        if (editForm) editForm.style.display = "block";
      }
    }
  }

  if (target.classList.contains("close-edit-form")) {
    const editForm = target.closest(".edit-form");
    if (editForm) {
      editForm.style.display = "none";
      clearAllErrors();
    }
  }
});

////// Helper function to close all edit forms //////
function closeAllEditForms() {
  const allEditForms = document.querySelectorAll(".edit-form");
  allEditForms.forEach((form) => {
    form.style.display = "none";
  });
}

////// Topic edit form validation //////
const topicEditForm = document.querySelector("form.topic-edit-form");

if (topicEditForm) {
  topicEditForm.addEventListener("submit", function (e) {
    clearAllErrors();
    let hasError = false;

    const categorySelect = this.querySelector('select[name="category_id"]');
    const titleInput = this.querySelector('input[name="title"]');
    const contentInput = this.querySelector('textarea[name="content"]');
    const imageInput = this.querySelector('input[name="image_path"]');

    const category = categorySelect?.value;
    const title = titleInput?.value.trim();
    const content = contentInput?.value.trim();

    if (!category) {
      showError("error-topic-category", "Category is required");
      hasError = true;
    }

    if (!title) {
      showError("error-topic-title", "Title is required");
      hasError = true;
    } else if (title.length < 3) {
      showError("error-topic-title", "Title must be at least 3 characters");
      hasError = true;
    }

    if (!content) {
      showError("error-topic-content", "Content is required");
      hasError = true;
    } else if (content.length < 10) {
      showError(
        "error-topic-content",
        "Content must be at least 10 characters"
      );
      hasError = true;
    }

    if (imageInput && imageInput.files.length > 0) {
      const file = imageInput.files[0];
      const allowedTypes = ["image/jpeg", "image/png", "image/gif"];
      const maxSize = 2 * 1024 * 1024; // 2MB

      if (!allowedTypes.includes(file.type)) {
        showError(
          "error-topic-image",
          "Only JPG, PNG, or GIF images are allowed"
        );
        hasError = true;
      } else if (file.size > maxSize) {
        showError("error-topic-image", "Image must be smaller than 2MB");
        hasError = true;
      }
    }

    if (hasError) {
      e.preventDefault();
    }
  });

  // Clear errors on user input
  topicEditForm
    .querySelector('select[name="category_id"]')
    ?.addEventListener("change", () => clearError("error-topic-category"));

  topicEditForm
    .querySelector('input[name="title"]')
    ?.addEventListener("input", () => clearError("error-topic-title"));

  topicEditForm
    .querySelector('textarea[name="content"]')
    ?.addEventListener("input", () => clearError("error-topic-content"));

  topicEditForm
    .querySelector('input[name="image_path"]')
    ?.addEventListener("change", () => clearError("error-topic-image"));
}

////// Create Comment form validation //////
const createCommentForm = document.querySelector(
  'form[action="/comments/create"]'
);

if (createCommentForm) {
  createCommentForm.addEventListener("submit", function (e) {
    clearAllErrors();
    let hasError = false;

    const contentInput = this.querySelector('textarea[name="content"]');
    const content = contentInput?.value.trim();

    if (!content) {
      showError("error-comment-content", "Comment is required");
      hasError = true;
    } else if (content.length < 3) {
      showError(
        "error-comment-content",
        "Comment must be at least 3 characters"
      );
      hasError = true;
    }

    if (hasError) {
      e.preventDefault();
    }
  });
}

////// Edit Comment form validation //////
document.querySelectorAll('form[action="/comments/edit"]').forEach((form) => {
  form.addEventListener("submit", function (e) {
    const contentInput = this.querySelector('textarea[name="content"]');
    const content = contentInput?.value.trim();
    const commentId = this.getAttribute("data-comment-id");
    const errorId = `error-edit-comment-${commentId}`;

    clearError(errorId);

    if (!content) {
      showError(errorId, "Comment is required");
      e.preventDefault();
      return false;
    } else if (content.length < 3) {
      showError(errorId, "Comment must be at least 3 characters");
      e.preventDefault();
      return false;
    }
  });

  const commentId = form.getAttribute("data-comment-id");
  const errorId = `error-edit-comment-${commentId}`;

  form
    .querySelector('textarea[name="content"]')
    ?.addEventListener("input", () => clearError(errorId));
});

// Confirm Delete Actions
document
  .querySelector('form[action="/topics/delete"]')
  ?.addEventListener("submit", (e) => {
    if (
      !confirm(
        "Are you sure you want to delete this topic? This action cannot be undone."
      )
    ) {
      e.preventDefault();
    }
  });

document.querySelectorAll('form[action="/comments/delete"]').forEach((form) => {
  form.addEventListener("submit", (e) => {
    if (
      !confirm(
        "Are you sure you want to delete this comment? This action cannot be undone."
      )
    ) {
      e.preventDefault();
    }
  });
});
// ////// Reaction box click triggers button click //////
// document.querySelectorAll(".reaction-box").forEach((box) => {
//   box.addEventListener("click", function () {
//     const btn = this.querySelector("button");
//     if (btn) btn.click();
//   });
// });
