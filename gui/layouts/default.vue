<script setup>
import { ref, onMounted, watch } from 'vue';
import Container from '@/models/Container';

const containers = ref([]);
const links = ref([
  {
    id: "dashboard",
    label: "Dashboard",
    icon: "i-heroicons-rectangle-group",
    to: "/",
  },
]);

const containersLinks = ref([
  {
    label: "Containers",
    icon: "i-heroicons-cube",
    children: [
      {
        id: "containers",
        label: "Overview",
        to: "/containers",
      },
    ],
  },
]);

const projectsLinks = [
  {
    id: "projects",
    label: "Projects",
    icon: "i-heroicons-folder",
    to: "/projects",
  },
];

const databasesLinks = [
  {
    id: "databases",
    label: "Databases",
    icon: "i-heroicons-circle-stack",
    to: "/databases",
  },
];

const loadContainers = async () => {
  containers.value = await Container.all();
};

// Watch the containers ref and update containersLinks accordingly
watch(containers, (newContainers) => {
  containersLinks.value = [
    {
      label: "Containers",
      icon: "i-heroicons-cube",
      children: [
        {
          id: "containers",
          label: "Overview",
          to: "/containers",
        },
        ...newContainers.map((container) => ({
          id: container.getId(),
          label: container.getName(),
          to: `/containers/${container.getId()}`,
          badge: {
            color: container.isRunning() ? "green" : "red",
            label: container.getStatus(),
          },
        })),
      ],
    },
  ];
});

onMounted(loadContainers);
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
  
          <UNavigationAccordion :links="containersLinks" />
  
          <UNavigationAccordion :links="projectsLinks" />
  
          <UNavigationAccordion :links="databasesLinks" />
        </UDashboardSidebar>
      </UDashboardPanel>
  
      <UDashboardPage>
        <UDashboardPanel grow>
          <UDashboardPanelContent>
            <slot />
          </UDashboardPanelContent>
        </UDashboardPanel>
      </UDashboardPage>
    </UDashboardLayout>
  </template>
  