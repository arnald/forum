// Reaction box click triggers button click
document.querySelectorAll(".reaction-box").forEach((box) => {
  box.addEventListener("click", function () {
    const btn = this.querySelector("button");
    if (btn) btn.click();
  });
});

// Toggle add comment form
const addCommentBtn = document.querySelector(".btn-comment");
const commentForm = document.querySelector(".add-comment");
const closeBtn = document.querySelector(".close-comment-form");

if (addCommentBtn && commentForm) {
  addCommentBtn.addEventListener("click", () =>
    commentForm.classList.toggle("active")
  );
}
if (closeBtn && commentForm) {
  closeBtn.addEventListener("click", () =>
    commentForm.classList.remove("active")
  );
}

// Image upload (post + comment)
document.addEventListener("DOMContentLoaded", () => {
  initUploadFeature(
    "uploadBox",
    "image-upload",
    "imagePreview",
    ".upload-placeholder",
    "error-image"
  );

  initUploadFeature(
    "commentUploadBox",
    "commentImageUpload",
    "commentImagePreview",
    ".upload-placeholder"
  );
});

function initUploadFeature(
  boxId,
  inputId,
  previewId,
  placeholderSelector,
  errorId,
  removeBtnId
) {
  const uploadBox = document.getElementById(boxId);
  const fileInput = document.getElementById(inputId);
  const preview = document.getElementById(previewId);
  const uploadPlaceholder = uploadBox
    ? uploadBox.querySelector(placeholderSelector)
    : null;
  const errorEl = errorId ? document.getElementById(errorId) : null;

  // Early return if essential elements are missing
  if (!uploadBox || !fileInput || !preview || !uploadPlaceholder) {
    return;
  }

  let removeBtn = removeBtnId ? document.getElementById(removeBtnId) : null;
  if (!removeBtn) {
    removeBtn = document.createElement("button");
    removeBtn.type = "button";
    removeBtn.textContent = "✕ Remove";
    Object.assign(removeBtn.style, {
      display: "none",
      marginTop: "8px",
      background: "#e63946",
      color: "white",
      border: "none",
      borderRadius: "5px",
      padding: "4px 8px",
      cursor: "pointer",
    });
    uploadBox.appendChild(removeBtn);
  }

  uploadBox.addEventListener("click", () => fileInput.click());

  fileInput.addEventListener("change", () => {
    if (errorEl) errorEl.textContent = "";
    const file = fileInput.files[0];
    if (!file) return reset();

    const allowedTypes = ["image/jpeg", "image/png", "image/gif"];
    const maxSizeMB = 20;
    if (
      !allowedTypes.includes(file.type) ||
      file.size > maxSizeMB * 1024 * 1024
    ) {
      if (errorEl)
        errorEl.textContent = !allowedTypes.includes(file.type)
          ? "Only JPEG, PNG, or GIF images are allowed."
          : "Image is too big. Max size is 20 MB.";
      return reset();
    }

    const reader = new FileReader();
    reader.onload = (e) => {
      preview.src = e.target.result;
      preview.style.display = "block";
      removeBtn.style.display = "inline-block";
      uploadPlaceholder.style.display = "none";
    };
    reader.readAsDataURL(file);
  });

  removeBtn.addEventListener("click", (ev) => {
    ev.stopPropagation();
    reset();
  });

  function reset() {
    fileInput.value = "";
    preview.src = "";
    preview.style.display = "none";
    removeBtn.style.display = "none";
    uploadPlaceholder.style.display = "flex";
  }
}

// Buttons (edit/save/cancel/delete)
document.addEventListener("click", (e) => {
  const target = e.target;

  // --- EDIT ---
  if (target.classList.contains("btn-edit")) {
    const container =
      target.closest(".comment-content") ||
      target.closest(".topic-container")?.querySelector(".topic-content");
    const textEl = container.querySelector(".post-text, .comment-text");
    const imgEl = container.querySelector(".post-image, .comment-image");

    container.dataset.originalText = textEl.textContent.trim();
    container.dataset.originalImg = imgEl ? imgEl.src : "";

    textEl.innerHTML = `
      <textarea class="edit-text" style="width:100%; min-height:80px;">${
        container.dataset.originalText
      }</textarea>
      <div class="edit-image-box" id="editImageBox">
        <img id="editImagePreview" src="${container.dataset.originalImg || ""}" 
             style="max-width:100%; max-height:600px; border-radius:8px; margin-top:8px; ${
               container.dataset.originalImg ? "" : "display:none"
             }">
        <button type="button" id="editRemoveImage" 
                style="margin-top:8px; background:#e63946; color:white; border:none; border-radius:5px; padding:4px 8px; cursor:pointer; ${
                  container.dataset.originalImg ? "" : "display:none"
                }">
          ✕ Remove Image
        </button>
        <input type="file" id="edit-image-upload" accept="image/jpeg,image/png,image/gif" hidden />
        <div class="upload-placeholder" style="margin-top:8px; color:gray; cursor:pointer;">Click to upload new image</div>
      </div>
      <button class="btn-save action-btn" style="background:#068f56; color:white; margin-top:8px;">Save</button>
      <button class="btn-cancel action-btn" style="background:#ccc; margin-top:8px;">Cancel</button>
    `;

    initUploadFeature(
      "editImageBox",
      "edit-image-upload",
      "editImagePreview",
      ".upload-placeholder",
      null,
      "editRemoveImage"
    );
    return;
  }

  // --- SAVE ---
  if (target.classList.contains("btn-save")) {
    const container =
      target.closest(".comment-content") ||
      target.closest(".topic-container")?.querySelector(".topic-content");
    const newText = container.querySelector(".edit-text").value.trim();
    const imgPreview = container.querySelector("#editImagePreview");
    const imgSrc =
      imgPreview && imgPreview.style.display !== "none" ? imgPreview.src : "";

    const textEl = container.querySelector(".post-text, .comment-text");
    textEl.textContent = newText;

    let imgEl = container.querySelector(".comment-image, .post-image");
    let imgBox =
      container.querySelector(".img-box") ||
      (() => {
        const box = document.createElement("div");
        box.className = "img-box";
        textEl.insertAdjacentElement("afterend", box);
        return box;
      })();

    if (imgSrc) {
      if (!imgEl) {
        imgEl = document.createElement("img");
        imgEl.className = textEl.classList.contains("post-text")
          ? "post-image"
          : "comment-image";
        Object.assign(imgEl.style, {
          maxWidth: "100%",
          maxHeight: "600px",
          borderRadius: "8px",
          marginTop: "8px",
        });
        imgBox.appendChild(imgEl);
      }
      imgEl.src = imgSrc;
    } else if (imgEl) imgEl.remove();

    container.querySelector(".edit-image-box")?.remove();
    container.querySelector(".btn-save")?.remove();
    container.querySelector(".btn-cancel")?.remove();
    return;
  }

  // --- CANCEL ---
  if (target.classList.contains("btn-cancel")) {
    const container =
      target.closest(".comment-content") ||
      target.closest(".topic-container")?.querySelector(".topic-content");
    const textEl = container.querySelector(".post-text, .comment-text");
    const imgEl = container.querySelector(".post-image, .comment-image");

    textEl.textContent = container.dataset.originalText || "";

    if (container.dataset.originalImg) {
      if (!imgEl) {
        const img = document.createElement("img");
        img.className = textEl.classList.contains("post-text")
          ? "post-image"
          : "comment-image";
        img.src = container.dataset.originalImg;
        img.style.maxWidth = "100%";
        (
          container
            .querySelector(".topic-body, .comment-body")
            ?.querySelector(".img-box") ||
          container.appendChild(
            Object.assign(document.createElement("div"), {
              className: "img-box",
            })
          )
        ).appendChild(img);
      } else imgEl.src = container.dataset.originalImg;
    } else if (imgEl) imgEl.remove();
    return;
  }

  // --- DELETE ---
  if (target.classList.contains("btn-delete")) {
    if (confirm("Are you sure you want to delete this?")) {
      const commentContainer = target.closest(".comment-content");
      const postContainer = target.closest(".topic-container");
      if (commentContainer) commentContainer.remove();
      else if (postContainer) {
        postContainer.remove();
        window.location.href = "/";
      }
    }
  }
});
