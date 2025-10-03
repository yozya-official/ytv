<template>
  <!-- 头部装饰背景 -->
  <div class="w-full bg-gradient-to-r from-primary/10 via-chart-2/10 to-chart-3/10">
    <div class="text-center py-6">
      <!-- 图标/徽章 -->
      <div
        class="inline-flex items-center justify-center size-16 rounded-2xl bg-gradient-to-br from-primary to-chart-2 text-primary-foreground mb-6 shadow-xl hover-scale"
      >
        <slot name="icon">
          <!-- 默认图标 -->
          <svg class="w-8 h-8" fill="white" viewBox="0 0 20 20">
            <path
              d="M9 6a3 3 0 11-6 0 3 3 0 016 0zM17 6a3 3 0 11-6 0 3 3 0 016 0zM12.93 17c.046-.327.07-.66.07-1a6.97 6.97 0 00-1.5-4.33A5 5 0 0119 16v1h-6.07zM6 11a5 5 0 015 5v1H1v-1a5 5 0 015-5z"
            />
          </svg>
        </slot>
      </div>
      <!-- 标题和副标题 -->
      <h2 class="text-foreground mb-2">{{ title }}</h2>
      <p v-if="subtitle" class="text-muted-foreground text-sm">{{ subtitle }}</p>

      <!-- 装饰线 -->
      <div
        class="w-24 h-1 bg-gradient-to-r from-primary to-chart-2 mx-auto mt-4 rounded-full"
      ></div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title: string
  subtitle?: string
}>()
</script>
