<template>
    <div class="flex-1 p-4">
      <div class="flex justify-between items-center mb-4">
        <input 
          type="text" 
          placeholder="Search mail" 
          class="w-full px-4 py-2 border rounded focus:outline-none" />
      </div>
      <MailList :mails="mails" />
    </div>
  </template>
  
  <script>
  import MailList from "../Molecules/MailList.vue";
  
  export default {
    components: { MailList },
    data() {
      return {
        mails: [
          { sender: "Acme Inc.", subject: "Insights: The latest in industrial equipment", date: "Feb, 26" },
          { sender: "Travel Tales", subject: "Our latest Adventures and Destinations", date: "March, 26" },
          { sender: "Delta Enterprises", subject: "Delta Weekly News", date: "Jan, 26" }
        ]
      };
    }
  };
  </script>
  