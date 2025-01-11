<template>
    <button 
      :class="`px-4 py-2 rounded text-white ${color}`"
      @click="onClick">
      <slot />
    </button>
  </template>
  
  <script>
  export default {
    props: {
      color: { type: String, default: 'bg-blue-500' },
      onClick: { type: Function, default: () => {} }
    }
  };
  </script>
  