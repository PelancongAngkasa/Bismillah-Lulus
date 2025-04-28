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
        <p class="text-gray-600 font-bold">Subject: {{ mail.subject }}</p> <!-- Tambahkan ini -->
      </div>

      <!-- Body -->
      <div class="mb-6">
        <p class="text-lg text-gray-800 mb-4">{{ mail.message }}</p>
        <a
          v-if="mail.attachment"
          :href="mail.attachment"
          target="_blank"
          class="text-blue-500 underline"
        >
          Download Attachment
        </a>
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
    };
  },
  methods: {
    fetchMail(mailId) {
      fetch(`http://localhost:9092/api/mail?id=${mailId}`)
        .then((response) => {
          if (!response.ok) {
            throw new Error(`Failed to fetch mail details: ${response.statusText}`);
          }
          return response.json();
        })
        .then((data) => {
          this.mail = {
            id: data.id,
            sender: data.sender || "Unknown Sender",
            receiver: data.receiver || "Unknown Receiver",
            date: new Date(data.date).toLocaleString(),
            subject: data.subject || "No subject available", // Tambahkan ini
            message: data.content || "No message content available",
            attachment: data.attachment || null,
          };
        })
        .catch((error) => {
          console.error("Error fetching mail:", error);
          alert("Failed to load mail details. Please try again later.");
        });
    },
    replyMail() {
      this.$router.push("/compose");
    },
    deleteMail() {
      alert("Mail Deleted!");
    },
  },
  created() {
    const mailId = this.$route.params.id;
    this.fetchMail(mailId);
  },
};
</script>