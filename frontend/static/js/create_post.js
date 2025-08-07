document.addEventListener("DOMContentLoaded", () => {
  (function () {
    const ms = document.getElementById("categorySelect");
    const multiOptions = document.getElementById("multiOptions");
    const chipsContainer = document.getElementById("chips");
    const placeholder = ms.querySelector(".placeholder");

    // open/close
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
      // if clicking checkbox or inside options, ignore
      if (e.target.closest(".options")) return;
      if (ms.classList.contains("open")) close();
      else open();
    });

    // clicking outside closes
    document.addEventListener("click", (e) => {
      if (!ms.contains(e.target)) close();
    });

    // Escape key
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

        // small remove icon inside chip
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

    // attach change listeners to checkboxes
    const checkboxes = ms.querySelectorAll('input[type="checkbox"]');
    checkboxes.forEach((cb) => {
      cb.addEventListener("change", () => {
        rebuildChips();
      });
    });

    // initialize
    rebuildChips();

    // keyboard navigation: open on Enter or Space when focused
    ms.addEventListener("keydown", (e) => {
      if (e.key === "Enter" || e.key === " ") {
        e.preventDefault();
        if (ms.classList.contains("open")) close();
        else open();
      }
    });
  })();

  (function () {
    const uploadBox = document.getElementById("uploadBox");
    const fileInput = document.getElementById("image-upload");
    const preview = document.getElementById("imagePreview");
    const uploadPlaceholder = document.querySelector(".upload-placeholder");
    const errorEl = document.getElementById("error-image");

    // Create "remove" button dynamically
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

    // Append remove button after preview
    uploadBox.appendChild(removeBtn);

    uploadBox.addEventListener("click", () => fileInput.click());

    fileInput.addEventListener("change", () => {
      errorEl.textContent = "";

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
        errorEl.textContent = !allowedTypes.includes(file.type)
          ? "Only JPEG, PNG, or GIF images are allowed."
          : "Image is too big. Max size is 20 MB.";
        fileInput.value = "";
        preview.style.display = "none";
        removeBtn.style.display = "none";
        uploadPlaceholder.style.display = "flex";
        return;
      }

      // Hide placeholder and show preview
      uploadPlaceholder.style.display = "none";

      // Show preview
      const reader = new FileReader();
      reader.onload = (e) => {
        preview.src = e.target.result;
        preview.style.display = "block";
        removeBtn.style.display = "inline-block"; // show remove button
      };
      reader.readAsDataURL(file);
    });

    // Remove button click
    removeBtn.addEventListener("click", (ev) => {
      ev.stopPropagation(); // avoid triggering file picker
      fileInput.value = "";
      preview.src = "";
      preview.style.display = "none";
      removeBtn.style.display = "none";
      uploadPlaceholder.style.display = "flex";
    });
  })();
});
