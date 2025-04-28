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

        <!-- Input untuk To -->
        <input
          v-model="subject"
          type="text"
          placeholder="subject"
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
          <div v-if="attachments.length" class="mt-4">
            <ul>
              <li
                v-for="(file, index) in attachments"
                :key="file.name"
                class="flex justify-between items-center mb-2"
              >
                <span class="truncate">{{ file.name }}</span>
                <button
                  class="text-red-500 hover:text-red-700"
                  @click="removeAttachment(index)"
                >
                  Remove
                </button>
              </li>
            </ul>
          </div>
        </div>

        <!-- Tombol Kirim -->
        <div class="flex justify-end">
          <button
            :disabled="loading"
            @click="sendMail"
            class="bg-blue-500 text-white px-4 py-2 rounded"
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
    toParty: "", // Penerima pesan
    subject: "", // Subjek pesan
    message: "", // Pesan utama
    attachments: [], // Daftar file lampiran
    loading: false, // Status loading
  };
},
methods: {
  triggerFileInput() {
    this.$refs.fileInput.click();
  },
  handleFileUpload(event) {
    const files = Array.from(event.target.files);
    files.forEach((file) => {
      if (!this.attachments.some((f) => f.name === file.name)) {
        this.attachments.push(file);
      }
    });
  },
  handleDrop(event) {
    const files = Array.from(event.dataTransfer.files);
    files.forEach((file) => {
      if (!this.attachments.some((f) => f.name === file.name)) {
          this.attachments.push(file);
        }
      });
    },
    removeAttachment(index) {
      this.attachments.splice(index, 1);
    },
  

    async sendMail() {
      if (!this.toParty || !this.message) {
        alert("To and Message fields are required.");
        return;
      }

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
        formData.append("attachments", file);
      });

      try {
        const response = await fetch("http://localhost:8081/api/as4/send", {
          method: "POST",
          body: formData,
        });

        if (!response.ok) throw new Error("Failed to send message");
          alert("Message sent successfully!");
          this.resetForm();
        } catch (error) {
          alert(`Error: ${error.message}`);
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
