class Container {
  constructor(containerData) {
    this.id = containerData.Id;
    this.created = containerData.Created;
    this.name = containerData.Name;
    this.status = containerData.State.Status;
    this.pid = containerData.State.Pid;
    this.startedAt = containerData.State.StartedAt;
    this.image = containerData.Image;
    this.ports = containerData.NetworkSettings.Ports;
    this.labels = containerData.Config.Labels;
    this.network = containerData.NetworkSettings.Networks;

    // Loading state for actions
    this.loading = {
      start: false,
      restart: false,
      stop: false,
      remove: false,
      rebuild: false,
    };
  }

  static async fetchContainer(id) {
    const containerData = await $fetch(`/api/containers/${id}/details`, {
      method: "GET",
    });
    return new Container(containerData);
  }

  static async all() {
    const data = await $fetch(`/api/containers`, { method: "GET" });
    const containers = await Promise.all(
      data.map(async (containerData) => {
        const model = await Container.fetchContainer(containerData.Id);
        if (model.getNetworks().devner && !model.getName().toLowerCase().includes("gui")) {
          return model;
        }
      })
    ).then((results) => results.filter((container) => container !== undefined));
    return containers;
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
  getPorts() {
    const portsArray = [];
    for (const [port, mappings] of Object.entries(this.ports)) {
      if (mappings) {
        mappings.forEach((mapping) => {
          portsArray.push(`${mapping.HostPort}/${port.split("/")[1]}`);
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
