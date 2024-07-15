<template>
  <UPageHeader
    title="Containers"
    class="sticky top-0 bg-white z-10 dark:bg-gray-800 mb-4 pb-4"
  />

  <UPageGrid v-if="containers.length && !loading">
    <UPageCard v-for="container in containers" :key="container.getId()">
      <template #header>
        {{ container.getName() }} -
        <UBadge 
          :label="container.getStatus()"
          :color="container.isRunning() ? 'green' : 'red'"
          variant="soft"
        />
      </template>

      <div class="flex flex-wrap gap-2">
        <UButton
          :to="`/containers/${container.getId()}`"
          icon="i-heroicons-eye"
          size="sm"
          color="white"
          square
          variant="solid"
          label="Details"
        />
        <UButton
            icon="i-heroicons-envelope"
            size="sm"
            color="white"
            square
            variant="solid"
            label="Go to Mailpit"
            to="https://mailpit.localhost"
            v-if="container.getName().toLowerCase() === 'mailpit'"
          />
          <UButton
            icon="i-heroicons-circle-stack"
            size="sm"
            color="white"
            square
            variant="solid"
            label="Go to Adminer"
            to="https://adminer.localhost"
            v-if="container.getName().toLowerCase() === 'adminer'"
          />
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
    </UPageCard>
  </UPageGrid>
  <div v-else>
    <UButton 
      size="xl"
      color="white"
      square
      block
      variant="ghost"
      loading>
        Loading containers...
    </UButton>
  </div>
</template>

<script setup>
import Container from "@/models/Container";

const containers = ref([]);
const loading = ref(true);

const loadContainers = async () => {
  try {
    containers.value = (await Container.all()).sort((a, b) => a.getName().localeCompare(b.getName()));
    loading.value = false;
  } catch (error) {
    console.error("Failed to load containers:", error);
  }
};

const performContainerAction = async (container, action) => {
  try {
    await container[action]();
    loading.value = true;
    loadContainers();
    loading.value = false;
  } catch (error) {
    console.error(`Failed to ${action} container:`, error);
  }
};

onMounted(loadContainers);

// Watch and update the containers ref
watch(containers, (newContainers) => {
  containers.value = newContainers;
});

watch(loading, (newLoading) => {
  loading.value = newLoading;
});
</script>
