<template>
  <div v-if="mail" class="flex h-screen">
    <!-- Sidebar -->
    <Sidebar />

    <!-- Mail Content -->
    <div class="flex-1 p-6 bg-gray-50">
      <!-- Header -->
      <div class="border-b pb-4 mb-4">
        <div class="flex justify-between items-center">
          <div class="text-gray-500">{{ mail.date }}</div>
        </div>
        <p class="text-gray-600">From: {{ mail.sender }}</p>
        <p class="text-gray-600">To: {{ mail.receiver }}</p>
        <p class="text-gray-600 font-bold">Subject: {{ mail.subject }}</p>
      </div>

      <!-- Body -->
      <div class="mb-6">
        <p class="text-lg text-gray-800 mb-4">{{ mail.message }}</p>

        <!-- Attachments Table -->
        <div v-if="mail.attachments && mail.attachments.length > 0" class="mt-6">
          <h3 class="text-lg font-semibold mb-2">Attachments:</h3>
          <table class="min-w-full bg-white border border-gray-200 rounded-lg">
            <thead>
              <tr class="bg-gray-100 text-left text-sm font-medium text-gray-700">
                <th class="p-2 border">File Name</th>
                <th class="p-2 border">Type</th>
                <th class="p-2 border">Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="attachment in mail.attachments" :key="attachment.path" class="text-sm text-gray-700">
                <td class="p-2 border">{{ attachment.name }}</td>
                <td class="p-2 border">{{ attachment.mimeType }}</td>
                <td class="p-2 border space-x-2">
                  <button v-if="canPreview(attachment.mimeType)"
                          @click="openPreview(attachment)"
                          class="text-blue-500 hover:text-blue-700">
                    Preview
                  </button>
                  <a :href="`http://localhost:8083/download?id=${attachment.path}&name=${attachment.name}`"
                     class="text-green-500 hover:text-green-700"
                     download>
                    Download
                  </a>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- File Preview Modal -->
        <div v-if="showPreview" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4">
          <div class="bg-white rounded-lg p-4 max-w-4xl w-full">
            <div class="flex justify-between mb-4">
              <h3 class="text-lg font-semibold">{{ previewFile.name }}</h3>
              <button @click="showPreview = false" class="text-gray-500 hover:text-gray-700">
                Close
              </button>
            </div>
            <!-- Displaying the file based on its type -->
            <div class="max-h-[80vh] overflow-auto">
              <img v-if="isImage(previewFile.mimeType)" :src="getFileUrl(previewFile)" class="max-w-full">
              <iframe v-else-if="isPDF(previewFile.mimeType)" 
                      :src="getFileUrl(previewFile)"
                      class="w-full h-[70vh]">
              </iframe>
            </div>
          </div>
        </div>
      </div>

      <!-- Action Buttons -->
      <div class="flex space-x-4">
        <Button color="bg-blue-500" @click="replyMail">Reply</Button>
        <Button color="bg-red-500" @click="deleteMail">Delete</Button>
      </div>
    </div>
  </div>

  <div v-else class="flex h-screen items-center justify-center">
    <p>Loading mail...</p>
  </div>
</template>


<script>
import Sidebar from "@/components/Organisms/Sidebar.vue";
import Button from "@/components/Atom/Button.vue";

export default {
  components: { Sidebar, Button },
  data() {
    return {
      mail: null,
      showPreview: false,
      previewFile: null,
    };
  },

  methods: {
    fetchMail(mailId) {
      fetch(`http://localhost:8083/api/mail?id=${mailId}`)
        .then((response) => {
          if (!response.ok) {
            throw new Error(`Failed to fetch mail details: ${response.statusText}`);
          }
          return response.json();
        })
        .then((data) => {
          // Membuat array attachments jika ada attachment
          let attachments = [];
          if (data.attachments && Array.isArray(data.attachments)) {
            attachments = data.attachments.map(att => ({
              name: att.fileName,
              path: att.fileName,
              mimeType: att.mimeType,
              url: att.url
            }));
          }

          this.mail = {
            id: data.id,
            sender: data.sender || "Unknown Sender",
            receiver: data.receiver || "Unknown Receiver",
            date: new Date(data.date).toLocaleString(),
            subject: data.subject || "No subject available",
            message: data.content || "No message content available",
            attachments: attachments,
          };
        })
        .catch((error) => {
          console.error("Error fetching mail:", error);
          alert("Failed to load mail details. Please try again later.");
        });
    },

    // Tambahkan method baru untuk menentukan MIME type
    getMimeType(filename) {
      if (filename.toLowerCase().endsWith('.pdf')) {
        return 'application/pdf';
      } else if (filename.toLowerCase().match(/\.(jpg|jpeg|png|gif)$/)) {
        return 'image/' + filename.toLowerCase().split('.').pop();
      }
      return 'application/octet-stream';
    },
    replyMail() {
      this.$router.push("/compose");
    },
    deleteMail() {
      alert("Mail Deleted!");
    },

    getFileIcon(mimeType) {
      if (mimeType.startsWith('image/')) return 'fas fa-image';
      if (mimeType === 'application/pdf') return 'fas fa-file-pdf';
      return 'fas fa-file';
    },

    canPreview(mimeType) {
      return this.isImage(mimeType) || this.isPDF(mimeType);
    },

    isImage(mimeType) {
      return mimeType.startsWith('image/');
    },

    isPDF(mimeType) {
      return mimeType === 'application/pdf';
    },

    openPreview(file) {
      this.previewFile = file;
      this.showPreview = true;
    },

    getFileUrl(file) {
      return `http://localhost:8083/attachments/${file.path}`;
    }
  },
  created() {
    const mailId = this.$route.params.id;
    this.fetchMail(mailId);
  },
};
</script>