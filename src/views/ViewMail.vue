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

        <!-- Security Information Section -->
        <div v-if="mail.securityInfo" class="mt-6 mb-6">
          <h3 class="text-lg font-semibold mb-4 text-blue-600">
            <i class="fas fa-shield-alt mr-2"></i>
            Security Information
          </h3>
          <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
            <div class="space-y-2 text-sm">
              <div><span class="font-medium">Keystore Alias:</span> {{ mail.securityInfo.keystoreAlias || 'Not specified' }}</div>
              <div v-if="mail.securityInfo.dname">
                <span class="font-medium">Distinguished Name (DName):</span>
                <span>{{ mail.securityInfo.dname }}</span>
              </div>
              <div v-else class="text-gray-500 text-sm">
                <i class="fas fa-exclamation-triangle mr-1"></i>
                DName not available
              </div>
            </div>
          </div>
        </div>

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
                <td class="p-2 border flex items-center space-x-2">
                  <i :class="getFileIcon(attachment.mimeType)" class="text-gray-500"></i>
                  <span>{{ attachment.name }}</span>
                </td>
                <td class="p-2 border">{{ attachment.mimeType }}</td>
                <td class="p-2 border space-x-2">
                  <button v-if="canPreview(attachment.mimeType)"
                          @click="openPreview(attachment)"
                          class="text-blue-500 hover:text-blue-700">
                    Preview
                  </button>
                  <a :href="`http://localhost:8083/download?name=${attachment.path}&name=${attachment.name}`"
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
              <img v-if="isImage(previewFile.mimeType)" :src="getFileUrl(previewFile)" class="max-w-full mx-auto">
              <iframe v-else-if="isPDF(previewFile.mimeType)" 
                      :src="getFileUrl(previewFile)"
                      class="w-full h-[70vh]">
              </iframe>
              <p v-else class="text-gray-500 text-center">Preview not available for this file type.</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Action Buttons -->
      <div class="flex space-x-4">
        <Button color="bg-blue-500" @click="replyMail">Reply</Button>
        
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
      showCertInfo: false,
      certificateInfo: null,
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
          let attachments = [];
          if (data.attachments && Array.isArray(data.attachments)) {
            attachments = data.attachments.map(att => ({
              name: att.fileName || att.name || att,
              path: att.fileName || att.name || att,
              mimeType: att.mimeType || this.getMimeType(att.fileName || att.name || att),
              url: att.url || this.getFileUrl({ path: att.fileName || att.name || att })
            }));
          }

          this.mail = {
            id: data.id,
            sender: data.sender || "Unknown Sender",
            receiver: data.receiver || "Unknown Receiver",
            date: data.date ? new Date(data.date).toLocaleString() : "",
            subject: data.subject || "No subject available",
            message: data.content || "No message content available",
            attachments: attachments,
            securityInfo: data.securityInfo || null,
          };
        })
        .catch((error) => {
          console.error("Error fetching mail:", error);
          alert("Failed to load mail details. Please try again later.");
        });
    },

    copyCertificate() {
      if (this.mail.securityInfo && this.mail.securityInfo.publicKeyCertificate) {
        navigator.clipboard.writeText(this.mail.securityInfo.publicKeyCertificate)
          .then(() => {
            alert('Certificate copied to clipboard!');
          })
          .catch(() => {
            alert('Failed to copy certificate to clipboard');
          });
      }
    },

    showCertificateInfo() {
      this.showCertInfo = true;
      this.certificateInfo = null; // Clear previous info
      this.fetchCertificateInfo();
    },

         fetchCertificateInfo() {
       fetch(`http://localhost:8083/api/certificates`)
         .then((response) => {
           if (!response.ok) {
             throw new Error(`Failed to fetch certificate info: ${response.statusText}`);
           }
           return response.json();
         })
         .then((data) => {
           this.certificateInfo = data;
         })
         .catch((error) => {
           console.error("Error fetching certificate info:", error);
           alert("Failed to load certificate information. Please try again later.");
         });
     },

    formatFileSize(bytes) {
      if (bytes === 0) return '0 Bytes';
      const k = 1024;
      const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
      const i = Math.floor(Math.log(bytes) / Math.log(k));
      return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    },

    getMimeType(filename) {
      if (!filename) return 'application/octet-stream';
      const ext = filename.toLowerCase().split('.').pop();
      switch (ext) {
        case 'pdf':
          return 'application/pdf';
        case 'jpg':
        case 'jpeg':
        case 'png':
        case 'gif':
        case 'bmp':
        case 'webp':
          return 'image/' + ext;
        default:
          return 'application/octet-stream';
      }
    },

    replyMail() {
      this.$router.push("/compose");
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
    },
    // extractDName method dihapus karena backend sudah mengirim dname langsung
  },

  created() {
    const mailId = this.$route.params.id;
    this.fetchMail(mailId);
  },
};
</script>
