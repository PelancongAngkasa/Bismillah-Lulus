<template>
  <div class="min-h-screen flex flex-col">
    <Navbar @toggle-sidebar="toggleSidebar" />
    
    <!-- Spacer untuk navbar fixed -->
    <div class="h-16"></div>
    
    <div class="flex flex-1">
      <!-- Sidebar -->
      <Sidebar v-show="isSidebarVisible" class="min-h-[calc(100vh-64px)] bg-gray-800" />
      
      <!-- Main content area with scroll -->
      <div class="flex-1 bg-gray-100 overflow-auto transition-all duration-300">
        <ViewPartnerLimit />
      </div>
    </div>
  </div>
</template>

<script>
import Sidebar from '@/components/Organisms/Sidebar.vue';
import Navbar from '@/components/Molecules/Navbar.vue';
import ViewPartnerLimit from '@/components/Organisms/ViewPartnerLimit.vue';

export default {
  components: {
    Sidebar,
    Navbar,
    ViewPartnerLimit
  },
  data() {
    return {
      isSidebarVisible: true
    }
  },
  methods: {
    toggleSidebar() {
      this.isSidebarVisible = !this.isSidebarVisible;
    }
  }
}
</script>