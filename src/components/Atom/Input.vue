<template>
    <input
      :type="type"
      :placeholder="placeholder"
      class="w-full px-4 py-2 border rounded focus:outline-none"
    />
  </template>
  
  <script>
  export default {
    props: {
      type: { type: String, default: 'text' },
      placeholder: { type: String, default: '' }
    }
  };
  </script>
  