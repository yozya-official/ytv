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
          class="group relative flex items-center space-x-3 transition-all duration-200 hover:opacity-90"
        >
          <!-- 站点标题 -->
          <div class=" ">
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
        <div class="relative flex">
          <ThemeToggle class="btn-hover btn btn-icon btn-ghost"></ThemeToggle>

          <button class="btn btn-ghost btn-icon btn-hover" @click="showHistory = true">
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
