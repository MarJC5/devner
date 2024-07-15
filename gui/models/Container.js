import stripAnsi from 'strip-ansi';

class Container {
  constructor(containerData) {
    this.id = containerData.Id || 0;
    this.created = containerData.Created || "";
    this.name = containerData.Name || "";
    this.status = containerData.State.Status || "";
    this.pid = containerData.State.Pid || 0;
    this.startedAt = containerData.State.StartedAt || "";
    this.image = containerData.Image || "";
    this.ports = containerData.NetworkSettings.Ports || {};
    this.labels = containerData.Config.Labels || {};
    this.network = containerData.NetworkSettings.Networks || {};
    this.project = this.labels["com.docker.compose.project"] || "";

    // Loading state for actions
    this.loading = {
      start: false,
      restart: false,
      stop: false,
      remove: false,
      rebuild: false,
    };
  }

  static async logs(id, options = {}) {
    const logStream = await $fetch(`/api/containers/${id}/logs`, {
      method: "POST",
      body: JSON.stringify(options),
      headers: { "Content-Type": "application/json" },
    });

    logStream.on("data", (chunk) => {
      return chunk.toString('utf8');
    });

    logStream.on("end", () => {
      console.log("Stream ended");
    });

    logStream.on("error", (error) => {
      console.error("Stream error:", error);
    });
  }

  static async fetchContainer(id) {
    const containerData = await $fetch(`/api/containers/${id}/details`, {
      method: "GET",
    });
    return new Container(containerData);
  }

  static async fetchContainerByName(containerName) {
    try {
      const containers = await Container.all();
      return containers.find((container) =>
        container.name.includes(`/${containerName}`)
      );
    } catch (error) {
      console.error("Error fetching container by name:", error);
      throw error;
    }
  }

  static async all() {
    const data = await $fetch(`/api/containers`, { method: "GET" });
    const containers = await Promise.all(
      data.map(async (containerData) => {
        const model = await Container.fetchContainer(containerData.Id);
        if (
          model.getNetworks().devner &&
          !model.getName().toLowerCase().includes("gui")
        ) {
          return model;
        }
      })
    ).then((results) => results.filter((container) => container !== undefined));

    return containers;
  }

  static async allOthers() {
    const data = await $fetch(`/api/containers`, { method: "GET" });
    const containers = await Promise.all(
      data.map(async (containerData) => {
        const model = await Container.fetchContainer(containerData.Id);
        if (!model.getNetworks().devner) {
          return model;
        }
      })
    ).then((results) => results.filter((container) => container !== undefined));
  
    // Group the containers by project name
    const groupedContainers = containers.reduce((acc, container) => {
      const project = container.project;
      if (!acc[project]) {
        acc[project] = [];
      }
      acc[project].push(container);
      return acc;
    }, {});
  
    // Transform the grouped containers into the desired format
    const formattedContainers = Object.entries(groupedContainers).map(([label, containers]) => ({
      label: label.charAt(0).toUpperCase() + label.slice(1),
      icon: "i-heroicons-cube-transparent",
      containers,
    }));
  
    return formattedContainers;
  }  

  formatDate(isoString) {
    const date = new Date(isoString);
    return date.toLocaleString();
  }

  /**
   * Perform an action on the container
   *
   * @param {string} action
   */

  async performAction(action) {
    this.loading[action] = true;
    console.log(`Performing action ${action} on container with ID: ${this.id}`);
    await $fetch(`/api/containers/${this.id}/${action}`, { method: "POST" });
    this.loading[action] = false;
    await this.updateStatus();
  }

  async start() {
    await this.performAction("start");
  }

  async restart() {
    await this.performAction("restart");
  }

  async stop() {
    await this.performAction("stop");
  }

  async remove() {
    await this.performAction("remove");
  }

  async rebuild() {
    await this.performAction("rebuild");
  }

  async updateStatus() {
    const updatedContainer = await Container.fetchContainer(this.id);
    this.status = updatedContainer.status;
    this.pid = updatedContainer.pid;
    this.startedAt = updatedContainer.startedAt;
  }

  /**
   * Exec a command in the container
   */
  async cmd(command, parseAsTable = false, clean = true) {
    try {
      const response = await $fetch(`/api/containers/${this.id}/exec`, {
        method: "POST",
        body: JSON.stringify({ command }),
        headers: { "Content-Type": "application/json" },
      });

      if (response.status === "error") {
        throw new Error(response.message);
      }

      // Split the raw output into lines
      const rawOutput = response.stdout.split("\n");

      // If parsing as table, handle the parsing
      if (parseAsTable) {
        return this.parseTableOutput(rawOutput);
      }

      // Clean up the output by removing non-printable characters if 'clean' is true
      if (clean) {
        return rawOutput
          .map((line) => stripAnsi(line).replace(/[^\x20-\x7E]/g, ''))
          .filter((line) => line);
      }

      return rawOutput;
    } catch (error) {
      console.error("Error executing command in container:", error);
      throw error;
    }
  }

  parseTableOutput(output) {
    const [header, ...rows] = output;

    // Split header into columns
    const headers = header.split("\t").map((h) => h.trim());

    return rows.map((row) => {
      const columns = row.split("\t").map((c) => c.trim()); // Split by tab characters and trim
      let rowObject = {};
      headers.forEach((header, index) => {
        rowObject[stripAnsi(header).replace(/[^\x20-\x7E]/g, '')] = columns[index] || "";
      });
      return rowObject;
    });
  }

  /**
   * Get Directories
   */
  async getDirectories(path = "/") {
    const filesAndDirs = await this.cmd(`ls -a ${path}`);
    return { filesAndDirs, path, count: filesAndDirs.length };
  }

  /**
   * Get id of the container
   *
   * @returns {string}
   */
  getId() {
    return this.id;
  }

  /**
   * Get creation date of the container
   *
   * @returns {string}
   */
  getCreated() {
    return this.formatDate(this.created);
  }

  /**
   * Get name of the container
   *
   * @returns {string}
   */
  getName() {
    let friendlyName = this.name
      .replace("/", "")
      .replace(/[-_]/g, " ")
      .replace(/\b\w/g, (c) => c.toUpperCase());
    friendlyName = friendlyName.replace(/ Devner$/, "");
    return friendlyName;
  }

  /**
   * Get status of the container
   *
   * @returns {string}
   */
  getStatus(format = "capitalize") {
    if (format === "capitalize") {
      return this.status.charAt(0).toUpperCase() + this.status.slice(1);
    }
    return this.status;
  }

  /**
   * Get image of the container
   *
   * @returns {string}
   */
  getImage() {
    return this.image;
  }

  /**
   * Get ports of the container
   *
   * @returns {object}
   */
  getPorts(full = true) {
    const portsArray = [];
    for (const [port, mappings] of Object.entries(this.ports)) {
      if (mappings) {
        mappings.forEach((mapping) => {
          if (full) {
            portsArray.push(`${mapping.HostPort}/${port.split("/")[1]}`);
          } else {
            portsArray.push(mapping.HostPort);
          }
        });
      }
    }
    return portsArray;
  }

  /**
   * Get labels of the container
   *
   * @returns {object}
   */
  getLabels() {
    return this.labels;
  }

  /**
   * Get NetworkSettings of the container
   *
   * @returns {object}
   */
  getNetworks() {
    return this.network;
  }

  /**
   * Is running
   */
  isRunning() {
    return this.status === "running";
  }

  /**
   * Actions status
   */
  isLoading(action) {
    return this.loading[action];
  }
}

export default Container;
