// scanner.js — 真实摄像头扫码组件 (BarcodeDetector API + QuaggaJS fallback)
// 用法: <input type="text" class="scan-input" data-scanner="true">
//        <button class="scan-camera-btn">📷</button>
(function(){
  'use strict';

  // 检测 BarcodeDetector 支持
  const hasNativeBarcode = 'BarcodeDetector' in window;
  const detector = hasNativeBarcode ? new BarcodeDetector({formats:['code_128','ean_13','ean_8','code_39','upc_a','itf','qr_code']}) : null;

  let stream = null, scanning = false, videoEl = null, overlayEl = null, currentInput = null;

  function createOverlay(){
    if(overlayEl) return overlayEl;
    overlayEl = document.createElement('div');
    overlayEl.id = 'scanner-overlay';
    overlayEl.innerHTML = `<div style="position:fixed;top:0;left:0;width:100%;height:100%;background:rgba(0,0,0,0.9);z-index:9999;display:flex;flex-direction:column;align-items:center;justify-content:center">
      <video id="scanner-video" style="max-width:100%;max-height:70vh;border:3px solid #38bdf8;border-radius:12px" autoplay playsinline></video>
      <div style="margin-top:16px;color:#fff;text-align:center">
        <div id="scanner-status" style="font-size:16px">🔍 扫描中...</div>
        <div id="scanner-result" style="font-size:14px;color:#38bdf8;margin-top:8px"></div>
      </div>
      <button id="scanner-close" style="margin-top:16px;background:#ef4444;color:#fff;border:none;border-radius:8px;padding:10px 32px;font-size:16px;cursor:pointer">关闭摄像头</button>
    </div>`;
    document.body.appendChild(overlayEl);
    videoEl = overlayEl.querySelector('#scanner-video');
    overlayEl.querySelector('#scanner-close').onclick = stopScan;
    return overlayEl;
  }

  async function startScan(inputEl){
    if(scanning) return;
    currentInput = inputEl;
    createOverlay();
    overlayEl.style.display = 'block';
    scanning = true;

    try {
      stream = await navigator.mediaDevices.getUserMedia({video:{facingMode:'environment',width:{ideal:1280},height:{ideal:720}}});
      videoEl.srcObject = stream;
      videoEl.play();
      document.getElementById('scanner-status').textContent = '🔍 扫描中...';
      scanLoop();
    } catch(e){
      document.getElementById('scanner-status').textContent = '❌ 摄像头不可用: ' + e.message;
      scanning = false;
    }
  }

  function scanLoop(){
    if(!scanning || !videoEl) return;
    const w = videoEl.videoWidth, h = videoEl.videoHeight;
    if(w === 0 || h === 0){requestAnimationFrame(scanLoop); return;}

    const canvas = document.createElement('canvas');
    canvas.width = w; canvas.height = h;
    canvas.getContext('2d').drawImage(videoEl, 0, 0, w, h);

    if(detector){
      detector.detect(canvas).then(barcodes => {
        if(barcodes.length > 0){
          const val = barcodes[0].rawValue;
          onBarcode(val);
          return;
        }
        requestAnimationFrame(scanLoop);
      }).catch(()=>requestAnimationFrame(scanLoop));
    } else {
      // Fallback: use Quagga2 if loaded
      if(typeof Quagga !== 'undefined'){
        Quagga.decodeSingle({src: canvas.toDataURL(), numOfWorkers:0, decoder:{readers:['code_128_reader','ean_reader','ean_8_reader','code_39_reader','upc_reader']}},
          function(result){
            if(result && result.codeResult){
              onBarcode(result.codeResult.code);
            } else {
              requestAnimationFrame(scanLoop);
            }
          });
      } else {
        document.getElementById('scanner-status').textContent = '⚠️ 浏览器不支持扫码，请手动输入';
        scanning = false;
      }
    }
  }

  function onBarcode(val){
    if(!scanning || !currentInput) return;
    document.getElementById('scanner-result').textContent = '✅ 识别: ' + val;
    currentInput.value = val;
    currentInput.dispatchEvent(new Event('input',{bubbles:true}));
    currentInput.dispatchEvent(new Event('change',{bubbles:true}));
    // Auto-trigger HTMX or form submit
    if(currentInput.hasAttribute('hx-post') || currentInput.hasAttribute('hx-get')){
      currentInput.dispatchEvent(new Event('keydown',{key:'Enter',bubbles:true}));
    }
    setTimeout(stopScan, 800);
  }

  function stopScan(){
    scanning = false;
    if(stream){stream.getTracks().forEach(t=>t.stop()); stream = null;}
    if(videoEl) videoEl.srcObject = null;
    if(overlayEl) overlayEl.style.display = 'none';
    currentInput = null;
  }

  // Init: attach camera buttons to all .scan-input fields
  document.addEventListener('DOMContentLoaded', function(){
    document.querySelectorAll('.scan-input').forEach(inp => {
      // Only add camera button if not already added
      const parent = inp.parentElement;
      if(parent && !parent.querySelector('.scan-camera-btn')){
        const btn = document.createElement('button');
        btn.type = 'button';
        btn.className = 'scan-camera-btn btn btn-outline-info btn-sm';
        btn.innerHTML = '📷';
        btn.title = '打开摄像头扫码';
        btn.style.cssText = 'position:absolute;right:4px;top:50%;transform:translateY(-50%);padding:2px 8px;z-index:5';
        inp.style.paddingRight = '40px';
        inp.parentElement.style.position = 'relative';
        btn.onclick = function(e){e.preventDefault(); startScan(inp);};
        inp.parentElement.appendChild(btn);
      }
    });
  });

  window.scannerStart = startScan;
  window.scannerStop = stopScan;
})();
