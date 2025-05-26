<template>
  <div class="flex h-screen">
    <Sidebar />

    <div class="flex-1 p-8">
      <div class="max-w-3xl mx-auto bg-white shadow rounded p-6">
        <h2 class="text-lg font-bold mb-4">Compose Mail</h2>

        <!-- Input untuk To -->
        <input
          v-model="toParty"
          type="text"
          placeholder="To"
          class="mb-2 w-full p-2 border rounded"
        />

        <!-- Input untuk Subject -->
        <input
          v-model="subject"
          type="text"
          placeholder="Subject"
          class="mb-2 w-full p-2 border rounded"
        />

        <!-- Input untuk Message -->
        <textarea
          v-model="message"
          placeholder="Write your message here"
          class="w-full p-2 border rounded h-32 mb-4"
        ></textarea>

        <!-- Area Upload File -->
        <div
          @drop.prevent="handleDrop"
          @dragover.prevent
          class="border-2 border-dashed border-gray-400 p-6 text-center rounded mb-4"
        >
          <p class="text-gray-500">Drag and drop files here or click to upload</p>
          <input
            type="file"
            ref="fileInput"
            class="hidden"
            @change="handleFileUpload"
            multiple
          />
          <button
            class="mt-2 bg-gray-200 px-4 py-2 rounded hover:bg-gray-300"
            @click="triggerFileInput"
          >
            Choose Files
          </button>
          
          <!-- Enhanced Attachment Preview -->
          <div v-if="attachments.length" class="mt-4">
            <h3 class="text-sm font-semibold text-gray-700 mb-2">Attachments:</h3>
            <ul class="space-y-2">
              <li
                v-for="(file, index) in attachments"
                :key="file.name + index"
                class="flex items-center justify-between p-2 bg-gray-50 rounded"
              >
                <div class="flex items-center truncate">
                  <span class="text-sm font-medium truncate">{{ file.name }}</span>
                  <span class="text-xs text-gray-500 ml-2">
                    ({{ formatFileSize(file.size) }})
                  </span>
                </div>
                <button
                  class="text-red-500 hover:text-red-700 ml-2"
                  @click="removeAttachment(index)"
                  title="Remove attachment"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
                  </svg>
                </button>
              </li>
            </ul>
            <p class="text-xs text-gray-500 mt-2">
              Total: {{ attachments.length }} file(s), {{ formatFileSize(totalAttachmentSize) }}
            </p>
          </div>
        </div>

        <!-- Tombol Kirim -->
        <div class="flex justify-end">
          <button
            :disabled="loading || !canSend"
            @click="sendMail"
            class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 disabled:bg-blue-300 disabled:cursor-not-allowed"
          >
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
      toParty: "",
      subject: "",
      message: "",
      attachments: [],
      loading: false,
    };
  },
  computed: {
    totalAttachmentSize() {
      return this.attachments.reduce((total, file) => total + file.size, 0);
    },
    canSend() {
      return this.toParty && this.message && !this.loading;
    }
  },
  methods: {
    triggerFileInput() {
      this.$refs.fileInput.click();
    },
    handleFileUpload(event) {
      this.addFiles(Array.from(event.target.files));
      this.$refs.fileInput.value = ''; // Reset input
    },
    handleDrop(event) {
      this.addFiles(Array.from(event.dataTransfer.files));
    },
    addFiles(files) {
      files.forEach((file) => {
        if (!this.attachments.some(f => f.name === file.name && f.size === file.size)) {
          this.attachments.push(file);
        }
      });
    },
    removeAttachment(index) {
      this.attachments.splice(index, 1);
    },
    formatFileSize(bytes) {
      if (bytes === 0) return '0 Bytes';
      const k = 1024;
      const sizes = ['Bytes', 'KB', 'MB', 'GB'];
      const i = Math.floor(Math.log(bytes) / Math.log(k));
      return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    },
    async sendMail() {
      if (!this.canSend) return;

      this.loading = true;

      const formData = new FormData();
      formData.append("fromParty", "org:holodeckb2b:example:company:A");
      formData.append("toParty", this.toParty);
      formData.append("service", "Examples");
      formData.append("action", "StoreMessage");
      formData.append("messageId", `msg-${Date.now()}`);
      formData.append("payload", this.message);
      formData.append("subject", this.subject);

      this.attachments.forEach((file) => {
        formData.append("attachments", file, file.name);
      });

      try {
        const response = await fetch("http://localhost:9091/api/as4/send", {
          method: "POST",
          body: formData,
        });

        if (!response.ok) {
          const errorText = await response.text();
          throw new Error(errorText || "Failed to send message");
        }
        
        alert('Message sent successfully!');
        this.resetForm();
      } catch (error) {
        console.error('Error sending message:', error);
        alert(error.message || 'Failed to send message');
      } finally {
        this.loading = false;
      }
    },
    resetForm() {
      this.toParty = "";
      this.subject = "";
      this.message = "";
      this.attachments = [];
    },
  },
};
</script>


<style scoped>
/* Optional: Add some transitions for better UX */
button {
  transition: background-color 0.2s ease;
}
</style>