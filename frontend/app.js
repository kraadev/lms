document.addEventListener('DOMContentLoaded', () => {
  const statusIndicator = document.querySelector('.status-indicator');
  const statusText = document.getElementById('status-text');

  async function checkServerHealth() {
    try {
      const res = await fetch('http://localhost:8080/health');
      if (res.ok) {
        const data = await res.json();
        statusIndicator.className = 'status-indicator online';
        statusText.textContent = `Backend Terhubung (${data.message})`;
      } else {
        throw new Error('Server returned non-200');
      }
    } catch {
      statusIndicator.className = 'status-indicator offline';
      statusText.textContent = 'Backend tidak terhubung (port 8080)';
    }
  }

  checkServerHealth();
});
