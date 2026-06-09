<template>
  <div class="error-page">
    <div class="error-bg-orb error-bg-orb--1" />
    <div class="error-bg-orb error-bg-orb--2" />

    <div class="error-content">
      <!-- 装饰性几何图形 -->
      <div class="error-figure">
        <div class="figure-grid">
          <div class="grid-cell" v-for="n in 16" :key="n"
            :class="{ dim: [1,2,5,6,7,10,11,12,13,14,15,16].includes(n) }" />
        </div>
        <div class="figure-code">404</div>
      </div>

      <h1 class="error-title">页面未找到</h1>
      <p class="error-desc">您访问的页面不存在，可能已被移除或地址有误</p>

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
  --primary: #0d7377;
  --primary-light: #14a0a5;

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
  background: radial-gradient(circle, rgba(13, 115, 119, 0.08), transparent 70%);
  top: -100px;
  left: -100px;
}
.error-bg-orb--2 {
  width: 300px;
  height: 300px;
  background: radial-gradient(circle, rgba(196, 69, 54, 0.06), transparent 70%);
  bottom: -50px;
  right: -50px;
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

/* ===== 装饰性网格图形 ===== */
.error-figure {
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 32px;
}

.figure-grid {
  display: grid;
  grid-template-columns: repeat(4, 28px);
  grid-template-rows: repeat(4, 28px);
  gap: 4px;
  opacity: 0.15;
}

.grid-cell {
  border-radius: 4px;
  background: var(--primary);
  transition: opacity 0.3s;
}

.grid-cell.dim {
  opacity: 0.2;
}

.figure-code {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 72px;
  font-weight: 800;
  letter-spacing: 4px;
  color: var(--primary);
  line-height: 1;
  text-shadow: 0 2px 24px rgba(13, 115, 119, 0.15);
}

/* ===== 文字 ===== */
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

/* ===== 按钮 ===== */
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
  box-shadow: 0 4px 16px rgba(13, 115, 119, 0.25);
}

.btn-icon {
  flex-shrink: 0;
}
</style>
