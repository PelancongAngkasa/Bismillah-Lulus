<template>
  <div class="flex h-screen">
    <!-- Sidebar -->
    <Sidebar />

    <!-- Mail Content -->
    <div class="flex-1">
      <!-- Search Bar -->
      <div class="p-4 border-b">
        <input 
          type="text" 
          placeholder="Search mail" 
          class="w-full px-4 py-2 border rounded focus:outline-none" />
      </div>
      <!-- Mail List -->
      <MailList :mails="mails" @view-mail="viewMail" />
    </div>
  </div>
</template>

<script>
import Sidebar from "@/components/Organisms/Sidebar.vue";
import MailList from "@/components/Molecules/MailList.vue";

export default {
  components: { Sidebar, MailList },
  methods: {
  viewMail(mailId) {
    this.$router.push(`/view/${mailId}`);
    },
  },
  data() {
    return {
      mails: [
        { sender: "Acme Inc.", subject: "Insights: Latest equipment", date: "Feb, 26" },
        { sender: "Travel Tales", subject: "Adventures & Destinations", date: "March, 26" },
        { sender: "Delta Enterprises", subject: "Delta Weekly News", date: "Jan, 26" },
      ],
    };
  },
};
</script>
