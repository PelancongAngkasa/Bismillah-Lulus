<template>
  <ul>
    <li
      v-for="(item, index) in menuItems"
      :key="index"
      class="text-gray-300 hover:text-white mb-2 cursor-pointer"
    >
      <a :href="item.link">{{ item.name }}</a>
    </li>
  </ul>
</template>

<script>
export default {
  props: {
    menuItems: { type: Array, required: true },
  },
};
</script>
