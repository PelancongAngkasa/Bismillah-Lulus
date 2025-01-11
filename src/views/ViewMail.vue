<template>
    <div class="flex h-screen">
      <!-- Sidebar -->
      <Sidebar />
  
      <!-- Mail Content -->
      <div class="flex-1 p-6 bg-gray-50">
        <!-- Header -->
        <div class="border-b pb-4 mb-4">
          <div class="flex justify-between items-center">
            <h1 class="text-2xl font-bold">{{ mail.subject }}</h1>
            <div class="text-gray-500">{{ mail.date }}</div>
          </div>
          <p class="text-gray-600">From: {{ mail.from }}</p>
        </div>
  
        <!-- Body -->
        <div class="mb-6">
          <p class="text-lg text-gray-800 mb-4">{{ mail.message }}</p>
          <img
            v-if="mail.image"
            :src="mail.image"
            alt="Attachment"
            class="rounded shadow w-1/2"
          />
        </div>
  
        <!-- Action Buttons -->
        <div class="flex space-x-4">
          <Button color="bg-blue-500" @click="replyMail">Reply</Button>
          <Button color="bg-red-500" @click="deleteMail">Delete</Button>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  import Sidebar from "@/components/Organisms/Sidebar.vue";
  import Button from "@/components/Atom/Button.vue";
  
  export default {
    components: { Sidebar, Button },
    data() {
      return {
        mail: {
          id: this.$route.params.id, // Ambil ID dari URL
          from: "Acme Inc.",
          subject: "Check This Logo",
          date: "Feb 26, 2024",
          message: "Please check the attached logo for your reference.",
          image: "/path-to-image/logo-example.jpg", // Path gambar
        },
      };
    },
    methods: {
      replyMail() {
        this.$router.push("/compose");
      },
      deleteMail() {
        alert("Mail Deleted!");
      },
    },
  };
  </script>
  