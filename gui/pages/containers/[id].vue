<template>
  <UPageHeader
    :headline="container ? '' : 'Loading...'"
    class="mb-4"
  >
    <template #title>
        <h1 v-if="container" class="flex">{{ container.getName() }}
        <UBadge 
          :label="container.getStatus()"
          :color="container.isRunning() ? 'green' : 'red'"
          class="ml-2 h-6"
          variant="soft"/>
        </h1>
    </template>
    <template #description>
      <div class="flex flex-wrap gap-4" v-if="container">
        <UButton
          @click="() => performContainerAction(container, 'start')"
          v-if="!container.isRunning()"
          :icon="container.isLoading('start') ? 'i-heroicons-spinner' : 'i-heroicons-play'"
          size="sm"
          color="white"
          square
          variant="solid"
          :loading="container.isLoading('start')"
          label="Start"
        />
        <UButton
          @click="() => performContainerAction(container, 'restart')"
          v-if="container.isRunning()"
          :icon="container.isLoading('restart') ? 'i-heroicons-spinner' : 'i-heroicons-arrow-path'"
          size="sm"
          color="white"
          square
          variant="solid"
          :loading="container.isLoading('restart')"
          label="Restart"
        />
        <UButton
          @click="() => performContainerAction(container, 'stop')"
          v-if="container.isRunning()"
          :icon="container.isLoading('stop') ? 'i-heroicons-spinner' : 'i-heroicons-stop'"
          size="sm"
          color="white"
          square
          variant="solid"
          :loading="container.isLoading('stop')"
          label="Stop"
        />
        <UButton
          @click="() => performContainerAction(container, 'remove')"
          :icon="container.isLoading('remove') ? 'i-heroicons-spinner' : 'i-heroicons-trash'"
          size="sm"
          color="white"
          square
          variant="solid"
          :loading="container.isLoading('remove')"
          label="Remove"
        />
        <UButton
          @click="() => performContainerAction(container, 'rebuild')"
          :icon="container.isLoading('rebuild') ? 'i-heroicons-spinner' : 'i-heroicons-arrow-path-rounded-square'"
          size="sm"
          color="white"
          square
          variant="solid"
          :loading="container.isLoading('rebuild')"
          label="Rebuild"
        />
      </div>
    </template>
  </UPageHeader>
  <div v-if="container">
    <div>
      <p>Image: {{ container.getImage() }}</p>
      <p>Created: {{ container.getCreated() }}</p>
      <p>
        Ports:
        <ul>
          <li v-for="port in container.getPorts()" :key="port">
            {{ port }}
          </li>
        </ul>
      </p>
    </div>
    <div>
      <LogViewer :containerId="container.getId()" />
    </div>
  </div>
  <div v-else>
    <UButton 
      size="xl"
      color="white"
      square
      block
      variant="ghost"
      loading>
        Loading container...
    </UButton>
  </div>
</template>

<script setup>
import Container from "@/models/Container";
import LogViewer from '@/components/LogViewer.client.vue';

const route = useRoute();
const container = ref(null);

const loadContainer = async () => {
  try {
    container.value = await Container.fetchContainer(route.params.id);
  } catch (error) {
    console.error("Failed to load container details:", error);
  }
};

const performContainerAction = async (container, action) => {
  try {
    await container[action]();
    loadContainer();
  } catch (error) {
    console.error(`Failed to ${action} container:`, error);
  }
};

onMounted(loadContainer);

// Watch and update the containers ref
watch(container, (newContainer) => {
  container.value = newContainer;
});
</script>
