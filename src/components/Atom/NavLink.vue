<template>
    <a :href="href" class="text-gray-700 hover:text-blue-500 px-3 py-2">{{ label }}</a>
  </template>
  
  <script>
  export default {
    props: {
      href: { type: String, required: true },
      label: { type: String, required: true }
    }
  };
  </script>
  