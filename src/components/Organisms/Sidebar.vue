<template>
  <aside class="bg-black text-white w-64 h-screen p-6">
    <h1 class="text-xl font-bold mb-6">Gmail</h1>
    <button 
      @click="$router.push('/compose')" 
      class="mb-4 w-full bg-blue-500 text-white py-2 px-4 rounded">
      New
    </button>
    <SidebarMenu :menuItems="menuItems" />
  </aside>
</template>

<script>
import SidebarMenu from "@/components/Molecules/SidebarMenu.vue";

export default {
  components: { SidebarMenu },
  data() {
    return {
      menuItems: [
        { name: "Inbox", link: "/" },
        { name: "Starred", link: "#" },
        { name: "Snoozed", link: "#" },
        { name: "Sent", link: "#" },
        { name: "Drafts", link: "#" },
      ],
    };
  },
};
</script>
