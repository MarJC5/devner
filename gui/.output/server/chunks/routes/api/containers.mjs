import { d as defineEventHandler } from '../../runtime.mjs';
import { d as dockerService } from '../../_/docker.mjs';
import 'node:http';
import 'node:https';
import 'fs';
import 'path';
import 'node:fs';
import 'node:url';
import 'dockerode';

const containers = defineEventHandler(async (event) => {
  if (event.request.method === "GET") {
    const containers = await dockerService.getContainers();
    return containers;
  }
});

export { containers as default };
//# sourceMappingURL=containers.mjs.map
