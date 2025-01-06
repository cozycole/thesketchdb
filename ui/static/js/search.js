document.body.addEventListener('htmx:configRequest', function (evt) {
    //console.log(evt.detail.headers);
    if (evt.detail.elt.classList.contains('htmxSearchPage')) {
        const searchType = document.querySelector('input[name="type"]:checked').value;
        evt.detail.parameters['type'] = searchType;
    }
});

document.addEventListener('DOMContentLoaded', () => {
  const toggles = document.getElementsByClassName('filterMenuToggle');
  const menu = document.getElementById('filterMenu');

  Array.from(toggles).forEach((e) => {
        e.addEventListener('click', () => {
            menu.classList.toggle('hidden');
        }
        )}
  );

  const typeDropToggleUp = document.getElementById('typeDropUp');
  const typeDropToggleDown = document.getElementById('typeDropDown');
  const typeMenu = document.getElementById('typeMenu');

  typeDropToggleUp.addEventListener('click', () => {
        typeDropToggleUp.classList.toggle('hidden');
        typeDropToggleDown.classList.toggle('hidden');
        typeMenu.classList.toggle('hidden');
  });

  typeDropToggleDown.addEventListener('click', () => {
        typeDropToggleUp.classList.toggle('hidden');
        typeDropToggleDown.classList.toggle('hidden');
        typeMenu.classList.toggle('hidden');
  });


  // Optional: close the dropdown when clicking outside of it
  document.addEventListener('click', (event) => {
    if (!document.querySelector('button.filterMenuToggle').contains(event.target) && !menu.contains(event.target)) {
      menu.classList.add('hidden');
    }
  });
});
