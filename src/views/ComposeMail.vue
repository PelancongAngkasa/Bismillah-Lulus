<template>
  <div class="flex h-screen">
    <Sidebar />

    <div class="flex-1 p-8">
      <div class="max-w-3xl mx-auto bg-white shadow rounded p-6">
        <h2 class="text-lg font-bold mb-4">Compose Mail</h2>

        <!-- Input untuk To -->
        <input v-model="toParty" type="text" placeholder="To" class="mb-2 w-full p-2 border rounded" />

        <!-- Input untuk Message -->
        <textarea v-model="message" placeholder="Write your message here" class="w-full p-2 border rounded h-32 mb-4"></textarea>

        <!-- Area Upload File -->
        <div @drop.prevent="handleDrop" @dragover.prevent class="border-2 border-dashed border-gray-400 p-6 text-center rounded mb-4">
          <p class="text-gray-500">Drag and drop a file here or click to upload</p>
          <input type="file" ref="fileInput" class="hidden" @change="handleFileUpload" />
          <button class="mt-2 bg-gray-200 px-4 py-2 rounded hover:bg-gray-300" @click="triggerFileInput">
            Choose File
          </button>
          <div v-if="attachment" class="mt-2">
            <p class="truncate">{{ attachment.name }}</p>
            <button class="text-red-500 hover:text-red-700" @click="removeAttachment">Remove</button>
          </div>
        </div>

        <!-- Tombol Kirim -->
        <div class="flex justify-end">
          <button :disabled="loading" @click="sendMail" class="bg-blue-500 text-white px-4 py-2 rounded">
            {{ loading ? "Sending..." : "Send" }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Sidebar from "@/components/Organisms/Sidebar.vue";

export default {
  components: { Sidebar },
  data() {
    return {
      toParty: "", // Penerima pesan
      message: "", // Pesan utama
      attachment: null, // File lampiran
      loading: false, // Status loading
    };
  },
  methods: {
    triggerFileInput() {
      this.$refs.fileInput.click();
    },
    handleFileUpload(event) {
      this.attachment = event.target.files[0];
    },
    handleDrop(event) {
      this.attachment = event.dataTransfer.files[0];
    },
    removeAttachment() {
      this.attachment = null;
    },
    async sendMail() {
      if (!this.toParty || !this.message) {
        alert("To and Message fields are required.");
        return;
      }

      this.loading = true;

      // Membuat FormData untuk backend
      const formData = new FormData();
      formData.append("fromParty", "Company A"); // FromParty default
      formData.append("toParty", this.toParty);
      formData.append("service", "EmailService");
      formData.append("action", "SendMail");
      formData.append("messageId", `msg-${Date.now()}`); // ID pesan unik
      formData.append("payload", this.message);

      // Tambahkan file lampiran jika ada
      if (this.attachment) {
        formData.append("attachment", this.attachment);
      }

      try {
        // Mengirim data ke backend
        const response = await fetch("http://localhost:8081/api/as4/send", {
          method: "POST",
          body: formData,
        });

        if (!response.ok) throw new Error("Failed to send message");

        alert("Message sent successfully!");
        this.resetForm(); // Mengosongkan form setelah berhasil
      } catch (error) {
        alert(`Error: ${error.message}`);
      } finally {
        this.loading = false;
      }
    },
    resetForm() {
      this.toParty = "";
      this.message = "";
      this.attachment = null;
    },
  },
};
</script>

<style scoped>
/* Tambahkan styling jika diperlukan */
</style>
