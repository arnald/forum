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
  addCommentBtn.addEventListener("click", () => {
    commentForm.classList.toggle("active");
  });
}

if (closeBtn && commentForm) {
  closeBtn.addEventListener("click", () => {
    commentForm.classList.remove("active");
  });
}

// Image Input, Image preview(both for post creation and comment creation)
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
    "comment-image-upload",
    "commentImagePreview",
    ".upload-placeholder"
  );
});

function initUploadFeature(
  boxId,
  inputId,
  previewId,
  placeholderSelector,
  errorId
) {
  const uploadBox = document.getElementById(boxId);
  const fileInput = document.getElementById(inputId);
  const preview = document.getElementById(previewId);
  const uploadPlaceholder = uploadBox.querySelector(placeholderSelector);
  const errorEl = errorId ? document.getElementById(errorId) : null;

  const removeBtn = document.createElement("button");
  removeBtn.type = "button";
  removeBtn.textContent = "✕ Remove";
  removeBtn.style.display = "none";
  removeBtn.style.marginTop = "8px";
  removeBtn.style.background = "#e63946";
  removeBtn.style.color = "white";
  removeBtn.style.border = "none";
  removeBtn.style.borderRadius = "5px";
  removeBtn.style.padding = "4px 8px";
  removeBtn.style.cursor = "pointer";
  uploadBox.appendChild(removeBtn);

  uploadBox.addEventListener("click", () => fileInput.click());

  fileInput.addEventListener("change", () => {
    if (errorEl) errorEl.textContent = "";
    const file = fileInput.files[0];
    if (!file) {
      preview.style.display = "none";
      removeBtn.style.display = "none";
      uploadPlaceholder.style.display = "flex";
      return;
    }
    const allowedTypes = ["image/jpeg", "image/png", "image/gif"];
    const maxSizeMB = 20;
    if (
      !allowedTypes.includes(file.type) ||
      file.size > maxSizeMB * 1024 * 1024
    ) {
      if (errorEl) {
        errorEl.textContent = !allowedTypes.includes(file.type)
          ? "Only JPEG, PNG, or GIF images are allowed."
          : "Image is too big. Max size is 20 MB.";
      }
      fileInput.value = "";
      preview.style.display = "none";
      removeBtn.style.display = "none";
      uploadPlaceholder.style.display = "flex";
      return;
    }
    uploadPlaceholder.style.display = "none";
    const reader = new FileReader();
    reader.onload = (e) => {
      preview.src = e.target.result;
      preview.style.display = "block";
      removeBtn.style.display = "inline-block";
    };
    reader.readAsDataURL(file);
  });

  removeBtn.addEventListener("click", (ev) => {
    ev.stopPropagation();
    fileInput.value = "";
    preview.src = "";
    preview.style.display = "none";
    removeBtn.style.display = "none";
    uploadPlaceholder.style.display = "flex";
  });
}

// Buttons Functionality
document.addEventListener("click", (e) => {
  const target = e.target;

  // --- EDIT ---
  if (target.classList.contains("btn-edit")) {
    const container =
      target.closest(".comment-content") ||
      target.closest(".topic-container")?.querySelector(".topic-content");
    const textEl = container.querySelector(".post-text, .comment-text");
    const imgEl = container.querySelector(".post-image, .comment-image");

    // Store original
    container.dataset.originalText = textEl.textContent.trim();
    container.dataset.originalImg = imgEl ? imgEl.src : "";

    // Replace text with textarea
    textEl.innerHTML = `
      <textarea class="edit-text" style="width:100%; min-height:80px;">${container.dataset.originalText}</textarea>
    `;

    // Image upload HTML
    let imgHtml = "";
    if (container.dataset.originalImg) {
      imgHtml = `<img id="editImagePreview" src="${container.dataset.originalImg}" style="max-width:100%; border-radius:8px; margin-top:8px;">`;
    }

    textEl.insertAdjacentHTML(
      "beforeend",
      `
      <div class="edit-image-box" id="editImageBox">
        ${imgHtml}
        <input type="file" id="edit-image-upload" accept="image/jpeg,image/png,image/gif" hidden />
        <div class="upload-placeholder" style="margin-top:8px; color:gray; cursor:pointer;">
          Click to upload new image
        </div>
      </div>
      <button class="btn-save action-btn" style="background:#068f56; color:white; margin-top:8px;">Save</button>
      <button class="btn-cancel action-btn" style="background:#ccc; margin-top:8px;">Cancel</button>
    `
    );

    // Init upload
    initUploadFeature(
      "editImageBox",
      "edit-image-upload",
      "editImagePreview",
      ".upload-placeholder"
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

    // Update text
    const textEl = container.querySelector(".post-text, .comment-text");
    textEl.textContent = newText;

    // Update image
    let imgEl = container.querySelector(".post-image, .comment-image");
    if (imgSrc) {
      if (!imgEl) {
        imgEl = document.createElement("img");
        imgEl.className = textEl.classList.contains("post-text")
          ? "post-image"
          : "comment-image";
        imgEl.style.maxWidth = "100%";

        // Append image in correct place for topic vs comment
        const bodyContainer = container.querySelector(
          ".topic-body, .comment-body"
        );
        const imgBox =
          bodyContainer.querySelector(".img-box") ||
          document.createElement("div");
        if (!imgBox.classList.contains("img-box")) {
          imgBox.className = "img-box";
          bodyContainer.appendChild(imgBox);
        }
        imgBox.appendChild(imgEl);
      }
      imgEl.src = imgSrc;
    } else if (imgEl) {
      imgEl.remove();
    }

    return;
  }

  // --- CANCEL ---
  if (target.classList.contains("btn-cancel")) {
    const container =
      target.closest(".comment-content") ||
      target.closest(".topic-container")?.querySelector(".topic-content");
    const textEl = container.querySelector(".post-text, .comment-text");
    const imgEl = container.querySelector(".post-image, .comment-image");

    // Restore original text
    textEl.textContent = container.dataset.originalText || "";

    // Restore original image
    if (container.dataset.originalImg) {
      if (!imgEl) {
        const img = document.createElement("img");
        img.className = textEl.classList.contains("post-text")
          ? "post-image"
          : "comment-image";
        img.src = container.dataset.originalImg;
        img.style.maxWidth = "100%";

        const bodyContainer = container.querySelector(
          ".topic-body, .comment-body"
        );
        const imgBox =
          bodyContainer.querySelector(".img-box") ||
          document.createElement("div");
        if (!imgBox.classList.contains("img-box")) {
          imgBox.className = "img-box";
          bodyContainer.appendChild(imgBox);
        }
        imgBox.appendChild(img);
      } else {
        imgEl.src = container.dataset.originalImg;
      }
    } else if (imgEl) {
      imgEl.remove();
    }

    return;
  }

  // --- DELETE ---
  if (target.classList.contains("btn-delete")) {
    if (confirm("Are you sure you want to delete this?")) {
      const container = target.closest(".comment-content, .topic-content");
      if (container) {
        container.remove();
      }
    }
    return;
  }
});
