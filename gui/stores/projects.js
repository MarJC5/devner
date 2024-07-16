import { defineStore } from 'pinia';
import Project from "~/utils/Project.js";
import Database from "~/utils/Database.js";

export const useProjectsStore = defineStore('projects', {
  state: () => ({
    projects: [],
    project: null,
  }),
  getters: {
    wordpress() {
      return this.projects.find(project => project.type.label === "Wordpress");
    },
    laravel() {
      return this.projects.find(project => project.type.label  === "Laravel");
    },
    nodejs() {
      return this.projects.find(project => project.type.label  === "Node.js");
    },
    unknown() {
      return this.projects.find(project => project.type.label  === "Unknown");
    }
  },
  actions: {
    findByName(name) {
      return this.projects.find(project => project.getName() === name);
    },
    async loadProjects() {
      if (import.meta.client) {
        const storedProjects = localStorage.getItem('projects');
        if (storedProjects) {
          const parsedProjects = JSON.parse(storedProjects);
          this.projects = parsedProjects.map(data => {
            const project = new Project(data);

            if (project.database) {
              project.database = new Database(project.database);
            }

            return project;
          });
        } else {
          try {
            this.projects = await Project.all();
            localStorage.setItem('projects', JSON.stringify(this.projects));
          } catch (error) {
            console.error("Failed to load projects:", error);
          }
        }
      }
    },
    async fetchProject(name) {
      this.loading = true;
      try {
        if (!this.projects.length) {
          await this.loadProjects();
        }
        this.project = this.projects.find(project => project.getName() === name);
      } catch (error) {
        console.error("Failed to load project details:", error);
      } finally {
        this.loading = false;
      }
    },
    async reloadProjects() {
      try {
        this.projects = await Project.all();
        localStorage.setItem('projects', JSON.stringify(this.projects));
      } catch (error) {
        console.error("Failed to reload projects:", error);
      }
    },
    async performProjectAction(project, action) {
      try {
        await project[action]();
        await this.loadProjects();
      } catch (error) {
        console.error(`Failed to ${action} project:`, error);
      }
    }
  },
});
