<template>
    <div class="p-6 border rounded shadow-md">
      <h2 class="font-bold mb-4">Compose Email</h2>
      <Input placeholder="To" class="mb-2" />
      <Input placeholder="CC" class="mb-2" />
      <Input placeholder="Subject" class="mb-2" />
      <textarea placeholder="Write your message" class="w-full p-2 border rounded h-32"></textarea>
      <div class="mt-4 flex justify-end">
        <Button color="bg-blue-500">Send</Button>
      </div>
    </div>
  </template>
  
  <script>
  import InputField from "../Atom/Input.vue";
  import Button from "../Atom/Button.vue";
  
  export default {
    components: { InputField, Button }
  };
  </script>
  