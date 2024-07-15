<template>
  <UPageHeader
    title="Other Containers"
    class="sticky top-0 bg-white z-10 dark:bg-gray-800 mb-4 pb-4"
  />

  <UAccordion
    multiple
    color="white"
    variant="solid"
    size="lg"
    :items="containers"
    v-if="containers.length"
  >
    <template #item="{ item }">
      <UPageGrid class="py-4">
      <UCard v-for="container in item.containers" :key="container.getId()">
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
            @click="() => performContainerAction(container, 'start')"
            v-if="!container.isRunning()"
            :icon="
              container.isLoading('start')
                ? 'i-heroicons-spinner'
                : 'i-heroicons-play'
            "
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
            :icon="
              container.isLoading('restart')
                ? 'i-heroicons-spinner'
                : 'i-heroicons-arrow-path'
            "
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
            :icon="
              container.isLoading('stop')
                ? 'i-heroicons-spinner'
                : 'i-heroicons-stop'
            "
            size="sm"
            color="white"
            square
            variant="solid"
            :loading="container.isLoading('stop')"
            label="Stop"
          />
          <UButton
            @click="() => performContainerAction(container, 'remove')"
            :icon="
              container.isLoading('remove')
                ? 'i-heroicons-spinner'
                : 'i-heroicons-trash'
            "
            size="sm"
            color="white"
            square
            variant="solid"
            :loading="container.isLoading('remove')"
            label="Remove"
          />
          <UButton
            @click="() => performContainerAction(container, 'rebuild')"
            :icon="
              container.isLoading('rebuild')
                ? 'i-heroicons-spinner'
                : 'i-heroicons-arrow-path-rounded-square'
            "
            size="sm"
            color="white"
            square
            variant="solid"
            :loading="container.isLoading('rebuild')"
            label="Rebuild"
          />
          <UButton
            :to="`/containers/${container.getId()}`"
            icon="i-heroicons-eye"
            size="sm"
            color="white"
            square
            variant="solid"
            label="Details"
          />
        </div>
      </UCard>
      </UPageGrid>
    </template>
  </UAccordion>
  <div v-else>
    <UButton size="xl" color="white" square block variant="ghost" loading>
      Loading containers...
    </UButton>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import Container from "@/models/Container";

const containers = ref([]);

const loadContainers = async () => {
  try {
    containers.value = await Container.allOthers();
    console.log("Loaded containers:", containers.value);
  } catch (error) {
    console.error("Failed to load containers:", error);
  }
};

const performContainerAction = async (container, action) => {
  try {
    await container[action]();
    loadContainers();
  } catch (error) {
    console.error(`Failed to ${action} container:`, error);
  }
};

onMounted(loadContainers);

// Watch and update the containers ref
watch(containers, (newContainers) => {
  containers.value = newContainers;
});
</script>
