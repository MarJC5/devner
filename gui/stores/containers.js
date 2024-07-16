import { defineStore } from 'pinia';
import Container from "~/utils/Container.js";

export const useContainersStore = defineStore('containers', {
  state: () => ({
    containers: [],
    container: null,
  }),
  getters: {
    runningContainers() {
      return this.containers.filter(container => container.getStatus() === "running");
    },
    stoppedContainers() {
      console.log("Stopped containers:", this.containers.filter(container => container.getStatus() === "stopped"));
      return this.containers.filter(container => container.getStatus() === "stopped");
    },
  },
  actions: {
    findByName(name) {
      return this.containers.find(container => container.getName() === name);
    },
    async loadContainers() {
      if (import.meta.client) {
        const storedContainers = localStorage.getItem('containers');
        if (storedContainers) {
          const parsedContainers = JSON.parse(storedContainers);
          this.containers = parsedContainers.map(data => new Container(data));
        } else {
          try {
            this.containers = await Container.all();
            localStorage.setItem('containers', JSON.stringify(this.containers));
          } catch (error) {
            console.error("Failed to load containers:", error);
          }
        }
      }
    },
    async fetchContainer(id) {
      this.loading = true;
      try {
        this.container = await Container.fetchContainer(id);
      } catch (error) {
        console.error("Failed to load container details:", error);
      } finally {
        this.loading = false;
      }
    },
    async reloadContainers() {
      try {
        this.containers = await Container.all();
        localStorage.setItem('containers', JSON.stringify(this.containers));
      } catch (error) {
        console.error("Failed to reload containers:", error);
      }
    },
    async performContainerAction(container, action) {
      try {
        await container[action]();
        await this.loadContainers();
      } catch (error) {
        console.error(`Failed to ${action} container:`, error);
      }
    }
  }
});
