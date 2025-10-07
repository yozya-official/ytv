<template>
  <footer class="mt-auto bg-gradient-to-b from-background to-muted/20">
    <div class="max-w-7xl mx-auto px-4 md:px-6 py-8 md:py-12">
      <!-- 免责声明 -->
      <div class="pt-6 border-t border-border">
        <div class="bg-muted/30 rounded-xl p-4 md:p-5 mb-6">
          <div class="flex items-start gap-3">
            <div class="flex-shrink-0 mt-0.5">
              <svg
                class="w-5 h-5 text-primary"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
            </div>
            <div class="flex-1 space-y-2 text-xs md:text-sm text-muted-foreground leading-relaxed">
              <h4 class="text-sm font-semibold text-foreground">免责声明</h4>
              <p>本站仅为视频搜索工具，不存储、上传或分发任何视频内容。</p>
              <p>所有视频均来自第三方API接口。请勿相信视频中的广告等内容。</p>
              <p class="">
                如您认为某内容涉及侵权，请发送邮件至
                <a
                  href="mailto:yuelioi1210@gmail.com"
                  class="text-primary hover:underline font-medium"
                  >yuelioi1210@gmail.com</a
                >
                ，我们将在核实后尽快处理。
              </p>
              <p>🙅‍♂️请不要在公共场合宣传本网站, 仅限🐶群友使用</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
// Footer组件，无需额外逻辑
</script>

<style scoped>
@keyframes pulse {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

.animate-pulse {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}
</style>
