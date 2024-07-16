import { defineStore } from 'pinia';
import { useProjectsStore } from '~/stores/projects';
import { useDatabasesStore } from '~/stores/databases';
import { useContainersStore } from '~/stores/containers';

export const useLinksStore = defineStore('links', {
  state: () => ({
    lastUpdate: null,
    links: [
      {
        id: "dashboard",
        label: "Dashboard",
        icon: "i-heroicons-rectangle-group",
        to: "/",
      },
    ],
    containersLinks: [
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
    ],
    otherContainersLinks: [
      {
        id: "other-containers",
        label: "Other Containers",
        icon: "i-heroicons-cube-transparent",
        children: [
          {
            id: "other-containers",
            label: "Overview",
            to: "/other-containers",
          },
        ],
      },
    ],
    projectsLinks: [
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
    ],
    databasesLinks: [
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
    ],
  }),
  actions: {
    updateContainersLinks(newContainers) {
      this.containersLinks[0].children = [
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
      ];
    },
    updateProjectsLinks(newProjects) {
      this.projectsLinks[0].children = [
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
      ];
    },
    updateDatabasesLinks(newDatabases) {
      this.databasesLinks[0].children = [
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
      ];
    },
    initialize() {
      // Load lastUpdate from local storage
      const savedLastUpdate = localStorage.getItem('lastUpdate');
      if (savedLastUpdate) {
        this.lastUpdate = new Date(savedLastUpdate);
      }
    },
    async autoUpdate(intervalSeconds) {
      this.updateLastUpdate();
      this.loadAllData();
      const lastUpdate = this.lastUpdate;

      setInterval(() => {
        if (new Date().getTime() - lastUpdate.getTime() > intervalSeconds * 1000) {
          this.loadAllData();
        }
      }, 1000);
    },
    updateLastUpdate() {
      this.lastUpdate = new Date();
      localStorage.setItem('lastUpdate', this.lastUpdate.toISOString());
    },
    async loadAllData() {
      const projectsStore = useProjectsStore();
      const databasesStore = useDatabasesStore();
      const containersStore = useContainersStore();

      await Promise.all([
        projectsStore.reloadProjects(),
        databasesStore.reloadDatabases(),
        containersStore.reloadContainers(),
      ]);

      this.updateLastUpdate();

      /* const toast = useToast();
      toast.add({
        title: 'Data updated', 
        type: 'success',
        color: 'green',
        id: this.lastUpdate.toISOString(),
        timeout: 2000
      }); */
    }
  },
});
