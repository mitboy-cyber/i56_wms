// I56 Framework — App Scripts
document.addEventListener('DOMContentLoaded', function() {
  // Auto-dismiss alerts
  setTimeout(function() {
    document.querySelectorAll('.alert').forEach(function(el) {
      el.style.transition = 'opacity 0.5s';
      el.style.opacity = '0';
      setTimeout(function() { el.remove(); }, 500);
    });
  }, 3000);
});

// HTMX response target for API results
document.addEventListener('htmx:afterRequest', function(evt) {
  var target = document.getElementById('api-result');
  if (target && evt.detail.target === target) {
    var resp = evt.detail.xhr.responseText;
    target.innerHTML = '<div class="alert alert-info small p-2 m-0">' +
      '<button class="btn-close btn-sm float-end" onclick="this.parentElement.remove()"></button>' +
      '<pre class="mb-0" style="max-height:200px;overflow:auto">' + resp + '</pre></div>';
  }
});
