<template>
  <div class="p-8 max-w-5xl mx-auto">
    <h2 class="text-2xl font-bold mb-4">Daftar File Konfigurasi</h2>
    <ul class="mb-6 space-y-2">
      <li
        v-for="file in files"
        :key="file"
        class="flex items-center"
      >
        <button
          class="flex items-center gap-2 text-base font-medium text-blue-700 hover:text-blue-900 transition-colors px-3 py-2 rounded hover:bg-blue-50 focus:outline-none"
          @click="selectFile(file)"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17v-2a2 2 0 012-2h2a2 2 0 012 2v2m-6 4h6a2 2 0 002-2V7a2 2 0 00-2-2h-1.5a1 1 0 01-.7-.3l-1.5-1.4a1 1 0 00-.7-.3H9a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
          {{ file }}
        </button>
      </li>
    </ul>

    <div v-if="selectedFile">
      <h3 class="text-xl font-semibold mb-2">Edit: {{ selectedFile }}</h3>
      <textarea
        v-model="fileContent"
        class="w-full h-96 border rounded p-3 font-mono text-base"
      ></textarea>
      <div class="mt-3 flex gap-2">
        <button
          @click="saveFile"
          class="px-4 py-2 bg-blue-600 text-white rounded"
        >Simpan</button>
        <span v-if="saveMsg" class="text-green-600">{{ saveMsg }}</span>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "PmodeEdit",
  data() {
    return {
      files: [],
      selectedFile: "",
      fileContent: "",
      saveMsg: ""
    }
  },
  mounted() {
    this.fetchFiles()
  },
  methods: {
    async fetchFiles() {
      const res = await fetch('http://localhost:8081/api/pmode/list')
      this.files = await res.json()
    },
    async selectFile(file) {
      this.selectedFile = file
      this.saveMsg = ""
      const res = await fetch(`http://localhost:8081/api/pmode/get?name=${encodeURIComponent(file)}`)
      this.fileContent = await res.text()
    },
    async saveFile() {
      const res = await fetch("http://localhost:8081/api/pmode/save", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          name: this.selectedFile,
          content: this.fileContent
        })
      })
      if (res.ok) {
        this.saveMsg = "Berhasil disimpan!"
        // Tutup tampilan editor setelah berhasil simpan
        setTimeout(() => {
          this.selectedFile = ""
          this.fileContent = ""
          this.saveMsg = ""
        }, 1000)
      } else {
        this.saveMsg = "Gagal menyimpan file."
      }
    }
  }
}
</script>