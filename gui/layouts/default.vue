<script setup>
import { ref, onMounted, watch } from "vue";
import Container from "@/models/Container";
import Project from "@/models/Project";
import Database from "@/models/Database";

const containers = ref([]);
const projects = ref([]);
const databases = ref([]);

const containersCollapsed = ref(false);
const projectsCollapsed = ref(false);
const databasesCollapsed = ref(false);

const loadCollapsedState = (key, ref) => {
  const savedState = localStorage.getItem(key);
  ref.value = savedState === "true";
};

const saveCollapsedState = (key, newValue) => {
  localStorage.setItem(key, newValue);
};

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

const projectsLinks = ref([
  {
    label: "Projects",
    icon: "i-heroicons-folder",
    children: [
      {
        id: "projects",
        label: "Overview",
        to: "/projects",
      },
    ],
  },
]);

const databasesLinks = ref([
  {
    label: "Databases",
    icon: "i-heroicons-circle-stack",
    children: [
      {
        id: "databases",
        label: "Overview",
        to: "/databases",
      },
    ],
  },
]);

const loadContainers = async () => {
  containers.value = await Container.all();
  projects.value = await Project.all();
  databases.value = await Database.all("mysql8");
};

watch(containers, (newContainers) => {
  // Watch the containers ref and update containersLinks accordingly
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

watch(projects, (newProjects) => {
  projectsLinks.value = [
    {
      id: "projects",
      label: "Projects",
      icon: "i-heroicons-folder",
      children: [
        {
          id: "projects",
          label: "Overview",
          to: "/projects",
        },
        ...newProjects.map((project) => ({
          id: project.getPath(),
          label: project.getName(),
          to: `/projects/${project.getName()}`,
          badge: {
            color: project.getType().color,
            label: project.getType().label,
          },
        })),
      ],
    },
  ];
});

watch(databases, (newDatabases) => {
  databasesLinks.value = [
    {
      label: "Databases",
      icon: "i-heroicons-circle-stack",
      children: [
        {
          id: "databases",
          label: "Overview",
          to: "/databases",
        },
        ...newDatabases.map((database) => ({
          id: database.getName(),
          label: database.getName(),
          to: `/databases/${database.getContainerName()}/${database.getName()}`,
          badge: {
            color: database.getType().color,
            label: database.getType().label,
          },
        })),
      ],
    },
  ];
});

watch(containersCollapsed, (newValue) => {
  saveCollapsedState("containers-collapsed", newValue);
});

watch(projectsCollapsed, (newValue) => {
  saveCollapsedState("projects-collapsed", newValue);
});

watch(databasesCollapsed, (newValue) => {
  saveCollapsedState("databases-collapsed", newValue);
});

onMounted(() => {
  loadCollapsedState('containers-collapsed', containersCollapsed);
  loadCollapsedState('projects-collapsed', projectsCollapsed);
  loadCollapsedState('databases-collapsed', databasesCollapsed);
  loadContainers();
});
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

        <UNavigationAccordion
          :links="containersLinks"
          :defaultOpen="containersCollapsed"
        />

        <UNavigationAccordion
          :links="projectsLinks"
          :defaultOpen="projectsCollapsed"
        />

        <UNavigationAccordion
          :links="databasesLinks"
          :defaultOpen="databasesCollapsed"
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
