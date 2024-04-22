function paging(elm, currentPage, pageCount, buttonCount, clickHandler, options = {}) {
  if(elm) {
    var child = elm.lastElementChild;
    while(child) {
      elm.removeChild(child);
      child = elm.lastElementChild;
    }

    var btnCount = Math.min(pageCount, buttonCount);
    var pageFrom = Math.min(Math.max(currentPage - Math.floor(btnCount / 2), 1), pageCount - (btnCount - 1));
    var enbPrev = currentPage > 1;
    var enbNext = currentPage < pageCount;
    var clsPrev = enbPrev ? "" : " pg-disabled";
    var clsNext = enbNext ? "" : " pg-disabled";
    var btnHome = document.createElement("span");
    var btnPrev = document.createElement("span");
    var btnNext = document.createElement("span");
    var btnLast = document.createElement("span");

    btnHome.innerHTML = "first" in options ? options.first : "First";
    btnHome.className = `pg-btn${clsPrev}`;

    btnPrev.innerHTML = "prev" in options ? options.prev : "Prev";
    btnPrev.className = `pg-btn${clsPrev}`;

    btnNext.innerHTML = "next" in options ? options.next : "Next";
    btnNext.className = `pg-btn${clsNext}`;

    btnLast.innerHTML = "last" in options ? options.last : "Last";
    btnLast.className = `pg-btn${clsNext}`;

    if(enbPrev && (pageCount > 0)) {
      btnHome.addEventListener('click', function(event) {if(clickHandler != null) {clickHandler(1);}});
      btnPrev.addEventListener('click', function(event) {if(clickHandler != null) {clickHandler(currentPage - 1);}});
    }

    if(enbNext && (pageCount > 0)) {
      btnNext.addEventListener('click', function(event) {if(clickHandler != null) {clickHandler(currentPage + 1);}});
      btnLast.addEventListener('click', function(event) {if(clickHandler != null) {clickHandler(pageCount);}});
    }

    elm.appendChild(btnHome);
    elm.appendChild(btnPrev);

    for(var i = 0; i < btnCount; i++) {
      let pageNum = pageFrom + i;
      var btnPage = document.createElement("span");
      btnPage.innerHTML = pageNum;
      btnPage.className = `pg-btn${pageNum == currentPage ? " pg-active" : ""}`;
      if(pageNum != currentPage) {
        btnPage.addEventListener('click', function(event) {
          if(clickHandler != null) {
            clickHandler(pageNum);
          }
        });
      }
      elm.appendChild(btnPage);
    }

    elm.appendChild(btnNext);
    elm.appendChild(btnLast);
  }
}
