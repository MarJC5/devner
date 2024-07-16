<template>
  <UPageHeader
    :headline="database ? '' : 'Loading...'"
    class="mb-4 sticky top-0 bg-white z-10 dark:bg-gray-900 pb-4"
  >
    <template #title>
      <h1 v-if="database" class="flex">{{ database.getName() }}
        <UBadge 
          :label="database.getType().label"
          :color="database.getType().color"
          class="ml-2 h-6"
          variant="soft"
        />
      </h1>
    </template>
    <template #description>
      <div class="flex flex-wrap gap-2" v-if="database">
        <UButton
          v-if="database.getProject()"
          :to="`/projects/${database.getProject().getName()}`"
          icon="i-heroicons-folder-open"
          size="sm"
          color="white"
          square
          variant="solid"
          label="View Project"
        />
        <UButton
          icon="i-heroicons-circle-stack"
          size="sm"
          color="white"
          square
          variant="solid"
          label="Go to Adminer"
          :to="`https://adminer.localhost/?~/server=${database.getTypeUrl()}&username=${database.getName()}&db=${database.getName()}`"
        />
      </div>
    </template>
  </UPageHeader>
  <div v-if="database">
    <!-- Add more database-related content here -->
  </div>
  <div v-else>
    <UButton 
      size="xl"
      color="white"
      square
      block
      variant="ghost"
      loading>
      Loading database...
    </UButton>
  </div>
</template>

<script setup>
import { useDatabasesStore } from '~/stores/databases';

const route = useRoute();
const databasesStore = useDatabasesStore();

onMounted(() => {
  databasesStore.fetchDatabase(route.params.type, route.params.name);
});

const database = computed(() => databasesStore.database);
const loading = computed(() => databasesStore.loading);

const performDatabaseAction = async (database, action) => {
  await databasesStore.performDatabaseAction(database, action);
};
</script>
