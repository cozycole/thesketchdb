export function initViewSeason() {
  document.getElementById('seasonDropdown').addEventListener('change', function() {
    let selectedValue = this.value;
    let dropdown = document.getElementById('seasonDropdown');
    console.log(dropdown);
    if (selectedValue) {
      console.log(dropdown.dataset.url + `/${this.value}`);
      window.location.href = dropdown.dataset.url + `/${this.value}`;
    }
  });
}
