import Docker from 'dockerode';

const docker = new Docker({ socketPath: "/var/run/docker.sock" });
const dockerService = {
  /**
   * Get all containers related to devner
   * 
   * @returns {Promise<Docker.ContainerInfo[]>}
   */
  getContainers() {
    const containers = docker.listContainers({ all: true });
    return containers;
  },
  /**
   * Get logs of a container by its id
   * 
   * @param {string} containerId
   * @param {object} options
   * @returns {Promise<string>}
   */
  getContainerLogs(containerId, options = { follow: true, stdout: true, stderr: true, since: 0 }) {
    const container = docker.getContainer(containerId);
    return container.logs(options);
  },
  /**
   * Get a container by its id
   * 
   * @param {string} containerId
   * @returns {Docker.Container}
   */
  startContainer(containerId) {
    const container = docker.getContainer(containerId);
    return container.start();
  },
  /**
   * Restart a container by its id
   * 
   * @param {string} containerId
   * @returns {Docker.Container}
   */
  restartContainer(containerId) {
    const container = docker.getContainer(containerId);
    return container.restart();
  },
  /**
   * Stop a container by its id
   * 
   * @param {string} containerId
   * @returns {Docker.Container}
   */
  stopContainer(containerId) {
    const container = docker.getContainer(containerId);
    return container.stop();
  },
  /**
   * Stop all containers
   * 
   * @returns {Promise<Docker.Container>}
   * @returns {Docker.Container}
   */
  stopAllContainers() {
    const containers = docker.listContainers({ all: true });
    return containers.map((container) => {
      const c = docker.getContainer(container.Id);
      return c.stop();
    });
  },
  /**
   * Create a new container
   * 
   * @param {object} data
   * @returns {Docker.Container}
   */
  createContainer(data) {
    return docker.createContainer(data);
  },
  /**
   * Remove a container by its id
   * 
   * @param {string} containerId
   * @returns {Docker.Container}
   */
  removeContainer(containerId) {
    const container = docker.getContainer(containerId);
    return container.remove();
  },
  /**
   * Rebuild a container by its id
   * 
   * @param {string} containerId
   * @returns {Docker.Container}
   */
  rebuildContainer(containerId) {
    const container = docker.getContainer(containerId);
    return container.rebuild();
  },
  /**
   * Get a container details by its id
   * 
   * @param {string} containerId
   * @returns {Docker.ContainerInfo}
   */
  getContainerDetail(containerId) {
    const container = docker.getContainer(containerId);
    return container.inspect();
  },
  /**
   * Get a container by its id
   * 
   * @param {string} containerId
   * @returns {Docker.Container}
   */
  getContainer(containerId) {
    const container = docker.getContainer(containerId);
    return container;
  }
};

export { dockerService as d };
//# sourceMappingURL=docker.mjs.map
