import Container from "./Container.js";

class Project {
  constructor(projectData) {
    this.name = projectData.name;
    this.path = projectData.path;
    this.type = projectData.type;
    this.directories = projectData.directories || [];

    this.loading = {
        create: false,
        delete: false,
        code: false,
    };
  }

  /**
   * Get all projects
   *
   * @returns {Array<Project>}
   */

  static async all() {
    const container = await Container.fetchContainerByName("frankenphp_devner");
    const filesAndDirs = await container.cmd("ls -a /var/www/html/*/");
    const projects = this.groupProjects(filesAndDirs);

    return projects;
  }

  /**
   * Get a project by its name
   *
   * @returns {Project}
   */

  static async fetchProject(name) {
    const container = await Container.fetchContainerByName("frankenphp_devner");
    const filesAndDirs = await container.cmd(`ls -a /var/www/html/${name}/`);
    const project = this.groupProject(name, filesAndDirs);

    return project;
  }

   /**
   * Perform an action on the container
   *
   * @param {string} action
   */

   async performAction(action, body = {}) {
    this.loading[action] = true;
    console.log(`Performing action ${action} on project with name: ${this.name}`);
    await $fetch(`/api/project/${this.name}/${action}`, { method: "POST", body });
    this.loading[action] = false;
  }

  async code() {
    await this.performAction("code");
  }

  async delete() {
    await this.performAction("delete");
  }

  async create(name, type) {
    await this.performAction("create", { name, type });
  }

  /**
   * Get project root directory name
   *
   * @returns {string}
   */

  getName() {
    return this.name;
  }

  /**
   * Get project path
   *
   * @returns {string}
   */

  getPath() {
    return this.path;
  }

  /**
   * Get project type
   *
   * @returns {string}
   */

  getType() {
    return this.type;
  }

  /**
   * Group projects by root directory
   *
   * @param {*} filesAndDirs
   * @returns
   */

  static groupProjects(filesAndDirs) {
    const directories = {};
    const projects = [];
    let currentProject = "";

    filesAndDirs.forEach((item) => {
      if (item.endsWith(":")) {
        currentProject = item.replace(":", "");
        currentProject = currentProject.replace("/var/www/html/", "");
        currentProject = currentProject.replace("/", "");
        directories[currentProject] = [];
      } else if (currentProject) {
        directories[currentProject].push(item);
      }
    });

    for (const directory in directories) {
      if (directory === "." || directory === ".." || directory === "gui") {
        continue;
      }
      projects.push(
        new Project({
          name: directory,
          path: `/var/www/html/${directory}`,
          type: this.identifyProject(directories[directory]),
          directories: directories[directory],
        })
      );
    }

    return projects;
  }

  static groupProject(name, filesAndDirs) {
    const directories = filesAndDirs.map(item => item.trim());
    return new Project({
      name: name,
      path: `/var/www/html/${name}`,
      type: this.identifyProject(directories),
      directories: directories,
    });
  }

  /**
   * Identify project type
   *
   * @param {*} filesAndDirs
   * @returns
   */
  static identifyProject(filesAndDirs) {
    if (filesAndDirs.includes("artisan")) {
      return { label: "Laravel", color: "red" };
    } else if (filesAndDirs.includes("wp-config.php")) {
      return { label: "WordPress", color: "blue" };
    } else if (filesAndDirs.includes("package.json")) {
      return { label: "Node.js", color: "green" };
    } else {
      return { label: "Unknown", color: "gray" };
    }
  }

  /**
   * Actions status
   */
  isLoading(action) {
    return this.loading[action];
  }
}

export default Project;
