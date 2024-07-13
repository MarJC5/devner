<template>
  <UPageHeader
    :headline="container ? '' : 'Loading...'"
    class="sticky top-0 bg-white z-10 dark:bg-gray-800 pb-4"
  >
    <template #title>
      <h1 v-if="container" class="flex">
        {{ container.getName() }}
        <UBadge
          :label="container.getStatus()"
          :color="container.isRunning() ? 'green' : 'red'"
          class="ml-2 h-6"
          variant="soft"
        />
      </h1>
    </template>
    <template #description>
      <div class="flex flex-wrap gap-2" v-if="container">
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
    <UDashboardToolbar class="pl-0 py-4">
      <template #left>
        <USelectMenu 
          :icon="selected === 'General' ? 'i-heroicons-information-circle' : 'i-heroicons-clipboard-document-list'"
          v-model="selected" 
          :options="status" 
          color="white"
          v-if="selected !== 'Shell'"/>
        <UButton
          icon="i-heroicons-command-line"
          size="sm"
          color="white"
          square
          variant="solid"
          label="Terminal"
          v-if="selected !== 'Shell'"
          @click="enableShell()"
        />
        <UButton
          icon="i-heroicons-arrow-right-on-rectangle"
          size="sm"
          color="white"
          square
          variant="solid"
          label="Quit"
          v-if="selected === 'Shell'"
          @click="leaveShell()"
        />
      </template>
    </UDashboardToolbar>
        <div v-if="selected === 'General'"
          class="mt-4">
          <UPageGrid
            :ui="{
              wrapper: 'grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-2 gap-8',
            }"
          >
            <UPageCard title="General">
              <p>Name: {{ container.getName() }}</p>
              <p class="truncate">Image: {{ container.getImage() }}</p>
              <p>Created: {{ container.getCreated() }}</p>
            </UPageCard>
            <UPageCard title="Network">
              Ports:
              <ul>
                <li v-for="port in container.getPorts()" :key="port">
                  {{ port }}
                </li>
              </ul>
            </UPageCard>
          </UPageGrid>
        </div>
        <div v-if="selected === 'Logs'"
          class="mt-4">
          <LogViewer :containerId="container.getId()" />
        </div>
        <div v-if="selected === 'Shell'"
          class="mt-4">
          <LogTerminal :containerId="container.getId()" v-if="selected === 'Shell'"/>
        </div>
  </div>
  <div v-else>
    <UButton size="xl" color="white" square block variant="ghost" loading>
      Loading container...
    </UButton>
  </div>
</template>

<script setup>
import Container from "@/models/Container";
import LogViewer from "@/components/LogViewer.vue";
import LogTerminal from "@/components/LogTerminal.vue";

const route = useRoute();
const container = ref(null);
const status = ['General', 'Logs']

const selected = ref(status[0])

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

const enableShell = () => {
  selected.value = 'Shell'
}

const leaveShell = () => {
  selected.value = 'General'
}

onMounted(loadContainer);

// Watch and update the containers ref
watch(container, (newContainer) => {
  container.value = newContainer;
});
</script>
