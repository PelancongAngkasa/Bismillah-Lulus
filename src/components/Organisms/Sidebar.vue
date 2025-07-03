<template>
  <aside class="bg-black text-white w-64 h-screen p-6">
    <h1 class="text-xl font-bold mb-6">SanapatiLink</h1>
    <button 
      @click="$router.push('/compose')" 
      class="mb-4 w-full bg-blue-500 text-white py-2 px-4 rounded">
      Compose
    </button>
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
      @click="goToAdmin"
      class="mt-8 w-full bg-yellow-500 text-black py-2 px-4 rounded font-semibold hover:bg-yellow-400 transition"
    >
      Masuk sebagai Admin
    </button>
  </aside>
</template>

<script>
export default {
  data() {
    return {
      menuItems: [
        {
          name: "Inbox",
          link: "/",
          icon: `<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V7a2 2 0 00-2-2H6a2 2 0 00-2 2v6m16 0v6a2 2 0 01-2 2H6a2 2 0 01-2-2v-6m16 0H4" /></svg>`
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