package web

const indexHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>namer - 姓名阴阳五行评分</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI","PingFang SC","Hiragino Sans GB","Microsoft YaHei",sans-serif;background:#f0f2f5;color:#333;min-height:100vh}
.header{background:linear-gradient(135deg,#1a1a2e 0%,#16213e 50%,#0f3460 100%);color:#fff;padding:24px 0;text-align:center}
.header h1{font-size:28px;font-weight:700;letter-spacing:2px}
.header p{margin-top:6px;font-size:14px;color:rgba(255,255,255,.7)}
.container{max-width:900px;margin:0 auto;padding:20px}
.tabs{display:flex;gap:0;margin-bottom:20px;background:#fff;border-radius:10px;overflow:hidden;box-shadow:0 2px 8px rgba(0,0,0,.06)}
.tab{flex:1;padding:14px;text-align:center;cursor:pointer;font-size:15px;font-weight:500;border:none;background:none;color:#666;transition:all .2s}
.tab:hover{color:#333;background:#f8f9fa}
.tab.active{color:#fff;background:linear-gradient(135deg,#e94560,#c23152)}
.card{background:#fff;border-radius:10px;padding:24px;box-shadow:0 2px 8px rgba(0,0,0,.06);margin-bottom:20px}
.card h2{font-size:18px;font-weight:600;margin-bottom:16px;padding-bottom:10px;border-bottom:2px solid #f0f2f5}
.form-row{display:flex;gap:12px;margin-bottom:14px;flex-wrap:wrap}
.form-group{display:flex;flex-direction:column;flex:1;min-width:120px}
.form-group label{font-size:13px;color:#888;margin-bottom:4px;font-weight:500}
.form-group input{padding:10px 12px;border:1px solid #e0e0e0;border-radius:6px;font-size:15px;transition:border-color .2s;outline:none}
.form-group input:focus{border-color:#e94560}
.btn{padding:12px 32px;border:none;border-radius:8px;font-size:15px;font-weight:600;cursor:pointer;transition:all .2s}
.btn-primary{background:linear-gradient(135deg,#e94560,#c23152);color:#fff}
.btn-primary:hover{transform:translateY(-1px);box-shadow:0 4px 12px rgba(233,69,96,.3)}
.btn-primary:disabled{opacity:.6;cursor:not-allowed;transform:none;box-shadow:none}
.btn-row{display:flex;justify-content:center;gap:12px;margin-top:8px}
.hidden{display:none}
.result-header{display:flex;justify-content:space-between;align-items:center;margin-bottom:16px}
.score-big{font-size:48px;font-weight:800;background:linear-gradient(135deg,#e94560,#ff6b6b);-webkit-background-clip:text;-webkit-text-fill-color:transparent}
.score-label{font-size:14px;color:#888}
.score-bar-wrap{margin-bottom:10px}
.score-bar-label{display:flex;justify-content:space-between;font-size:13px;margin-bottom:4px}
.score-bar-label span:first-child{font-weight:500}
.score-bar-label span:last-child{color:#e94560;font-weight:600}
.score-bar{height:8px;background:#f0f2f5;border-radius:4px;overflow:hidden}
.score-bar-fill{height:100%;border-radius:4px;transition:width .6s ease}
.detail-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(200px,1fr));gap:12px;margin-top:16px}
.detail-item{background:#f8f9fa;border-radius:8px;padding:12px}
.detail-item .label{font-size:12px;color:#888;margin-bottom:4px}
.detail-item .value{font-size:14px;font-weight:600}
.rank-list{width:100%}
.rank-item{display:flex;align-items:center;padding:12px 16px;border-bottom:1px solid #f0f2f5;cursor:pointer;transition:background .15s}
.rank-item:hover{background:#f8f9fa}
.rank-num{width:28px;height:28px;border-radius:50%;display:flex;align-items:center;justify-content:center;font-size:13px;font-weight:700;margin-right:14px;flex-shrink:0}
.rank-num.top1{background:#ffd700;color:#333}
.rank-num.top2{background:#c0c0c0;color:#333}
.rank-num.top3{background:#cd7f32;color:#fff}
.rank-num.other{background:#f0f2f5;color:#888}
.rank-name{flex:1;font-size:16px;font-weight:600}
.rank-score{font-size:18px;font-weight:700;color:#e94560}
.loading{text-align:center;padding:40px;color:#888}
.loading .spinner{display:inline-block;width:32px;height:32px;border:3px solid #f0f2f5;border-top-color:#e94560;border-radius:50%;animation:spin .8s linear infinite}
@keyframes spin{to{transform:rotate(360deg)}}
.empty{text-align:center;padding:40px;color:#aaa;font-size:14px}
@media(max-width:600px){.form-row{flex-direction:column}.score-big{font-size:36px}}
</style>
</head>
<body>
<div class="header">
  <h1>namer</h1>
  <p>中国姓名阴阳五行评分工具</p>
</div>
<div class="container">
  <div class="tabs">
    <button class="tab active" onclick="switchTab('score')">单名评估</button>
    <button class="tab" onclick="switchTab('batch')">批量生成</button>
  </div>

  <!-- 单名评估 -->
  <div id="tab-score">
    <div class="card">
      <h2>输入信息</h2>
      <div class="form-row">
        <div class="form-group"><label>姓氏</label><input id="s-last" placeholder="如：王" maxlength="2"></div>
        <div class="form-group"><label>名字</label><input id="s-first" placeholder="如：明轩" maxlength="4"></div>
      </div>
      <div class="form-row">
        <div class="form-group"><label>出生年</label><input id="s-year" type="number" placeholder="2024" min="1900" max="2100"></div>
        <div class="form-group"><label>月</label><input id="s-month" type="number" placeholder="3" min="1" max="12"></div>
        <div class="form-group"><label>日</label><input id="s-day" type="number" placeholder="15" min="1" max="31"></div>
        <div class="form-group"><label>时</label><input id="s-hour" type="number" placeholder="10" min="0" max="23"></div>
        <div class="form-group"><label>分</label><input id="s-min" type="number" placeholder="30" min="0" max="59"></div>
      </div>
      <div class="btn-row"><button class="btn btn-primary" id="btn-score" onclick="doScore()">评估</button></div>
    </div>
    <div id="score-result" class="hidden"></div>
  </div>

  <!-- 批量生成 -->
  <div id="tab-batch" class="hidden">
    <div class="card">
      <h2>输入信息</h2>
      <div class="form-row">
        <div class="form-group"><label>姓氏</label><input id="b-last" placeholder="如：王" maxlength="2"></div>
        <div class="form-group" style="flex:3"><label>备选字（逗号分隔）</label><input id="b-words" placeholder="如：明,轩,浩,然,宇,泽,博,文"></div>
      </div>
      <div class="form-row">
        <div class="form-group"><label>出生年</label><input id="b-year" type="number" placeholder="2024" min="1900" max="2100"></div>
        <div class="form-group"><label>月</label><input id="b-month" type="number" placeholder="3" min="1" max="12"></div>
        <div class="form-group"><label>日</label><input id="b-day" type="number" placeholder="15" min="1" max="31"></div>
        <div class="form-group"><label>时</label><input id="b-hour" type="number" placeholder="10" min="0" max="23"></div>
        <div class="form-group"><label>分</label><input id="b-min" type="number" placeholder="30" min="0" max="59"></div>
      </div>
      <div class="btn-row"><button class="btn btn-primary" id="btn-batch" onclick="doBatch()">批量评分</button></div>
    </div>
    <div id="batch-result" class="hidden"></div>
  </div>
</div>

<script>
function switchTab(name) {
  document.querySelectorAll('.tab').forEach((t,i) => {
    t.classList.toggle('active', (name==='score'?i===0:i===1));
  });
  document.getElementById('tab-score').classList.toggle('hidden', name!=='score');
  document.getElementById('tab-batch').classList.toggle('hidden', name!=='batch');
}

function barColor(pct) {
  if (pct >= 80) return '#52c41a';
  if (pct >= 60) return '#1890ff';
  if (pct >= 40) return '#faad14';
  return '#ff4d4f';
}

function renderScoreResult(d, container) {
  const pct = d.total.toFixed(1);
  const dims = [
    {name:'五格数理', score:d.wuge, max:30},
    {name:'三才配置', score:d.sancai, max:25},
    {name:'喜用神匹配', score:d.xiyong, max:20},
    {name:'内部五行', score:d.wuxing, max:15},
    {name:'阴阳平衡', score:d.yinyang, max:10},
  ];
  const dt = d.detail;
  const chars = d.name.split('');
  const strokeStr = chars.map((c,i) => c+'('+dt.strokes[i]+')').join(' ');
  container.innerHTML = '<div class="card">' +
    '<div class="result-header"><div><div class="score-big">'+pct+'</div><div class="score-label">总分 / 100</div></div><div style="font-size:24px;font-weight:700">'+d.name+'</div></div>' +
    dims.map(function(dm){
      var p = (dm.score/dm.max*100).toFixed(0);
      return '<div class="score-bar-wrap"><div class="score-bar-label"><span>'+dm.name+' (满分'+dm.max+')</span><span>'+dm.score.toFixed(1)+'</span></div><div class="score-bar"><div class="score-bar-fill" style="width:'+p+'%;background:'+barColor(p)+'"></div></div></div>';
    }).join('') +
    '<div class="detail-grid">' +
    '<div class="detail-item"><div class="label">康熙笔画</div><div class="value">'+strokeStr+'</div></div>' +
    '<div class="detail-item"><div class="label">五格</div><div class="value">天'+dt.tian_ge+' 人'+dt.ren_ge+' 地'+dt.di_ge+' 总'+dt.zong_ge+' 外'+dt.wai_ge+'</div></div>' +
    '<div class="detail-item"><div class="label">三才配置</div><div class="value">'+dt.sancai_desc+' → '+dt.sancai_jx+'</div></div>' +
    '<div class="detail-item"><div class="label">八字</div><div class="value">'+dt.bazi+'</div></div>' +
    '<div class="detail-item"><div class="label">喜用神</div><div class="value">'+dt.xiyong_shen+'</div></div>' +
    '<div class="detail-item"><div class="label">字五行</div><div class="value">'+dt.char_wuxing.join(' ')+'</div></div>' +
    '<div class="detail-item"><div class="label">阴阳</div><div class="value">'+dt.yinyang_pat+'</div></div>' +
    '</div></div>';
  container.classList.remove('hidden');
}

async function doScore() {
  const btn = document.getElementById('btn-score');
  const last = document.getElementById('s-last').value.trim();
  const first = document.getElementById('s-first').value.trim();
  if (!last || !first) { alert('请输入姓氏和名字'); return; }
  btn.disabled = true; btn.textContent = '评估中...';
  const container = document.getElementById('score-result');
  container.innerHTML = '<div class="loading"><div class="spinner"></div><p style="margin-top:12px">正在计算...</p></div>';
  container.classList.remove('hidden');
  try {
    const resp = await fetch('/api/score', {
      method:'POST', headers:{'Content-Type':'application/json'},
      body: JSON.stringify({
        last_name: last, first_name: first,
        year: parseInt(document.getElementById('s-year').value)||2024,
        month: parseInt(document.getElementById('s-month').value)||1,
        day: parseInt(document.getElementById('s-day').value)||1,
        hour: parseInt(document.getElementById('s-hour').value)||0,
        minute: parseInt(document.getElementById('s-min').value)||0,
      })
    });
    if (!resp.ok) { const t=await resp.text(); throw new Error(t); }
    const d = await resp.json();
    renderScoreResult(d, container);
  } catch(e) { container.innerHTML='<div class="card empty">评估失败: '+e.message+'</div>'; }
  finally { btn.disabled=false; btn.textContent='评估'; }
}

async function doBatch() {
  const btn = document.getElementById('btn-batch');
  const last = document.getElementById('b-last').value.trim();
  const words = document.getElementById('b-words').value.trim();
  if (!last || !words) { alert('请输入姓氏和备选字'); return; }
  btn.disabled = true; btn.textContent = '评分中...';
  const container = document.getElementById('batch-result');
  container.innerHTML = '<div class="loading"><div class="spinner"></div><p style="margin-top:12px">正在批量评分...</p></div>';
  container.classList.remove('hidden');
  try {
    const resp = await fetch('/api/batch', {
      method:'POST', headers:{'Content-Type':'application/json'},
      body: JSON.stringify({
        last_name: last, key_words: words,
        year: parseInt(document.getElementById('b-year').value)||2024,
        month: parseInt(document.getElementById('b-month').value)||1,
        day: parseInt(document.getElementById('b-day').value)||1,
        hour: parseInt(document.getElementById('b-hour').value)||0,
        minute: parseInt(document.getElementById('b-min').value)||0,
      })
    });
    if (!resp.ok) { const t=await resp.text(); throw new Error(t); }
    const d = await resp.json();
    let html = '<div class="card"><h2>排名 Top '+d.results.length+'</h2><div class="rank-list">';
    d.results.forEach(function(r, i){
      const cls = i===0?'top1':i===1?'top2':i===2?'top3':'other';
      html += '<div class="rank-item" onclick="showDetail('+i+')"><div class="rank-num '+cls+'">'+(i+1)+'</div><div class="rank-name">'+r.name+'</div><div class="rank-score">'+r.total.toFixed(1)+'</div></div>';
    });
    html += '</div></div><div id="batch-detail"></div>';
    container.innerHTML = html;
    container.classList.remove('hidden');
    window._batchResults = d.results;
    if (d.results.length > 0) showDetail(0);
  } catch(e) { container.innerHTML='<div class="card empty">评分失败: '+e.message+'</div>'; }
  finally { btn.disabled=false; btn.textContent='批量评分'; }
}

function showDetail(idx) {
  const d = window._batchResults[idx];
  const container = document.getElementById('batch-detail');
  renderScoreResult(d, container);
  document.querySelectorAll('.rank-item').forEach(function(el, i){
    el.style.background = i===idx ? '#fff5f5' : '';
  });
}
</script>
</body>
</html>`
