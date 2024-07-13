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
            variant="soft"/>
          </h1>
      </template>
      <template #description>
        <div class="flex flex-wrap gap-4" v-if="project">
          
        </div>
      </template>
    </UPageHeader>
    <div v-if="project">
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
  import Project from "@/models/Project";
  
  const route = useRoute();
  const project = ref(null);
  
  const loadProject = async () => {
    try {
        project.value = await Project.fetchProject(route.params.name);
    } catch (error) {
      console.error("Failed to load project details:", error);
    }
  };
  
  const performProjectAction = async (project, action) => {
    try {
      await project[action]();
      loadProject();
    } catch (error) {
      console.error(`Failed to ${action} project:`, error);
    }
  };
  
  onMounted(loadProject);
  
  // Watch and update the containers ref
  watch(project, (newProject) => {
    project.value = newProject;
  });
  </script>
  