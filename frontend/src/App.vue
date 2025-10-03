<script setup lang="ts">
import { RouterView } from 'vue-router'
import { ToastContainer } from '@yuelioi/toast'
import TheHeader from './components/TheHeader.vue'

import TheFooter from '@yuelioi/components/the-footer'
import '@yuelioi/components/the-footer.css'
import '@yuelioi/toast/dist/toast.css'
</script>

<template>
  <ToastContainer />
  <TheHeader site-title="YTV" />

  <RouterView> </RouterView>

  <TheFooter />
</template>
