<template>
    <div>
      <h1>Create Docker Container</h1>
      <form @submit.prevent="createContainer">
        <div>
          <label for="name">Name</label>
          <input v-model="containerData.name" id="name" required>
        </div>
        <div>
          <label for="image">Image</label>
          <input v-model="containerData.image" id="image" required>
        </div>
        <button type="submit">Create</button>
      </form>
    </div>
  </template>
  
  <script setup>
  import { ref } from 'vue';
  import { useRouter } from 'vue-router';
  
  const router = useRouter();
  const containerData = ref({ name: '', image: '' });
  
  const createContainer = async () => {
    await $fetch('/api/containers/create', {
      method: 'POST',
      body: JSON.stringify(containerData.value),
      headers: { 'Content-Type': 'application/json' },
    });
    router.push('/containers');
  };
  </script>
  