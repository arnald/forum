document.addEventListener("DOMContentLoaded", () => {
  // Multi-select category dropdown
  (function () {
    const ms = document.getElementById("categorySelect");
    if (!ms) return;

    const multiOptions = document.getElementById("multiOptions");
    const chipsContainer = document.getElementById("chips");
    const placeholder = ms.querySelector(".placeholder");

    function open() {
      ms.classList.add("open");
      ms.setAttribute("aria-expanded", "true");
      multiOptions.focus();
    }
    function close() {
      ms.classList.remove("open");
      ms.setAttribute("aria-expanded", "false");
    }
    ms.addEventListener("click", (e) => {
      if (e.target.closest(".options")) return;
      if (ms.classList.contains("open")) close();
      else open();
    });

    // clicking outside closes, clicking "Esc" closes
    document.addEventListener("click", (e) => {
      if (!ms.contains(e.target)) close();
    });

    document.addEventListener("keydown", (e) => {
      if (e.key === "Escape") close();
    });

    // update chips based on checkboxes
    function rebuildChips() {
      chipsContainer.innerHTML = "";
      const checked = ms.querySelectorAll('input[type="checkbox"]:checked');
      if (checked.length === 0) {
        placeholder.style.display = "inline";
      } else {
        placeholder.style.display = "none";
      }
      checked.forEach((cb) => {
        const label =
          cb.parentElement.querySelector(".option-label")?.textContent ||
          cb.value;
        const chip = document.createElement("span");
        chip.className = "chip";
        chip.textContent = label;

        const btn = document.createElement("button");
        btn.type = "button";
        btn.setAttribute("aria-label", `Remove ${label}`);
        btn.style.background = "transparent";
        btn.style.border = "none";
        btn.style.marginLeft = "8px";
        btn.style.cursor = "pointer";
        btn.textContent = "✕";
        btn.addEventListener("click", (ev) => {
          ev.stopPropagation();
          cb.checked = false;
          cb.dispatchEvent(new Event("change", { bubbles: true }));
        });

        chip.appendChild(btn);
        chipsContainer.appendChild(chip);
      });
    }

    const checkboxes = ms.querySelectorAll('input[type="checkbox"]');
    checkboxes.forEach((cb) => {
      cb.addEventListener("change", () => {
        rebuildChips();
      });
    });

    // initialize
    rebuildChips();

    ms.addEventListener("keydown", (e) => {
      if (e.key === "Enter" || e.key === " ") {
        e.preventDefault();
        if (ms.classList.contains("open")) close();
        else open();
      }
    });
  })();

  // File upload handling
  (function () {
    const fileInput = document.getElementById("image-upload");
    const uploadBox = document.getElementById("uploadBox");
    const fileNameDisplay = document.getElementById("file-name");
    const errorEl = document.getElementById("error-image");

    if (!fileInput || !uploadBox) return;

    uploadBox.addEventListener("click", () => fileInput.click());

    fileInput.addEventListener("change", () => {
      errorEl.textContent = "";
      fileNameDisplay.textContent = "";

      const file = fileInput.files[0];
      if (!file) {
        return;
      }

      const allowedTypes = ["image/jpeg", "image/png", "image/gif"];
      const maxSizeMB = 20;

      if (!allowedTypes.includes(file.type)) {
        errorEl.textContent = "Only JPEG, PNG, or GIF images are allowed.";
        fileInput.value = "";
        return;
      }

      if (file.size > maxSizeMB * 1024 * 1024) {
        errorEl.textContent = "Image is too large. Maximum size is 20 MB.";
        fileInput.value = "";
        return;
      }

      fileNameDisplay.textContent = `Selected: ${file.name}`;
      fileNameDisplay.style.color = "#068f56";
      fileNameDisplay.style.marginTop = "0.5rem";
      fileNameDisplay.style.fontSize = "0.9rem";
    });
  })();

  // Form validation and reset handling
  (function () {
    const form = document.getElementById("createPostForm");
    if (!form) return;

    // Form reset handler - clears all fields and error messages
    form.addEventListener("reset", () => {
      document.getElementById("error-categories").textContent = "";
      document.getElementById("error-title").textContent = "";
      document.getElementById("error-content").textContent = "";
      document.getElementById("error-image").textContent = "";
      document.getElementById("file-name").textContent = "";

      const checkboxes = document.querySelectorAll(
        '#categorySelect input[type="checkbox"]'
      );
      checkboxes.forEach((cb) => {
        cb.checked = false;
      });

      const event = new Event("change", { bubbles: true });
      if (checkboxes.length > 0) {
        checkboxes[0].dispatchEvent(event);
      }

      const fileInput = document.getElementById("image-upload");
      if (fileInput) {
        fileInput.value = "";
      }
    });

    // Client-side validation before submit
    form.addEventListener("submit", (e) => {
      let hasError = false;

      document.getElementById("error-categories").textContent = "";
      document.getElementById("error-title").textContent = "";
      document.getElementById("error-content").textContent = "";
      document.getElementById("error-image").textContent = "";

      const selectedCategories = document.querySelectorAll(
        '#categorySelect input[type="checkbox"]:checked'
      );
      if (selectedCategories.length === 0) {
        document.getElementById("error-categories").textContent =
          "Please select at least one category";
        hasError = true;
      }

      const title = document.getElementById("title").value.trim();
      if (!title) {
        document.getElementById("error-title").textContent =
          "Title is required";
        hasError = true;
      } else if (title.length < 3) {
        document.getElementById("error-title").textContent =
          "Title must be at least 3 characters";
        hasError = true;
      } else if (title.length > 200) {
        document.getElementById("error-title").textContent =
          "Title must not exceed 200 characters";
        hasError = true;
      }

      const content = document.getElementById("content").value.trim();
      if (!content) {
        document.getElementById("error-content").textContent =
          "Content is required";
        hasError = true;
      } else if (content.length < 10) {
        document.getElementById("error-content").textContent =
          "Content must be at least 10 characters";
        hasError = true;
      } else if (content.length > 5000) {
        document.getElementById("error-content").textContent =
          "Content must not exceed 5000 characters";
        hasError = true;
      }

      if (hasError) {
        e.preventDefault();
      }
    });
  })();
});
