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
            variant="soft"/>
          </h1>
      </template>
      <template #description>
        <div class="flex flex-wrap gap-4" v-if="database">
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
        </div>
      </template>
    </UPageHeader>
    <div v-if="database">
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
  import Database from "@/models/Database";
  import Project from "@/models/Project";
  
  const route = useRoute();
  const database = ref(null);
  const project = ref(null);
  
  const loadDatabase = async () => {
    try {
      const [db, pjct] = await Promise.all([
          Database.fetchDatabase(route.params.type, route.params.name),
          Project.fetchProject(route.params.name),
      ]);

      database.value = db;
      database.value.setProject(pjct);
    } catch (error) {
      console.error("Failed to load database details:", error);
    }
  };
  
  const performDatabaseAction = async (database, action) => {
    try {
      await database[action]();
      loadDatabase();
    } catch (error) {
      console.error(`Failed to ${action} database:`, error);
    }
  };
  
  onMounted(loadDatabase);
  
  // Watch and update the containers ref
  watch(database, (newDatabase) => {
    database.value = newDatabase;
  });
  </script>
  