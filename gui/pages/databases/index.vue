<template>
  <UPageHeader
    title="Databases"
    class="sticky top-0 bg-white z-10 dark:bg-gray-800 mb-4 pb-4"
  >
    <template #description>
      <div class="flex flex-wrap gap-2">
        <UButton
          @click="() => console.log('New Database')"
          icon="i-heroicons-plus"
          size="sm"
          color="white"
          square
          variant="solid"
          label="New Database"
        />
      </div>
    </template>
  </UPageHeader>
  <UPageGrid v-if="databases && !loading">
    <UPageCard v-for="database in databases" :key="database.getName()">
      <template #header>
        {{ database.getName() }} -
        <UBadge
          :label="database.getType().label"
          :color="database.getType().color"
          class="ml-1 h-6"
          variant="soft"
        />
      </template>
      <div class="flex flex-wrap gap-2">
        <UButton
          :to="`/databases/${database.getTypeUrl()}/${database.getName()}`"
          icon="i-heroicons-eye"
          size="sm"
          color="white"
          square
          variant="solid"
          label="Details"
        />
        <UButton
          v-if="database.getProject()"
          :to="`/projects/${database.getProject().getName()}`"
          icon="i-heroicons-folder-open"
          size="sm"
          color="white"
          square
          variant="solid"
          label="Project"
        />
        <UButton
          @click="() => performDatabaseAction(database, 'remove')"
          :icon="database.isLoading('remove') ? 'i-heroicons-spinner' : 'i-heroicons-trash'"
          size="sm"
          color="white"
          square
          variant="solid"
          :loading="database.isLoading('remove')"
          label="Remove"
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
      Loading databases...
    </UButton>
  </div>
</template>

<script setup>
import { useDatabasesStore } from '~/stores/databases';

const databasesStore = useDatabasesStore();

onMounted(() => {
  databasesStore.loadDatabases();
});

const databases = computed(() => databasesStore.databases);
const loading = computed(() => databasesStore.loading);

const performDatabaseAction = async (database, action) => {
  await databasesStore.performDatabaseAction(database, action);
};
</script>
