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

  toggleCheckbox.addEventListener("change", function () {
    if (toggleCheckbox.checked) {
      passwordInput.type = "text";
    } else {
      passwordInput.type = "password";
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
