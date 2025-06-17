<template>
  <div :class="['bg-white rounded-lg shadow p-6', $attrs.class]">
    <h2 v-if="title" class="mb-4 text-xl font-bold text-gray-800">{{ title }}</h2>
    <slot />
  </div>
</template>

<script setup>
defineProps({
  title: String
})
</script>