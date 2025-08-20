<template>
  <aside class="bg-gray-800 w-64 flex-shrink-0 min-h-full">
    <nav class="p-4 text-white h-full">
      <ul class="space-y-2">
        <li v-for="item in menuItems" :key="item.link">
          <router-link
            :to="item.link"
            class="flex items-center gap-3 px-4 py-2 rounded hover:bg-blue-700 transition-colors font-medium"
            active-class="bg-blue-900"
          >
            <span v-if="item.icon" v-html="item.icon"></span>
            <span>{{ item.name }}</span>
          </router-link>
        </li>
      </ul>
    </nav>
  </aside>
</template>

<script>
export default {
  data() {
    return {
      menuItems: [
        {
          name: "Kirim Pesan",
          link: "/kirim",
          icon: `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-send-icon lucide-send"><path d="M14.536 21.686a.5.5 0 0 0 .937-.024l6.5-19a.496.496 0 0 0-.635-.635l-19 6.5a.5.5 0 0 0-.024.937l7.93 3.18a2 2 0 0 1 1.112 1.11z"/><path d="m21.854 2.147-10.94 10.939"/></svg>`
        },
        {
          name: "Pesan Masuk",
          link: "/",
          icon: `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-inbox-icon lucide-inbox"><polyline points="22 12 16 12 14 15 10 15 8 12 2 12"/><path d="M5.45 5.11 2 12v6a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2v-6l-3.45-6.89A2 2 0 0 0 16.76 4H7.24a2 2 0 0 0-1.79 1.11z"/></svg>`
        },
        {
          name: "Daftar Partner",
          link: "/partner",
          icon: `<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" /></svg>`
        },
      ],
    };
  },
  methods: {
    goToAdmin() {
      this.$router.push('/admin/partner/add');
    }
  }
};
</script>