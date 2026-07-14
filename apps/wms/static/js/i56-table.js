/**
 * I56 Table Component — BFT56 Interaction Patterns
 * All features work client-side without Go backend changes.
 * Extracted from templates/generic_list.html (Phase 0 consolidation).
 */
(function() {
  // Read page identifier from URL path: /admin/<page> or /admin/<group>/<page>
  var fullPath = window.location.pathname.replace(/^\/+|\/+$/g, '').replace(/^admin\//, '');
  var pathParts = window.location.pathname.replace(/^\/+|\/+$/g, '').split('/');
  var page = pathParts[pathParts.length - 1] || '';
  // For nested paths like /admin/system/api-couriers, use the full sub-path
  if (pathParts.length > 3) {
    page = pathParts.slice(2).join('/');
  }

  var STATE = {
    page: page,
    sortCol: -1,
    sortDir: 'asc',
    filterActive: false
  };

  window.I56Table = {
    /* ── Add Form ── */
    openAddForm: function(url) {
      var modal = document.getElementById('modal-container');
      if (!modal) return;
      modal.style.display = 'flex';
      modal.innerHTML = '<div class="modal-overlay"><div class="modal-content"><div class="modal-body" style="text-align:center;padding:40px"><p>⚡ 加载中...</p></div></div></div>';
      fetch(url).then(function(r) { if (!r.ok) throw new Error('HTTP ' + r.status); return r.text(); }).then(function(html) {
        modal.innerHTML = html;
        // Add submit handler validation
        setTimeout(function() {
          var form = modal.querySelector('form');
          if (form) {
            form.addEventListener('submit', function(e) {
              var required = form.querySelectorAll('[required]');
              var valid = true;
              required.forEach(function(field) {
                if (!field.value.trim()) {
                  field.style.borderColor = 'red';
                  valid = false;
                } else {
                  field.style.borderColor = '';
                }
              });
              if (!valid) {
                e.preventDefault();
                alert('请填写所有必填字段');
              }
            });
          }
        }, 100);
      }).catch(function(e) {
        modal.innerHTML = '<div class="modal-overlay" onclick="closeI56Modal()"><div class="modal-content"><div class="modal-header"><span class="modal-title">错误</span><button class="modal-close" onclick="this.closest(\'div.modal-overlay\').remove()">&times;</button></div><div class="modal-body"><p>加载添加表单失败: ' + e.message + '</p></div></div></div>';
      });
    },
    /* ── Search / Filter ── */
    filter: function() {
      var q = document.getElementById('table-search').value.toLowerCase().trim();
      var rows = document.querySelectorAll('#data-table tbody tr');
      var visible = 0;
      rows.forEach(function(tr) {
        var checkTd = tr.querySelector('td.clickable-cell');
        if (!checkTd) return;
        var text = tr.textContent.toLowerCase();
        var match = !q || text.indexOf(q) !== -1;
        tr.style.display = match ? '' : 'none';
        if (match) visible++;
      });
      var empty = document.getElementById('filter-empty');
      if (visible === 0 && q) {
        if (!empty) {
          empty = document.createElement('tr');
          empty.id = 'filter-empty';
          var colSpan = document.querySelectorAll('#data-table thead th').length;
          empty.innerHTML = '<td colspan="' + colSpan + '" class="empty-state"><div class="empty-state-icon">🔍</div><div class="empty-state-text">未找到匹配结果</div></td>';
          document.querySelector('#data-table tbody').appendChild(empty);
        }
        empty.style.display = '';
      } else if (empty) {
        empty.style.display = 'none';
      }
    },

    /* ── Sort ── */
    sort: function(colIdx) {
      var dir = STATE.sortCol === colIdx ? (STATE.sortDir === 'asc' ? 'desc' : 'asc') : 'asc';
      STATE.sortCol = colIdx;
      STATE.sortDir = dir;

      // Update sort indicators on all headers
      document.querySelectorAll('#data-table th.sortable').forEach(function(th, i) {
        var arrows = th.querySelector('.sort-arrows');
        if (arrows) {
          arrows.className = 'sort-arrows' + (i === colIdx ? ' active' : '');
          arrows.setAttribute('data-dir', i === colIdx ? dir : '');
        }
      });

      // Sort rows
      var tbody = document.querySelector('#data-table tbody');
      var rows = Array.from(tbody.querySelectorAll('tr:not(#filter-empty)'));
      rows.sort(function(a, b) {
        var aCell = a.querySelectorAll('td.clickable-cell')[colIdx];
        var bCell = b.querySelectorAll('td.clickable-cell')[colIdx];
        if (!aCell || !bCell) return 0;
        var aVal = aCell.textContent.trim();
        var bVal = bCell.textContent.trim();
        // Try numeric sort
        var aNum = parseFloat(aVal.replace(/[^0-9.-]/g, ''));
        var bNum = parseFloat(bVal.replace(/[^0-9.-]/g, ''));
        if (!isNaN(aNum) && !isNaN(bNum)) {
          return dir === 'asc' ? aNum - bNum : bNum - aNum;
        }
        return dir === 'asc' ? aVal.localeCompare(bVal, 'zh-CN') : bVal.localeCompare(bVal, 'zh-CN');
      });
      rows.forEach(function(r) { tbody.appendChild(r); });
    },

    /* ── Navigate to Detail ── */
    navigateToDetail: function(cell) {
      var cellText = cell.textContent.trim();
      // Navigate to /admin/<page>/<id>
      var base = '/admin/' + STATE.page;
      window.location.href = base + '?detail=' + encodeURIComponent(cellText);
    },

    /* ── Row Actions ── */
    editRow: function(btn) {
      var id = btn.getAttribute('data-id');
      var base = '/admin/' + STATE.page;
      var url = base + '/edit-form?id=' + encodeURIComponent(id);
      var modal = document.getElementById('modal-container');
      modal.style.display = 'flex';
      modal.innerHTML = '<div class="modal-overlay" onclick="closeI56Modal()"><div class="modal-content"><div class="modal-body" style="text-align:center;padding:40px"><div style="font-size:20px">⚡</div><p>加载中...</p></div></div></div>';
      fetch(url).then(function(r) { if (!r.ok) throw new Error('HTTP ' + r.status); return r.text(); }).then(function(html) {
        modal.innerHTML = html;
      }).catch(function(e) {
        modal.innerHTML = '<div class="modal-overlay" onclick="closeI56Modal()"><div class="modal-content"><div class="modal-header"><span class="modal-title">错误</span><button class="modal-close" onclick="this.closest(\'div.modal-overlay\').remove()">&times;</button></div><div class="modal-body"><p>加载编辑表单失败: ' + e.message + '</p></div></div></div>';
      });
    },
    deleteRow: function(btn) {
      var id = btn.getAttribute('data-id');
      if (confirm('确认删除 "' + id + '" ？此操作不可撤销。')) {
        fetch('/admin/delete?page=' + STATE.page + '&id=' + encodeURIComponent(id), { method: 'POST' })
          .then(function(r) {
            if (r.ok) window.location.reload();
            else alert('删除失败');
          })
          .catch(function() { alert('删除失败'); });
      }
    },

    /* ── Batch Operations ── */
    toggleSelectAll: function(cb) {
      var checked = cb.checked;
      document.querySelectorAll('.row-checkbox').forEach(function(c) { c.checked = checked; });
      I56Table.updateBatchBar();
    },
    updateBatchBar: function() {
      var checked = document.querySelectorAll('.row-checkbox:checked');
      var count = checked.length;
      document.getElementById('batch-count').textContent = count;
      document.getElementById('batch-bar').style.display = count > 0 ? '' : 'none';
      // Update select-all state
      var allVisible = document.querySelectorAll('#data-table tbody tr:not(#filter-empty) .row-checkbox');
      var allChecked = document.querySelectorAll('#data-table tbody tr:not(#filter-empty) .row-checkbox:checked');
      var selAll = document.getElementById('select-all');
      selAll.checked = allChecked.length === allVisible.length && allVisible.length > 0;
      selAll.indeterminate = allChecked.length > 0 && allChecked.length < allVisible.length;
    },
    batchSelectAll: function() {
      document.querySelectorAll('.row-checkbox').forEach(function(c) { c.checked = true; });
      document.getElementById('select-all').checked = true;
      document.getElementById('select-all').indeterminate = false;
      I56Table.updateBatchBar();
    },
    batchClear: function() {
      document.querySelectorAll('.row-checkbox').forEach(function(c) { c.checked = false; });
      document.getElementById('select-all').checked = false;
      document.getElementById('select-all').indeterminate = false;
      I56Table.updateBatchBar();
    },
    batchDelete: function() {
      var checked = document.querySelectorAll('.row-checkbox:checked');
      var ids = [];
      checked.forEach(function(cb) {
        var row = cb.closest('tr');
        var firstCell = row.querySelector('td.clickable-cell');
        if (firstCell) ids.push(firstCell.textContent.trim());
      });
      if (ids.length === 0) return;
      if (!confirm('确认删除 ' + ids.length + ' 条记录？此操作不可撤销。')) return;

      Promise.all(ids.map(function(id) {
        return fetch('/admin/' + STATE.page + '?delete=' + encodeURIComponent(id), { method: 'POST' });
      }))
        .then(function() { window.location.reload(); })
        .catch(function() { alert('部分删除失败'); });
    },

    /* ── Export CSV ── */
    exportCSV: function() {
      var rows = [];
      // Header
      var headers = ['✔'];
      document.querySelectorAll('#data-table thead th.sortable .sort-label').forEach(function(h) {
        headers.push(h.textContent.trim());
      });
      rows.push(headers.join(','));

      // Data rows (visible only)
      document.querySelectorAll('#data-table tbody tr').forEach(function(tr) {
        if (tr.style.display === 'none' || tr.id === 'filter-empty') return;
        var cells = [];
        tr.querySelectorAll('td.clickable-cell').forEach(function(td) {
          cells.push('"' + td.textContent.replace(/"/g, '""').trim() + '"');
        });
        if (cells.length > 0) rows.push(cells.join(','));
      });

      var blob = new Blob(['\uFEFF' + rows.join('\n')], { type: 'text/csv;charset=utf-8' });
      var url = URL.createObjectURL(blob);
      var a = document.createElement('a');
      a.href = url;
      a.download = STATE.page + '_export.csv';
      a.click();
      URL.revokeObjectURL(url);
    },

    /* ── Filter Toggle ── */
    toggleFilter: function() {
      STATE.filterActive = !STATE.filterActive;
      var badge = document.getElementById('filter-count');
      var btn = document.getElementById('btn-filter');
      if (STATE.filterActive) {
        btn.classList.add('active');
        badge.style.display = '';
        badge.textContent = '1';
      } else {
        btn.classList.remove('active');
        badge.style.display = 'none';
      }
    },

    /* ── Export Dropdown ── */
    toggleExportMenu: function(e) {
      e.stopPropagation();
      var menu = document.getElementById('export-menu');
      menu.style.display = menu.style.display === 'none' ? '' : 'none';
      // Close on outside click
      if (menu.style.display !== 'none') {
        setTimeout(function() {
          document.addEventListener('click', function handler(ev) {
            document.getElementById('export-menu').style.display = 'none';
            document.removeEventListener('click', handler);
          });
        }, 0);
      }
    },

    /* ── Export Declaration (PDF template) ── */
    exportDeclaration: function() {
      document.getElementById('export-menu').style.display = 'none';
      // Gather visible rows
      var headers = [];
      document.querySelectorAll('#data-table thead th.sortable .sort-label').forEach(function(h) {
        headers.push(h.textContent.trim());
      });
      var rows = [];
      document.querySelectorAll('#data-table tbody tr').forEach(function(tr) {
        if (tr.style.display === 'none' || tr.id === 'filter-empty') return;
        var cells = [];
        tr.querySelectorAll('td.clickable-cell').forEach(function(td) {
          cells.push(td.textContent.trim());
        });
        if (cells.length > 0) rows.push(cells);
      });

      // Generate a simple HTML declaration template and trigger print
      var html = '<!DOCTYPE html><html lang="zh-CN"><head><meta charset="UTF-8"><title>申报单导出 - I56</title>' +
        '<style>body{font-family:sans-serif;padding:20px;color:#222}' +
        '.hd{text-align:center;margin-bottom:20px}.hd h2{font-size:18px;margin:0 0 4px}.hd span{font-size:11px;color:#666}' +
        'table{width:100%;border-collapse:collapse;font-size:10px}th,td{border:1px solid #333;padding:4px 6px;text-align:left}th{background:#eee}' +
        '.ft{font-size:10px;color:#666;margin-top:20px;text-align:center}' +
        '</style></head><body>' +
        '<div class="hd"><h2>📋 申报单</h2><span>导出时间: ' + new Date().toLocaleString('zh-CN') + ' | 页面: ' + STATE.page + '</span></div>' +
        '<table><thead><tr>' + headers.map(function(h){return '<th>'+h+'</th>'}).join('') + '</tr></thead><tbody>' +
        rows.map(function(r){return '<tr>'+r.map(function(c){return '<td>'+c+'</td>'}).join('')+'</tr>'}).join('') +
        '</tbody></table>' +
        '<div class="ft">I56 WMS 系统 · ' + new Date().toISOString().split('T')[0] + '</div>' +
        '</body></html>';

      var w = window.open('', '_blank', 'width=900,height=700');
      w.document.write(html);
      w.document.close();
      setTimeout(function() { w.print(); }, 300);
    },

    /* ── Order Status Transition ── */
    transitionOrder: function(orderNo, newStatus, btn) {
      if (!confirm('确认将订单 ' + orderNo + ' 状态更新为 ' + newStatus + '？')) return;
      var origText = btn.textContent;
      btn.textContent = '...';
      btn.disabled = true;
      fetch('/admin/orders/' + encodeURIComponent(orderNo) + '/status', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ status: newStatus })
      }).then(function(r) {
        return r.json().then(function(data) { return { ok: r.ok, data: data }; });
      }).then(function(result) {
        if (result.ok) {
          // Reload page to show updated status and new transition buttons
          window.location.reload();
        } else {
          alert('状态更新失败: ' + (result.data.error || result.data.message || '未知错误'));
          btn.textContent = origText;
          btn.disabled = false;
        }
      }).catch(function() {
        alert('网络错误，请稍后重试');
        btn.textContent = origText;
        btn.disabled = false;
      });
    }
  };
})();
