<script setup>
import { useContainersStore } from '~/stores/containers';
import { useProjectsStore } from '~/stores/projects';
import { useDatabasesStore } from '~/stores/databases';
import { useLinksStore } from '~/stores/links';

const containersStore = useContainersStore();
const projectsStore = useProjectsStore();
const databasesStore = useDatabasesStore();
const linksStore = useLinksStore();

onMounted(() => {
  // initialize links store
  linksStore.initialize();

  // load containers, projects, and databases
  containersStore.loadContainers().then(() => {
    linksStore.updateContainersLinks(containersStore.containers);
  });
  projectsStore.loadProjects().then(() => {
    linksStore.updateProjectsLinks(projectsStore.projects);
  });
  databasesStore.loadDatabases().then(() => {
    linksStore.updateDatabasesLinks(databasesStore.databases);
  });

  // auto update data every x minute
  linksStore.autoUpdate(60 * 3);
});

const links = computed(() => linksStore.links);
const containersLinks = computed(() => linksStore.containersLinks);
const otherContainersLinks = computed(() => linksStore.otherContainersLinks);
const projectsLinks = computed(() => linksStore.projectsLinks);
const databasesLinks = computed(() => linksStore.databasesLinks);
</script>


<template>
  <UDashboardLayout>
    <UDashboardPanel
      :width="250"
      :resizable="{ min: 200, max: 300 }"
      collapsible
    >
      <UDashboardNavbar class="!border-transparent" :ui="{ left: 'flex-1' }">
        <template #left>
          <UButton
            to="/"
            icon="i-heroicons-square-3-stack-3d"
            size="xl"
            color="white"
            square
            variant="ghost"
            label="Devner"
            class="pl-0 uppercase"
          />
        </template>
      </UDashboardNavbar>

      <UDashboardSidebar>
        <UNavigationLinks :links="links" />

        <UDivider class="my-4" />

        <UNavigationAccordion
          :links="containersLinks"
        />

        <UNavigationAccordion
          :links="projectsLinks"
        />

        <UNavigationAccordion
          :links="databasesLinks"
        />

        <div class="flex-1" />

        <UNavigationAccordion
          :links="otherContainersLinks"
        />
      </UDashboardSidebar>
    </UDashboardPanel>

    <UDashboardPage>
      <UDashboardPanel grow>
        <UDashboardPanelContent class="pt-0">
          <slot />
        </UDashboardPanelContent>
      </UDashboardPanel>
    </UDashboardPage>
  </UDashboardLayout>
</template>
