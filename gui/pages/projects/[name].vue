<template>
  <UPageHeader
    :headline="project ? '' : 'Loading...'"
    class="sticky top-0 bg-white z-10 dark:bg-gray-800 mb-4 pb-4"
  >
    <template #title>
      <h1 v-if="project" class="flex">{{ project.getName() }}
        <UBadge 
          :label="project.getType().label"
          :color="project.getType().color"
          class="ml-2 h-6"
          variant="soft"
        />
      </h1>
    </template>
    <template #description>
      <div class="flex flex-wrap gap-2" v-if="project">
        <UButton
          v-if="project.getDatabase()"
          :to="`/databases/${project.getDatabase().getTypeUrl()}/${project.database.getName()}`"
          icon="i-heroicons-circle-stack"
          size="sm"
          color="white"
          square
          variant="solid"
          label="Database"
        />
      </div>
    </template>
  </UPageHeader>
  <div v-if="project">
    <!-- Add more project-related content here -->
  </div>
  <div v-else>
    <UButton 
      size="xl"
      color="white"
      square
      block
      variant="ghost"
      loading>
      Loading project...
    </UButton>
  </div>
</template>

<script setup>
import { useProjectsStore } from '~/stores/projects';

const route = useRoute();
const projectsStore = useProjectsStore();

onMounted(() => {
  projectsStore.fetchProject(route.params.name);
});

const project = computed(() => projectsStore.project);
const loading = computed(() => projectsStore.loading);

const performProjectAction = async (project, action) => {
  await projectsStore.performProjectAction(project, action);
};
</script>
