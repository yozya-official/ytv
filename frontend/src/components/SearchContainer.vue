<template>
  <div v-if="loading" class="space-y-8">
    <div class="skeleton h-6 w-32"></div>
    <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
      <div v-for="j in 24" :key="j" class="space-y-2">
        <div class="skeleton aspect-[2/3] w-full"></div>
        <div class="skeleton h-4 w-full"></div>
        <div class="skeleton h-3 w-2/3"></div>
      </div>
    </div>
  </div>

  <div v-else>
    <div v-if="results.length">
      <SearchResult ref="filterResultRef" :vods="results"></SearchResult>
    </div>
    <!-- 空状态 -->
  </div>
</template>

<script setup lang="ts">
defineProps<{ results: VodItem[]; loading: boolean }>()
</script>
