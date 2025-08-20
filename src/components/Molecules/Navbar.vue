<template>
  <nav class="bg-white border-b border-gray-200 h-16 flex items-center px-4 fixed w-full top-0 z-50">
    <!-- Left side with toggle button -->
    <div class="flex items-center">
      <button 
        @click="$emit('toggle-sidebar')" 
        class="p-2 hover:bg-gray-100 rounded-lg"
      >
        <svg 
          xmlns="http://www.w3.org/2000/svg" 
          class="h-6 w-6 text-gray-600" 
          fill="none" 
          viewBox="0 0 24 24" 
          stroke="currentColor"
        >
          <path 
            stroke-linecap="round" 
            stroke-linejoin="round" 
            stroke-width="2" 
            d="M4 6h16M4 12h16M4 18h16"
          />
        </svg>
      </button>
      <h1 class="text-xl font-semibold text-gray-800">SanapatiLink</h1>
    </div>

    <!-- Right side with profile -->
    <div class="relative ml-auto">
      <button 
        @click="isProfileOpen = !isProfileOpen" 
        class="flex items-center space-x-3 focus:outline-none"
      >
        <span class="text-gray-700">Sanapati Organization</span>
        <img 
          src="@/assets/logo.jpeg"
          alt="Profile" 
          class="h-8 w-8 rounded-full object-cover"
        />
      </button>

      <!-- Dropdown menu -->
      <div 
        v-if="isProfileOpen"
        class="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-xl py-2"
      >
        <div class="px-4 py-2 text-sm text-gray-700 border-b">
          Sanapati Organization
        </div>
        <router-link
          to="/admin/partner/add"
          class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
          @click="isProfileOpen = false"
        >
          Login sebagai Admin
        </router-link>
      </div>
    </div>
  </nav>
</template>

<script>
export default {
  name: 'Navbar',
  data() {
    return {
      isProfileOpen: false
    }
  },
  created() {
    document.addEventListener('click', this.handleClickOutside)
  },
  beforeUnmount() {
    document.removeEventListener('click', this.handleClickOutside)
  },
  methods: {
    handleClickOutside(event) {
      const dropdown = this.$el.querySelector('.relative')
      if (dropdown && !dropdown.contains(event.target)) {
        this.isProfileOpen = false
      }
    }
  }
}
</script>