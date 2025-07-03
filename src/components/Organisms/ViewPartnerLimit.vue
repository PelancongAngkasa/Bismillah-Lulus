<template>
  <Card class="w-full max-w-7xl mx-4 mt-10">
    <h3 class="text-lg font-bold mb-2">Daftar Partner</h3>
    <table class="min-w-full border text-sm">
      <thead>
        <tr class="bg-gray-100">
          <th class="border px-2 py-1">Party ID</th>
          <th class="border px-2 py-1">Name</th>
          <th class="border px-2 py-1">Endpoint URL</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="p in paginatedPartners" :key="p.partyid">
          <td class="border px-2 py-1">{{ p.partyid }}</td>
          <td class="border px-2 py-1">{{ p.name }}</td>
          <td class="border px-2 py-1">{{ p.endpoint_url }}</td>
        </tr>
      </tbody>
    </table>
    <!-- Pagination Controls -->
    <div class="flex justify-center mt-4 gap-2">
      <button
        @click="currentPage--"
        :disabled="currentPage === 1"
        class="px-3 py-1 rounded bg-gray-200 disabled:opacity-50"
      >Prev</button>
      <span>Page {{ currentPage }} / {{ totalPages }}</span>
      <button
        @click="currentPage++"
        :disabled="currentPage === totalPages"
        class="px-3 py-1 rounded bg-gray-200 disabled:opacity-50"
      >Next</button>
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