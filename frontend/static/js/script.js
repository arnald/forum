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
