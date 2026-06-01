export function showToast(message) {
  let container = document.querySelector("#toast-container");

  if (!container) {
    container = document.createElement("div");
    container.id = "toast-container";
    container.className =
      "fixed top-4 right-4 z-50 flex flex-col gap-2 pointer-events-none";
    document.body.appendChild(container);
  }

  const toast = document.createElement("div");
  toast.textContent = message;
  toast.className =
    "pointer-events-auto rounded-lg bg-orange-500 px-4 py-2 text-white shadow-lg " +
    "translate-x-full opacity-0 transition-all duration-300 ease-out";

  container.appendChild(toast);

  requestAnimationFrame(() => {
    toast.classList.remove("translate-x-full", "opacity-0");
    toast.classList.add("translate-x-0", "opacity-100");
  });

  setTimeout(() => {
    toast.classList.remove("translate-x-0", "opacity-100");
    toast.classList.add("translate-x-full", "opacity-0");

    setTimeout(() => toast.remove(), 300);
  }, 3000);
}
