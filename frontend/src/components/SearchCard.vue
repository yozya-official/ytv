<template>
  <div class="flex flex-nowrap h-40">
    <!-- 图片容器 -->
    <div class="w-28 relative overflow-hidden flex-shrink-0">
      <img
        class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
        :src="vod.vod_pic_thumb || vod.vod_pic"
        :alt="vod.vod_name"
        @error="handleImageError"
      />

      <!-- 悬浮操作按钮 -->
      <div class="image-mask flex items-center justify-center">
        <div class="bg-background/90 backdrop-blur-sm rounded-full p-3 shadow-lg">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="size-6 text-primary"
            viewBox="0 0 24 24"
            fill="currentColor"
          >
            <path d="M8 5v14l11-7z" />
          </svg>
        </div>
      </div>
      <!-- 年份标签 -->
      <div class="absolute top-2 right-2" v-if="vod.vod_year != '0'">
        <span
          class="px-2 py-0.5 rounded text-xs font-semibold bg-muted text-muted-foreground backdrop-blur-sm"
        >
          {{ vod.vod_year }}
        </span>
      </div>
    </div>

    <!-- 内容区域 -->
    <div class="flex flex-col p-1.5 md:p-2 flex-1 gap-2 min-w-0">
      <!-- 标题 -->
      <span
        class="text-base font-semibold text-card-foreground line-clamp-2 group-hover:text-primary transition-colors leading-snug"
      >
        {{ vod.vod_name }}
      </span>

      <!-- 类型标签 -->
      <div class="flex flex-wrap gap-1.5">
        <span
          class="inline-flex items-center gap-1 px-2 py-0.5 rounded-md text-xs font-medium bg-primary/10 text-primary border border-primary/20"
        >
          {{ vod.type_name }}
        </span>
        <span
          class="inline-flex items-center gap-1 px-2 py-0.5 rounded-md text-xs font-medium bg-primary/10 text-destructive/75 border border-primary/20"
        >
          <span class="font-sm">共{{ vod.episodes.length }}集</span>
        </span>
      </div>

      <!-- 信息区域 -->
      <div class="mt-auto space-y-1.5 text-sm">
        <!-- 备注 -->
        <div v-if="vod.vod_remarks" class="flex items-start gap-2 text-muted-foreground">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-4 w-4 flex-shrink-0 mt-0.5"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <span class="line-clamp-1">{{ vod.vod_remarks }}</span>
        </div>

        <!-- 来源 -->
        <div class="flex items-center gap-2 text-foreground/70">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-4 w-4 flex-shrink-0"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01"
            />
          </svg>
          <span class="font-medium truncate">{{ vod.source_name }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { handleImageError } from '@/utils'

defineProps<{ vod: VodItem }>()
</script>
