<template>
  <Card class="mx-auto my-10 max-w-7xl transition-all duration-300 px-6">
    <h3 class="text-lg font-bold mb-4">Daftar Partner</h3>
    <div class="overflow-x-auto">
      <table class="w-full border text-sm">
        <thead>
          <tr class="bg-gray-100">
            <th class="border px-4 py-2">ID Partner</th>
            <th class="border px-4 py-2">Nama Partner</th>
            <th class="border px-4 py-2">Alamat URL</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="p in paginatedPartners" :key="p.partyid" class="hover:bg-gray-50">
            <td class="border px-4 py-2">{{ p.partyid }}</td>
            <td class="border px-4 py-2">{{ p.name }}</td>
            <td class="border px-4 py-2">{{ p.endpoint_url }}</td>
          </tr>
        </tbody>
      </table>
    </div>
    <!-- Pagination Controls -->
    <div class="flex justify-center mt-6 gap-2">
      <button
        @click="currentPage--"
        :disabled="currentPage === 1"
        class="px-4 py-2 rounded bg-gray-200 hover:bg-gray-300 disabled:opacity-50 disabled:hover:bg-gray-200"
      >
        Prev
      </button>
      <span class="flex items-center">Page {{ currentPage }} / {{ totalPages }}</span>
      <button
        @click="currentPage++"
        :disabled="currentPage === totalPages"
        class="px-4 py-2 rounded bg-gray-200 hover:bg-gray-300 disabled:opacity-50 disabled:hover:bg-gray-200"
      >
        Next
      </button>
    </div>
  </Card>
</template>

<script>
import Card from '@/components/Molecules/Card.vue'

export default {
  name: 'ViewPartnerUser',
  components: { Card },
  data() {
    return {
      partners: [],
      currentPage: 1,
      pageSize: 5
    }
  },
  computed: {
    totalPages() {
      return Math.ceil(this.partners.length / this.pageSize) || 1
    },
    paginatedPartners() {
      const start = (this.currentPage - 1) * this.pageSize
      return this.partners.slice(start, start + this.pageSize)
    }
  },
  mounted() {
    this.fetchPartners()
  },
  methods: {
    async fetchPartners() {
      const res = await fetch('http://localhost:8081/api/partner')
      this.partners = await res.json()
      if (this.currentPage > this.totalPages) this.currentPage = this.totalPages
    }
  }
}
</script>