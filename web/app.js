// Beeyard dashboard
async function fetchJSON(url) {
  const res = await fetch(url);
  return res.json();
}

async function refresh() {
  try {
    const hives = await fetchJSON('/api/apiaries/1/hives');
    document.getElementById('hive-count').textContent = hives ? hives.length : 0;
    const alerts = await fetchJSON('/api/apiaries/1/alerts');
    document.getElementById('alert-count').textContent = alerts ? alerts.length : 0;
    document.getElementById('inspect-count').textContent = '-';
  } catch (e) {
    console.error('refresh error', e);
  }
}

refresh();
setInterval(refresh, 30000);
