// ==== Homepage ==== //
// Close the details when clicking outside
document.addEventListener("click", function (event) {
  const details = document.querySelector(".category-details");

  if (!details) return;

  if (!details.contains(event.target)) {
    details.removeAttribute("open");
  }
});

// Add active class to a button
document.addEventListener("DOMContentLoaded", function () {
  const buttons = document.querySelectorAll(".nav-categories-btn");
  const path = window.location.pathname;

  buttons.forEach((button) => {
    button.classList.remove("active");
    if (button.getAttribute("href") === path) {
      button.classList.add("active");
    }
  });
});

// ==== Signup ==== //
/*--- Show/Hide Password ---*/
document.addEventListener("DOMContentLoaded", function () {
  const passwordInput = document.getElementById("password");
  const toggleCheckbox = document.getElementById("togglePassword");
  const showIcon = document.getElementById("eye-icon");
  const hiddenIcon = document.getElementById("hidden-icon");

  // Initially hide the hidden icon
  hiddenIcon.style.display = "none";

  toggleCheckbox.addEventListener("change", function () {
    if (toggleCheckbox.checked) {
      passwordInput.type = "text";
      showIcon.style.display = "none";
      hiddenIcon.style.display = "block";
    } else {
      passwordInput.type = "password";
      showIcon.style.display = "block";
      hiddenIcon.style.display = "none";
    }
  });
});

/*--- Reset button removes preserved values ---*/
document.addEventListener("DOMContentLoaded", function () {
  const resetButtons = document.querySelectorAll(".btn-reset-form");

  resetButtons.forEach((button) => {
    button.addEventListener("click", function (e) {
      e.preventDefault();

      const form = this.closest("form");

      const inputs = form.querySelectorAll(
        'input[type="text"], input[type="email"]'
      );
      inputs.forEach((input) => {
        input.value = "";
        input.classList.remove("input-error");
      });

      const passwordInput = form.querySelector('input[type="password"]');
      if (passwordInput) {
        passwordInput.value = "";
        passwordInput.classList.remove("input-error");
      }

      const errorMessages = form.querySelectorAll(".error-message");
      errorMessages.forEach((msg) => (msg.textContent = ""));

      // Focus on first field
      if (inputs.length > 0) {
        inputs[0].focus();
      }
    });
  });
});

// ==== Signin ==== //
/*--- Remove errors after new input ---*/
document.addEventListener("DOMContentLoaded", () => {
  const inputs = document.querySelectorAll(".form-input");

  inputs.forEach((input) => {
    input.addEventListener("input", () => {
      if (input.classList.contains("input-error")) {
        input.classList.remove("input-error");

        const errorSpan = input
          .closest(".input-box")
          .querySelector(".error-message");

        if (errorSpan) {
          errorSpan.textContent = "";
        }
      }
    });
  });
});
