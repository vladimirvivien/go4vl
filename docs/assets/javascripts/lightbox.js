(function () {
  "use strict";

  var EXCLUDED_LOGOS = [
    "go4vl-logo-light.png",
    "go4vl-logo-dark.png",
    "go4vl-mark.png",
    "favicon.png",
    "go4vl-logo-small.png"
  ];

  function isExcludedImage(img) {
    if (!img || img.classList.contains("tx-hero__logo")) {
      return true;
    }
    var src = img.getAttribute("src") || "";
    for (var i = 0; i < EXCLUDED_LOGOS.length; i++) {
      if (src.indexOf(EXCLUDED_LOGOS[i]) !== -1) {
        return true;
      }
    }
    return false;
  }

  function ensureLightbox() {
    var existing = document.getElementById("go4vl-lightbox");
    if (existing) {
      return existing;
    }

    var overlay = document.createElement("div");
    overlay.id = "go4vl-lightbox";
    overlay.className = "go4vl-lightbox";
    overlay.setAttribute("role", "dialog");
    overlay.setAttribute("aria-modal", "true");
    overlay.setAttribute("aria-hidden", "true");
    overlay.setAttribute("aria-label", "Enlarged diagram");

    overlay.innerHTML =
      '<div class="go4vl-lightbox__backdrop"></div>' +
      '<div class="go4vl-lightbox__dialog">' +
        '<button type="button" class="go4vl-lightbox__close" aria-label="Close enlarged diagram" title="Close (Esc)">&times;</button>' +
        '<div class="go4vl-lightbox__stage">' +
          '<img class="go4vl-lightbox__img" src="" alt="">' +
        '</div>' +
        '<div class="go4vl-lightbox__bar">' +
          '<span class="go4vl-lightbox__caption"></span>' +
          '<a class="go4vl-lightbox__raw" href="" target="_blank" rel="noopener">Open full size &#8599;</a>' +
        '</div>' +
      '</div>';

    document.body.appendChild(overlay);
    return overlay;
  }

  var lastTrigger = null;

  function openLightbox(img) {
    var overlay = ensureLightbox();
    var modalImg = overlay.querySelector(".go4vl-lightbox__img");
    var captionEl = overlay.querySelector(".go4vl-lightbox__caption");
    var rawLink = overlay.querySelector(".go4vl-lightbox__raw");
    var closeBtn = overlay.querySelector(".go4vl-lightbox__close");

    var src = img.currentSrc || img.src;
    var alt = img.getAttribute("alt") || "Architecture diagram";

    lastTrigger = img;
    modalImg.src = src;
    modalImg.alt = alt;
    captionEl.textContent = alt;
    rawLink.href = src;

    overlay.classList.add("is-open");
    overlay.setAttribute("aria-hidden", "false");
    document.documentElement.style.overflow = "hidden";
    document.body.style.overflow = "hidden";

    if (closeBtn) {
      closeBtn.focus({ preventScroll: true });
    }
  }

  function closeLightbox() {
    var overlay = document.getElementById("go4vl-lightbox");
    if (!overlay || !overlay.classList.contains("is-open")) {
      return;
    }
    overlay.classList.remove("is-open");
    overlay.setAttribute("aria-hidden", "true");
    document.documentElement.style.overflow = "";
    document.body.style.overflow = "";

    if (lastTrigger && typeof lastTrigger.focus === "function") {
      lastTrigger.focus({ preventScroll: true });
    }
    lastTrigger = null;
  }

  function bindGlobalListeners() {
    if (window.__go4vlLightboxGlobalBound) {
      return;
    }
    window.__go4vlLightboxGlobalBound = true;

    document.addEventListener("keydown", function (e) {
      if (e.key === "Escape") {
        closeLightbox();
      }
    });

    document.addEventListener("click", function (e) {
      var overlay = document.getElementById("go4vl-lightbox");
      if (!overlay || !overlay.classList.contains("is-open")) {
        return;
      }
      if (
        e.target.classList.contains("go4vl-lightbox__backdrop") ||
        e.target.classList.contains("go4vl-lightbox__close") ||
        e.target.classList.contains("go4vl-lightbox__img") ||
        e.target === overlay
      ) {
        closeLightbox();
      }
    });
  }

  function attachDiagrams() {
    bindGlobalListeners();
    ensureLightbox();

    var candidates = document.querySelectorAll('.diagram-frame img, .md-content img[src*="assets/images/"]');
    for (var i = 0; i < candidates.length; i++) {
      var img = candidates[i];
      if (isExcludedImage(img) || img.dataset.lightboxBound === "1") {
        continue;
      }
      img.dataset.lightboxBound = "1";
      img.classList.add("go4vl-zoomable");
      if (!img.hasAttribute("tabindex")) {
        img.setAttribute("tabindex", "0");
      }
      img.setAttribute("role", "button");
      img.setAttribute("title", "Click to enlarge");
      var alt = img.getAttribute("alt") || "Diagram";
      img.setAttribute("aria-label", "Enlarge diagram: " + alt);

      (function (targetImg) {
        targetImg.addEventListener("click", function (e) {
          e.preventDefault();
          openLightbox(targetImg);
        });
        targetImg.addEventListener("keydown", function (e) {
          if (e.key === "Enter" || e.key === " ") {
            e.preventDefault();
            openLightbox(targetImg);
          }
        });
      })(img);
    }
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", attachDiagrams);
  } else {
    attachDiagrams();
  }

  if (typeof window.document$ !== "undefined" && typeof window.document$.subscribe === "function") {
    window.document$.subscribe(function () {
      attachDiagrams();
    });
  }
})();
