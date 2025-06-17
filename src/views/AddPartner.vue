<template>
  <div class="flex h-screen">
    <Sidebar />
    <main class="flex-1 bg-gray-100 px-0 pt-0 overflow-auto">
      <AddPartnerSection />
      <ViewPartner />
    </main>
  </div>
</template>

<script setup>
import Sidebar from '@/components/Organisms/Sidebar.vue'
import AddPartnerSection from '@/components/Organisms/AddPartnerSection.vue'
import ViewPartner from '@/components/Organisms/ViewPartner.vue'
</script>