function msgbox(elementId, title, message, options = {}) {
  var elm = document.getElementById(elementId);
  if(elm) {
    var child = elm.lastElementChild;
    while(child) {
      elm.removeChild(child);
      child = elm.lastElementChild;
    }

    var onclick = options.onclick ? options.onclick : function() {};
    var buttons = options.buttons ? options.buttons : ["OK", "Cancel"];
    var d_class = options.class ? " " + options.class : "";
    var htmlBody = [
      '<div class="modal-dialog' + d_class + '">', 
      '  <div class="modal-content">', 
      '    <div class="modal-header">', 
      '      <h5 class="modal-title">' + title + '</h5>', 
      '      <button type="button" class="close" data-dismiss="modal" aria-label="Close">', 
      '        <span aria-hidden="true">&times;</span>', 
      '      </button>', 
      '    </div>', 
      '    <div class="modal-body">', 
      '      <div>' + message + '</div>', 
      '    </div>', 
      '    <div id="msgbox-footer" class="modal-footer"></div>', 
      '  </div>', 
      '</div>', 
    ];

    elm.className = "modal";
    elm.setAttribute("tabindex", "-1");
    elm.setAttribute("style", "z-index: 2000;");
    elm.innerHTML = htmlBody.join("");

    buttons.forEach(function(b) {
      var btn = document.createElement("button");
      btn.className = "btn btn-sm btn-w-normal btn-primary";
      btn.setAttribute("type", "button");
      btn.setAttribute("data-dismiss", "modal");
      btn.innerHTML = b;
      btn.addEventListener('click', function(event) {onclick(b);});
      document.getElementById("msgbox-footer").appendChild(btn);
    });

    $('#' + elementId).modal('show');
  }
}
