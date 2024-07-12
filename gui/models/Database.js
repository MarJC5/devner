import Container from "./Container.js";

class Database {
  constructor(databaseData) {
    this.name = databaseData.name;
    this.type = databaseData.type;
    this.port = databaseData.port;
    this.host = databaseData.host;
    this.containerName = databaseData.containerName;
    this.tables = databaseData.tables || [];

    this.loading = {
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
      const rawDatabases = await container.cmd(
        `mysql --defaults-file=/.my.cnf -e "SHOW DATABASES;"`,
        true,
        false
      );

      // Remove BDatabase if information_schema, performance_schema, mysql, and sys or empty
      const cleanedDatabases = rawDatabases.map((db) =>
        db.BDatabase.replace(/[^\x20-\x7E]/g, "").trim()
      );

      // Remove unwanted databases
      const databases = cleanedDatabases.filter(
        (dbName) =>
          ![
            "information_schema",
            "performance_schema",
            "mysql",
            "sys",
            "",
          ].includes(dbName)
      );

      // Filter out warning and empty lines, and map to database names
      return await this.identifyDatabase(container, databases, name);
    } catch (error) {
      console.error("Error fetching databases:", error);
      return [];
    }
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
    if (this.type.includes("postgres")) {
      return {
        label: "PostgreSQL",
        color: "blue",
      };
    }

    return {
      label: "MySQL",
      color: "orange",
    };
  }

  static findType(containerName) {
    if (containerName.includes("postgres")) {
      return {
        label: "PostgreSQL",
        color: "blue",
      };
    }

    return {
      label: "MySQL",
      color: "orange",
    };
  }

  /**
   * idenfy database
   *
   * @param {Container} container
   * @param {Array<string>} rawDatabases
   * @param {string} name
   * @returns {Array<Database>}
   */
  static async identifyDatabase(container, rawDatabases, name) {
    // Filter out warning and empty lines, and map to database names
    try {
      // Use Promise.all to wait for all database objects to be created
      const databaseObjects = await Promise.all(
        rawDatabases.map(async (dbName) => {
          const tables = await container.cmd(
            `mysql --defaults-file=/.my.cnf -e "SHOW TABLES;" ${dbName};`
          );

          // Remove first line of output as it is the table header
          tables.shift();

          return new Database({
            name: dbName,
            type: this.findType(`${name}_devner`).label,
            port: container.getPorts(false)[0],
            host: name,
            containerName: `${name}`,
            tables: tables.map((table) => table.trim()),
          });
        })
      );

      return databaseObjects;
    } catch (error) {
      console.error("Error fetching databases:", error);
      return [];
    }
  }

  /**
   * Inscpect a table
   *
   * @param {string} tableName
   * @returns {string}
   */
  async inspectTable(tableName) {
    const container = await Container.fetchContainerByName(this.containerName);
    const tableDescription = await container.cmd(
      `mysql --defaults-file=/.my.cnf -e "DESCRIBE ${tableName}" ${this.name};`,
      true,
      false
    );
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
      const container = await Container.fetchContainerByName(
        `${containerName}_devner`
      );
      await container.cmd(
        `mysql --defaults-file=/.my.cnf -e "CREATE DATABASE ${dbName};"`
      );
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
      const container = await Container.fetchContainerByName(
        `${containerName}_devner`
      );
      await container.cmd(
        `mysql --defaults-file=/.my.cnf -e "DROP DATABASE ${dbName};"`
      );
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
    try {
      await container.cmd(
        `mysql --defaults-file=/.my.cnf -e "CREATE TABLE ${this.name}.${tableName} (${tableSchema});"`
      );
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
   * @param {string} tableSchema
   */
  async deleteTable(tableName) {
    const container = await Container.fetchContainerByName(this.containerName);
    try {
      await container.cmd(
        `mysql --defaults-file=/.my.cnf -e "DROP TABLE ${this.name}.${tableName};"`
      );
      return await this.inspectTable(tableName);
    } catch (error) {
      console.error("Error deleting table:", error);
      return [];
    }
  }
}

export default Database;
