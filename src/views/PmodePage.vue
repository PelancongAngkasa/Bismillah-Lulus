<template>
  <div class="flex h-screen">
    <SidebarAdmin />
    <main class="flex-1 bg-gray-100 px-0 pt-0 overflow-auto">
      <PmodeEdit />
    </main>
  </div>
</template>

<script setup>
import PmodeEdit from '@/components/Organisms/PmodeEdit.vue';
import SidebarAdmin from '@/components/Organisms/SidebarAdmin.vue';
</script>