<template>
  <div class="flex h-screen">
    <Sidebar />
    <main class="flex-1 bg-gray-100 px-0 pt-0 overflow-auto">
      <ViewPartnerLimit />
    </main>
  </div>
</template>

<script setup>
import Sidebar from '@/components/Organisms/Sidebar.vue';
import ViewPartnerLimit from '@/components/Organisms/ViewPartnerLimit.vue';

</script>