<template>
  <header
    class="sticky top-0 z-50 w-full border-b border-border/40 bg-background/80 backdrop-blur-xl supports-[backdrop-filter]:bg-background/60 transition-all duration-300"
  >
    <div class="container mx-auto flex h-16 max-w-screen-2xl items-center justify-between px-6">
      <!-- Logo 区域 -->
      <div class="flex items-center space-x-3">
        <a
          href="/"
          aria-label="首页"
          class="group flex items-center space-x-3 transition-all duration-200 hover:opacity-90"
        >
          <div class="relative">
            <!-- 装饰光晕 -->
            <div
              class="absolute inset-0 rounded-xl bg-gradient-to-br from-primary to-chart-2 opacity-0 group-hover:opacity-20 transition-opacity duration-300 blur-md -z-10"
            ></div>
          </div>
          <!-- 站点标题 -->
          <div class="hidden sm:block">
            <h1
              class="text-2xl sm:text-3xl font-extrabold bg-clip-text text-transparent bg-gradient-to-r from-purple-500 via-pink-500 to-yellow-400 transition-all duration-300 hover:scale-105"
              style="font-family: 'Poppins', 'Segoe UI', sans-serif"
            >
              {{ siteTitle }}
            </h1>
          </div>
        </a>
      </div>

      <!-- 右侧功能区 -->
      <div class="flex items-center">
        <!-- 主题切换器+设置 -->
        <div class="relative flex gap-4">
          <ThemeToggle class="btn-scale"></ThemeToggle>

          <button
            class="cursor-pointer btn btn-ghost btn-icon btn-scale"
            @click="showHistory = true"
          >
            <svg class="size-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
          </button>
        </div>
      </div>
    </div>

    <PlayerHistory v-model:show-history="showHistory"></PlayerHistory>
  </header>
</template>

<script setup lang="ts">
import ThemeToggle from '@yuelioi/components/theme-toggle'
import '@yuelioi/components/theme-toggle.css'

const showHistory = ref(false)

defineProps(['siteTitle'])

// 滚动效果 - 添加阴影和背景模糊
function initScrollEffect() {
  const header = document.querySelector('header')
  let lastScrollY = window.scrollY

  const onScroll = () => {
    const scrollY = window.scrollY

    // 滚动阴影效果
    if (scrollY > 10) {
      header?.classList.add('shadow-lg', 'bg-background/90')
      header?.classList.remove('bg-background/80')
    } else {
      header?.classList.remove('shadow-lg', 'bg-background/90')
      header?.classList.add('bg-background/80')
    }

    // 滚动方向检测（可用于隐藏/显示导航栏）
    if (scrollY > lastScrollY && scrollY > 100) {
      // 向下滚动 - 可以添加隐藏逻辑
      header?.classList.add('transform', '-translate-y-1')
    } else {
      // 向上滚动
      header?.classList.remove('transform', '-translate-y-1')
    }

    lastScrollY = scrollY
  }

  window.addEventListener('scroll', onScroll, { passive: true })
  return () => window.removeEventListener('scroll', onScroll)
}

let cleanupScroll: (() => void) | undefined

onMounted(() => {
  cleanupScroll = initScrollEffect()
})

onBeforeUnmount(() => {
  cleanupScroll?.()
  // 确保页面可滚动
  document.body.style.overflow = ''
})
</script>

<style scoped lang="css">
/* Logo 悬停效果 */
.group:hover .w-10 {
  transform: scale(1.05) rotate(2deg);
}

/* 导航链接活跃状态动画 */
nav a {
  position: relative;
  overflow: hidden;
}

nav a::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    90deg,
    transparent,
    var(--primary-rgb, 84 114 183) / 0.1,
    transparent
  );
  transition: left 0.5s;
}

nav a:hover::before {
  left: 100%;
}

/* 汉堡菜单动画优化 */
#mobile-menu-toggle span {
  transform-origin: center;
}

/* 移动端菜单项悬停效果 */
#mobile-menu a {
  backdrop-filter: blur(8px);
}

/* 头部毛玻璃效果增强 */
header {
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
}

/* 响应式字体大小 */
@media (max-width: 640px) {
  h1 {
    font-size: 1rem;
  }
}
</style>
