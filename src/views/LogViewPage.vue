<template>
  <div class="flex h-screen">
    <SidebarAdmin />
    <main class="flex-1 bg-gray-100 px-0 pt-0 overflow-auto">
      <LogView />
    </main>
  </div>
</template>

<script setup>
import LogView from '@/components/Organisms/LogView.vue'
import SidebarAdmin from '@/components/Organisms/SidebarAdmin.vue'
</script>