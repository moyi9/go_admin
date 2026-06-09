<template>
  <div class="error-page">
    <div class="error-bg-orb error-bg-orb--1" />
    <div class="error-bg-orb error-bg-orb--2" />

    <div class="error-content">
      <!-- 装饰性锁形几何图形 -->
      <div class="error-figure">
        <svg width="120" height="140" viewBox="0 0 120 140" fill="none" class="figure-svg">
          <!-- 锁体 -->
          <rect x="20" y="60" width="80" height="65" rx="6" stroke="currentColor" stroke-width="2.5" fill="none" />
          <rect x="32" y="72" width="56" height="20" rx="2" stroke="currentColor" stroke-width="1.5" fill="none" opacity="0.4" />
          <!-- 锁梁 -->
          <path d="M35 60V40c0-13.8 11.2-25 25-25s25 11.2 25 25v20" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" fill="none" />
          <!-- X 标记 -->
          <line x1="52" y1="102" x2="68" y2="118" stroke="#c44536" stroke-width="3" stroke-linecap="round" opacity="0.7" />
          <line x1="68" y1="102" x2="52" y2="118" stroke="#c44536" stroke-width="3" stroke-linecap="round" opacity="0.7" />
        </svg>
      </div>

      <h1 class="error-title">访问被拒绝</h1>
      <p class="error-desc">您没有权限访问此页面，请联系管理员获取权限</p>

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
  --primary: #c44536;
  --primary-light: #e05a4a;

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
  background: radial-gradient(circle, rgba(196, 69, 54, 0.08), transparent 70%);
  top: -100px;
  right: -100px;
}
.error-bg-orb--2 {
  width: 300px;
  height: 300px;
  background: radial-gradient(circle, rgba(13, 115, 119, 0.05), transparent 70%);
  bottom: -50px;
  left: -50px;
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
  filter: drop-shadow(0 4px 24px rgba(196, 69, 54, 0.12));
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
  box-shadow: 0 4px 16px rgba(196, 69, 54, 0.25);
}

.btn-icon {
  flex-shrink: 0;
}
</style>
