<template>
  <div class="flex h-screen">
    <Sidebar />
    <div class="flex-1 bg-gray-100 px-4 md:px-12 py-8 overflow-auto">
      <Card class="w-full max-w-7xl mx-auto">
        <div class="flex justify-between items-center mb-6">
          <h2 class="text-3xl font-bold">Inbox</h2>
          <input
            type="text"
            placeholder="Search mail"
            v-model="searchQuery"
            class="px-6 py-3 border rounded text-lg focus:outline-none"
            style="min-width: 300px"
          />
        </div>
        <table class="min-w-full text-lg">
          <thead>
            <tr class="bg-gray-100">
              <th class="px-6 py-3 border-transparent">Subject</th>
              <th class="px-6 py-3 border-transparent">Sender</th>
              <th class="px-6 py-3 border-transparent">Date</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="mail in paginatedMails"
              :key="mail.id"
              class="hover:bg-blue-50 cursor-pointer"
              @click="viewMail(mail.id)"
            >
              <td class="px-6 py-3 border-transparent font-semibold">{{ mail.subject }}</td>
              <td class="px-6 py-3 border-transparent">{{ mail.sender }}</td>
              <td class="px-6 py-3 border-transparent">{{ mail.date }}</td>
            </tr>
            <tr v-if="paginatedMails.length === 0">
              <td colspan="3" class="text-center py-6 text-gray-400 text-lg">No mail found.</td>
            </tr>
          </tbody>
        </table>
        <!-- Pagination Controls -->
        <div class="flex justify-center gap-4 pt-6" v-if="totalPages > 1">
          <button
            @click="currentPage--"
            :disabled="currentPage === 1"
            class="px-5 py-2 rounded bg-gray-200 text-lg disabled:opacity-50"
          >Prev</button>
          <span class="text-lg">Page {{ currentPage }} / {{ totalPages }}</span>
          <button
            @click="currentPage++"
            :disabled="currentPage === totalPages"
            class="px-5 py-2 rounded bg-gray-200 text-lg disabled:opacity-50"
          >Next</button>
        </div>
      </Card>
    </div>
  </div>
</template>

<script>
import Sidebar from "@/components/Organisms/Sidebar.vue";
import Card from "@/components/Molecules/Card.vue";

export default {
  components: { Sidebar, Card },
  data() {
    return {
      mails: [],
      searchQuery: "",
      currentPage: 1,
      pageSize: 25,
    };
  },
  computed: {
    filteredMails() {
      if (!this.searchQuery) return this.mails;
      const q = this.searchQuery.toLowerCase();
      return this.mails.filter(
        mail =>
          (mail.subject && mail.subject.toLowerCase().includes(q)) ||
          (mail.sender && mail.sender.toLowerCase().includes(q))
      );
    },
    totalPages() {
      return Math.ceil(this.filteredMails.length / this.pageSize) || 1;
    },
    paginatedMails() {
      const start = (this.currentPage - 1) * this.pageSize;
      return this.filteredMails.slice(start, start + this.pageSize);
    },
  },
  created() {
    fetch("http://localhost:8082/api/mails")
      .then((response) => response.json())
      .then((data) => {
        this.mails = data;
      })
      .catch((error) => console.error("Error fetching mails:", error));
  },
  methods: {
    viewMail(mailId) {
      this.$router.push(`/view/${mailId}`);
    },
  },
  watch: {
    filteredMails() {
      if (this.currentPage > this.totalPages) this.currentPage = this.totalPages;
    }
  }
};
</script>