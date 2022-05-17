// For showing/hiding sidebar on medium and small screens
(function () {
  let hideSidebar = document.getElementById("close-sidebar"),
    sidebar = document.getElementById("sidebar");
    shade = document.getElementById("shade")

  moveToMain("pageTitle")
  moveToMain("ghAccount")

  function moveToMain(id) {
    let dom = document.getElementById(id)
    let masterDom = document.getElementById(id + "Placeholder")
    if(dom != null && masterDom != null){
      masterDom.innerHTML = dom.innerHTML
      dom.remove()
    }
  }

  function toggleSidebar() {
    let classesToApplyForSidebar = {
      active: [],
      inactive: ["transform", "-translate-x-full"],
    },
    classesToApplyForShade = {
      active: [],
      inactive: ["hidden"]
    };

    let isMenuActive = hideSidebar.getAttribute("data-menu-active") === "true";

    if (isMenuActive) {
      hideSidebar.setAttribute("data-menu-active", null);

      sidebar.classList.remove(...classesToApplyForSidebar.active);
      sidebar.classList.add(...classesToApplyForSidebar.inactive);

      shade.classList.remove(...classesToApplyForShade.active);
      shade.classList.add(...classesToApplyForShade.inactive)
    } else {
      hideSidebar.setAttribute("data-menu-active", true);

      sidebar.classList.add(...classesToApplyForSidebar.active);
      sidebar.classList.remove(...classesToApplyForSidebar.inactive);

      shade.classList.add(...classesToApplyForShade.active);
      shade.classList.remove(...classesToApplyForShade.inactive)
    }
  }

  document.getElementById("close-sidebar").addEventListener("click", (e) => {
    e.preventdefault;
    toggleSidebar();
  });

  document.getElementById("show-sidebar").addEventListener("click", (e) => {
    e.preventdefault;
    toggleSidebar();
  });

  toggleSidebar();
  setTimeout(() => {
    sidebar.classList.add("duration-300");
  }, 300);
})();
