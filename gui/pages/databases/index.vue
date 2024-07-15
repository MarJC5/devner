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
          :to="`/databases/${database.getType().label.toLowerCase()}/${database.getName()}`"
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
import Database from "@/models/Database";
import Project from "@/models/Project";

const databases = ref([]);
const loading = ref(true);

const loadDatabases = async () => {
  try {
    // Fetch containers, projects, and databases in parallel
    const [mysqlDatabases, postgresDatabases, projects] =
      await Promise.all([
        Database.all("mysql"),
        Database.all("postgres"),
        Project.all(),
      ]);

      // Merge both MySQL & Postgress databases results and sort them by name
      databases.value = [...mysqlDatabases, ...postgresDatabases].sort((a, b) => a.getName().localeCompare(b.getName()));

      // Add projects
      databases.value.forEach((database) => {
        database.setProject(projects.find((project) => project.getName() === database.getName()));
      });

      loading.value = false;
  } catch (error) {
    console.error("Failed to load databases:", error);
  }
};

const performDatabaseAction = async (database, action) => {
  try {
    await database[action]();
    loading.value = true;
    loadDatabases();
    loading.value = false;
  } catch (error) {
    console.error(`Failed to ${action} database:`, error);
  }
};

onMounted(loadDatabases);

// Watch and update the databases ref
watch(databases, (newProjects) => {
  databases.value = newProjects;
});

watch(loading, (newLoading) => {
  loading.value = newLoading;
});
</script>
