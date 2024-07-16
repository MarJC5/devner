import Container from "./Container.js";

export default class Database {
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
      delete: false,
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

      if (dbType === 'MySQL') {
        rawDatabases = await container.cmd(
          `mysql --defaults-file=/.my.cnf -e "SHOW DATABASES;"`,
          true,
          false
        );
        // Remove unwanted databases and clean the output
        rawDatabases = rawDatabases
          .map((db) => db.Database.replace(/[^\x20-\x7E]/g, "").trim())
          .filter((dbName) => !["information_schema", "performance_schema", "mysql", "sys", ""].includes(dbName));
      } else {
        rawDatabases = await container.cmd(
          `psql -U devner -c "\\l"`,
          true,
          false
        );

        // Parse PostgreSQL output
        const parsedDatabases = this.parsePostgresOutput(rawDatabases);

        // Filter out template databases
        rawDatabases = parsedDatabases
          .map((db) => db.Name)
          .filter((dbName) => !["postgres", "template0", "template1", "devner"].includes(dbName));
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
            currentRow[header] = (currentRow[header] || '') + " " + values[index];
          }
        });
      }
    });

    if (Object.keys(currentRow).length > 0) {
      rows.push(currentRow);
    }

    return rows.filter(row => row.Name && row.Name !== "");
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
      color: "cyan",
    } : {
      label: "MySQL",
      color: "orange",
    };
  }

  getTypeUrl() {
    return this.type.includes("PostgreSQL") ? 'postgres' : 'mysql';
  }

  /**
   * Get project
   *
   * @returns void
   */
  getProject() {
    return this.project;
  };

  /**
   * Set project
   *
   * @returns void
   */
  setProject(project) {
    this.project = project;
  };

  static findType(containerName) {
    return containerName.includes("postgres") ? {
      label: "PostgreSQL",
      color: "cyan",
    } : {
      label: "MySQL",
      color: "orange",
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

          if (dbType === 'MySQL') {
            tables = await container.cmd(
              `mysql --defaults-file=/.my.cnf -e "SHOW TABLES;" ${dbName};`
            );
            tables.shift();
          } else {
            tables = await container.cmd(
              `psql -U devner -d ${dbName} -c "\\dt"`
            );
            tables = tables.map(item => item["List of databases"]).map(line => line.trim()).filter(line => line).slice(3, -1).map(table => table.split('|')[1].trim());
          }

          return new Database({
            name: dbName,
            type: dbType,
            port: container.getPorts(false)[0],
            host: name,
            containerName: `${name}`,
            tables: tables,
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
    
    if (dbType === 'MySQL') {
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

      if (dbType === 'MySQL') {
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

      if (dbType === 'MySQL') {
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
      if (dbType === 'MySQL') {
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
      if (dbType === 'MySQL') {
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
