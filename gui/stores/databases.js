import { defineStore } from "pinia";
import Database from "~/utils/Database.js";
import Project from "~/utils/Project.js";

export const useDatabasesStore = defineStore("databases", {
  state: () => ({
    databases: [],
    database: null,
  }),
  getters: {
    mysqlDatabases() {
      return this.databases.filter(
        (database) => database.getType() === "mysql"
      );
    },
    postgresDatabases() {
      return this.databases.filter(
        (database) => database.getType() === "postgres"
      );
    },
  },
  actions: {
    findByName(name) {
      return this.databases.find((database) => database.getName() === name);
    },
    async loadDatabases() {
      if (import.meta.client) {
        const storedDatabases = localStorage.getItem("databases");
        const storedProjects = localStorage.getItem("projects");
        if (storedDatabases) {
          const parsedDatabases = JSON.parse(storedDatabases);
          this.databases = parsedDatabases.map((data) => {
            const database = new Database(data);

            if (storedProjects) {
              const parsedProjects = JSON.parse(storedProjects);
              parsedProjects.map((project) => {
                if (
                  project.database &&
                  project.database.name === database.getName() &&
                  project.database.type === database.getType().label
                ) {
                  database.setProject(new Project(project));
                }
              });
            }

            return database;
          });
        } else {
          try {
            const [mysqlDatabases, postgresDatabases] = await Promise.all([
              Database.all("mysql"),
              Database.all("postgres"),
            ]);
            this.databases = [...mysqlDatabases, ...postgresDatabases];
            localStorage.setItem("databases", JSON.stringify(this.databases));
          } catch (error) {
            console.error("Failed to load databases:", error);
          }
        }
      }
    },
    async fetchDatabase(type, name) {
      this.loading = true;
      try {
        if (!this.databases.length) {
          await this.loadDatabases();
        }
        this.database = this.databases.find((database) => {
          return database.getTypeUrl() === type && database.getName() === name;
        });
      } catch (error) {
        console.error("Failed to load database details:", error);
      } finally {
        this.loading = false;
      }
    },
    async reloadDatabases() {
      try {
        const [mysqlDatabases, postgresDatabases] = await Promise.all([
          Database.all("mysql"),
          Database.all("postgres"),
        ]);
        this.databases = [...mysqlDatabases, ...postgresDatabases];
        localStorage.setItem("databases", JSON.stringify(this.databases));
      } catch (error) {
        console.error("Failed to reload databases:", error);
      }
    },
    async performDatabaseAction(database, action) {
      try {
        await database[action]();
        await this.loadDatabases();
      } catch (error) {
        console.error(`Failed to ${action} database:`, error);
      }
    },
  },
});
