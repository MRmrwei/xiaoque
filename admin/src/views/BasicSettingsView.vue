<template>
  <div class="adm-page">
    <header class="adm-page-head">
      <div>
        <p class="adm-kicker">BASIC SETTINGS</p>
        <h1 class="adm-title">基础设置</h1>
        <p class="adm-sub">管理后台的基础设置。当前开放「后台登录」配置。</p>
      </div>
    </header>

    <div class="category-grid">
      <router-link to="/settings/login" class="category-card">
        <span class="category-icon"><el-icon><Lock /></el-icon></span>
        <span class="category-body">
          <span class="category-head">
            <strong>后台登录</strong>
            <span class="category-tag">已开放</span>
          </span>
          <span class="category-desc">
            图形验证码、登录页 Logo、名称与展示文案。
          </span>
          <span class="category-foot">进入设置 <span>→</span></span>
        </span>
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Lock } from '@element-plus/icons-vue'
</script>

<style scoped>
.category-grid {
  display: grid;
  grid-template-columns: minmax(0, 560px);
  gap: 14px;
}

.category-card {
  display: flex;
  gap: 14px;
  min-width: 0;
  padding: 20px;
  border: 1px solid var(--adm-line);
  border-radius: var(--adm-radius);
  background: var(--adm-card);
  box-shadow: var(--adm-shadow);
  color: inherit;
  cursor: pointer;
  transition: border-color 0.18s ease, box-shadow 0.18s ease, transform 0.18s ease;
}

.category-card:hover {
  border-color: var(--adm-accent);
  box-shadow: 0 1px 2px rgba(20, 21, 24, 0.02), 0 14px 36px rgba(20, 21, 24, 0.08);
  transform: translateY(-1px);
}

.category-icon {
  width: 42px;
  height: 42px;
  flex: none;
  display: grid;
  place-items: center;
  border-radius: 11px;
  background: var(--adm-accent-soft);
  color: var(--adm-accent-deep);
  font-size: 19px;
}

.category-body {
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-width: 0;
  flex: 1;
}

.category-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.category-head strong {
  font-size: 14px;
  font-weight: 600;
  letter-spacing: -0.01em;
}

.category-tag {
  flex: none;
  padding: 3px 8px;
  border-radius: 999px;
  background: var(--adm-green-soft);
  color: var(--adm-green);
  font-size: 10.5px;
  font-weight: 600;
  white-space: nowrap;
}

.category-desc {
  font-size: 12px;
  line-height: 1.65;
  color: var(--adm-muted);
}

.category-foot {
  margin-top: auto;
  padding-top: 6px;
  font-size: 11.5px;
  font-weight: 600;
  color: var(--adm-accent-deep);
}

@media (max-width: 640px) {
  .category-grid {
    grid-template-columns: 1fr;
  }
}
</style>
