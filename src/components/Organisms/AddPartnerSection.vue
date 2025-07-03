<template>
  <div class="flex justify-center pt-10">
    <Card
      title="Tambah Partner Baru"
      class="w-full max-w-7xl mx-4"
    >
      <form @submit.prevent="onSubmit" class="space-y-4">
        <div class="flex flex-col gap-4">
          <input
            v-model="form.partyid"
            placeholder="Party ID"
            required
            class="p-3 border rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <input
            v-model="form.name"
            placeholder="Name"
            required
            class="p-3 border rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <input
            v-model="form.endpoint_url"
            placeholder="Endpoint URL"
            required
            class="p-3 border rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <Button type="submit" class="w-full mt-4">Tambah Partner</Button>
        <p v-if="success" class="text-green-600 mt-2">Partner berhasil ditambahkan!</p>
        <p v-if="error" class="text-red-600 mt-2">{{ error }}</p>
      </form>
    </Card>
  </div>
</template>

<script>
import Card from '@/components/Molecules/Card.vue'
import Button from '@/components/Atom/Button.vue'

export default {
  name: 'AddPartnerSection',
  components: { Card, Button },
  data() {
    return {
      form: {
        partyid: '',
        name: '',
        endpoint_url: ''
      },
      success: false,
      error: ''
    }
  },
  methods: {
    async onSubmit() {
      this.success = false
      this.error = ''

      try {
        const response = await fetch('http://localhost:8081/api/partner', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(this.form)
        })

        if (!response.ok) {
          const errText = await response.text()
          throw new Error(errText || 'Gagal menambah partner')
        }

        this.success = true
        this.form = { partyid: '', name: '', endpoint_url: '' }
      } catch (e) {
        this.error = e.message
      }
    }
  }
}
</script>