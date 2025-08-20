<template>
  <div class="min-h-screen flex flex-col">
    <Navbar @toggle-sidebar="toggleSidebar" />
    
    <!-- Spacer untuk navbar fixed -->
    <div class="h-16"></div>
    
    <div class="flex flex-1">
      <!-- Sidebar -->
      <Sidebar
        v-show="isSidebarVisible"
        class="fixed top-16 left-0 h-[calc(100vh-64px)] w-64 bg-gray-800"
      />
      
      <!-- Main content area with scroll -->
      <div class="flex-1 bg-gray-100 p-6 overflow-auto ml-64">
        <Card class="w-full max-w-7xl mx-auto">
          <div class="flex justify-between items-center mb-6">
            <h2 class="text-3xl font-bold">Pesan Masuk</h2>
            <input
              type="text"
              placeholder="cari pesan..."
              v-model="searchQuery"
              class="px-6 py-3 border rounded text-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              style="min-width: 300px"
            />
          </div>

          <!-- Message Count Info -->
          <div class="text-sm text-gray-600 mb-4">
            Showing {{ paginationInfo.from }}-{{ paginationInfo.to }} of {{ filteredMails.length }} messages
          </div>

          <table class="min-w-full text-lg">
            <thead>
              <tr class="bg-gray-100">
                <th class="px-6 py-3 border-transparent text-left">Subjek</th>
                <th class="px-6 py-3 border-transparent text-left">Pengirim</th>
                <th class="px-6 py-3 border-transparent text-left">Tanggal</th>
              </tr>
            </thead>
            <tbody>
              <template v-if="mails.length > 0">
                <tr
                  v-for="mail in paginatedMails"
                  :key="mail.id"
                  class="hover:bg-blue-50 cursor-pointer border-b border-gray-100"
                  @click="viewMail(mail.id)"
                >
                  <td class="px-6 py-3 border-transparent font-semibold">{{ mail.subject }}</td>
                  <td class="px-6 py-3 border-transparent">{{ mail.sender }}</td>
                  <td class="px-6 py-3 border-transparent">{{ mail.date }}</td>
                </tr>
              </template>
              <tr v-else>
                <td colspan="3" class="text-center py-6 text-gray-400 text-lg">
                  Loading messages...
                </td>
              </tr>
              <tr v-if="mails.length > 0 && paginatedMails.length === 0">
                <td colspan="3" class="text-center py-6 text-gray-400 text-lg">
                  {{ searchQuery ? 'No messages found matching your search.' : 'No messages available.' }}
                </td>
              </tr>
            </tbody>
          </table>

          <!-- Enhanced Pagination Controls -->
          <div class="flex justify-between items-center mt-6" v-if="totalPages > 1">
            <div class="flex items-center gap-2">
              <button
                @click="currentPage = 1"
                :disabled="currentPage === 1"
                class="px-3 py-1 rounded bg-gray-200 text-sm disabled:opacity-50 hover:bg-gray-300"
              >
                First
              </button>
              <button
                @click="currentPage--"
                :disabled="currentPage === 1"
                class="px-4 py-2 rounded bg-gray-200 disabled:opacity-50 hover:bg-gray-300"
              >
                Prev
              </button>
            </div>

            <!-- Page Numbers -->
            <div class="flex gap-1">
              <button
                v-for="pageNum in displayedPages"
                :key="pageNum"
                @click="currentPage = pageNum"
                :class="[
                  'px-3 py-1 rounded text-sm',
                  pageNum === '...' ? 'cursor-default' : 'cursor-pointer',
                  currentPage === pageNum 
                    ? 'bg-blue-500 text-white' 
                    : pageNum === '...' 
                      ? 'bg-transparent'
                      : 'bg-gray-200 hover:bg-gray-300'
                ]"
                :disabled="pageNum === '...'"
              >
                {{ pageNum }}
              </button>
            </div>

            <div class="flex items-center gap-2">
              <button
                @click="currentPage++"
                :disabled="currentPage === totalPages"
                class="px-4 py-2 rounded bg-gray-200 disabled:opacity-50 hover:bg-gray-300"
              >
                Next
              </button>
              <button
                @click="currentPage = totalPages"
                :disabled="currentPage === totalPages"
                class="px-3 py-1 rounded bg-gray-200 text-sm disabled:opacity-50 hover:bg-gray-300"
              >
                Last
              </button>
            </div>
          </div>
        </Card>
      </div>
    </div>
  </div>
</template>

<script>
import Sidebar from "@/components/Organisms/Sidebar.vue";
import Card from "@/components/Molecules/Card.vue";
import Navbar from "@/components/Molecules/Navbar.vue";

export default {
  components: { Sidebar, Card, Navbar },
  data() {
    return {
      mails: [],
      searchQuery: "",
      currentPage: 1,
      pageSize: 10,
      isSidebarVisible: true
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
    paginatedMails() {
      const start = (this.currentPage - 1) * this.pageSize;
      const end = start + this.pageSize;
      console.log(`Showing mails from ${start} to ${end}`); // Debug info
      return this.filteredMails.slice(start, end);
    },
    totalPages() {
      const total = Math.ceil(this.filteredMails.length / this.pageSize);
      console.log("Total pages:", total); // Debug info
      return total || 1;
    },
    displayedPages() {
      const total = this.totalPages;
      const current = this.currentPage;
      const delta = 2;
      
      let pages = [];
      pages.push(1);
      
      for (let i = current - delta; i <= current + delta; i++) {
        if (i > 1 && i < total) {
          pages.push(i);
        }
      }
      
      if (total !== 1) {
        pages.push(total);
      }
      
      pages = pages.reduce((acc, curr, i, arr) => {
        if (i > 0) {
          if (curr - arr[i-1] > 1) {
            acc.push('...');
          }
        }
        acc.push(curr);
        return acc;
      }, []);
      
      return pages;
    },
    paginationInfo() {
      const total = this.filteredMails.length;
      const from = total === 0 ? 0 : (this.currentPage - 1) * this.pageSize + 1;
      const to = Math.min(this.currentPage * this.pageSize, total);
      
      return { from, to, total };
    }
  },
  created() {
  fetch("http://localhost:8082/api/mails")
    .then((response) => response.json())
    .then((data) => {
      // Pastikan data adalah array
      this.mails = Array.isArray(data) ? data : [];
      console.log("Total mails:", this.mails.length); // Debug info
    })
    .catch((error) => {
      console.error("Error fetching mails:", error);
      this.mails = [];
    });
  },
  methods: {
    viewMail(mailId) {
      this.$router.push(`/view/${mailId}`);
    },
    toggleSidebar() {
      this.isSidebarVisible = !this.isSidebarVisible;
    }
  },
  watch: {
    currentPage(newPage) {
      console.log("Current page changed to:", newPage);
    },
    filteredMails(newMails) {
      console.log("Filtered mails length:", newMails.length);
      if (this.currentPage > this.totalPages) {
        this.currentPage = this.totalPages;
      }
    },
    searchQuery() {
      this.currentPage = 1;
    }
  }
};
</script>