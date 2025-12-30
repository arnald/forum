////// Toggle add comment form //////|
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
  });
}

////// Edit button handlers - Show edit forms //////
document.addEventListener("click", (e) => {
  const target = e.target;

  if (target.classList.contains("btn-edit")) {
    const dataType = target.getAttribute("data-type");

    closeAllEditForms();

    if (addCommentForm) {
      addCommentForm.classList.remove("active");
    }

    if (dataType === "topic") {
      const topicEditForm = document.querySelector(".edit-topic-form");
      if (topicEditForm) {
        topicEditForm.style.display = "block";
      }
    } else if (dataType === "comment") {
      const commentId = target.getAttribute("data-comment-id");
      const commentContent = document.querySelector(
        `.comment-content[data-comment-id="${commentId}"]`
      );
      if (commentContent) {
        const editForm = commentContent.querySelector(".edit-comment-form");
        if (editForm) {
          editForm.style.display = "block";
        }
      }
    }
  }

  if (target.classList.contains("close-edit-form")) {
    const editForm = target.closest(".edit-form");
    if (editForm) {
      editForm.style.display = "none";
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

// ////// Reaction box click triggers button click //////
// document.querySelectorAll(".reaction-box").forEach((box) => {
//   box.addEventListener("click", function () {
//     const btn = this.querySelector("button");
//     if (btn) btn.click();
//   });
// });
