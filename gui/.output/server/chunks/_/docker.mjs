import Docker from 'dockerode';

const docker = new Docker({ socketPath: "/var/run/docker.sock" });
const dockerService = {
  getContainers() {
    return docker.listContainers({ all: true });
  },
  startContainer(containerId) {
    const container = docker.getContainer(containerId);
    return container.start();
  },
  stopContainer(containerId) {
    const container = docker.getContainer(containerId);
    return container.stop();
  },
  createContainer(data) {
    return docker.createContainer(data);
  },
  removeContainer(containerId) {
    const container = docker.getContainer(containerId);
    return container.remove();
  }
  // Add more methods as needed
};

export { dockerService as d };
//# sourceMappingURL=docker.mjs.map
