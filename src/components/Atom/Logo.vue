<template>
    <img :src="src" :alt="alt" class="h-12 w-auto" />
  </template>
  
  <script>
  export default {
    props: {
      src: { type: String, required: true },
      alt: { type: String, default: "Logo" }
    }
  };
  </script>
  