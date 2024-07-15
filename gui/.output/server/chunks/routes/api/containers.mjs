import { d as defineEventHandler } from '../../runtime.mjs';
import { d as dockerService } from '../../_/docker.mjs';
import 'node:http';
import 'node:https';
import 'events';
import 'https';
import 'http';
import 'net';
import 'tls';
import 'crypto';
import 'stream';
import 'url';
import 'zlib';
import 'buffer';
import 'fs';
import 'path';
import 'engine.io';
import 'socket.io';
import 'dockerode';
import 'node:fs';
import 'node:url';

const containers = defineEventHandler(async (event) => {
  if (event.node.req.method === "GET") {
    const containers = await dockerService.getContainers();
    return containers;
  }
});

export { containers as default };
//# sourceMappingURL=containers.mjs.map
