<template>
  <div class="flex h-screen">
    <SidebarAdmin />
    <main class="flex-1 bg-gray-100 px-0 pt-0 overflow-auto">
      <AddPartnerSection />
      <ViewPartner />
    </main>
  </div>
</template>

<script setup>
import AddPartnerSection from '@/components/Organisms/AddPartnerSection.vue'
import SidebarAdmin from '@/components/Organisms/SidebarAdmin.vue';
import ViewPartner from '@/components/Organisms/ViewPartner.vue'
</script>