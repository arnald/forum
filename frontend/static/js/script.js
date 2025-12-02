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
/*--- Login Type Selector (Username vs Email) ---*/
document.addEventListener("DOMContentLoaded", function () {
  const loginTypeRadios = document.querySelectorAll(".login-type-radio");
  const usernameBox = document.getElementById("usernameBox");
  const emailBox = document.getElementById("emailBox");
  const usernameInput = document.getElementById("username");
  const emailInput = document.getElementById("email");

  loginTypeRadios.forEach((radio) => {
    radio.addEventListener("change", function () {
      if (this.value === "username") {
        usernameBox.style.display = "block";
        emailBox.style.display = "none";
        usernameInput.focus();
        // Clear email input and errors
        emailInput.value = "";
        emailInput.classList.remove("input-error");
        const emailError = emailBox.querySelector(".error-message");
        if (emailError) emailError.textContent = "";
      } else if (this.value === "email") {
        usernameBox.style.display = "none";
        emailBox.style.display = "block";
        emailInput.focus();
        // Clear username input and errors
        usernameInput.value = "";
        usernameInput.classList.remove("input-error");
        const usernameError = usernameBox.querySelector(".error-message");
        if (usernameError) usernameError.textContent = "";
      }
    });
  });
});

/*--- Reset button removes preserved values ---*/
document.addEventListener("DOMContentLoaded", function () {
  const resetButtons = document.querySelectorAll(".btn-reset-form");

  resetButtons.forEach((button) => {
    button.addEventListener("click", function (e) {
      e.preventDefault();

      const form = this.closest("form");

      // Reset all form fields with the form-input class
      const inputs = form.querySelectorAll(".form-input");
      inputs.forEach((input) => {
        input.value = "";
        input.classList.remove("input-error");
      });

      // Clear error messages
      const errorMessages = form.querySelectorAll(".error-message");
      errorMessages.forEach((msg) => (msg.textContent = ""));

      // Reset login type to Username by default
      const usernameRadio = form.querySelector("#loginTypeUsername");
      const emailBox = form.querySelector("#emailBox");
      const usernameBox = form.querySelector("#usernameBox");

      if (usernameRadio) usernameRadio.checked = true;

      // Show username box, hide email box
      if (usernameBox && emailBox) {
        usernameBox.style.display = "block";
        emailBox.style.display = "none";
      }

      // Focus username
      const usernameInput = form.querySelector("#username");
      if (usernameInput) usernameInput.focus();
    });
  });
});
