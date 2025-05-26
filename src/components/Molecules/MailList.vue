<template>
  <div>
    <div
      v-for="(mail, index) in mails"
      :key="index"
      class="p-4 border-b cursor-pointer hover:bg-gray-100"
      @click="$emit('view-mail', mail.id)"
    >
      <h3 class="font-bold">{{ mail.subject || 'No Subject' }}</h3>
      <p class="text-gray-500">From: {{ mail.sender }} - {{ mail.date }}</p>
    </div>
  </div>
</template>

<script>
export default {
  props: { mails: Array },
};
</script>
