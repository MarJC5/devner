import stripAnsi from 'strip-ansi';

class Container {
  constructor(containerData) {
    this.id = containerData.id || containerData.Id || 0;
    this.created = containerData.created || containerData.Created || "";
    this.name = containerData.name || containerData.Name || "";
    this.status = containerData.status || containerData.State.Status || "";
    this.pid = containerData.pid || containerData.State.Pid || 0;
    this.startedAt = containerData.startedAt || containerData.State.StartedAt || "";
    this.image = containerData.image || containerData.Image || "";
    this.ports = containerData.ports || containerData.NetworkSettings.Ports || {};
    this.labels = containerData.labels || containerData.Config.Labels || {};
    this.network = containerData.network || containerData.NetworkSettings.Networks || {};
    this.project = this.labels["com.docker.compose.project"] || "";
    this.loading = {
      start: false,
      restart: false,
      stop: false,
      remove: false,
      rebuild: false
    };
  }
  static async logs(id, options = {}) {
    const logStream = await $fetch(`/api/containers/${id}/logs`, {
      method: "POST",
      body: JSON.stringify(options),
      headers: { "Content-Type": "application/json" }
    });
    logStream.on("data", (chunk) => {
      return chunk.toString("utf8");
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
      method: "GET"
    });
    return new Container(containerData);
  }
  static async fetchContainerByName(containerName) {
    try {
      const containers = await Container.all();
      return containers.find(
        (container) => container.name.includes(`/${containerName}`)
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
        if (model.getNetworks().devner && !model.getName().toLowerCase().includes("gui")) {
          return model;
        }
      })
    ).then((results) => results.filter((container) => container !== void 0));
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
    ).then((results) => results.filter((container) => container !== void 0));
    const groupedContainers = containers.reduce((acc, container) => {
      const project = container.project;
      if (!acc[project]) {
        acc[project] = [];
      }
      acc[project].push(container);
      return acc;
    }, {});
    const formattedContainers = Object.entries(groupedContainers).map(([label, containers2]) => ({
      label: label.charAt(0).toUpperCase() + label.slice(1),
      icon: "i-heroicons-cube-transparent",
      containers: containers2
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
        headers: { "Content-Type": "application/json" }
      });
      if (response.status === "error") {
        throw new Error(response.message);
      }
      const rawOutput = response.stdout.split("\n");
      if (parseAsTable) {
        return this.parseTableOutput(rawOutput);
      }
      if (clean) {
        return rawOutput.map((line) => stripAnsi(line).replace(/[^\x20-\x7E]/g, "")).filter((line) => line);
      }
      return rawOutput;
    } catch (error) {
      console.error("Error executing command in container:", error);
      throw error;
    }
  }
  parseTableOutput(output) {
    const [header, ...rows] = output;
    const headers = header.split("	").map((h) => h.trim());
    return rows.map((row) => {
      const columns = row.split("	").map((c) => c.trim());
      let rowObject = {};
      headers.forEach((header2, index) => {
        rowObject[stripAnsi(header2).replace(/[^\x20-\x7E]/g, "")] = columns[index] || "";
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
    let friendlyName = this.name.replace("/", "").replace(/[-_]/g, " ").replace(/\b\w/g, (c) => c.toUpperCase());
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

class Database {
  constructor(databaseData) {
    this.name = databaseData.name;
    this.type = databaseData.type;
    this.port = databaseData.port;
    this.host = databaseData.host;
    this.containerName = databaseData.containerName;
    this.tables = databaseData.tables || [];
    this.project = databaseData.project || null;
    this.loading = {
      project: false,
      create: false,
      delete: false
    };
  }
  /**
   * Get all databases
   *
   * @returns {Array<Database>}
   */
  static async all(name) {
    try {
      const container = await Container.fetchContainerByName(`${name}_devner`);
      const dbType = this.findType(container.name).label;
      let rawDatabases;
      if (dbType === "MySQL") {
        rawDatabases = await container.cmd(
          `mysql --defaults-file=/.my.cnf -e "SHOW DATABASES;"`,
          true,
          false
        );
        rawDatabases = rawDatabases.map((db) => db.Database.replace(/[^\x20-\x7E]/g, "").trim()).filter((dbName) => !["information_schema", "performance_schema", "mysql", "sys", ""].includes(dbName));
      } else {
        rawDatabases = await container.cmd(
          `psql -U devner -c "\\l"`,
          true,
          false
        );
        const parsedDatabases = this.parsePostgresOutput(rawDatabases);
        rawDatabases = parsedDatabases.map((db) => db.Name).filter((dbName) => !["postgres", "template0", "template1", "devner"].includes(dbName));
      }
      return await this.identifyDatabase(container, rawDatabases, name);
    } catch (error) {
      console.error("Error fetching databases:", error);
      return [];
    }
  }
  /**
   * Parse PostgreSQL command output
   * @param {Array<Object>} output
   * @returns {Array<Object>}
   */
  static parsePostgresOutput(output) {
    const headers = ["Name", "Owner", "Encoding", "Collate", "Ctype", "Access privileges"];
    const rows = [];
    let currentRow = {};
    let isHeaderProcessed = false;
    let separatorLine = false;
    output.forEach((lineObj) => {
      const line = lineObj["List of databases"];
      if (!isHeaderProcessed) {
        if (line.includes("Name") && line.includes("Owner")) {
          isHeaderProcessed = true;
        }
        return;
      }
      if (!separatorLine) {
        if (/^-+\+-+\+-+\+-+\+-+\+-+$/.test(line)) {
          separatorLine = true;
          return;
        }
      }
      if (line.includes("rows)")) {
        return;
      }
      const values = line.split("|").map((value) => value.trim());
      if (values.length === headers.length) {
        if (Object.keys(currentRow).length > 0) {
          rows.push(currentRow);
        }
        currentRow = {};
        headers.forEach((header, index) => {
          currentRow[header] = values[index];
        });
      } else {
        headers.forEach((header, index) => {
          if (values[index]) {
            currentRow[header] = (currentRow[header] || "") + " " + values[index];
          }
        });
      }
    });
    if (Object.keys(currentRow).length > 0) {
      rows.push(currentRow);
    }
    return rows.filter((row) => row.Name && row.Name !== "");
  }
  /**
   * Get a database by its name
   *
   * @returns {Database}
   */
  static async fetchDatabase(containerName, name) {
    const databases = await this.all(containerName);
    return databases.find((database) => database.getName() === name);
  }
  /**
   * Get database name
   *
   * @param {string} name
   */
  getName() {
    return this.name;
  }
  /**
   * Get container name
   *
   * @param {string} name
   */
  getContainerName() {
    return this.containerName;
  }
  /**
   * Get db type
   *
   * @returns {string}
   */
  getType() {
    return this.type.includes("PostgreSQL") ? {
      label: "PostgreSQL",
      color: "cyan"
    } : {
      label: "MySQL",
      color: "orange"
    };
  }
  getTypeUrl() {
    return this.type.includes("PostgreSQL") ? "postgres" : "mysql";
  }
  /**
   * Get project
   *
   * @returns void
   */
  getProject() {
    return this.project;
  }
  /**
   * Set project
   *
   * @returns void
   */
  setProject(project) {
    this.project = project;
  }
  static findType(containerName) {
    return containerName.includes("postgres") ? {
      label: "PostgreSQL",
      color: "cyan"
    } : {
      label: "MySQL",
      color: "orange"
    };
  }
  /**
   * Identify database
   *
   * @param {Container} container
   * @param {Array<string>} rawDatabases
   * @param {string} name
   * @returns {Array<Database>}
   */
  static async identifyDatabase(container, rawDatabases, name) {
    try {
      const databaseObjects = await Promise.all(
        rawDatabases.map(async (dbName) => {
          const dbType = this.findType(`${name}_devner`).label;
          let tables;
          if (dbType === "MySQL") {
            tables = await container.cmd(
              `mysql --defaults-file=/.my.cnf -e "SHOW TABLES;" ${dbName};`
            );
            tables.shift();
          } else {
            tables = await container.cmd(
              `psql -U devner -d ${dbName} -c "\\dt"`
            );
            tables = tables.map((item) => item["List of databases"]).map((line) => line.trim()).filter((line) => line).slice(3, -1).map((table) => table.split("|")[1].trim());
          }
          return new Database({
            name: dbName,
            type: dbType,
            port: container.getPorts(false)[0],
            host: name,
            containerName: `${name}`,
            tables
          });
        })
      );
      return databaseObjects;
    } catch (error) {
      console.error("Error identifying databases:", error);
      return [];
    }
  }
  /**
   * Inspect a table
   *
   * @param {string} tableName
   * @returns {string}
   */
  async inspectTable(tableName) {
    const container = await Container.fetchContainerByName(this.containerName);
    const dbType = this.getType().label;
    let tableDescription;
    if (dbType === "MySQL") {
      tableDescription = await container.cmd(
        `mysql --defaults-file=/.my.cnf -e "DESCRIBE ${tableName}" ${this.name};`,
        true,
        false
      );
    } else {
      tableDescription = await container.cmd(
        `psql -U devner -d ${this.name} -c "\\d ${tableName}"`,
        true,
        false
      );
    }
    return tableDescription;
  }
  /**
   * Create a database
   *
   * @param {string} containerName
   * @param {string} dbName
   */
  static async createDatabase(containerName, dbName) {
    try {
      const container = await Container.fetchContainerByName(`${containerName}_devner`);
      const dbType = this.findType(container.name).label;
      if (dbType === "MySQL") {
        await container.cmd(
          `mysql --defaults-file=/.my.cnf -e "CREATE DATABASE ${dbName};"`
        );
      } else {
        await container.cmd(
          `psql -U devner -c "CREATE DATABASE ${dbName};"`
        );
      }
      return await this.all(containerName);
    } catch (error) {
      console.error("Error creating database:", error);
      return [];
    }
  }
  /**
   * Delete a database
   *
   * @param {string} containerName
   * @param {string} dbName
   */
  static async deleteDatabase(containerName, dbName) {
    try {
      const container = await Container.fetchContainerByName(`${containerName}_devner`);
      const dbType = this.findType(container.name).label;
      if (dbType === "MySQL") {
        await container.cmd(
          `mysql --defaults-file=/.my.cnf -e "DROP DATABASE ${dbName};"`
        );
      } else {
        await container.cmd(
          `psql -U devner -c "DROP DATABASE ${dbName};"`
        );
      }
      return await this.all(containerName);
    } catch (error) {
      console.error("Error deleting database:", error);
      return [];
    }
  }
  /**
   * Create a new table
   *
   * @param {string} tableName
   * @param {string} tableSchema
   */
  async createTable(tableName, tableSchema) {
    const container = await Container.fetchContainerByName(this.containerName);
    const dbType = this.getType().label;
    try {
      if (dbType === "MySQL") {
        await container.cmd(
          `mysql --defaults-file=/.my.cnf -e "CREATE TABLE ${this.name}.${tableName} (${tableSchema});"`
        );
      } else {
        await container.cmd(
          `psql -U devner -d ${this.name} -c "CREATE TABLE ${tableName} (${tableSchema});"`
        );
      }
      return await this.inspectTable(tableName);
    } catch (error) {
      console.error("Error creating table:", error);
      return [];
    }
  }
  /**
   * Delete a table
   *
   * @param {string} tableName
   */
  async deleteTable(tableName) {
    const container = await Container.fetchContainerByName(this.containerName);
    const dbType = this.getType().label;
    try {
      if (dbType === "MySQL") {
        await container.cmd(
          `mysql --defaults-file=/.my.cnf -e "DROP TABLE ${this.name}.${tableName};"`
        );
      } else {
        await container.cmd(
          `psql -U devner -d ${this.name} -c "DROP TABLE ${tableName};"`
        );
      }
      return await this.inspectTable(tableName);
    } catch (error) {
      console.error("Error deleting table:", error);
      return [];
    }
  }
  /**
   * Actions status
   */
  isLoading(action) {
    return this.loading[action];
  }
}

export { Container as C, Database as D };
//# sourceMappingURL=Database.mjs.map
