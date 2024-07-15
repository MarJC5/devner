<template>
  <UPageHeader
    title="Projects"
    class="sticky top-0 bg-white z-10 dark:bg-gray-800 mb-4 pb-4"
  >
    <template #description>
      <div class="flex flex-wrap gap-2">
        <UButton
          @click="() => console.log('New Project')"
          icon="i-heroicons-plus"
          size="sm"
          color="white"
          square
          variant="solid"
          label="New Project"
        />
      </div>
    </template>
  </UPageHeader>
  <UPageGrid v-if="projects && !loading">
    <UPageCard v-for="project in projects" :key="project.getName()">
      <template #header>
        {{ project.getName() }} -
        <UBadge
          :label="project.getType().label"
          :color="project.getType().color"
          class="ml-1 h-6"
          variant="soft"
        />
      </template>
      <div class="flex flex-wrap gap-2">
        <UButton
          :to="`/projects/${project.getName()}`"
          icon="i-heroicons-eye"
          size="sm"
          color="white"
          square
          variant="solid"
          label="Details"
        />
        <UButton
          @click="() => performProjectAction(project, 'code')"
          :icon="project.isLoading('code') ? 'i-heroicons-spinner' : 'i-heroicons-code-bracket-square'"
          size="sm"
          color="white"
          square
          variant="solid"
          :loading="project.isLoading('code')"
          label="VS Code"
        />
        <UButton
          v-if="project.getDatabase()"
          :to="`/databases/${project.getDatabase().getType().label.toLowerCase()}/${project.getDatabase().getName()}`"
          icon="i-heroicons-circle-stack"
          size="sm"
          color="white"
          square
          variant="solid"
          label="Database"
        />
        <UButton
          @click="() => performProjectAction(project, 'remove')"
          :icon="project.isLoading('remove') ? 'i-heroicons-spinner' : 'i-heroicons-trash'"
          size="sm"
          color="white"
          square
          variant="solid"
          :loading="project.isLoading('remove')"
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
        Loading projects...
    </UButton>
  </div>
</template>

<script setup>
import Project from "@/models/Project";

const projects = ref([]);
const loading = ref(true);

const loadProjects = async () => {
  try {
    projects.value = await Project.all();
    loading.value = false;
  } catch (error) {
    console.error("Failed to load projects:", error);
  }
};

const performProjectAction = async (project, action) => {
  try {
    await project[action]();
    loading.value = true;
    loadProjects();
    loading.value = false;
  } catch (error) {
    console.error(`Failed to ${action} project:`, error);
  }
};

onMounted(loadProjects);

// Watch and update the projects ref
watch(projects, (newProjects) => {
  projects.value = newProjects;
});

watch(loading, (newLoading) => {
  loading.value = newLoading;
});
</script>
