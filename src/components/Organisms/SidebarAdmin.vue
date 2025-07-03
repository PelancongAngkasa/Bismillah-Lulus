<template>
  <aside class="bg-black text-white w-64 h-screen p-6">
    <h1 class="text-xl font-bold mb-6">SanapatiLink</h1>
    <nav>
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
    <button 
      @click="$router.push('/compose')" 
      class="mb-4 w-full bg-red-500 text-white py-2 px-4 mt-2 rounded">
      Kembali
    </button>
  </aside>
</template>

<script>
export default {
  data() {
    return {
      menuItems: [
        {
          name: "Add Partner",
          link: "/admin/partner/add",
          icon: `<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" /></svg>`
        },
        {
          name: "Log Aplikasi",
          link: "/admin/log",
          icon: `<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17v-2a2 2 0 012-2h2a2 2 0 012 2v2m-6 4h6a2 2 0 002-2V7a2 2 0 00-2-2h-1.5a1 1 0 01-.7-.3l-1.5-1.4a1 1 0 00-.7-.3H9a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg>`
        },
        {
          name: "Konfigurasi Pesan",
          link: "/admin/pmode/edit",
          icon: `<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536M9 13h3l8-8a2.828 2.828 0 00-4-4l-8 8v3z" /></svg>`
        }
      ],
    };
  },
};
</script>