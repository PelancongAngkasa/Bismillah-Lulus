<template>
  <div class="p-4">
    <h2 class="text-lg font-bold mb-2">Aplikasi Log</h2>
    <pre class="bg-gray-900 text-green-300 p-4 rounded overflow-x-auto max-h-[60vh]">{{ log }}</pre>
    <button @click="fetchLog" class="mt-2 px-4 py-1 bg-blue-500 text-white rounded">Refresh</button>
  </div>
</template>

<script>
export default {
  name: 'LogView',
  data() {
    return {
      log: ''
    }
  },
  mounted() {
    this.fetchLog()
  },
  methods: {
    async fetchLog() {
      try {
        const res = await fetch('/api/log')
        this.log = await res.text()
      } catch (e) {
        this.log = 'Gagal mengambil log: ' + e.message
      }
    }
  }
}
</script>