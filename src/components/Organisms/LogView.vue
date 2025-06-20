<template>
  <div class="p-4">
    <h2 class="text-lg font-bold mb-2">Aplikasi Log</h2>
    <input
      v-model="filter"
      type="text"
      placeholder="Cari log (misal: berhasil, gagal, diterima)"
      class="mb-3 px-3 py-2 border rounded w-full text-gray-800"
    />
    <pre class="bg-gray-900 text-green-300 p-4 rounded overflow-x-auto max-h-[60vh]">
{{ filteredLog }}
    </pre>
    <button @click="fetchLog" class="mt-2 px-4 py-1 bg-blue-500 text-white rounded">Refresh</button>
  </div>
</template>

<script>
export default {
  name: 'LogView',
  data() {
    return {
      log: '',
      filter: ''
    }
  },
  mounted() {
    this.fetchLog()
  },
  computed: {
    filteredLog() {
      if (!this.filter) return this.log;
      // Filter baris log yang mengandung kata filter (case-insensitive)
      return this.log
        .split('\n')
        .filter(line => line.toLowerCase().includes(this.filter.toLowerCase()))
        .join('\n');
    }
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