<template>
  <div class="error-page">
    <div class="error-bg-orb error-bg-orb--1" />
    <div class="error-bg-orb error-bg-orb--2" />

    <div class="error-content">
      <!-- 装饰性警示几何图形 -->
      <div class="error-figure">
        <svg width="120" height="120" viewBox="0 0 120 120" fill="none" class="figure-svg">
          <!-- 齿轮外圈 -->
          <circle cx="60" cy="60" r="40" stroke="currentColor" stroke-width="2" fill="none" opacity="0.3" />
          <!-- 齿轮内圈 -->
          <circle cx="60" cy="60" r="25" stroke="currentColor" stroke-width="2" fill="none" />
          <!-- 齿 -->
          <g stroke="currentColor" stroke-width="3" stroke-linecap="round">
            <line x1="60" y1="12" x2="60" y2="22" />
            <line x1="60" y1="98" x2="60" y2="108" />
            <line x1="12" y1="60" x2="22" y2="60" />
            <line x1="98" y1="60" x2="108" y2="60" />
            <line x1="26" y1="26" x2="33" y2="33" />
            <line x1="87" y1="87" x2="94" y2="94" />
            <line x1="26" y1="94" x2="33" y2="87" />
            <line x1="87" y1="33" x2="94" y2="26" />
          </g>
          <!-- 叹号 -->
          <line x1="60" y1="48" x2="60" y2="66" stroke="#c44536" stroke-width="3.5" stroke-linecap="round" opacity="0.8" />
          <circle cx="60" cy="78" r="2.5" fill="#c44536" opacity="0.8" />
        </svg>
      </div>

      <h1 class="error-title">服务器错误</h1>
      <p class="error-desc">服务器遇到了意外错误，请稍后重试或联系管理员</p>

      <el-button class="error-btn" @click="$router.push('/')">
        <svg width="16" height="16" viewBox="0 0 16 16" fill="none" class="btn-icon">
          <path d="M8 1L1 7h3v7h8V7h3L8 1z" stroke="currentColor" stroke-width="1.2" fill="none" />
          <path d="M6 11h4" stroke="currentColor" stroke-width="1.2" stroke-linecap="round" />
        </svg>
        返回首页
      </el-button>
    </div>
  </div>
</template>

<style scoped>
.error-page {
  --primary: #e6a23c;
  --primary-light: #f0b854;

  min-height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  background: var(--el-bg-color-page);
  font-family: "PingFang SC", "Noto Sans SC", "Microsoft YaHei", sans-serif;
}

.error-bg-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(100px);
  pointer-events: none;
  z-index: 0;
}
.error-bg-orb--1 {
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, rgba(230, 162, 60, 0.08), transparent 70%);
  top: -80px;
  left: -80px;
}
.error-bg-orb--2 {
  width: 350px;
  height: 350px;
  background: radial-gradient(circle, rgba(196, 69, 54, 0.04), transparent 70%);
  bottom: -100px;
  right: -60px;
}

.error-content {
  position: relative;
  z-index: 1;
  text-align: center;
  padding: 40px;
  animation: fadeUp 0.6s cubic-bezier(0.16, 1, 0.3, 1) both;
}

@keyframes fadeUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.error-figure {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 32px;
}

.figure-svg {
  color: var(--primary);
  filter: drop-shadow(0 4px 24px rgba(230, 162, 60, 0.12));
}

.error-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--el-text-color-primary);
  margin: 0 0 12px;
}

.error-desc {
  font-size: 15px;
  color: var(--el-text-color-secondary);
  margin: 0 0 36px;
  line-height: 1.6;
}

.error-btn {
  height: 44px;
  padding: 0 28px;
  border-radius: 10px;
  font-size: 15px;
  font-weight: 500;
  border: 1px solid var(--primary);
  color: var(--primary);
  background: transparent;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  transition: all 0.25s ease;
  cursor: pointer;
}

.error-btn:hover {
  background: var(--primary);
  color: #fff;
  transform: translateY(-1px);
  box-shadow: 0 4px 16px rgba(230, 162, 60, 0.25);
}

.btn-icon {
  flex-shrink: 0;
}
</style>
